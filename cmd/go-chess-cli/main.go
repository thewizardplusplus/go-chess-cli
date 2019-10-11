package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
	"github.com/thewizardplusplus/go-chess-models/uci"
)

func printStorage(
	storage models.PieceStorage,
) {
	fmt.Print(" ")
	width := storage.Size().Width
	for i := 0; i < width; i++ {
		fmt.Print(string(i + 97))
	}
	fmt.Println()

	positions := storage.Size().Positions()
	previousRank := -1
	for _, position := range positions {
		if position.Rank != previousRank {
			previousRank = position.Rank
			fmt.Print(position.Rank + 1)
		}

		piece, ok := storage.Piece(position)
		if ok {
			fmt.Print(uci.EncodePiece(piece))
		} else {
			fmt.Print(".")
		}

		lastFile := storage.Size().Height - 1
		if position.File == lastFile {
			fmt.Println()
		}
	}
}

func printPrompt() {
	fmt.Print("> ")
}

func main() {
	boardInFEN := flag.String(
		"fen",
		"rnbqk/ppppp/5/PPPPP/RNBQK",
		"board in FEN",
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

	printStorage(storage)

	scanner := bufio.NewScanner(os.Stdin)
	printPrompt()

	for scanner.Scan() {
		command := scanner.Text()
		command = strings.TrimSpace(command)
		command = strings.ToLower(command)

		switch command {
		case "exit", "quit":
			os.Exit(0)
		default:
			move, err := uci.DecodeMove(command)
			if err != nil {
				log.Print(
					"unable to decode a move: ",
					err,
				)

				printPrompt()

				continue
			}

			err = storage.CheckMove(move)
			if err != nil {
				log.Print("incorrect move: ", err)

				printPrompt()

				continue
			}

			storage = storage.ApplyMove(move)
			printStorage(storage)
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
