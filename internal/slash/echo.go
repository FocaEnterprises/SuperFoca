package slash

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	Add(&discordgo.ApplicationCommand{
		Name:        "echo",
		Description: "Eco! (Eco!)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "mensagem",
				Description: "A mensagem a ser ecoada",
				Required:    true,
			},
		},
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := *ParseOptions(i.ApplicationCommandData().Options)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Enviado...",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		s.ChannelMessageSend(i.ChannelID, options["mensagem"].StringValue())
	})
}
