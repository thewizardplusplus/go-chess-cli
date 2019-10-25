package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	cli "github.com/thewizardplusplus/go-chess-cli"
	minimax "github.com/thewizardplusplus/go-chess-minimax"
	"github.com/thewizardplusplus/go-chess-minimax/caches"
	"github.com/thewizardplusplus/go-chess-minimax/evaluators"
	"github.com/thewizardplusplus/go-chess-minimax/terminators"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/games"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

var (
	generator models.MoveGenerator
	evaluator evaluators.MaterialEvaluator
)

func newChecker() minimax.MoveSearcher {
	// limit to the minimum value
	// at which all checks are performed
	terminator :=
		terminators.NewDeepTerminator(1)

	return minimax.NewAlphaBetaSearcher(
		generator,
		terminator,
		evaluator,
	)
}

func newSearcher(
	maxCacheSize int,
) minimax.MoveSearcher {
	cache := caches.NewParallelCache(
		caches.NewStringHashingCache(
			maxCacheSize,
			uci.EncodePieceStorage,
		),
	)

	return minimax.NewParallelSearcher(
		// terminator will be set
		// before each search
		nil,
		runtime.NumCPU(),
		func() minimax.MoveSearcher {
			innerSearcher :=
				minimax.NewAlphaBetaSearcher(
					generator,
					// terminator will be set
					// automatically
					// by the iterative searcher
					nil,
					evaluator,
				)

			// make and bind a cached searcher
			// to inner one
			minimax.NewCachedSearcher(
				innerSearcher,
				cache,
			)

			return minimax.NewIterativeSearcher(
				innerSearcher,
				// terminator will be set
				// automatically
				// by the parallel searcher
				nil,
			)
		},
	)
}

func newGameModel(
	storage models.PieceStorage,
	searcher minimax.MoveSearcher,
	searcherColor models.Color,
) (games.Manual, error) {
	checker := newChecker()

	return games.NewManual(
		storage,
		minimax.SearcherAdapter{checker},
		minimax.SearcherAdapter{searcher},
		searcherColor,
		models.White,
	)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	boardInFEN := flag.String(
		"fen",
		//"rnbqk/ppppp/5/PPPPP/RNBQK",
		"rnbqkbnr/pppppppp/8/8"+
			"/8/8/PPPPPPPP/RNBQKBNR",
		"board in FEN",
	)
	humanColor := flag.String(
		"color",
		"random",
		"human color "+
			"(allowed: black, white, random)",
	)
	maxDuration := flag.Duration(
		"duration",
		5*time.Second,
		"maximal duration",
	)
	maxCacheSize := flag.Int(
		"cache",
		1e6,
		"maximal cache size",
	)
	flag.Parse()

	storage, err := uci.DecodePieceStorage(
		*boardInFEN,
		pieces.NewPiece,
		models.NewBoard,
	)
	if err != nil {
		log.Fatal(
			"unable to decode a board: ",
			err,
		)
	}

	var searcherColor models.Color
	switch *humanColor {
	case "black":
		searcherColor = models.White
	case "white":
		searcherColor = models.Black
	case "random":
		searcherColor =
			models.Color(rand.Intn(2))
	default:
		log.Fatal("incorrect human color")
	}

	searcher := newSearcher(*maxCacheSize)
	gameModel, err := newGameModel(
		storage,
		searcher,
		searcherColor,
	)
	if err != nil {
		log.Fatal(
			"unable to start a game: ",
			err,
		)
	}

	encoder := cli.NewPieceStorageEncoder(
		uci.EncodePiece,
		"x",
		searcherColor,
	)
	game := cli.NewGame(
		gameModel,
		encoder.Encode,
		os.Stdin,
		os.Stdout,
		"exit",
	)
	firstMove := true
loop:
	for {
		if !firstMove ||
			searcherColor == models.Black {
			err := game.ReadMove("human>")
			switch err {
			case nil:
			case cli.ErrExit:
				break loop
			case games.ErrCheckmate,
				games.ErrDraw:
				const message = "game in state: " +
					"%s\n"
				fmt.Printf(message, err)
				break loop
			default:
				fmt.Printf("error: %s\n", err)
				continue loop
			}
		}

		searcher.SetTerminator(
			terminators.NewTimeTerminator(
				time.Now,
				*maxDuration,
			),
		)

		err := game.SearchMove("ai:")
		switch err {
		case nil:
		case cli.ErrExit:
			break loop
		case games.ErrCheckmate,
			games.ErrDraw:
			const message = "game in state: " +
				"%s\n"
			fmt.Printf(message, err)
			break loop
		default:
			fmt.Printf("error: %s\n", err)
			continue loop
		}

		firstMove = false
	}
}
