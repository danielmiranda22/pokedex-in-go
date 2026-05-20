package pokeapi

import (
	"net/http"
	"time"

	"github.com/danielmiranda22/pokedexcli/internal/pokecache"
)

type PokeAPIClient struct {
	httpClient *http.Client
	cache      *pokecache.Cache
	Next       *string
	Previous   *string
}

func NewPokeAPIClient() *PokeAPIClient {
	return &PokeAPIClient{
		httpClient: &http.Client{},
		cache:      pokecache.NewCache(5 * time.Minute),
	}
}
