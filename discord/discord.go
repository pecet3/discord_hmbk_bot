package discord

import (
	"log"
	"os"
	"time"
	"webscraping/scraper"

	"github.com/bwmarrin/discordgo"
)

type command struct {
}

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
		prefix := m.Content[0]
		if string(prefix) != "$" {
			return
		}
		command := string(m.Content[1:])

		switch command {
		case "szczytno":

			nsz, _ := scrap.GetCity("szczytno")

			articles := nsz.Scraper.GetArticles(nsz)

			for _, a := range articles[0:4] {

				img := &discordgo.MessageEmbedImage{
					URL:    a.Image,
					Width:  400,
					Height: 280,
				}

				emb := &discordgo.MessageEmbed{
					URL:         a.Link,
					Title:       a.Title + "\n " + a.Date,
					Description: a.Content,
					Image:       img,
				}
				s.ChannelMessageSendEmbed(m.ChannelID, emb)

			}
			scrap.SaveCity(nsz)

		case "nsz":
			emb := &discordgo.MessageEmbed{
				Title: "Test",
			}
			nsz, _ := scrap.GetCity("szczytno")
			s.ChannelMessageSendEmbed(m.ChannelID, emb)
			articles := nsz.Scraper.GetArticles(nsz)
			display := ""
			for _, a := range articles[:4] {
				display = display + "## " + a.Title + "\n" + "### " + a.Date + "\n" + a.Content + "\n"
			}
			log.Println(display)
			scrap.SaveCity(nsz)
			s.ChannelMessageSend(m.ChannelID, string(display))
		}

	})

	if err = discord.Open(); err != nil {
		log.Fatal(err)
		return
	}
	defer discord.Close()
	discord.Identify.Intents = discordgo.IntentsAll | discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	log.Println("[🟢] Discord Bot is online")
}
