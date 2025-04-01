package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cmolloy36/Pokedex/internal/PokeAPIInteractions"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	initialNext := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	initialPrevious := ""

	ptr := &config{
		Next:     &initialNext,
		Previous: &initialPrevious,
	}
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := scanner.Text()
			newText := cleanInput(text)
			if len(newText) == 0 {
				fmt.Println("Please input a command.")
			} else {
				if _, ok := commandsMap[newText[0]]; !ok {
					fmt.Println("Please input a valid command.")
					continue
				}
				err := commandsMap[newText[0]].callback(ptr)
				if err != nil {
					fmt.Println(fmt.Errorf("error: %w", err))
				}
			}
			// fmt.Printf("Your command was: %s\n", newText[0]) \\ for printing first word
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next     *string
	Previous *string
}

var commandsMap map[string]cliCommand

func init() {
	commandsMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Get more info",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays names of next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays names of previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func cleanInput(text string) []string {
	strippedString := strings.TrimSpace(strings.ToLower(text))
	substrings := strings.Fields(strippedString)
	return substrings
}

func commandExit(ptr *config) error { // does this need to return an error? When would exiting not work?
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(ptr *config) error { // does this need to return an error? When would exiting not work?
	usageStr := ""
	for key, val := range commandsMap { // sort this!! currently returns in various orders
		currStr := fmt.Sprintf("%s: %s\n", key, val.description)
		usageStr += currStr
	}
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n%s", usageStr)
	return nil
}

func commandMap(ptr *config) error {
	locations, err := PokeAPIInteractions.GetLocations(ptr.Next)
	if err != nil {
		return err
	}

	ptr.Next = locations.Next
	ptr.Previous = locations.Previous

	for i := 0; i < len(locations.Results); i++ {
		fmt.Println(locations.Results[i].Name)
	}

	return nil
}

func commandMapb(ptr *config) error {
	locations, err := PokeAPIInteractions.GetLocations(ptr.Previous)
	if err != nil {
		return err
	}

	ptr.Next = locations.Next
	ptr.Previous = locations.Previous

	for i := 0; i < len(locations.Results); i++ {
		fmt.Println(locations.Results[i].Name)
	}

	return nil
}
