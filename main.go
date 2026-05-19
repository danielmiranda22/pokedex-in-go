package main

import "github.com/danielmiranda22/pokedexcli/internal/pokeapi"

func main() {
	config := &pokeapi.Config{}
	startRepl(config)
}
