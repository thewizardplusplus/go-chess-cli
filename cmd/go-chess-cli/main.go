package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/thewizardplusplus/go-chess-cli/encoding/ascii"
	"github.com/thewizardplusplus/go-chess-cli/encoding/unicode"
	climodels "github.com/thewizardplusplus/go-chess-cli/models"
	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

// nolint: gochecknoglobals
var (
	baseWideMargins = ascii.Margins{
		Legend: ascii.LegendMargins{
			File: ascii.VerticalMargins{
				Top: 1,
			},
			Rank: ascii.HorizontalMargins{
				Right: 1,
			},
		},
		Board: ascii.VerticalMargins{
			Top:    1,
			Bottom: 1,
		},
	}
	widePieceMargins = ascii.PieceMargins{
		HorizontalMargins: ascii.HorizontalMargins{
			Left: 1,
		},
		VerticalMargins: ascii.VerticalMargins{
			Bottom: 1,
		},
	}
	extraWidePieceMargins = ascii.PieceMargins{
		HorizontalMargins: ascii.HorizontalMargins{
			Left:  1,
			Right: 1,
		},
		VerticalMargins: ascii.VerticalMargins{
			Top:    1,
			Bottom: 1,
		},
	}
)

type colorCodeGroup map[models.Color]int

func setTTYMode(mode int) string {
	return fmt.Sprintf("\x1b[%dm", mode)
}

func makeColorizer(colorsCodes colorCodeGroup) ascii.Colorizer {
	return func(text string, color models.Color) string {
		return setTTYMode(colorsCodes[color]) + text + setTTYMode(0)
	}
}

func search(
	cache caches.Cache,
	storage models.PieceStorage,
	color models.Color,
	terminator terminators.SearchTerminator,
) (moves.ScoredMove, error) {
	searcher := minimax.NewParallelSearcher(
		terminator,
		runtime.NumCPU(),
		func() minimax.MoveSearcher {
			innerSearcher := minimax.NewAlphaBetaSearcher(
				models.MoveGenerator{},
				nil, // terminator will be set automatically by the iterative searcher
				evaluators.MaterialEvaluator{},
			)

			if cache != nil {
				// make and bind a cached searcher to inner one
				minimax.NewCachedSearcher(innerSearcher, cache)
			}

			return minimax.NewIterativeSearcher(
				innerSearcher,
				nil, // terminator will be set automatically by the parallel searcher
			)
		},
	)

	return searcher.SearchMove(
		storage,
		color,
		0, // initial deep
		moves.NewBounds(),
	)
}

func check(storage models.PieceStorage, color models.Color) error {
	// minimal deep, at which a game state will be detected
	terminator := terminators.NewDeepTerminator(1)
	_, err := search(
		nil, // without a cache
		storage,
		color,
		terminator,
	)
	return err // don't wrap
}

func writePrompt(
	storageEncoder ascii.PieceStorageEncoder,
	storage models.PieceStorage,
	color models.Color,
	side climodels.Side,
) error {
	text := storageEncoder.EncodePieceStorage(storage)
	fmt.Println(text)

	if err := check(storage, color); err != nil {
		return err // don't wrap
	}

	var mark string
	if side == climodels.Searcher {
		mark = "(searching) "
	}

	text = ascii.EncodeColor(color)
	fmt.Printf("%s> %s", text, mark) // don't break the line

	return nil
}

func readMove(
	reader *bufio.Reader,
	storageEncoder ascii.PieceStorageEncoder,
	storage models.PieceStorage,
	color models.Color,
	side climodels.Side,
) (models.Move, error) {
	if err := writePrompt(storageEncoder, storage, color, side); err != nil {
		return models.Move{}, err // don't wrap
	}

	text, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return models.Move{}, fmt.Errorf("unable to read the move: %s", err)
	}

	text = strings.TrimSuffix(text, "\n")
	move, err := uci.DecodeMove(text)
	if err != nil {
		return models.Move{}, fmt.Errorf("unable to decode the move: %s", err)
	}

	if err = storage.CheckMove(move); err != nil {
		return models.Move{}, fmt.Errorf("incorrect move: %s", err)
	}

	if piece, _ := storage.Piece(move.Start); piece.Color() != color {
		return models.Move{}, errors.New("incorrect move: opponent piece")
	}

	nextStorage := storage.ApplyMove(move)
	nextColor := color.Negative()
	if err = check(nextStorage, nextColor); err == models.ErrKingCapture {
		return models.Move{}, errors.New("incorrect move: check")
	}

	return move, nil
}

func searchMove(
	cache caches.Cache,
	storageEncoder ascii.PieceStorageEncoder,
	storage models.PieceStorage,
	color models.Color,
	side climodels.Side,
	deep int,
	duration time.Duration,
) (models.Move, error) {
	if err := writePrompt(storageEncoder, storage, color, side); err != nil {
		return models.Move{}, err // don't wrap
	}

	terminator := terminators.NewGroupTerminator(
		terminators.NewDeepTerminator(deep),
		terminators.NewTimeTerminator(time.Now, duration),
	)
	move, _ := search(cache, storage, color, terminator) // nolint: gosec
	return move.Move, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fen := flag.String(
		"fen",
		"rnbqk/ppppp/5/PPPPP/RNBQK",
		"board in FEN (default: Gardner's minichess)",
	)
	humanColor := flag.String(
		"humanColor",
		"random",
		"human color (allowed: random, black, white)",
	)
	deep := flag.Int("deep", 5, "search deep")
	duration := flag.Duration(
		"duration",
		5*time.Second,
		"search duration (e.g. 72h3m0.5s)",
	)
	cacheSize := flag.Int("cacheSize", 1e6, "maximal cache size (in items)")
	useUnicode := flag.Bool("unicode", true, "use Unicode to display pieces")
	colorfulPieces := flag.Bool(
		"colorfulPieces",
		true,
		"use colors to display pieces",
	)
	pieceBlackColor := flag.Int(
		"pieceBlackColor",
		34, // blue
		"SGR parameter for ANSI escape sequences for setting a color of black pieces",
	)
	pieceWhiteColor := flag.Int(
		"pieceWhiteColor",
		31, // red
		"SGR parameter for ANSI escape sequences for setting a color of white pieces",
	)
	colorfulBoard := flag.Bool(
		"colorfulBoard",
		true,
		"use colors to display the board",
	)
	squareBlackColor := flag.Int(
		"squareBlackColor",
		40, // black
		"SGR parameter for ANSI escape sequences "+
			"for setting a color of black squares",
	)
	squareWhiteColor := flag.Int(
		"squareWhiteColor",
		47, // white
		"SGR parameter for ANSI escape sequences "+
			"for setting a color of white squares",
	)
	wide := flag.Bool("wide", true, "display the board wide")
	flag.Parse()

	storage, err := uci.DecodePieceStorage(*fen, pieces.NewPiece, models.NewBoard)
	if err != nil {
		log.Fatal("unable to decode the board: ", err)
	}

	parsedHumanColor, err := ascii.DecodeColor(*humanColor)
	switch {
	case err == nil:
	case *humanColor == "random":
		if rand.Intn(2) == 0 {
			parsedHumanColor = models.Black
		} else {
			parsedHumanColor = models.White
		}
	default:
		log.Fatal("unable to decode the color: ", err)
	}

	var pieceEncoder ascii.PieceEncoder
	var placeholder string
	if *useUnicode {
		pieceEncoder = unicode.EncodePiece
		placeholder = "\u00b7"
	} else {
		pieceEncoder = uci.EncodePiece
		placeholder = "."
	}
	if *colorfulPieces {
		pieceColorizer := makeColorizer(colorCodeGroup{
			models.Black: *pieceBlackColor,
			models.White: *pieceWhiteColor,
		})
		basePieceEncoder := pieceEncoder
		pieceEncoder = func(piece models.Piece) string {
			text := basePieceEncoder(piece)
			return pieceColorizer(text, piece.Color())
		}
	}
	if *colorfulBoard {
		placeholder = " "
	}

	var margins ascii.Margins
	if *wide {
		margins = baseWideMargins

		if *colorfulBoard {
			margins.Piece = extraWidePieceMargins
		} else {
			margins.Piece = widePieceMargins
		}
	}

	var squareColorizer ascii.OptionalColorizer
	if *colorfulBoard {
		baseSquareColorizer := makeColorizer(colorCodeGroup{
			models.Black: *squareBlackColor,
			models.White: *squareWhiteColor,
		})
		squareColorizer = ascii.NewOptionalColorizer(baseSquareColorizer)
	} else {
		squareColorizer = ascii.WithoutColor
	}

	side := climodels.NewSide(parsedHumanColor)
	reader := bufio.NewReader(os.Stdin)
	storageEncoder := ascii.NewPieceStorageEncoder(
		pieceEncoder,
		placeholder,
		margins,
		squareColorizer,
		parsedHumanColor.Negative(),
		1,
	)
	cache := caches.NewParallelCache(caches.NewStringHashingCache(
		*cacheSize,
		uci.EncodePieceStorage,
	))
loop:
	for {
		var move models.Move
		var err error
		switch side {
		case climodels.Human:
			move, err = readMove(reader, storageEncoder, storage, parsedHumanColor, side)
		case climodels.Searcher:
			move, err = searchMove(
				cache,
				storageEncoder,
				storage,
				parsedHumanColor.Negative(),
				side,
				*deep,
				*duration,
			)
			if err == nil {
				text := uci.EncodeMove(move)
				fmt.Println(text)
			}
		}
		switch err {
		case nil:
		case minimax.ErrCheckmate, minimax.ErrDraw:
			log.Print("game in the state: ", err)
			break loop
		default:
			log.Print("error: ", err)
			continue loop
		}

		storage = storage.ApplyMove(move)
		side = side.Invert()
	}
}
