package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
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

func newGame(
	storage models.PieceStorage,
	searcher minimax.MoveSearcher,
) (games.Manual, error) {
	checker := newChecker()

	return games.NewManual(
		storage,
		minimax.SearcherAdapter{checker},
		minimax.SearcherAdapter{searcher},
		models.Black,
		models.White,
	)
}

func printPrompt() {
	fmt.Print("> ")
}

func main() {
	boardInFEN := flag.String(
		"fen",
		//"rnbqk/ppppp/5/PPPPP/RNBQK",
		"rnbqkbnr/pppppppp/8/8"+
			"/8/8/PPPPPPPP/RNBQKBNR",
		"board in FEN",
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

	searcher := newSearcher(*maxCacheSize)
	game, err := newGame(storage, searcher)
	if err != nil {
		log.Fatal(
			"unable to start a game: ",
			err,
		)
	}
	if game.State() != nil {
		fmt.Println(
			"game in state: ",
			game.State(),
		)
		os.Exit(0)
	}

	encoder := cli.PieceStorageEncoder{
		PieceEncoder: uci.EncodePiece,
		Separator:    "x",
		TopColor:     models.Black,
	}
	fmt.Println(
		encoder.Encode(game.Storage()),
	)
	printPrompt()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		command = strings.TrimSpace(command)
		command = strings.ToLower(command)
		if command == "exit" {
			os.Exit(0)
		}

		move, err := uci.DecodeMove(command)
		if err != nil {
			log.Print(
				"unable to decode the move: ",
				err,
			)

			fmt.Println(
				encoder.Encode(game.Storage()),
			)
			printPrompt()

			continue
		}

		err = game.ApplyMove(move)
		if err != nil {
			log.Print(
				"unable to apply the move: ",
				err,
			)

			fmt.Println(
				encoder.Encode(game.Storage()),
			)
			printPrompt()

			continue
		}

		fmt.Println(
			encoder.Encode(game.Storage()),
		)
		if game.State() != nil {
			fmt.Println(
				"game in state: ",
				game.State(),
			)
			os.Exit(0)
		}

		searcher.SetTerminator(
			terminators.NewTimeTerminator(
				time.Now,
				*maxDuration,
			),
		)

		move, err = game.SearchMove()
		if err != nil {
			log.Print(
				"unable to search a move: ",
				err,
			)

			continue
		}

		printPrompt()
		fmt.Println(uci.EncodeMove(move))

		fmt.Println(
			encoder.Encode(game.Storage()),
		)
		if game.State() != nil {
			fmt.Println(
				"game in state: ",
				game.State(),
			)
			os.Exit(0)
		}

		printPrompt()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(
			"unable to read a command: ",
			err,
		)
	}
}
