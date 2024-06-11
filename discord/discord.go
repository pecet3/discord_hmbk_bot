package discord

import (
	"log"
	"os"
	"time"
	"webscraping/scraper"

	"github.com/bwmarrin/discordgo"
)

func Run() {
	discordToken := os.Getenv("DISCORD_TOKEN")
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal(err)
		return
	}
	scrap := scraper.New()

	scrap.CitiesMap["szczytno"] = &scraper.CityNews{
		Name:      "szczytno",
		ExpiresAt: time.Now(),
		Scraper:   scraper.NszScraper{},
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		log.Println(m.Content)
		log.Println(m.ChannelID)
		log.Println(m.Author)

		prefix := m.Content[0]
		if string(prefix) != "$" {
			return
		}
		command := string(m.Content[1:])

		switch command {
		case "szczytno":
			emb := &discordgo.MessageEmbed{
				Title: "Test",
			}
			nsz, _ := scrap.GetCity("szczytno")
			s.ChannelMessageSendEmbed(m.ChannelID, emb)
			articles := nsz.Scraper.GetArticles(nsz)
			display := ""
			for _, a := range articles[:16] {
				display = display + " ## " + a.Title + "\n"
			}
			log.Println(display)
			s.ChannelMessageSend(m.ChannelID, string(display))
		}

	})

	if err = discord.Open(); err != nil {
		log.Fatal(err)
		return
	}
	defer discord.Close()
	discord.Identify.Intents = discordgo.IntentsAll | discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	log.Println("[+] Discord Bot is online")
}
