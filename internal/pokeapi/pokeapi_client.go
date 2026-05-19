package pokeapi

import (
	"net/http"
	"time"

	"github.com/danielmiranda22/pokedexcli/internal/pokecache"
)

type PokeAPIClient struct {
	httpClient *http.Client
	Next       *string // nil = first page
	Previous   *string // nil = on first page
	cache      *pokecache.Cache
}

func NewPokeAPIClient() *PokeAPIClient {
	return &PokeAPIClient{
		httpClient: &http.Client{},
		cache:      pokecache.NewCache(5 * time.Minute), // 5min sensible
	}
}
