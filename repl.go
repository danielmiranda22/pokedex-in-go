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
	callback    func() error
	api         *pokeapi.Config
}

func startRepl(apiConfig *pokeapi.Config) {
	cmds := getCommands(apiConfig)
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

		// look up the command in the registry
		cmd, ok := cmds[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		// call the command — print error if it returns one
		if err := cmd.callback(); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil // never reached — os.Exit stops the program
}

func commandHelp(cmds map[string]cliCommand) func() error {
	return func() error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("Usage:\n")
		for _, cmd := range cmds {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		fmt.Println()
		return nil
	}
}

// commandMap fetches next 20 location areas
func commandMap(config *pokeapi.Config) func() error {
	return func() error {
		areas, err := pokeapi.GetLocationAreas(config)
		if err != nil {
			return err
		}
		for _, area := range areas {
			fmt.Println(area)
		}
		return nil
	}
}

// commandMapBack fetches previous 20 location areas
func commandMapBack(config *pokeapi.Config) func() error {
	return func() error {
		if config.Previous == nil {
			fmt.Println("you're on the first page")
			return nil
		}
		areas, err := pokeapi.GetPreviousLocationAreas(config)
		if err != nil {
			return err
		}
		for _, area := range areas {
			fmt.Println(area)
		}
		return nil
	}
}

func getCommands(apiConfig *pokeapi.Config) map[string]cliCommand {
	cmds := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
	cmds["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp(cmds),
	}
	cmds["map"] = cliCommand{
		name:        "map",
		description: "Displays a list of location areas",
		callback:    commandMap(apiConfig),
	}
	cmds["mapb"] = cliCommand{
		name:        "map-back",
		description: "Go back to the previous list of location areas",
		callback:    commandMapBack(apiConfig),
	}
	return cmds
}

func cleanInput(input string) []string {
	return strings.Fields(strings.ToLower(input))
}
