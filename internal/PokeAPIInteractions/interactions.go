package PokeAPIInteractions

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BatchInfo struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocations(url *string) (BatchInfo, error) { // Do I need to make a separate request for each location ID?
	if url == nil {
		return BatchInfo{}, fmt.Errorf("you're on the first page")
	}

	res, err := http.Get(*url)
	if err != nil {
		return BatchInfo{}, err
	}

	defer res.Body.Close()

	var locations BatchInfo

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&locations); err != nil {
		return BatchInfo{}, err
	}

	// fmt.Print(locations.Previous)
	// fmt.Printf(" is nil: %t\n", locations.Previous == nil)
	// fmt.Println(locations.Next)
	// fmt.Printf(" is nil: %t\n", locations.Next == nil)

	return locations, nil
}
