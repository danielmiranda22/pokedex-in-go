package main

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"

	"github.com/danielmiranda22/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

func startRepl(client *pokeapi.PokeAPIClient, pokedex *Pokedex) {
	cmds := getCommands(client, pokedex)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "Pokedex > ",
		HistoryFile:     "/tmp/pokedex_history.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Println("Error initializing readline:", err)
		return
	}
	defer rl.Close()

	fmt.Println("🤾🏻‍♀️ Welcome to the Pokedex! Type 'help' to see available commands.")

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		words := cleanInput(line)
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
