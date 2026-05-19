package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// GetLocationAreas fetches the next 20 location areas
func (c *PokeAPIClient) GetLocationAreas() ([]string, error) {
	// use Next URL if available, otherwise start from beginning
	url := locationAreaURL
	if c.Next != nil {
		url = *c.Next
	}

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

	var payload LocationAreaResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	// update pagination state on the client
	c.Next = payload.Next
	c.Previous = payload.Previous

	names := make([]string, len(payload.Results))
	for i, area := range payload.Results {
		names[i] = area.Name
	}
	return names, nil
}

// GetPreviousLocationAreas goes back one page
func (c *PokeAPIClient) GetPreviousLocationAreas() ([]string, error) {
	if c.Previous == nil {
		return nil, nil // caller checks for nil — means first page
	}
	// point Next at Previous so GetLocationAreas uses it
	c.Next = c.Previous
	return c.GetLocationAreas()
}
