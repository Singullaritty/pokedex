package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Singullaritty/pokedexcli/internal/pokapi"
)

type Command interface {
	RunCmd() error
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

type MapCommand struct {
	Name        string
	Description string
	Config      *Config
}

type MapBackCommand struct {
	Name        string
	Description string
	Config      *Config
}

func NewCli() map[string]Command {
	sharedConfig := &Config{}

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
		},
		"mapb": &MapBackCommand{
			Name:        "mapb",
			Description: "Print previous pokemon location areas",
			Config:      sharedConfig,
		},
	}

}

func (e ExitCommand) RunCmd() error {
	fmt.Println("Closing the Pokedex... Goodbye!")

	return nil
}

func (h HelpCommand) RunCmd() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("\nhelp: Displays a help message")
	fmt.Println("exit: Exits from program")
	// fmt.Println("\nmap: Print pokemon locations area")

	return nil
}

func (m *MapCommand) RunCmd() error {
	names := []string{}
	pokeNext := m.Config.NextUrl

	if pokeNext == "" {
		pokeNext = "https://pokeapi.co/api/v2/location-area/"
	}

	res, err := pokapi.GetAreas(pokeNext)
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

func (mb *MapBackCommand) RunCmd() error {
	names := []string{}
	pokePrevious := mb.Config.PreviousUrl

	if pokePrevious == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	res, err := pokapi.GetAreas(pokePrevious)
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

func StartRepl() {
	initRepl := NewCli()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := scanner.Text()

			cmd, exists := initRepl[text]
			if exists {
				if err := cmd.RunCmd(); err != nil {
					fmt.Println("Error executing command", err)
				}
				if text == "exit" {
					os.Exit(0)
				}
			} else {
				fmt.Println("No such command: ", text)
			}
		} else {
			fmt.Println("Error reading user input")
			break
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(text)
}
