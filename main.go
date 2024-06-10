package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	discordToken := os.Getenv("DISCORD_TOKEN")
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal(err)
		return
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		log.Println(m.Content)
		log.Println(m.ChannelID)
		log.Println(m.Author)
		if m.Content == "hello" {
			log.Println("user hi")
			s.ChannelMessageSend(m.ChannelID, "world")
		}
		s.ChannelMessageSend(m.ChannelID, m.Message.Content+"test")
		s.ChannelMessageSend(m.ChannelID, m.Author.Username)

	})

	if err = discord.Open(); err != nil {
		log.Fatal(err)
		return
	}
	defer discord.Close()
	discord.Identify.Intents = discordgo.IntentsAll | discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	log.Println("[*] HMBK Bot is online")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Loaded .env")
}
