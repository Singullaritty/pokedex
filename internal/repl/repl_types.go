package repl

import (
	"github.com/Singullaritty/pokedexcli/internal/pokapi"
	"github.com/Singullaritty/pokedexcli/internal/pokecache"
)

type Config struct {
	NextUrl     string
	PreviousUrl string
}

type HelpCommand struct {
	Name        string
	Description string
}

type ExitCommand struct {
	Name        string
	Description string
}

type ExploreCommand struct {
	Name        string
	Description string
	Config      *Config
	Cache       *pokecache.Cache
}

type CatchCommand struct {
	Name        string
	Description string
	Config      *Config
	Cache       *pokecache.Cache
	Pokemons    map[string]pokapi.Pokemon
}

type MapCommand struct {
	Name        string
	Description string
	Config      *Config
	Cache       *pokecache.Cache
}

type MapBackCommand struct {
	Name        string
	Description string
	Config      *Config
	Cache       *pokecache.Cache
}

type InspectCommand struct {
	Name        string
	Description string
	Config      *Config
	Pokemons    map[string]pokapi.Pokemon
}

type PokedexCommand struct {
	Name        string
	Description string
	Config      *Config
	Pokemons    map[string]pokapi.Pokemon
}

type Command interface {
	RunCmd(args []string) error
}
