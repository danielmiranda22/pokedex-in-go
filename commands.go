package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/danielmiranda22/pokedexcli/internal/pokeapi"
)

var (
	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
	colorCyan    = "\033[36m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorGray    = "\033[90m"
	colorRed     = "\033[31m"
)

func commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cmds map[string]cliCommand) func([]string) error {
	return func(args []string) error {
		fmt.Printf("\n%s🤾 Pokédex CLI Help%s\n\n", colorBlue, colorReset)

		seen := make(map[string]bool)

		for _, cmd := range cmds {
			if seen[cmd.name] {
				continue
			}

			seen[cmd.name] = true

			// Customize the formatting of the command name and description
			// Ansi color codes for cyan and reset
			// pad left to 15 characters for alignment
			fmt.Printf("%s  %-15s%s %s\n", colorCyan, cmd.name, colorReset, cmd.description)
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
			fmt.Printf(" - %s%s%s\n", colorGreen, area, colorReset)
		}
		return nil
	}
}

func commandMapBack(client *pokeapi.PokeAPIClient) func([]string) error {
	return func(args []string) error {
		if client.Previous == nil {
			fmt.Printf("%sYou're on the first page!%s\n", colorGray, colorReset)
			return nil
		}
		areas, err := client.GetPreviousLocationAreas()
		if err != nil {
			return err
		}
		for _, area := range areas {
			fmt.Printf(" - %s%s%s\n", colorGreen, area, colorReset)
		}
		return nil
	}
}

func commandExplore(client *pokeapi.PokeAPIClient) func([]string) error {
	return func(args []string) error {
		// validate — explore needs exactly one argument
		if len(args) < 2 {
			return fmt.Errorf("%susage: explore <area-name>%s", colorGray, colorReset)
		}

		areaName := args[1] // words[0]=explore, words[1]=area name
		fmt.Printf("Exploring %s%s%s...\n", colorCyan, areaName, colorReset)

		pokemons, err := client.ExploreArea(areaName)
		if err != nil {
			return err
		}

		fmt.Printf("%sFound Pokemon:%s\n", colorYellow, colorReset)
		for _, name := range pokemons {
			fmt.Printf(" - %s%s%s\n", colorGreen, name, colorReset)
		}
		return nil
	}
}

func commandCatch(client *pokeapi.PokeAPIClient, pokedex *Pokedex) func([]string) error {
	return func(args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("%susage: catch <pokemon-name>%s", colorGray, colorReset)
		}

		pokemonName := args[1] // words[0]=catch, words[1]=pokemon name
		fmt.Printf("Throwing a Pokeball at %s%s%s...\n", colorCyan, pokemonName, colorReset)

		pokemon, err := client.GetPokemon(pokemonName)
		if err != nil {
			return err
		}

		rnd := rand.Intn(pokemon.BaseExperience)
		if rnd < 40 {
			fmt.Printf("%s%s was caught!%s\n", colorGreen, pokemonName, colorReset)
			fmt.Printf("%sYou may now inspect it with the inspect command.%s\n", colorGray, colorReset)
			pokedex.CaughtPokemon[pokemonName] = pokemon
		} else {
			fmt.Printf("%s%s escaped!%s\n", colorRed, pokemonName, colorReset)
		}

		return nil
	}
}

func commandInspect(pokedex *Pokedex) func([]string) error {
	return func(args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("%susage: inspect <pokemon-name>%s", colorGray, colorReset)
		}

		pokemonName := args[1] // words[0]=inspect, words[1]=pokemon name
		pokemon, ok := pokedex.CaughtPokemon[pokemonName]
		if !ok {
			return fmt.Errorf("%s%s%s is not in your Pokedex", colorRed, pokemonName, colorReset)
		}

		fmt.Printf("\n%sName: %s%s\n", colorYellow, pokemon.Name, colorReset)
		fmt.Printf("%sHeight: %d%s\n", colorYellow, pokemon.Height, colorReset)
		fmt.Printf("%sWeight: %d%s\n", colorYellow, pokemon.Weight, colorReset)
		if len(pokemon.Stats) > 0 {
			fmt.Printf("%sStats:%s\n", colorYellow, colorReset)
			for _, stat := range pokemon.Stats {
				fmt.Printf("  %s-%s: %s%d%s\n", colorCyan, stat.Stat.Name, colorReset, stat.BaseStat, colorReset)
			}
		}
		if len(pokemon.Types) > 0 {
			fmt.Printf("%sTypes:%s\n", colorYellow, colorReset)
			for _, t := range pokemon.Types {
				fmt.Printf("  -%s%s%s\n", colorCyan, t.Type.Name, colorReset)
			}
		}
		return nil
	}
}

func commandPokedex(pokedex *Pokedex) func([]string) error {
	return func(args []string) error {
		if len(pokedex.CaughtPokemon) == 0 {
			fmt.Printf("%sYou haven't caught any Pokemon yet!%s\n", colorRed, colorReset)
			return nil
		}

		fmt.Printf("\n%sYour Pokedex:%s\n", colorBlue, colorReset)
		for name := range pokedex.CaughtPokemon {
			fmt.Printf(" - %s%s%s\n", colorGreen, name, colorReset)
		}
		return nil
	}
}

func getCommands(client *pokeapi.PokeAPIClient, pokedex *Pokedex) map[string]cliCommand {
	cmds := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	// Create the map command and add it under both "map" and "m"
	mapCmd := cliCommand{
		name:        "map | m",
		description: "Displays next 20 location areas",
		callback:    commandMap(client),
	}

	cmds["map"] = mapCmd
	cmds["m"] = mapCmd

	// Create the catch command and add it under both "catch" and "c"
	catchCmd := cliCommand{
		name:        "catch | c",
		description: "usage: catch <pokemon-name>",
		callback:    commandCatch(client, pokedex),
	}

	cmds["catch"] = catchCmd
	cmds["c"] = catchCmd

	cmds["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp(cmds),
	}
	cmds["mapb"] = cliCommand{
		name:        "mapb | mb",
		description: "Displays previous 20 location areas",
		callback:    commandMapBack(client),
	}
	cmds["explore"] = cliCommand{
		name:        "explore",
		description: "usage: explore <area-name>",
		callback:    commandExplore(client),
	}
	cmds["inspect"] = cliCommand{
		name:        "inspect",
		description: "usage: inspect <pokemon-name>",
		callback:    commandInspect(pokedex),
	}
	cmds["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "List all caught Pokemon",
		callback:    commandPokedex(pokedex),
	}
	return cmds
}
