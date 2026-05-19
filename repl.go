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
	// no api field — commands get client via closure
}

func startRepl(client *pokeapi.PokeAPIClient) {
	cmds := getCommands(client)
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

		if err := cmd.callback(); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
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

func commandMap(client *pokeapi.PokeAPIClient) func() error {
	return func() error {
		areas, err := client.GetLocationAreas()
		if err != nil {
			return err
		}
		for _, area := range areas {
			fmt.Println(area)
		}
		return nil
	}
}

func commandMapBack(client *pokeapi.PokeAPIClient) func() error {
	return func() error {
		if client.Previous == nil {
			fmt.Println("you're on the first page")
			return nil
		}
		areas, err := client.GetPreviousLocationAreas()
		if err != nil {
			return err
		}
		for _, area := range areas {
			fmt.Println(area)
		}
		return nil
	}
}

func getCommands(client *pokeapi.PokeAPIClient) map[string]cliCommand {
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
		description: "Displays next 20 location areas",
		callback:    commandMap(client),
	}
	cmds["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays previous 20 location areas",
		callback:    commandMapBack(client),
	}
	return cmds
}

func cleanInput(input string) []string {
	return strings.Fields(strings.ToLower(input))
}
