package discord

import (
	"log"
	"strings"
	"time"

	"github.com/pecet3/discord_hmbk_bot/paint"
	"github.com/pecet3/discord_hmbk_bot/scraper"

	"github.com/bwmarrin/discordgo"
)

const (
	PREFIX = "!"
	NSZ    = "nsz"
	DAY    = "dzien"
	IQ     = "iq"
	PAINT  = "paint"
)

func Run(discord *discordgo.Session, ps *paint.PaintSessions) {
	defer discord.Close()
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
		if len(m.Content) <= 0 {
			return
		}
		pfix := string(m.Content[:1])
		if pfix != PREFIX {
			return
		}
		if strings.Contains(m.Content[1:], NSZ) {
			logActivity(m, NSZ)
			handleNsz(s, m, scrap)
		}

		if strings.Contains(m.Content, DAY) {
			logActivity(m, DAY)
			handleDay(s, m, scrap)
		}

		if strings.Contains(m.Content, IQ) {
			logActivity(m, IQ)
			handleIq(s, m)
		}

		if strings.Contains(m.Content, PAINT) {
			logActivity(m, PAINT)
			handlePaint(s, m, ps)
		}

	})

	if err := discord.Open(); err != nil {
		log.Fatal(err)
		return
	}
	discord.Identify.Intents = discordgo.IntentsAll | discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	log.Println("[🟢] Discord Bot is online")
}

func logActivity(m *discordgo.MessageCreate, command string) {
	log.Printf("> %s [%s] used command: !%s", m.Author.GlobalName, m.Author.ID, command)
}
