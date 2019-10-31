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
	"strconv"
	"strings"
	"time"

	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	moves "github.com/thewizardplusplus/go-chess-minimax/models"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func decodeColor(text string) (
	models.Color,
	error,
) {
	var color models.Color
	switch text {
	case "black":
		color = models.Black
	case "white":
		color = models.White
	default:
		return 0, errors.New("incorrect color")
	}

	return color, nil
}

func encodeColor(
	color models.Color,
) string {
	var text string
	switch color {
	case models.Black:
		text = "black"
	case models.White:
		text = "white"
	}

	return text
}

func encodeStorage(
	storage models.PieceStorage,
) string {
	var ranks []string
	var currentRank string
	positions := storage.Size().Positions()
	for _, position := range positions {
		if len(currentRank) == 0 {
			currentRank +=
				strconv.Itoa(position.Rank + 1)
		}

		piece, ok := storage.Piece(position)
		if ok {
			currentRank += uci.EncodePiece(piece)
		} else {
			currentRank += "x"
		}

		lastFile := storage.Size().Height - 1
		if position.File == lastFile {
			ranks = append(ranks, currentRank)
			currentRank = ""
		}
	}

	legendRank := " "
	width := storage.Size().Width
	for i := 0; i < width; i++ {
		legendRank += string(i + 97)
	}
	ranks = append(ranks, legendRank)

	return strings.Join(ranks, "\n")
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
	text := encodeStorage(storage)
	fmt.Println(text)

	err := check(storage, color)
	if err != nil {
		return err // don't wrap
	}

	text = encodeColor(color)
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
		// "rnbqkbnr/pppppppp/8/8"+
		// "/8/8/PPPPPPPP/RNBQKBNR",
		"board in FEN",
	)
	color := flag.String(
		"color",
		"white",
		"human color (allowed: black, white)",
	)
	duration := flag.Duration(
		"duration",
		5*time.Second,
		"search duration",
	)
	cacheSize := flag.Int(
		"cacheSize",
		1e6,
		"maximal cache size",
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

	parsedColor, err := decodeColor(*color)
	if err != nil {
		log.Fatal(
			"unable to decode the color: ",
			err,
		)
	}

	var side string
	switch parsedColor {
	case models.Black:
		side = "searcher"
	case models.White:
		side = "human"
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
		case "human":
			move, err = readMove(
				reader,
				storage,
				parsedColor,
			)
		case "searcher":
			move, err = searchMove(
				cache,
				storage,
				parsedColor.Negative(),
				*duration,
			)
		}
		switch err {
		case nil:
		case minimax.ErrCheckmate,
			minimax.ErrDraw:
			fmt.Printf(
				"game in state: %s\n",
				err,
			)
			break loop
		default:
			fmt.Printf("error: %s\n", err)
			continue loop
		}

		storage = storage.ApplyMove(move)
		switch side {
		case "human":
			side = "searcher"
		case "searcher":
			text := uci.EncodeMove(move)
			fmt.Println(text)

			side = "human"
		}
	}
}
