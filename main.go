package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"webscraping/discord"
	"webscraping/paint"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("[👍] Loaded .env")
}

func main() {
	loadEnv()
	discordToken := os.Getenv("DISCORD_TOKEN")
	dc, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal(err)
		return
	}

	ps := paint.NewPaintSessions()

	discord.Run(dc, ps)
	paint.Run(ps)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dc.Close()
	log.Println("Interupted")
}
