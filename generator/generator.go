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

	//go:embed data/party_perks.json
	perksRaw []byte

	//go:embed data/intros.json
	introsRaw []byte

	//go:embed data/adjectives.json
	adjectivesRaw []byte
)

type Location struct {
	Name     string `json:"name"`
	Site     string `json:"site"`
	Location string `json:"location"`
	Type     string `json:"type"`
}

type Chaplin string

type Perk string

type Intro string

type Adjective string

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

func RandomPerk() (string, error) {
	var perks []Perk

	if err := json.Unmarshal(perksRaw, &perks); err != nil {
		return "", err
	}

	rng := rand.Intn(len(perks))
	perk := string(perks[rng])

	return perk, nil
}

func RandomIntro() (string, error) {
	var intros []Intro

	if err := json.Unmarshal(introsRaw, &intros); err != nil {
		return "", err
	}

	rng := rand.Intn(len(intros))
	intro := string(intros[rng])

	return intro, nil
}

func RandomAdjective() (string, error) {
	var adjectives []Adjective

	if err := json.Unmarshal(adjectivesRaw, &adjectives); err != nil {
		return "", err
	}

	rng := rand.Intn(len(adjectives))
	adjective := string(adjectives[rng])

	return adjective, nil
}
