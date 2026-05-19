package main

import "github.com/danielmiranda22/pokedexcli/internal/pokeapi"

func main() {
	client := pokeapi.NewPokeAPIClient()
	startRepl(client)
}
