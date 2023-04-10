package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/sansaid/sponty/generator"
)

// To run, do `GUILD_ID=<GUILD ID> BOT_TOKEN=<TOKEN HERE> go run .`

const (
	CHAPLIN_ROLE_NAME = "Party Chaplins"
)

var commands = []api.CreateCommandData{
	{
		Name:        "rng-party",
		Description: "Create a new randomly generated party",
		Options: []discord.CommandOption{
			&discord.BooleanOption{
				OptionName:  "generate_chaplin",
				Description: fmt.Sprintf("(default: true) Whether to randomly generate a party chaplin for your party. This will randomly draw from people who have the role %s", CHAPLIN_ROLE_NAME),
				Required:    false,
			},
			&discord.StringOption{
				OptionName:  "location_type",
				Description: "(default: \"pub\") Which location type to host your party",
				Required:    false,
				Choices: []discord.StringChoice{
					{
						Name:  "pub",
						Value: "pub",
					},
					{
						Name:  "park",
						Value: "park",
					},
				},
			},
		},
	},
}

type handler struct {
	*cmdroute.Router
	state   *state.State
	guildID discord.GuildID
}

func newHandler(s *state.State, g discord.GuildID) *handler {
	h := &handler{
		state:   s,
		guildID: g,
	}

	h.Router = cmdroute.NewRouter()
	h.state.AddInteractionHandler(h)

	// Automatically defer handles if they're slow.
	h.Use(cmdroute.Deferrable(s, cmdroute.DeferOpts{}))
	h.initIntents()
	h.initHandlers()
	h.overwriteCommands(commands)

	return h
}

func (h *handler) overwriteCommands(newCommands []api.CreateCommandData) {
	app, err := h.state.CurrentApplication()

	if err != nil {
		log.Fatalln("Failed to get application ID:", err)
	}

	oldCommands, err := h.state.GuildCommands(app.ID, h.guildID)
	if err != nil {
		log.Fatalln("failed to get guild commands:", err)
	}

	for _, command := range oldCommands {
		log.Println("Existing command", command.Name, "found. Deleting..")
		h.state.DeleteGuildCommand(app.ID, h.guildID, command.ID)
	}

	for _, command := range newCommands {
		_, err := h.state.CreateGuildCommand(app.ID, h.guildID, command)
		if err != nil {
			log.Fatalln("failed to create guild command:", err)
		}
	}
}

func (h *handler) initIntents() {
	h.state.AddIntents(gateway.IntentGuilds)
	h.state.AddIntents(gateway.IntentGuildMessages)
	h.state.AddIntents(gateway.IntentGuildMembers)
}

func (h *handler) initHandlers() {
	// InteractionCreateEvent type: https://pkg.go.dev/github.com/diamondburned/arikawa/v3@v3.0.0-rc.4/gateway#InteractionCreateEvent
	h.state.AddHandler(func(*gateway.ReadyEvent) {
		me, _ := h.state.Me()
		log.Println("Connected to the gateway as ", me.Tag())
	})

	h.AddFunc("rng-party", h.cmdRngParty)
}

func (h *handler) run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := h.state.Connect(ctx); err != nil {
		log.Fatalln("failed to connect:", err)
	}
}

func main() {
	guildID := discord.GuildID(mustSnowflakeEnv("GUILD_ID"))
	token := os.Getenv("BOT_TOKEN")

	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}

	// The State type is also a Session, which is also a Client (so it will inherit the interfaces of those
	// two nested types - that explains why we can't find CurrentApplication and RespondInteraction as a
	// method for State; they're methods for Client and/or Session)
	h := newHandler(state.New("Bot "+token), guildID)

	h.run()
}

// Move to utils
func mustSnowflakeEnv(env string) discord.Snowflake {
	s, err := discord.ParseSnowflake(os.Getenv(env))
	if err != nil {
		log.Fatalf("Invalid snowflake for $%s: %v", env, err)
	}
	return s
}

func (h *handler) cmdRngParty(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
	var msgArray []string
	var chaplin discord.UserID
	var options struct {
		GenerateChaplin bool   `discord:"generate_chaplin"`
		LocationType    string `discord:"location_type"`
	}

	if err := data.Options.Unmarshal(&options); err != nil {
		return errorResponse(err)
	}

	if options.LocationType == "" {
		options.LocationType = "pub"
	}

	intro, err := generator.RandomIntro()
	if err != nil {
		log.Println("failed to get intro:", err)
	}
	msgArray = append(msgArray, intro)

	if options.GenerateChaplin {
		var err error
		chaplin, err = generator.RandomChaplin(CHAPLIN_ROLE_NAME, h.guildID, h.state)

		if err != nil {
			log.Println("failed to get party chaplins:", err)

			return persistentResponse(fmt.Sprintf("**What!?** No party chaplins!?"+
				" Make yourself useful and create a **%s** role with members.", CHAPLIN_ROLE_NAME))
		}

		msgArray = append(msgArray, fmt.Sprintf(":levitate_tone1: Your party chaplin is <@%s>", chaplin.String()))
	}

	location, err := generator.RandomLocation(options.LocationType)
	if err != nil {
		return errorResponse(fmt.Errorf("failed to get location: %w", err))
	}
	msgArray = append(msgArray, fmt.Sprintf("🍾 The adventure begins at **%s** (%s)", location.Name, location.Location))

	perk, err := generator.RandomPerk(options.LocationType)
	if err != nil {
		return errorResponse(fmt.Errorf("failed to get perk: %w", err))
	}
	msgArray = append(msgArray, fmt.Sprintf("📖 Tonight's golden rule: %s", perk))

	msg := strings.Join(msgArray, "\n\r")

	return persistentResponse(msg)
}

func errorResponse(err error) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content:         option.NewNullableString("**Error:** " + err.Error()),
		Flags:           discord.EphemeralMessage,
		AllowedMentions: &api.AllowedMentions{ /* none */ },
	}
}

func persistentResponse(msg string) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content: option.NewNullableString(msg),
	}
}
