package bot

import (
	"fmt"
	"math"
	"strings"
	"superfoca/internal/database"

	"github.com/bwmarrin/discordgo"
)

func GetIQIncrease(iq float64) float64 {
	base := 0.1
	k := 0.1
	e := math.Pow(math.E, -k*(iq+1))
	return base*e + 0.001
}

func SemPutariaHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	userId := m.Author.ID

	if userId == s.State.User.ID {
		return
	}

	if m.Author.Bot {
		return
	}

	prohibitedWords := []string{
		"sexo",
		"puta",
		"vagabunda",
		"gozar",
		"punheta",
		"bronha",
		"masturbar",
		"masturbacao",
		"pornografia",
		"punhetao",
		"punhetão",
		"pornô",
		"porno",
		"pinto",
		"penis",
		"buceta",
		"xereca",
		"sex",
	}

	lower := strings.ToLower(m.Content)
	for _, v := range prohibitedWords {
		if strings.Contains(lower, v) {
			s.ChannelMessageSendReply(m.ChannelID, m.Author.Mention()+" SEM PUTARIA!!! :fire::fire::fire::speaking_head::speaking_head::speaking_head:", m.Reference())
			return
		}
	}
}

func IQIncreaseHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	userId := m.Author.ID
	guildId := m.GuildID

	if userId == s.State.User.ID {
		return
	}

	if m.Author.Bot {
		return
	}

	rank := database.FindRankByUser(userId, guildId)

	if rank == nil {
		rank = database.CreateRank(userId, guildId)
	}

	currentIQ := rank.IQ
	newIQ := currentIQ + GetIQIncrease(currentIQ)

	database.UpdateRank(*rank, newIQ)

	title := database.ReadTitleFromRank(*rank)

	if int(newIQ) != int(currentIQ) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**%s** agora tem **%.02f** de QI! Seu ranking é **%s**.", m.Author.Mention(), newIQ, title.Title))
	}
}
