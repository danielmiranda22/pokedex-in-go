package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		words := cleanInput(input)

		if len(words) > 0 {
			fmt.Println("Your command was:", words[0])
		}
	}
}

// split the input into "words"
func cleanInput(input string) []string {
	return strings.Fields(strings.ToLower(input))
}
