package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2/location-area/"

type Config struct {
	Next     *string
	Previous *string
}

type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`     // nullable — *string not string
	Previous *string `json:"previous"` // nullable — first page is null
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationAreas(config *Config) ([]string, error) {
	url := baseURL
	if config.Next != nil {
		url = *config.Next
	}

	// make the GET request
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // always close body when done

	// read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// check status code AFTER reading body
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("bad status: %d", res.StatusCode)
	}

	// decode JSON into struct
	var payload LocationAreaResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	// update pagination state for next/previous calls
	config.Next = payload.Next
	config.Previous = payload.Previous

	// extract just the names
	names := make([]string, len(payload.Results))
	for i, area := range payload.Results {
		names[i] = area.Name
	}
	return names, nil
}

func GetPreviousLocationAreas(config *Config) ([]string, error) {
	if config.Previous == nil {
		return nil, nil // signal: no previous page
	}
	// swap next to previous URL so GetLocationAreas uses it
	config.Next = config.Previous
	return GetLocationAreas(config)
}
