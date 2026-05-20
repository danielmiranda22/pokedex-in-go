package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/danielmiranda22/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

func startRepl(client *pokeapi.PokeAPIClient, pokedex *Pokedex) {
	cmds := getCommands(client, pokedex)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		cmd, ok := cmds[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(words); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func cleanInput(input string) []string {
	return strings.Fields(strings.ToLower(input))
}
