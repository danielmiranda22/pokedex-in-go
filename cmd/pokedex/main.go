package main

import (
	"github.com/danielmiranda22/pokedexcli/internal/cli"
	"github.com/danielmiranda22/pokedexcli/internal/commands"
	"github.com/danielmiranda22/pokedexcli/internal/domain"
	"github.com/danielmiranda22/pokedexcli/internal/pokeapi"
)

func main() {
	client := pokeapi.NewPokeAPIClient()
	pokedex := domain.NewPokedex()
	cmds := commands.GetCommands(client, pokedex)
	cli.StartRepl(cmds)
}
