package main

import "github.com/danielmiranda22/pokedexcli/internal/pokeapi"

func main() {
	client := pokeapi.NewPokeAPIClient()
	pokedex := NewPokedex()
	startRepl(client, pokedex)
}
