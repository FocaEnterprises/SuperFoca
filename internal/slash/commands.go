package slash

import (
	"log"
	"superfoca/internal/bot"

	"github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	Command *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
	Commands map[string]*SlashCommand
)

func init() {
	Commands = make(map[string]*SlashCommand)
}

func Init(session *discordgo.Session) {
	bot.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if c, ok := Commands[i.ApplicationCommandData().Name]; ok {
			c.Handler(s, i)
		}
	})
}

func Add(command *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	if _, added := Commands[command.Name]; !added {
		Commands[command.Name] = &SlashCommand{
			Command: command,
			Handler: handler,
		}
	}
}

func Register() error {
	for _, v := range Commands {
		cmd, err := bot.Session.ApplicationCommandCreate(bot.Session.State.User.ID, bot.GuildId, v.Command)

		if err != nil {
			return err
		}

		Commands[cmd.Name].Command = cmd
	}

	return nil
}

func Clear() {
	for _, v := range Commands {
		err := bot.Session.ApplicationCommandDelete(bot.Session.State.User.ID, v.Command.GuildID, v.Command.ID)

		if err != nil {
			log.Printf("error deleting application command: %s", err)
		}
	}
}

func RespondInteractionEphemeralString(i *discordgo.InteractionCreate, message string) {
	bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func RespondInteraction(i *discordgo.InteractionCreate, r *discordgo.InteractionResponse) {
	bot.Session.InteractionRespond(i.Interaction, r)
}

func RespondInteractionString(i *discordgo.InteractionCreate, message string) {
	bot.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func ParseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) *map[string]*discordgo.ApplicationCommandInteractionDataOption {
	opts := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)

	for _, opt := range options {
		opts[opt.Name] = opt
	}

	return &opts
}
