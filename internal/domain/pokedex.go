package domain

import "github.com/danielmiranda22/pokedexcli/internal/pokeapi"

type Pokedex struct {
	CaughtPokemon map[string]pokeapi.Pokemon
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		CaughtPokemon: make(map[string]pokeapi.Pokemon),
	}
}
