package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/thewizardplusplus/go-chess-cli/encoding/ascii"
	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

type side int

const (
	searcher side = iota
	human
)

func (side side) invert() side {
	if side == searcher {
		return human
	}

	return searcher
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
			innerSearcher :=
				minimax.NewAlphaBetaSearcher(
					models.MoveGenerator{},
					// terminator will be set
					// automatically
					// by the iterative searcher
					nil,
					evaluators.MaterialEvaluator{},
				)

			if cache != nil {
				// make and bind a cached searcher
				// to inner one
				minimax.NewCachedSearcher(
					innerSearcher,
					cache,
				)
			}

			return minimax.NewIterativeSearcher(
				innerSearcher,
				// terminator will be set
				// automatically
				// by the parallel searcher
				nil,
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

func check(
	storage models.PieceStorage,
	color models.Color,
) error {
	// minimal deep, at which a game state
	// will be detected
	terminator :=
		terminators.NewDeepTerminator(1)
	_, err := search(
		nil, // without a cache
		storage,
		color,
		terminator,
	)
	return err // don't wrap
}

func writePrompt(
	storage models.PieceStorage,
	color models.Color,
) error {
	encoder := ascii.NewPieceStorageEncoder(
		uci.EncodePiece,
		".",
		models.White,
	)
	text :=
		encoder.EncodePieceStorage(storage)
	fmt.Println(text)

	err := check(storage, color)
	if err != nil {
		return err // don't wrap
	}

	text = ascii.EncodeColor(color)
	// don't break the line
	fmt.Print(text + "> ")

	return nil
}

func readMove(
	reader *bufio.Reader,
	storage models.PieceStorage,
	color models.Color,
) (models.Move, error) {
	err := writePrompt(storage, color)
	if err != nil {
		return models.Move{}, err // don't wrap
	}

	text, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return models.Move{}, fmt.Errorf(
			"unable to read the move: %s",
			err,
		)
	}

	text = strings.TrimSuffix(text, "\n")
	move, err := uci.DecodeMove(text)
	if err != nil {
		return models.Move{}, fmt.Errorf(
			"unable to decode the move: %s",
			err,
		)
	}

	err = storage.CheckMove(move)
	if err != nil {
		return models.Move{}, fmt.Errorf(
			"incorrect move: %s",
			err,
		)
	}

	piece, _ := storage.Piece(move.Start)
	if piece.Color() != color {
		return models.Move{}, errors.New(
			"incorrect move: opponent piece",
		)
	}

	nextStorage := storage.ApplyMove(move)
	nextColor := color.Negative()
	err = check(nextStorage, nextColor)
	if err == models.ErrKingCapture {
		return models.Move{}, errors.New(
			"incorrect move: check",
		)
	}

	return move, nil
}

func searchMove(
	cache caches.Cache,
	storage models.PieceStorage,
	color models.Color,
	duration time.Duration,
) (models.Move, error) {
	err := writePrompt(storage, color)
	if err != nil {
		return models.Move{}, err // don't wrap
	}

	terminator :=
		terminators.NewTimeTerminator(
			time.Now,
			duration,
		)
	move, _ := search(
		cache,
		storage,
		color,
		terminator,
	)
	return move.Move, nil
}

func main() {
	fen := flag.String(
		"fen",
		"rnbqk/ppppp/5/PPPPP/RNBQK",
		"board in FEN "+
			"(default: Gardner's minichess)",
	)
	color := flag.String(
		"color",
		"white",
		"human color (allowed: black, white)",
	)
	duration := flag.Duration(
		"duration",
		5*time.Second,
		"search duration (e.g. 72h3m0.5s)",
	)
	cacheSize := flag.Int(
		"cacheSize",
		1e6,
		"maximal cache size (in items)",
	)
	flag.Parse()

	storage, err := uci.DecodePieceStorage(
		*fen,
		pieces.NewPiece,
		models.NewBoard,
	)
	if err != nil {
		log.Fatal(
			"unable to decode the board: ",
			err,
		)
	}

	parsedColor, err :=
		ascii.DecodeColor(*color)
	if err != nil {
		log.Fatal(
			"unable to decode the color: ",
			err,
		)
	}

	var side side
	// detect an initial side
	switch parsedColor {
	case models.Black:
		side = searcher
	case models.White:
		side = human
	}

	reader := bufio.NewReader(os.Stdin)
	cache := caches.NewParallelCache(
		caches.NewStringHashingCache(
			*cacheSize,
			uci.EncodePieceStorage,
		),
	)
loop:
	for {
		var move models.Move
		var err error
		switch side {
		case human:
			move, err = readMove(
				reader,
				storage,
				parsedColor,
			)
		case searcher:
			move, err = searchMove(
				cache,
				storage,
				parsedColor.Negative(),
				*duration,
			)
			if err == nil {
				text := uci.EncodeMove(move)
				fmt.Println(text)
			}
		}
		switch err {
		case nil:
		case minimax.ErrCheckmate,
			minimax.ErrDraw:
			log.Print("game in the state: ", err)
			break loop
		default:
			log.Print("error: ", err)
			continue loop
		}

		storage = storage.ApplyMove(move)
		side = side.invert()
	}
}
