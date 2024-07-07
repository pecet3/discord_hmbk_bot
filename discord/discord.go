package discord

import (
	"log"
	"strings"
	"sync"
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

type BotUserSessions struct {
	Sessions map[string]*Session
	Mutex    sync.Mutex
}

type Session struct {
	UserId    string
	ExpiresAt time.Time
}

func NewSessions() *BotUserSessions {
	return &BotUserSessions{
		Sessions: make(map[string]*Session),
	}
}

func (bus *BotUserSessions) AddSession(userId string) {
	bus.Mutex.Lock()
	defer bus.Mutex.Unlock()
	bus.Sessions[userId] = &Session{
		UserId:    userId,
		ExpiresAt: time.Now().Add(10 * time.Second),
	}
}

func (bus *BotUserSessions) GetSession(id string) (*Session, bool) {
	bus.Mutex.Lock()
	defer bus.Mutex.Unlock()
	session, exists := bus.Sessions[id]
	return session, exists
}

func (bus *BotUserSessions) RemoveSession(id string) {
	bus.Mutex.Lock()
	defer bus.Mutex.Unlock()
	delete(bus.Sessions, id)
}

func Run(discord *discordgo.Session, ps *paint.PaintSessions) {
	defer discord.Close()
	scrap := scraper.New()
	sessions := NewSessions()

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
		us, exists := sessions.GetSession(m.Author.ID)
		if exists {
			if !us.ExpiresAt.Before(time.Now()) {
				log.Printf("<SPAM PROTECTION> [!] Blocked user: %s with ID: %s", m.Author.Username, m.Author.ID)
				return
			}
			log.Printf("<SPAM PROTECTION> New Session,  user: %s with ID: %s", m.Author.Username, m.Author.ID)

			sessions.RemoveSession(m.Author.ID)
			sessions.AddSession(m.Author.ID)
		} else {
			sessions.AddSession(m.Author.ID)
			log.Println("<SPAM PROTECTION> Added a session:", m.Author.ID)
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
