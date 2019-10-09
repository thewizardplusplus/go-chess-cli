package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/thewizardplusplus/go-chess-models/uci"
)

func prompt() {
	fmt.Print("> ")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	prompt()

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

				prompt()

				continue
			}

			fmt.Println(move)
		}

		prompt()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(
			"unable to read a command: ",
			err,
		)
	}
}
