package repl

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Singullaritty/pokedexcli/internal/pokapi"
	"github.com/Singullaritty/pokedexcli/internal/pokecache"
	"golang.org/x/term"
)

const (
	ClearLine    = "\r\033[K"
	ControlC     = byte(3)
	ControlD     = byte(4)
	KeyEscape    = byte(27)
	KeyBackspace = byte(8)
	KeyDelete    = byte(127)
	KeyUp        = byte(65)
	KeyDown      = byte(66)
	KeyEnter     = byte(13)
)

var ErrExit = errors.New("exit")

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
	fmt.Print("Closing the Pokedex... Goodbye!")
	return ErrExit
}

func (h HelpCommand) RunCmd(args []string) error {
	fmt.Println("\rWelcome to the Pokedex!")
	fmt.Println("\rUsage:")
	fmt.Println("\r    help: Displays a help message")
	fmt.Println("\r    exit: Exits from program")
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
	fmt.Print(strings.Join(names, "\r\n"))

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
	fmt.Println(strings.Join(names, "\r\n"))
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
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", args[0])
	res, err := pokapi.ExploreArea(url, e.Cache)
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
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", args[0])
	res, err := pokapi.GetPokemonInfo(url, c.Cache)
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
		chance := rand.Intn(3)
		if chance == 2 {
			fmt.Printf("%s was caught!\n", pokName)
			fmt.Println("You may now inspect it with the inspect command.")
			c.Pokemons[pokName] = res
		} else {
			fmt.Printf("%s escaped!\n", pokName)
		}
	case exp > 100:
		chance := rand.Intn(5)
		if chance == 4 {
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
	fd := int(os.Stdin.Fd())

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)

	var lineBuffer []byte
	var b = make([]byte, 1)
	var esc = make([]byte, 2)
	var history []string
	var historyIndex int
	for {
		fmt.Print(ClearLine)
		fmt.Print("Pokedex > ")
		for {
			_, err := os.Stdin.Read(b)
			if err != nil {
				fmt.Printf("Error reading stdin: %s", err)
				return
			}
			if b[0] >= 32 && b[0] <= 126 {
				lineBuffer = append(lineBuffer, b[0])
				fmt.Print(string(b[0]))
			} else if b[0] == KeyEscape {
				_, err := os.Stdin.Read(esc)
				if err != nil {
					fmt.Printf("Error reading stdin: %s", err)
					return
				}
				if len(history) == 0 {
					continue
				}
				if esc[0] == '[' && esc[1] == KeyUp {
					if historyIndex >= 0 {
						lineBuffer = []byte(history[historyIndex])
						fmt.Print(ClearLine)
						fmt.Print("Pokedex > ")
						fmt.Print(history[historyIndex])
						if historyIndex > 0 {
							historyIndex--
						}
					}
				}
				if esc[0] == '[' && esc[1] == KeyDown {
					if historyIndex == len(history)-1 {
						lineBuffer = nil
						fmt.Print(ClearLine)
						fmt.Print("Pokedex > ")
						continue
					}
					lineBuffer = []byte(history[historyIndex])
					fmt.Print(ClearLine)
					fmt.Print("Pokedex > ")
					fmt.Print(history[historyIndex+1])
					historyIndex++
				}
			} else if b[0] == KeyEnter {
				fmt.Println("\r")
				args := strings.Fields(string(lineBuffer))
				if len(args) == 0 {
					historyIndex = len(history) - 1
					break
				}
				cmd, exists := initRepl[args[0]]
				if !exists && args[0] != "" {
					fmt.Printf("Command %s doesn't exist\r\n", args[0])
					lineBuffer = nil
					break
				}
				err := cmd.RunCmd(args[1:])
				if errors.Is(err, ErrExit) {
					return
				}
				if err != nil {
					fmt.Printf("Error executing command: %s", err)
					return
				}
				history = append(history, args[0])
				historyIndex = len(history) - 1
				lineBuffer = nil
				fmt.Print("\r")
				break
			} else if b[0] == KeyBackspace || b[0] == KeyDelete {
				if len(lineBuffer) > 0 {
					lineBuffer = lineBuffer[:len(lineBuffer)-1]
					fmt.Print("\b \b")
				}
			} else if b[0] == ControlC || b[0] == ControlD {
				return
			}
		}
	}
}
