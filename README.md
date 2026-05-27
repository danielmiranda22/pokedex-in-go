# Pokédex CLI

Command-line Pokédex in Go. Uses [PokéAPI](https://pokeapi.co) to explore locations, catch Pokémon and manage your collection.

> Initially built as part of [Boot.dev's](https://boot.dev). Now being actively extended with new features, improvements, and explorations into Go CLI development.

---

## Run

```bash
git clone https://github.com/danielmiranda22/pokedexcli.git
cd pokedexcli
go run .
```

**Requires:** Go 1.21+

---

## Commands

| Command             | Description                    |
| ------------------- | ------------------------------ |
| `map` / `mapb`      | Navigate location areas        |
| `explore <area>`    | List Pokémon in an area        |
| `catch <pokemon>`   | Attempt to catch a Pokémon     |
| `inspect <pokemon>` | View stats of a caught Pokémon |
| `pokedex`           | List all caught Pokémon        |
| `help` / `exit`     | Help and exit                  |

---

## Stack

- Pure Go stdlib — no frameworks
- `net/http` — HTTP client
- `encoding/json` — JSON parsing
- `sync.Mutex` + goroutines — thread-safe in-memory cache with TTL
