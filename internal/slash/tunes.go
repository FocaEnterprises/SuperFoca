package slash

import (
	"encoding/json"
	"fmt"
	"log"
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

type FocaTunesVideoResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    *youtube.Video `json:"data"`
}

type FocaTunesPlaylistsResponse struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    map[string]*youtube.Playlist `json:"data"`
}

type FocaTunesPlaylistResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    *youtube.Playlist `json:"data"`
}

func tunesList(i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
}

func tunesAdd(i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	o := *ParseOptions(options)

	playlist := o["playlist"].StringValue()

	var songId string
	urlStruct, _ := url.Parse(o["música"].StringValue())

	songId = o["música"].StringValue()

	if slices.Contains([]string{"music.youtube.com", "youtube.com"}, urlStruct.Host) {
		songId = urlStruct.Query().Get("v")

		if songId == "" {
			RespondInteractionEphemeralString(i, "Não consegui retirar um ID válido.")
			return
		}
	}

	videoJSON, err := http.Get(fmt.Sprintf("http://localhost:8080/videos?id=%s", songId))

	if err != nil {
		RespondInteractionEphemeralString(i, "Não consegui me comunicar com a API.")
		return
	}

	defer videoJSON.Body.Close()

	var video FocaTunesVideoResponse

	err = json.NewDecoder(videoJSON.Body).Decode(&video)

	if err != nil {
		log.Printf("couldn't decode snippet: %s", err)
		return
	}

	if videoJSON.StatusCode != http.StatusOK {
		log.Printf("carai")
		return
	}

	_, err = http.Post(fmt.Sprintf("http://localhost:8080/playlists/%s", playlist), "application/json", strings.NewReader(fmt.Sprintf("{ \"resourceId\": %q }", video.Data.Id)))

	if err != nil {
		log.Printf("couldn't post song: %s", err)
		RespondInteractionEphemeralString(i, "Não consegui me comunicar com a API.")
		return
	}

	RespondInteraction(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Author: &discordgo.MessageEmbedAuthor{
						Name: video.Data.Snippet.ChannelTitle,
						URL:  "https://music.youtube.com/" + i.ChannelID,
					},
					Image: &discordgo.MessageEmbedImage{
						URL: video.Data.Snippet.Thumbnails.Maxres.Url,
					},
					URL:   "https://music.youtube.com/watch?v=" + video.Data.Id,
					Color: 0xFFFF00,
					Title: video.Data.Snippet.Title,
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

func init() {
	response, err := http.Get("http://localhost:8080/playlists")

	if err != nil {
		log.Printf("failed retrieving playlists: %s", err)
		return
	}

	playlists := FocaTunesPlaylistsResponse{}

	err = json.NewDecoder(response.Body).Decode(&playlists)

	if err != nil {
		log.Printf("failed decoding playlists: %s", err)
		return
	}

	if playlists.Status != "success" {
		log.Println("failed")
		return
	}

	playlistChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(playlists.Data))

	for v, k := range playlists.Data {
		playlistChoices = append(playlistChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  k.Snippet.Title,
			Value: v,
		})
	}

	Add(&discordgo.ApplicationCommand{
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
	}, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ChannelID != "1097228819087237240" {
			RespondInteractionEphemeralString(i, fmt.Sprintf("Esse comando só pode ser usado no canal da Foca Tunes!"))
			return
		}

		if i.Member.User.Bot == true {
			RespondInteractionEphemeralString(i, "Você não tem permissão para isso!")
			return
		}

		focaTunesRole := "1308938228589400095"

		if !slices.Contains(i.Member.Roles, focaTunesRole) {
			RespondInteractionEphemeralString(i, "Você não tem permissão para isso!")
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
	})
}
