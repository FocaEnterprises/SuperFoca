package slash

import (
	"fmt"
	"superfoca/internal/database"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Add(&discordgo.ApplicationCommand{
		Name:        "iq",
		Description: "Consulta o QI de um usuário",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "usuário",
				Description: "O usuário a ser consultado",
				Required:    false,
			},
		},
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := *ParseOptions(i.ApplicationCommandData().Options)

		var user *discordgo.User

		if options["usuário"] != nil {
			user = options["usuário"].UserValue(s)
		} else {
			user = i.Member.User
		}

		rank := database.FindRankByUser(user.ID, i.GuildID)

		if rank == nil {
			RespondInteractionString(i, fmt.Sprintf("O usuário %s nem sequer tem QI...", user.Username))
			return
		}

		title := database.ReadTitleFromRank(*rank)

		RespondInteractionString(i, fmt.Sprintf("%s tem QI de %.02f! Ranking %s", user.Username, rank.IQ, title.Title))
	})
}
