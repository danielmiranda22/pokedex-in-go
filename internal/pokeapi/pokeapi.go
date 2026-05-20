package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *PokeAPIClient) ExploreArea(areaName string) ([]string, error) {
	url := locationAreaURL + areaName

	// check cache first
	if cached, ok := c.cache.Get(url); ok {
		return parsePokemonNames(cached)
	}

	// cache miss — fetch from API
	body, err := c.fetch(url)
	if err != nil {
		return nil, err
	}

	// store in cache
	c.cache.Add(url, body)

	return parsePokemonNames(body)
}

func parsePokemonNames(body []byte) ([]string, error) {
	var payload ExploreAreaResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	names := make([]string, len(payload.PokemonEncounters))
	for i, encounter := range payload.PokemonEncounters {
		names[i] = encounter.Pokemon.Name
	}
	return names, nil
}

func (c *PokeAPIClient) GetLocationAreas() ([]string, error) {
	url := locationAreaURL
	if c.Next != nil {
		url = *c.Next
	}

	// check cache — avoid network call if we have it
	if cached, ok := c.cache.Get(url); ok {
		fmt.Println("(cache hit)")
		return c.parseLocationNames(cached)
	}

	// cache miss — make the HTTP request
	fmt.Println("(cache miss — fetching from API)")
	body, err := c.fetch(url)
	if err != nil {
		return nil, err
	}

	// store in cache for next time
	c.cache.Add(url, body)

	return c.parseLocationNames(body)
}

// GetPreviousLocationAreas goes back one page
func (c *PokeAPIClient) GetPreviousLocationAreas() ([]string, error) {
	if c.Previous == nil {
		return nil, nil
	}
	c.Next = c.Previous
	return c.GetLocationAreas()
}

func (c *PokeAPIClient) GetPokemon(pokemonName string) (Pokemon, error) {
	url := pokemonURL + pokemonName

	// check cache first
	if cached, ok := c.cache.Get(url); ok {
		var pokemon Pokemon
		if err := json.Unmarshal(cached, &pokemon); err != nil {
			return Pokemon{}, err
		}
		return pokemon, nil
	}

	// cache miss — fetch from API
	body, err := c.fetch(url)
	if err != nil {
		return Pokemon{}, err
	}

	// store in cache
	c.cache.Add(url, body)

	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}

// fetch makes a GET request and returns the raw body bytes
func (c *PokeAPIClient) fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("bad status: %d body: %s", res.StatusCode, body)
	}

	return body, nil
}

func (c *PokeAPIClient) parseLocationNames(body []byte) ([]string, error) {
	var payload LocationAreaResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	// update pagination for next/previous calls
	c.Next = payload.Next
	c.Previous = payload.Previous

	names := make([]string, len(payload.Results))
	for i, area := range payload.Results {
		names[i] = area.Name
	}
	return names, nil
}
