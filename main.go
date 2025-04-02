package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
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

	pkmnCaught := make(map[string]pokecache.PokemonInfo)

	// interval, _ := time.ParseDuration("1m30s") // make configurable?
	cache := pokecache.NewCache(5 * time.Second)

	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			text := scanner.Text()
			newText := cleanInput(text)
			if len(newText) == 0 {
				fmt.Println("Please enter a command.")
			} else {
				if _, ok := commandsMap[newText[0]]; !ok {
					fmt.Println("Please enter a valid command.")
					continue
				}

				param := ""
				if len(newText) > 1 {
					param = newText[1]
				}

				err := commandsMap[newText[0]].callback(scanner, ptr, cache, pkmnCaught, param)
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
	callback    func(*bufio.Scanner, *config, *pokecache.Cache, map[string]pokecache.PokemonInfo, string) error
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
			description: "Attempt to catch a Pokemon",
			callback:    commandCatch,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "See which Pokemon are in the location",
			callback:    commandExplore,
		},
		"help": {
			name:        "help",
			description: "Learn about the commands",
			callback:    commandHelp,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect the Pokemon",
			callback:    commandInspect,
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
		"pokedex": {
			name:        "pokedex",
			description: "List all caught Pokemon",
			callback:    commandPokedex,
		},
		"release": {
			name:        "release",
			description: "Release a caught Pokemon",
			callback:    commandRelease,
		},
	}
}

func cleanInput(text string) []string {
	strippedString := strings.TrimSpace(strings.ToLower(text))
	substrings := strings.Fields(strippedString)
	return substrings
}

func determineCatch(experience int) bool {
	const maxExperience = 635
	num := rand.Intn(maxExperience + 200)
	fmt.Printf("experience: %d\n", experience)
	fmt.Printf("number: %d\n", num)
	return num >= experience
}

func commandCatch(scanner *bufio.Scanner, ptr *config, cache *pokecache.Cache, pkmnCaught map[string]pokecache.PokemonInfo, param string) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", param)
	err := PokeAPIInteractions.GetPokemonInfo(&url, cache)
	if err != nil {
		return err
	}

	cachedData, ok := cache.Get(url)
	if !ok {
		return fmt.Errorf("could not get entry from catch")
	}

	var entry pokecache.PokemonInfo
	err = json.Unmarshal(cachedData, &entry)
	if err != nil {
		return err
	}
	pkmnName := entry.Forms[0].Name

	fmt.Printf("Throwing a Pokeball at the wild %s", pkmnName)
	time.Sleep(1 * time.Second)
	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\n")

	caught := determineCatch(entry.BaseExperience)
	if !caught {
		fmt.Printf("%s escaped!\n", pkmnName) // hardcoded, entry.Forms[0] because I think typing specific forms still works
		return nil
	}

	fmt.Printf("%s was caught!\n", pkmnName)

	if _, exist := pkmnCaught[pkmnName]; !exist {
		fmt.Printf("%s was added to the Pokedex!\n", pkmnName)
		pkmnCaught[pkmnName] = entry
	}

	return nil

}

func commandExit(scanner *bufio.Scanner, ptr *config, cache *pokecache.Cache, pkmnCaught map[string]pokecache.PokemonInfo, param string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandExplore(scanner *bufio.Scanner, ptr *config, cache *pokecache.Cache, pkmnCaught map[string]pokecache.PokemonInfo, param string) error {
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
		pkmnName := pkmn.Pokemon.Name
		_, exist := pkmnCaught[pkmnName]
		if exist {
			fmt.Printf(" - %s (caught)\n", pkmnName)
		} else {
			fmt.Printf(" - %s\n", pkmnName)
		}
	}

	return nil

}

func commandHelp(scanner *bufio.Scanner, ptr *config, cache *pokecache.Cache, pkmnCaught map[string]pokecache.PokemonInfo, param string) error {
	usageStr := ""
	for key, val := range commandsMap { // sort this!! currently returns in various orders
		currStr := fmt.Sprintf("%s: %s\n", key, val.description)
		usageStr += currStr
	}
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n%s", usageStr)
	return nil
}

func commandInspect(scanner *bufio.Scanner, ptr *config, cache *pokecache.Cache, pkmnCaught map[string]pokecache.PokemonInfo, param string) error {
	pkmn, exist := pkmnCaught[param]
	if !exist {
		return fmt.Errorf("you haven't caught that Pokemon")
	}

	fmt.Printf("Name: %s\n", pkmn.Forms[0].Name)
	fmt.Printf("Height: %d\n", pkmn.Height)
	fmt.Printf("Weight: %d\n", pkmn.Weight)
	fmt.Println("Stats:")
	for _, stat := range pkmn.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pkmnType := range pkmn.Types {
		fmt.Printf("  -%s\n", pkmnType.Type.Name)
	}

	return nil
}

func commandMap(scanner *bufio.Scanner, ptr *config, cache *pokecache.Cache, pkmnCaught map[string]pokecache.PokemonInfo, param string) error {
	err := PokeAPIInteractions.GetLocations(ptr.Next, cache)
	if err != nil {
		return err
	}

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

func commandMapb(scanner *bufio.Scanner, ptr *config, cache *pokecache.Cache, pkmnCaught map[string]pokecache.PokemonInfo, param string) error {
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

func commandPokedex(scanner *bufio.Scanner, ptr *config, cache *pokecache.Cache, pkmnCaught map[string]pokecache.PokemonInfo, param string) error {
	count := 0
	for key, _ := range pkmnCaught {
		count++
		if count == 1 {
			fmt.Println("Your Pokedex:")
		}
		fmt.Printf(" - %s\n", key)
	}
	if count == 0 {
		return fmt.Errorf("you haven't caught any pokemon")
	}

	return nil
}

func commandRelease(scanner *bufio.Scanner, ptr *config, cache *pokecache.Cache, pkmnCaught map[string]pokecache.PokemonInfo, param string) error {
	pkmn, exist := pkmnCaught[param]
	if !exist {
		return fmt.Errorf("you haven't caught that Pokemon")
	}

	pkmnName := pkmn.Forms[0].Name

	for {
		fmt.Printf("Are you sure you want to release %s? > ", pkmnName)
		if scanner.Scan() {
			text := scanner.Text()
			newText := cleanInput(text)
			if newText[0] == "yes" {
				delete(pkmnCaught, pkmnName)
				fmt.Printf("%s was released.\n", pkmnName)
				return nil
			} else if newText[0] == "no" {
				fmt.Printf("%s was not released.\n", pkmnName)
				return nil
			} else {
				fmt.Println("Please enter a valid command.")
			}
		}
	}
}
