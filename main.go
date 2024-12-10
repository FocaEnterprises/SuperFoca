package main

import (
	"log"
	"os"
	"os/signal"
	"superfoca/internal/bot"
	"superfoca/internal/database"
	"superfoca/internal/slash"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	database.Init()

	bot.Init()
	slash.Init(bot.Session)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	if err := bot.Session.Close(); err != nil {
		log.Printf("failed closing bot.Session: %s", err)
	}
}
