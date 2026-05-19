package pokeapi

import "net/http"

type PokeAPIClient struct {
	httpClient *http.Client
	Next       *string // nil = first page
	Previous   *string // nil = on first page
}

func NewPokeAPIClient() *PokeAPIClient {
	return &PokeAPIClient{
		httpClient: &http.Client{},
	}
}
