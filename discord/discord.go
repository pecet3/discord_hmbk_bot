package discord

import (
	"log"
	"strings"
	"time"
	"webscraping/scraper"

	"github.com/bwmarrin/discordgo"
)

const (
	PREFIX = "!"
	NSZ    = "nsz"
	DAY    = "dzien"
	IQ     = "iq"
)

func Run(discord *discordgo.Session) {

	scrap := scraper.New()

	scrap.PagesMap["szczytno"] = &scraper.Page{
		Name:      "szczytno",
		ExpiresAt: time.Now(),
		Scraper:   scraper.NszScraper{},
	}
	scrap.PagesMap["day"] = &scraper.Page{
		Name:      "day",
		ExpiresAt: time.Now(),
		Scraper:   scraper.DayScraper{},
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		pfix := string(m.Content[:1])
		if pfix != PREFIX {
			return
		}
		if strings.Contains(m.Content[1:], NSZ) {
			handleNsz(s, m, scrap)
		}

		if strings.Contains(m.Content, DAY) {
			log.Println(1)
			handleDay(s, m, scrap)
		}

	})

	if err := discord.Open(); err != nil {
		log.Fatal(err)
		return
	}
	discord.Identify.Intents = discordgo.IntentsAll | discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	log.Println("[🟢] Discord Bot is online")
}
