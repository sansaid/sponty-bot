package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/sansaid/sponty/generator"
)

// To run, do `GUILD_ID="GUILD ID" BOT_TOKEN="TOKEN HERE" go run .`

func main() {
	guildID := discord.GuildID(mustSnowflakeEnv("GUILD_ID"))

	token := os.Getenv("BOT_TOKEN")

	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}

	// The State type is also a Session, which is also a Client (so it will inherit the interfaces of those
	// two nested types - that explains why we can't find CurrentApplication and RespondInteraction as a
	// method for State; they're methods for Client and/or Session)
	s := state.New("Bot " + token)

	app, err := s.CurrentApplication()

	if err != nil {
		log.Fatalln("Failed to get application ID:", err)
	}

	// InteractionCreateEvent type: https://pkg.go.dev/github.com/diamondburned/arikawa/v3@v3.0.0-rc.4/gateway#InteractionCreateEvent
	s.AddHandler(func(e *gateway.InteractionCreateEvent) {
		var data api.InteractionResponse

		if e.Data.InteractionType() == discord.CommandInteractionType {
			if e.InteractionEvent.Data.(*discord.CommandInteraction).Name == "rng-party" {
				data = api.InteractionResponse{
					Type: api.MessageInteractionWithSource,
					Data: &api.InteractionResponseData{
						Content: generatePartyData(),
					},
				}

				if err := s.RespondInteraction(e.ID, e.Token, data); err != nil {
					log.Println("failed to send interaction callback:", err)
				}
			}
		}
	})

	s.AddIntents(gateway.IntentGuilds)
	s.AddIntents(gateway.IntentGuildMessages)

	if err := s.Open(context.Background()); err != nil {
		log.Fatalln("failed to open:", err)
	}
	defer s.Close()

	log.Println("Gateway connected. Getting all guild commands.")

	commands, err := s.GuildCommands(app.ID, guildID)
	if err != nil {
		log.Fatalln("failed to get guild commands:", err)
	}

	for _, command := range commands {
		log.Println("Existing command", command.Name, "found. Deleting..")
		s.DeleteGuildCommand(app.ID, guildID, command.ID)
	}

	newCommands := []api.CreateCommandData{
		{
			Name:        "rng-party",
			Description: "Create a new randomly generated party",
		},
	}

	for _, command := range newCommands {
		_, err := s.CreateGuildCommand(app.ID, guildID, command)
		if err != nil {
			log.Fatalln("failed to create guild command:", err)
		}
	}

	// Block forever.
	select {}
}

func mustSnowflakeEnv(env string) discord.Snowflake {
	s, err := discord.ParseSnowflake(os.Getenv(env))
	if err != nil {
		log.Fatalf("Invalid snowflake for $%s: %v", env, err)
	}
	return s
}

func generatePartyData() *option.NullableStringData {
	chaplin, err := generator.RandomChaplin()
	if err != nil {
		log.Println("failed to get chaplin", err)
	}

	location, err := generator.RandomLocation()
	if err != nil {
		log.Println("failed to get location", err)
	}

	adjective, err := generator.RandomAdjective()
	if err != nil {
		log.Println("failed to get location", err)
	}

	msg := fmt.Sprintf("🥳⚠️🎉SPONTANEOUS NITE OUT BROADCAST🎉⚠️🥳"+
		"\nThe party chaplin is the %s %s"+
		"\nThe adventure begins at %s", adjective, chaplin, location)
	return option.NewNullableString(msg)
}
