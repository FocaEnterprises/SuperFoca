package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

var (
	GuildId string
	Token   string

	session *discordgo.Session
	db      *sql.DB
)

func init() {
	flag.Parse()

	var err error

	err = godotenv.Load(".env")

	if err != nil {
		log.Fatalf("couldn't read .env: %s", err)
	}

	if GuildId = os.Getenv("DISCORD_GUILD"); GuildId == "" {
		log.Fatal("failed reading DISCORD_GUILD")
	}

	if Token = os.Getenv("DISCORD_TOKEN"); Token == "" {
		log.Fatal("failed reading DISCORD_TOKEN")
	}

	session, err = discordgo.New("Bot " + Token)

	if err != nil {
		log.Fatalf("failed to authorize bot: %s", err)
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("POSTGRES_PORT")

	db, err = sql.Open("postgres", fmt.Sprintf("host=127.0.0.1 port=%s user=%s password=%s dbname=%s sslmode=disable", dbPort, dbUser, dbPassword, dbName))

	if err != nil {
		log.Fatalf("failed to establish database connection: %s", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to reach database: %s", err)
	}

	log.Println("succesfully established database connection")
}

func main() {
	session.Identify.Intents = discordgo.IntentsAll

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("logged in as %s", session.State.User.String())
	})

	session.AddHandler(iqIncreaseHandler)

	if err := session.Open(); err != nil {
		log.Fatalf("failed to start bot client: %s", err)
	}

	registerSlashCommands()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	clearSlashCommands()

	if err := session.Close(); err != nil {
		log.Printf("failed closing session: %s", err)
	}
}
