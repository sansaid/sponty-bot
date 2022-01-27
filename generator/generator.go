package generator

import (
	_ "embed"
	"encoding/json"
	"math/rand"
)

// TODO: implement filters
// Random location generator

type Location struct {
	Name     string `json:"name"`
	Site     string `json:"site"`
	Location string `json:"location"`
	Type     string `json:"type"`
}

type Chaplin string

func RandomLocation() (string, error) {
	//go:embed ./../data/locations.json
	var locationsRaw []byte
	var locations []Location

	if err := json.Unmarshal(locationsRaw, &locations); err != nil {
		return "", err
	}

	rng := rand.Intn(len(locations))
	location := locations[rng].Name

	return location, nil
}

func RandomChaplin() (string, error) {
	//go:embed ./../data/party_chaplins.json
	var chaplinsRaw []byte
	var chaplins []Chaplin

	if err := json.Unmarshal(chaplinsRaw, &chaplins); err != nil {
		return "", err
	}

	rng := rand.Intn(len(chaplins))
	chaplin := string(chaplins[rng])

	return chaplin, nil
}
