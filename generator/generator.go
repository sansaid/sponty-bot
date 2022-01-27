package generator

import (
	_ "embed"
	"encoding/json"
	"math/rand"
)

// TODO: implement filters
// Random location generator

var (
	//go:embed data/party_chaplins.json
	chaplinsRaw []byte

	//go:embed data/locations.json
	locationsRaw []byte
)

type Location struct {
	Name     string `json:"name"`
	Site     string `json:"site"`
	Location string `json:"location"`
	Type     string `json:"type"`
}

type Chaplin string

func RandomLocation() (string, error) {
	var locations []Location

	if err := json.Unmarshal(locationsRaw, &locations); err != nil {
		return "", err
	}

	rng := rand.Intn(len(locations))
	location := locations[rng].Name

	return location, nil
}

func RandomChaplin() (string, error) {
	var chaplins []Chaplin

	if err := json.Unmarshal(chaplinsRaw, &chaplins); err != nil {
		return "", err
	}

	rng := rand.Intn(len(chaplins))
	chaplin := string(chaplins[rng])

	return chaplin, nil
}
