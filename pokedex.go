package main

import "github.com/danielmiranda22/pokedexcli/internal/pokeapi"

// Pokedex holds the user's game state — caught pokemon
type Pokedex struct {
	CaughtPokemon map[string]pokeapi.Pokemon
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		CaughtPokemon: make(map[string]pokeapi.Pokemon),
	}
}
