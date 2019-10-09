package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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

		fmt.Println(command)
		prompt()
	}
	if err := scanner.Err(); err != nil {
		const msg = "unable to read a command: "
		log.Fatal(msg, err)
	}
}
