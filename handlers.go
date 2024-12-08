package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/api/youtube/v3"
)

var (
	registeredCommands []*discordgo.ApplicationCommand
	commands           []*discordgo.ApplicationCommand
)

func initCommands() {
	response, err := http.Get("http://localhost:8080/playlists")

	if err != nil {
		log.Printf("failed retrieving playlists: %s", err)
		return
	}

	playlists := make(map[string]any)

	err = json.NewDecoder(response.Body).Decode(&playlists)

	if err != nil {
		log.Printf("failed decoding playlists: %s", err)
		return
	}

	playlistChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0)

	for k := range playlists {
		response, err := http.Get("http://localhost:8080/playlists/" + k)

		if err != nil {
			log.Printf("failed retrieving playlist title: %s", err)
			return
		}

		var playlistData map[string]any

		err = json.NewDecoder(response.Body).Decode(&playlistData)

		if err != nil {
			log.Printf("failed decoding playlist data: %s", err)
			return
		}

		defer response.Body.Close()

		snippet := playlistData["snippet"].(map[string]any)

		playlistChoices = append(playlistChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  snippet["title"].(string),
			Value: k,
		})
	}

	commands = []*discordgo.ApplicationCommand{
		{
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
		},
		{
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
		},
		{
			Name:        "tunes",
			Description: "Manipulação das playlists da Foca Tunes",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "list",
					Description: "Lista as músicas de uma playlist.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "playlist",
							Description: "A playlist que vai ser consultada.",
							Choices:     playlistChoices,
							Required:    true,
						},
					},
				},
				{
					Name:        "add",
					Description: "Adiciona uma música à playlist.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "playlist",
							Description: "A playlist em que a música vai ser adicionada.",
							Required:    true,
							Choices:     playlistChoices,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "música",
							Description: "Um link ou ID de música no YouTube.",
							Required:    true,
						},
					},
				},
			},
		},
	}

	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))

	registerSlashCommands()
}

func parseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) *map[string]*discordgo.ApplicationCommandInteractionDataOption {
	opts := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)

	for _, opt := range options {
		opts[opt.Name] = opt
	}

	return &opts
}

func interactionRespondEphemeral(i *discordgo.InteractionCreate, message string) {
	session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func interactionRespond(i *discordgo.InteractionCreate, message string) {
	session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func tunesList(i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	//o := *parseOptions(options)

}

func tunesAdd(i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	o := *parseOptions(options)

	playlist := o["playlist"].StringValue()

	var songId string
	urlStruct, _ := url.Parse(o["música"].StringValue())

	songId = o["música"].StringValue()

	if slices.Contains([]string{"music.youtube.com", "youtube.com"}, urlStruct.Host) {
		songId = urlStruct.Query().Get("v")

		if songId == "" {
			interactionRespondEphemeral(i, "Não consegui retirar um ID válido.")
			return
		}
	}

	videoJSON, err := http.Get(fmt.Sprintf("http://localhost:8080/videos?id=%s", songId))

	if err != nil {
		interactionRespondEphemeral(i, "Não consegui me comunicar com a API.")
		return
	}

	defer videoJSON.Body.Close()

	var video youtube.Video

	err = json.NewDecoder(videoJSON.Body).Decode(&video)

	if err != nil {
		log.Printf("couldn't decode snippet: %s", err)
		return
	}

	_, err = http.Post(fmt.Sprintf("http://localhost:8080/playlists/%s/songs", playlist), "application/json", videoJSON.Body)

	if err != nil {
		log.Printf("couldn't post song: %s", err)
		interactionRespondEphemeral(i, "Não consegui me comunicar com a API.")
		return
	}

	session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						Name: video.Snippet.ChannelTitle,
						URL:  "https://music.youtube.com/" + i.ChannelID,
					},
					Image: &discordgo.MessageEmbedImage{
						URL: video.Snippet.Thumbnails.Maxres.Url,
					},
					URL:   "https://music.youtube.com/watch?v=" + video.Id,
					Color: 0xFFFF00,
					Title: video.Snippet.Title,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Playlist",
							Value: cases.Title(language.BrazilianPortuguese).String(playlist),
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    fmt.Sprintf("Adicionado por %s", i.Member.User.Username),
						IconURL: i.Member.User.AvatarURL(""),
					},
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				},
			},
		},
	})
}

var slashHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ranking": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ranking := getRanking(0, 10)

		if ranking == nil {
			interactionRespondEphemeral(i, "Não consegui pegar o ranking...")
			return
		}

		buffer := &strings.Builder{}

		buffer.WriteString("```md\n")

		for n, v := range ranking {
			title := readTitleFromRank(*v)

			member, err := session.GuildMember(v.GuildId, v.UserId)

			if err != nil {
				interactionRespondEphemeral(i, "Não consegui pegar o ranking...")
				return
			}

			buffer.WriteString(fmt.Sprintf("%3v. QI %6.2f · %s, %s\n", n+1, v.IQ, member.User.Username, title.Title))
		}

		buffer.WriteString("```")

		interactionRespond(i, buffer.String())
	},
	"iq": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := *parseOptions(i.ApplicationCommandData().Options)

		var user *discordgo.User

		if options["usuário"] != nil {
			user = options["usuário"].UserValue(s)
		} else {
			user = i.Member.User
		}

		rank := findRank(user.ID, i.GuildID)

		if rank == nil {
			interactionRespond(i, fmt.Sprintf("O usuário %s nem sequer tem QI...", user.Username))
			return
		}

		title := readTitleFromRank(*rank)

		interactionRespond(i, fmt.Sprintf("%s tem QI de %.02f! Ranking %s", user.Username, rank.IQ, title.Title))
	},
	"tunes": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ChannelID != "1097228819087237240" {
			interactionRespondEphemeral(i, fmt.Sprintf("Esse comando só pode ser usado no canal da Foca Tunes!"))
			return
		}

		if i.Member.User.Bot == true {
			interactionRespondEphemeral(i, "Você não tem permissão para isso!")
			return
		}

		focaTunesRole := "1308938228589400095"

		if !slices.Contains(i.Member.Roles, focaTunesRole) {
			interactionRespondEphemeral(i, "Você não tem permissão para isso!")
			return
		}

		options := i.ApplicationCommandData().Options

		switch options[0].Name {
		case "add":
			tunesAdd(i, options[0].Options)
			break
		case "list":
			tunesList(i, options[0].Options)
			break
		default:
			break
		}

	},
	"echo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := *parseOptions(i.ApplicationCommandData().Options)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Enviado...",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		s.ChannelMessageSend(i.ChannelID, options["mensagem"].StringValue())
	},
}

func registerSlashCommands() {
	for _, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, GuildId, v)

		if err != nil {
			log.Fatalf("failed to register command: %s", err)
		}

		registeredCommands = append(registeredCommands, cmd)
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := slashHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func clearSlashCommands() {
	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, v.GuildID, v.ID)

		if err != nil {
			log.Printf("error deleting application command: %s", err)
		}
	}
}

func getIQIncrease(iq float64) float64 {
	base := 0.1
	k := 0.1
	e := math.Pow(math.E, -k*(iq+1))
	return base*e + 0.001
}

func semPutariaHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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

func iqIncreaseHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	userId := m.Author.ID
	guildId := m.GuildID

	if userId == s.State.User.ID {
		return
	}

	if m.Author.Bot {
		return
	}

	rank := findRank(userId, guildId)

	if rank == nil {
		rank = createRank(userId, guildId)
	}

	currentIQ := rank.IQ
	newIQ := currentIQ + getIQIncrease(currentIQ)

	updateRank(*rank, newIQ)

	title := readTitleFromRank(*rank)

	if int(newIQ) != int(currentIQ) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**%s** agora tem **%.02f** de QI! Seu ranking é **%s**.", m.Author.Mention(), newIQ, title.Title))
	}
}
