package bot

import (
	"log"
	"os"
	"slices"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildId string
	Token   string

	Session *discordgo.Session
)

func Init() {
	var err error
	GuildId = os.Getenv("DISCORD_GUILD")
	Token = os.Getenv("DISCORD_TOKEN")

	if slices.Contains([]string{GuildId, Token}, "") {
		log.Fatal("FAILED RETRIEVING ENVIRONMENT VARIABLES")
	}

	Session, err = discordgo.New("Bot " + Token)

	if err != nil {
		log.Fatalf("failed to authorize bot: %s", err)
	}

	Session.Identify.Intents = discordgo.IntentsAll

	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("logged in as %s", Session.State.User.String())
	})

	if err := Session.Open(); err != nil {
		log.Fatalf("failed to start bot client: %s", err)
	}

	Session.UpdateCustomStatus("vai mengão!!!!")

	Session.AddHandler(IQIncreaseHandler)
	Session.AddHandler(SemPutariaHandler)
}
