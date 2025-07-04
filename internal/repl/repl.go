package repl

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Singullaritty/pokedexcli/internal/pokapi"
	"github.com/Singullaritty/pokedexcli/internal/pokecache"
)

type Command interface {
	RunCmd(args []string) error
}

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

func NewCli() map[string]Command {
	sharedConfig := &Config{}
	cache := pokecache.NewCache(5 * time.Second)
	pokemons := make(map[string]pokapi.Pokemon)

	return map[string]Command{
		"help": HelpCommand{
			Name:        "help",
			Description: "Displays a help message",
		},
		"exit": ExitCommand{
			Name:        "exit",
			Description: "Exits from program",
		},
		"map": &MapCommand{
			Name:        "map",
			Description: "Print pokemon location areas",
			Config:      sharedConfig,
			Cache:       cache,
		},
		"mapb": &MapBackCommand{
			Name:        "mapb",
			Description: "Print previous pokemon location areas",
			Config:      sharedConfig,
			Cache:       cache,
		},
		"explore": &ExploreCommand{
			Name:        "explore",
			Description: "Explore pokemons in the location are",
			Config:      sharedConfig,
			Cache:       cache,
		},
		"catch": &CatchCommand{
			Name:        "catch",
			Description: "Catch pokemons",
			Config:      sharedConfig,
			Cache:       cache,
			Pokemons:    pokemons,
		},
		"inspect": &InspectCommand{
			Name:        "inspect",
			Description: "Inspect pokemons",
			Config:      sharedConfig,
			Pokemons:    pokemons,
		},
		"pokedex": &PokedexCommand{
			Name:        "pokedex",
			Description: "list of all the names of the Pokemons that were caught",
			Config:      sharedConfig,
			Pokemons:    pokemons,
		},
	}

}

func (e ExitCommand) RunCmd(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (h HelpCommand) RunCmd(args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("\nhelp: Displays a help message")
	fmt.Println("exit: Exits from program")
	return nil
}

func (m *MapCommand) RunCmd(args []string) error {
	names := []string{}
	pokeNext := m.Config.NextUrl
	if pokeNext == "" {
		pokeNext = "https://pokeapi.co/api/v2/location-area/"
	}
	res, err := pokapi.GetAreas(pokeNext, m.Cache)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	for _, n := range res.Results {
		names = append(names, n.Name)
	}
	fmt.Println(strings.Join(names, "\n"))

	if res.Next != nil {
		m.Config.NextUrl = *res.Next
	} else {
		fmt.Println("No areas to explore!")
		m.Config.NextUrl = ""
	}
	if res.Previous != nil {
		m.Config.PreviousUrl = *res.Previous
	} else {
		m.Config.PreviousUrl = ""
	}
	return nil
}

func (mb *MapBackCommand) RunCmd(args []string) error {
	names := []string{}
	pokePrevious := mb.Config.PreviousUrl
	if pokePrevious == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	res, err := pokapi.GetAreas(pokePrevious, mb.Cache)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	for _, n := range res.Results {
		names = append(names, n.Name)
	}
	fmt.Println(strings.Join(names, "\n"))
	if res.Previous != nil {
		mb.Config.PreviousUrl = *res.Previous
	} else {
		mb.Config.PreviousUrl = ""
	}

	if res.Next != nil {
		mb.Config.NextUrl = *res.Next
	}
	return nil
}

func (e *ExploreCommand) RunCmd(args []string) error {
	pokemons := []string{}
	res, err := pokapi.ExploreArea(args[0], e.Cache)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	for _, p := range res.PokemonEncounters {
		pokemons = append(pokemons, p.Pokemon.Name)
	}
	fmt.Printf("Exploring %s...", args[0])
	fmt.Print("\nFound pokemon:")
	for _, p := range pokemons {
		fmt.Printf("\n - %s", p)
	}
	fmt.Println("")
	return nil
}

func (c *CatchCommand) RunCmd(args []string) error {
	res, err := pokapi.GetPokemonInfo(args[0], c.Cache)
	pokName := res.Name
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	if pokName == "" {
		fmt.Printf("Pokemon %s doesn't exist\n", pokName)
	}
	if _, ok := c.Pokemons[pokName]; ok {
		fmt.Printf("%s already caught!\n", pokName)
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokName)
	exp := res.BaseExperience
	switch {
	case exp <= 50:
		fmt.Printf("%s was caught!\n", pokName)
		fmt.Println("You may now inspect it with the inspect command.")
		c.Pokemons[pokName] = res
	case exp > 50 && exp <= 100:
		chance := rand.Intn(4)
		if chance == 3 {
			fmt.Printf("%s was caught!\n", pokName)
			fmt.Println("You may now inspect it with the inspect command.")
			c.Pokemons[pokName] = res
		} else {
			fmt.Printf("%s escaped!\n", pokName)
		}
	case exp > 100:
		chance := rand.Intn(6)
		if chance == 5 {
			fmt.Printf("%s was caught!\n", pokName)
			fmt.Println("You may now inspect it with the inspect command.")
			c.Pokemons[pokName] = res
		} else {
			fmt.Printf("%s escaped!\n", pokName)
		}
	}
	return nil
}

func (i InspectCommand) RunCmd(args []string) error {
	pokemonName := args[0]
	val, ok := i.Pokemons[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", val.Name)
	fmt.Printf("Height: %d\n", val.Height)
	fmt.Printf("Weight: %d\n", val.Weight)
	fmt.Println("Stats:")
	for _, v := range val.Stats {
		fmt.Printf("  -%s: %d\n", v.Stat.Name, v.BaseStat)
	}
	fmt.Println("Types:")
	for _, v := range val.Types {
		fmt.Printf("  - %s\n", v.Type.Name)
	}
	return nil
}

func (p PokedexCommand) RunCmd(args []string) error {
	if len(p.Pokemons) == 0 {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for k := range p.Pokemons {
		fmt.Printf("  - %s\n", k)
	}
	return nil
}

func StartRepl() {
	initRepl := NewCli()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			args := strings.Split(scanner.Text(), " ")
			cmd, exists := initRepl[args[0]]
			if !exists && args[0] != "" {
				fmt.Println("No such command: ", args[0])
				continue
			}
			if args[0] == "" {
				continue
			}
			err := cmd.RunCmd(args[1:])
			if err != nil {
				fmt.Printf("Error executing command: %s", err)

			}
		} else {
			fmt.Println("Error reading user input")
			break
		}
	}
}
