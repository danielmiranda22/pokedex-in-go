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
