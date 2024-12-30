package discord

import (
	"log"
	"strings"
	"time"

	"github.com/pecet3/discord_hmbk_bot/pkg/paint"
	"github.com/pecet3/discord_hmbk_bot/pkg/scraper"
	"github.com/pecet3/discord_hmbk_bot/pkg/sessions"

	"github.com/bwmarrin/discordgo"
)

const (
	PREFIX      = "!"
	NSZ         = "nsz"
	DAY         = "dzien"
	IQ          = "iq"
	PAINT       = "paint"
	HUJ         = "huj"
	FONTANNA_ID = "408025348199022593"
)

func Run(discord *discordgo.Session, ps *paint.PaintSessions) {
	defer discord.Close()
	scrap := scraper.New()
	sessions := sessions.New()

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

	go loopRandom(sessions, discord)

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID || m.Author.Bot {
			return
		}

		if len(m.Content) <= 0 {
			return
		}

		_, rSexists := sessions.GetRandomSession(m.Author.ID)
		if !rSexists {
			log.Println("Added to random session, user: ", m.Author)
			sessions.AddRandomSession(m.Author)
		}

		// praiseCh <- m.Author

		pfix := string(m.Content[:1])
		if pfix != PREFIX {
			return
		}
		// spam protection
		us, sSexists := sessions.GetSpamSession(m.Author.ID)
		if sSexists {
			if !us.ExpiresAt.Before(time.Now()) {
				log.Printf("<SPAM PROTECTION> [!] Blocked user: %s with ID: %s", m.Author.Username, m.Author.ID)
				ch := m.ChannelID
				err := s.ChannelMessagesBulkDelete(ch, []string{m.Message.ID})
				if err != nil {
					log.Println(err, "HANDLER ERROR")
				}
				return
			}
			sessions.RemoveSpamSession(m.Author.ID)
			sessions.AddSpamSession(m.Author.ID)

		} else {
			sessions.AddSpamSession(m.Author.ID)
			log.Println("<SPAM PROTECTION> Added a session:", m.Author.ID)
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

		if strings.Contains(m.Content, HUJ) {
			logActivity(m, HUJ)
			handleHuj(s, m)
		}

		if strings.Contains(m.Content, PAINT) {
			logActivity(m, PAINT)
			handlePaint(s, m, ps)
			ch := m.ChannelID
			err := s.ChannelMessagesBulkDelete(ch, []string{m.Message.ID})
			if err != nil {
				log.Println(err, "HANDLER ERROR")
			}
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
