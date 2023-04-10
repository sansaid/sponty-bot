package generator

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/sansaid/sponty/utils"
	funk "github.com/thoas/go-funk"
)

var (
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
	Location string `json:"location"`
}

type Chaplin string

type Perk string

type Intro string

type Adjective string

// Move to another package - helpers
// Dummy interface for the State struct:
// https://pkg.go.dev/github.com/diamondburned/arikawa/v3@v3.0.0-rc.4/state
type DiscordState interface {
	Roles(guildID discord.GuildID) ([]discord.Role, error)
	Members(guildID discord.GuildID) ([]discord.Member, error)
}

// Move to another package - helpers
func GetRole(roleName string, guildId discord.GuildID, client DiscordState) (discord.Role, error) {
	roles := utils.Must(client.Roles(guildId)).([]discord.Role)

	for _, role := range roles {
		if role.Name == roleName {
			return role, nil
		}
	}

	return discord.Role{}, fmt.Errorf("role %s not found", roleName)
}

// Move to another package - helpers
func GetRoleMembers(role discord.Role, guildId discord.GuildID, client DiscordState) ([]discord.Member, error) {
	// Need to reset internal cache, otherwise, for some freaky reason, the cache persists between restarts
	// (code suggests the cache is just a map, so not sure how it's possible that it persists between restarts)
	client.(*state.State).Cabinet.MemberStore.Reset()

	// Paginates at 1000 members per request by default - we're not going to page through because
	// we'll have other problems if the guild exceeds 1000 members
	members := utils.Must(client.Members(guildId)).([]discord.Member)
	roleMembers := []discord.Member{}

	for _, member := range members {
		// Yes, I downloaded a whole package just to check membership because of how much of a pain in the ass
		// it was to get Golang to DeepEquals various different types - this person seems to have done the job
		// for us. Maybe I could have spared us by just nesting another for loop in here. Optimisation for
		// another time.
		if funk.Contains(member.RoleIDs, role.ID) {
			roleMembers = append(roleMembers, member)
		}
	}

	return roleMembers, nil
}

func RandomLocation(locationType string) (Location, error) {
	var locationMap map[string][]Location

	if err := json.Unmarshal(locationsRaw, &locationMap); err != nil {
		return Location{}, err
	}

	locations, ok := locationMap[locationType]

	if !ok {
		return Location{}, fmt.Errorf("location type unrecognised: %s", locationType)
	}

	rng := rand.Intn(len(locations))

	return locations[rng], nil
}

func RandomChaplin(roleName string, guildId discord.GuildID, client DiscordState) (discord.UserID, error) {
	chaplinRole, err := GetRole(roleName, guildId, client)

	// TODO: error handle this better - shouldn't panic when we can't find the role.
	// Should send a proper error message instead to Discord users.
	chaplinMembers := utils.Must(GetRoleMembers(chaplinRole, guildId, client)).([]discord.Member)

	if len(chaplinMembers) == 0 {
		return discord.NullUserID, fmt.Errorf("no members with role %s", roleName)
	}

	if err != nil {
		return discord.NullUserID, err
	}

	rng := rand.Intn(len(chaplinMembers))
	chaplin := chaplinMembers[rng].User.ID

	return chaplin, nil
}

func RandomPerk(locationType string) (string, error) {
	var perkMap map[string][]Perk

	if err := json.Unmarshal(perksRaw, &perkMap); err != nil {
		return "", err
	}

	perks, ok := perkMap[locationType]

	if !ok {
		return "", fmt.Errorf("location type unrecognised: %s", locationType)
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
