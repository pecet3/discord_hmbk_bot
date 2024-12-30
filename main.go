package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/pecet3/discord_hmbk_bot/discord"
	"github.com/pecet3/discord_hmbk_bot/pkg/message"
	"github.com/pecet3/discord_hmbk_bot/pkg/paint"
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
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	ps := paint.NewPaintSessions()

	discord.Run(dc, ps)
	paint.Run(mux, ps)
	message.Run(mux, dc)
	http.ListenAndServe(":8080", mux)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	dc.Close()
	log.Println("Interupted")
}
