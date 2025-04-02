package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cmolloy36/Pokedex/internal/PokeAPIInteractions"
	"github.com/cmolloy36/Pokedex/internal/pokecache"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	initialNext := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	var initialPrevious *string = nil

	ptr := &config{
		Next:     &initialNext,
		Previous: initialPrevious,
	}

	// interval, _ := time.ParseDuration("1m30s") // make configurable?
	cache := pokecache.NewCache(5 * time.Second)

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

				param := ""
				if len(newText) > 1 {
					param = newText[1]
				}

				err := commandsMap[newText[0]].callback(ptr, cache, param)
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
	callback    func(*config, *pokecache.Cache, string) error
}

type config struct {
	Next     *string
	Previous *string
}

var commandsMap map[string]cliCommand

func init() {
	commandsMap = map[string]cliCommand{
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "See which pokemon are in the location",
			callback:    commandExplore,
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

func commandCatch(ptr *config, cache *pokecache.Cache, param string) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", param)
	err := PokeAPIInteractions.GetPokemonInfo(&url, cache)
	if err != nil {
		return err
	}

	cachedData, ok := cache.Get(url)
	if !ok {
		return fmt.Errorf("could not get entry from catch")
	}

	var entry pokecache.LocationPokemonInfo
	err = json.Unmarshal(cachedData, &entry)
	if err != nil {
		return err
	}

	for _, pkmn := range entry.PokemonEncounters {
		fmt.Println(pkmn.Pokemon.Name)
	}
	// fmt.Printf("%d", len(cache.Cache))

	return nil

}

func commandExit(ptr *config, cache *pokecache.Cache, param string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(ptr *config, cache *pokecache.Cache, param string) error {
	usageStr := ""
	for key, val := range commandsMap { // sort this!! currently returns in various orders
		currStr := fmt.Sprintf("%s: %s\n", key, val.description)
		usageStr += currStr
	}
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n%s", usageStr)
	return nil
}

func commandExplore(ptr *config, cache *pokecache.Cache, param string) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", param)
	err := PokeAPIInteractions.GetLocationPokemon(&url, cache)
	if err != nil {
		return err
	}

	cachedData, ok := cache.Get(url)
	if !ok {
		return fmt.Errorf("could not get entry from explore")
	}

	var entry pokecache.LocationPokemonInfo
	err = json.Unmarshal(cachedData, &entry)
	if err != nil {
		return err
	}

	for _, pkmn := range entry.PokemonEncounters {
		fmt.Println(pkmn.Pokemon.Name)
	}
	// fmt.Printf("%d", len(cache.Cache))

	return nil

}

func commandMap(ptr *config, cache *pokecache.Cache, param string) error {
	err := PokeAPIInteractions.GetLocations(ptr.Next, cache)
	if err != nil {
		return err
	}

	// cache.Mu.Lock()

	// defer cache.Mu.Unlock()

	// fmt.Println("getlocations finished")

	cachedData, ok := cache.Get(*ptr.Next)
	if !ok {
		return fmt.Errorf("could not get entry from map")
	}

	var entry pokecache.BatchInfo
	err = json.Unmarshal(cachedData, &entry)
	if err != nil {
		return err
	}

	ptr.Next = entry.Next
	ptr.Previous = entry.Previous

	for i := 0; i < len(entry.Results); i++ {
		fmt.Println(entry.Results[i].Name)
	}
	// fmt.Printf("cache length: %d\n", len(cache.Cache))

	return nil
}

func commandMapb(ptr *config, cache *pokecache.Cache, param string) error {
	err := PokeAPIInteractions.GetLocations(ptr.Previous, cache)
	if err != nil {
		return err
	}

	cachedData, ok := cache.Get(*ptr.Previous)
	if !ok {
		return fmt.Errorf("could not get entry from mapb")
	}

	var entry pokecache.BatchInfo
	err = json.Unmarshal(cachedData, &entry)
	if err != nil {
		return err
	}

	ptr.Next = entry.Next
	ptr.Previous = entry.Previous

	for i := 0; i < len(entry.Results); i++ {
		fmt.Println(entry.Results[i].Name)
	}

	// fmt.Printf("cache length: %d\n", len(cache.Cache))

	return nil
}
