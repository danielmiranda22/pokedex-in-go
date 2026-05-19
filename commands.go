package main

import (
	"fmt"
	"os"

	"github.com/danielmiranda22/pokedexcli/internal/pokeapi"
)

func commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cmds map[string]cliCommand) func([]string) error {
	return func(args []string) error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("Usage:")
		for _, cmd := range cmds {
			fmt.Printf("%s: %s\n", cmd.name, cmd.description)
		}
		fmt.Println()
		return nil
	}
}

func commandMap(client *pokeapi.PokeAPIClient) func([]string) error {
	return func(args []string) error {
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

func commandMapBack(client *pokeapi.PokeAPIClient) func([]string) error {
	return func(args []string) error {
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

func commandExplore(client *pokeapi.PokeAPIClient) func([]string) error {
	return func(args []string) error {
		// validate — explore needs exactly one argument
		if len(args) < 2 {
			return fmt.Errorf("usage: explore <area-name>")
		}

		areaName := args[1] // words[0]=explore, words[1]=area name
		fmt.Printf("Exploring %s...\n", areaName)

		pokemon, err := client.ExploreArea(areaName)
		if err != nil {
			return err
		}

		fmt.Println("Found Pokemon:")
		for _, name := range pokemon {
			fmt.Printf(" - %s\n", name)
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
	cmds["explore"] = cliCommand{
		name:        "explore",
		description: "Explore a location area — usage: explore <area-name>",
		callback:    commandExplore(client),
	}
	return cmds
}
