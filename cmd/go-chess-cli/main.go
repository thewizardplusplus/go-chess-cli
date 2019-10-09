package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		fmt.Println(command)
	}
	if err := scanner.Err(); err != nil {
		const msg = "unable to read a command: "
		log.Fatal(msg, err)
	}
}
