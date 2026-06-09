package commands

import (
	"fmt"
	"math/rand"
	"os"
	"sort"

	"github.com/danielmiranda22/pokedexcli/internal/domain"
	"github.com/danielmiranda22/pokedexcli/internal/pokeapi"
)

type Command struct {
	Name        string
	Description string
	Callback    func(args []string) error
	Order       int
}

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

func commandHelp(cmds map[string]Command) func([]string) error {
	return func(args []string) error {
		orderedCmds := sortCommandByOrder(cmds)

		fmt.Printf("\n%s🤾 Pokédex CLI Help%s\n\n", colorBlue, colorReset)
		for _, cmd := range orderedCmds {
			fmt.Printf("%s  %-15s%s %s\n", colorCyan, cmd.Name, colorReset, cmd.Description)
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
		if err := validateArgs(args, 2, "usage: explore <area-name>"); err != nil {
			return fmt.Errorf("%s%v%s", colorGray, err, colorReset)
		}

		areaName := args[1]
		fmt.Printf("Exploring %s%s%s...\n", colorCyan, areaName, colorReset)

		pokemons, err := client.ExploreArea(areaName)
		if err != nil {
			return fmt.Errorf("%s%v%s", colorGray, err, colorReset)
		}

		fmt.Printf("%sFound Pokemon:%s\n", colorYellow, colorReset)
		for _, name := range pokemons {
			fmt.Printf(" - %s%s%s\n", colorGreen, name, colorReset)
		}
		return nil
	}
}

func commandCatch(client *pokeapi.PokeAPIClient, pokedex *domain.Pokedex) func([]string) error {
	return func(args []string) error {
		if err := validateArgs(args, 2, "usage: catch <pokemon-name>"); err != nil {
			return fmt.Errorf("%s%v%s", colorGray, err, colorReset)
		}

		pokemonName := args[1]
		fmt.Printf("Throwing a Pokeball at %s%s%s...\n", colorCyan, pokemonName, colorReset)

		pokemon, err := client.GetPokemon(pokemonName)
		if err != nil {
			return fmt.Errorf("%s%v%s", colorGray, err, colorReset)
		}

		if getCatchResult(pokemon.BaseExperience) {
			fmt.Printf("%s%s was caught!%s\n", colorGreen, pokemonName, colorReset)
			fmt.Printf("%sYou may now inspect it with the inspect command.%s\n", colorGray, colorReset)
			pokedex.CaughtPokemon[pokemonName] = pokemon
		} else {
			fmt.Printf("%s%s escaped!%s\n", colorRed, pokemonName, colorReset)
		}

		return nil
	}
}

func commandInspect(pokedex *domain.Pokedex) func([]string) error {
	return func(args []string) error {
		if err := validateArgs(args, 2, "usage: inspect <pokemon-name>"); err != nil {
			return fmt.Errorf("%s%v%s", colorGray, err, colorReset)
		}

		pokemonName := args[1]
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

func commandPokedex(pokedex *domain.Pokedex) func([]string) error {
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

func GetCommands(client *pokeapi.PokeAPIClient, pokedex *domain.Pokedex) map[string]Command {
	cmds := map[string]Command{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
			Order:       8,
		},
	}

	cmds["pokedex"] = Command{
		Name:        "pokedex",
		Description: "List all caught Pokemon",
		Callback:    commandPokedex(pokedex),
		Order:       7,
	}

	cmds["inspect"] = Command{
		Name:        "inspect",
		Description: "usage: inspect <pokemon-name>",
		Callback:    commandInspect(pokedex),
		Order:       6,
	}

	catchCmd := Command{
		Name:        "catch | c",
		Description: "usage: catch <pokemon-name>",
		Callback:    commandCatch(client, pokedex),
		Order:       5,
	}
	cmds["catch"] = catchCmd
	cmds["c"] = catchCmd

	cmds["explore"] = Command{
		Name:        "explore",
		Description: "usage: explore <area-name>",
		Callback:    commandExplore(client),
		Order:       4,
	}

	cmds["mapb"] = Command{
		Name:        "mapb",
		Description: "Displays previous 20 location areas",
		Callback:    commandMapBack(client),
		Order:       3,
	}

	mapCmd := Command{
		Name:        "map | m",
		Description: "Displays next 20 location areas",
		Callback:    commandMap(client),
		Order:       2,
	}
	cmds["map"] = mapCmd
	cmds["m"] = mapCmd

	cmds["help"] = Command{
		Name:        "help",
		Description: "Displays a help message",
		Callback:    commandHelp(cmds),
		Order:       1,
	}

	return cmds
}

// Private

func validateArgs(args []string, expected int, expectedMsg string) error {
	if len(args) < expected {
		return fmt.Errorf("%s", expectedMsg)
	}
	return nil
}

func getCatchResult(xp int) bool {
	randomVal := rand.Intn(xp)
	return shouldCatchPokemon(randomVal)
}

func shouldCatchPokemon(rnd int) bool {
	return rnd > 40
}

func sortCommandByOrder(cmds map[string]Command) []Command {
	seen := make(map[string]bool)
	orderedCmds := make([]Command, 0, len(cmds))

	for _, cmd := range cmds {
		if seen[cmd.Name] {
			continue
		}
		seen[cmd.Name] = true
		orderedCmds = append(orderedCmds, cmd)
	}

	sort.Slice(orderedCmds, func(i, j int) bool {
		return orderedCmds[i].Order < orderedCmds[j].Order
	})

	return orderedCmds

}
