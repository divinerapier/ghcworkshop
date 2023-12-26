package main

import (
	"encoding/json"
	"os"
	"sort"
)

// loadBeersFromFile loads a list of beers from a file and returns a pointer to the Beers struct and an error, if any.
// The file path is specified by the 'name' parameter.
// It opens the file, decodes the JSON data into a slice of Beer structs, sorts the beers, and returns the sorted Beers struct.
// If there is an error while opening the file or decoding the JSON data, it returns nil and the error.
func loadBeersFromFile(name string) (*Beers, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	beers := make([]Beer, 0)
	if err := json.NewDecoder(file).Decode(&beers); err != nil {
		return nil, err
	}

	bs := NewBeers(beers)
	sort.Sort(bs)

	return bs, nil
}
