package main

import (
	"fmt"
	"os"

	"github.com/danielmiranda22/pokedexcli/internal/pokeapi"
)

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cmds map[string]cliCommand) func() error {
	return func() error {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("Usage:")
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
