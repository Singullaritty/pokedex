package repl

import (
	"bufio"
	"fmt"
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

func NewCli() map[string]Command {
	sharedConfig := &Config{}
	cache := pokecache.NewCache(5 * time.Second)

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
			}
			if exists {
				err := cmd.RunCmd(args[1:])
				if err != nil {
					fmt.Println("Error executing command", err)
				}
			}
		} else {
			fmt.Println("Error reading user input")
			break
		}
	}
}
