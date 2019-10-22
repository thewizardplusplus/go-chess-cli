package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
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

func printPrompt(side string) {
	fmt.Printf("%s> ", side)
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
	userColor := flag.String(
		"color",
		"random",
		"user color "+
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
	switch *userColor {
	case "black":
		searcherColor = models.White
	case "white":
		searcherColor = models.Black
	case "random":
		searcherColor =
			models.Color(rand.Intn(2))
	default:
		log.Fatal("incorrect user color")
	}

	searcher := newSearcher(*maxCacheSize)
	game, err := newGame(
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

	encoder := cli.PieceStorageEncoder{
		PieceEncoder:     uci.EncodePiece,
		PiecePlaceholder: "x",
		TopColor:         searcherColor,
	}
	firstMove := true
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println(
			encoder.Encode(game.Storage()),
		)
		if game.State() != nil {
			fmt.Println(
				"game in state: ",
				game.State(),
			)

			break
		}

		if !firstMove ||
			searcherColor == models.Black {
			printPrompt("user")

			ok := scanner.Scan()
			if !ok {
				break
			}

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

				continue
			}

			err = game.ApplyMove(move)
			if err != nil {
				log.Print(
					"unable to apply the move: ",
					err,
				)

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

				break
			}
		}

		printPrompt("searcher")

		searcher.SetTerminator(
			terminators.NewTimeTerminator(
				time.Now,
				*maxDuration,
			),
		)

		move := game.SearchMove()
		fmt.Println(uci.EncodeMove(move))

		firstMove = false
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(
			"unable to read a command: ",
			err,
		)
	}
}
