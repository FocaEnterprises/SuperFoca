package slash

import (
	"fmt"
	"strings"
	"superfoca/internal/bot"
	"superfoca/internal/database"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Add(&discordgo.ApplicationCommand{
			Name:        "ranking",
			Description: "Consulta o QI de um usuário",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "usuário",
					Description: "O usuário a ser consultado",
					Required:    false,
				},
			},
		},
	func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ranking := database.GetRanking(0, 10)

		if ranking == nil {
			RespondInteractionEphemeralString(i, "Não consegui pegar o ranking...")
			return
		}

		buffer := &strings.Builder{}

		buffer.WriteString("```\n")

		for n, v := range ranking {
			title := database.ReadTitleFromRank(*v)

			member, err := bot.Session.GuildMember(v.GuildId, v.UserId)

			if err != nil {
				RespondInteractionEphemeralString(i, "Não consegui pegar o ranking...")
				return
			}

			buffer.WriteString(fmt.Sprintf("%3v. QI %6.2f · %s, %s\n", n+1, v.IQ, member.User.Username, title.Title))
		}

		buffer.WriteString("```")

		RespondInteractionString(i, buffer.String())
	})
}
