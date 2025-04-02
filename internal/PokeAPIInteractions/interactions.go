package PokeAPIInteractions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cmolloy36/Pokedex/internal/pokecache"
)

func GetLocations(url *string, c *pokecache.Cache) error {
	if url == nil || *url == "" {
		return fmt.Errorf("you're on the first page")
	}

	_, exist := c.Get(*url)

	if !exist {
		// fmt.Println("adding key...")
		res, err := http.Get(*url)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		var locations pokecache.BatchInfo

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&locations); err != nil {
			return err
		}

		jsonBytes, err := json.Marshal(locations)
		if err != nil {
			return err
		}

		// fmt.Println(jsonBytes)

		c.Add(*url, jsonBytes)
		// fmt.Println(c.Get(*url))
	}

	// fmt.Print("using cache!")
	// elem, _ := c.Get(*url)
	// fmt.Printf("%v", elem.Results)

	// fmt.Print(locations.Previous)
	// fmt.Printf(" is nil: %t\n", locations.Previous == nil)
	// fmt.Println(locations.Next)
	// fmt.Printf(" is nil: %t\n", locations.Next == nil)

	return nil
}
