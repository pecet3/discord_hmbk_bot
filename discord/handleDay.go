package discord

import (
	"fmt"
	"log"
	"time"
	"webscraping/scraper"

	"github.com/bwmarrin/discordgo"
)

var polishMonths = []string{
	"stycznia", "lutego", "marca", "kwietnia", "maja", "czerwca",
	"lipca", "sierpnia", "września", "października", "listopada", "grudnia",
}
var polishDaysOfWeek = []string{
	"niedziela", "poniedziałek", "wtorek", "środa", "czwartek", "piątek", "sobota",
}

func formatDatePolish(t time.Time) string {
	day := t.Day()
	month := polishMonths[t.Month()-1]
	year := t.Year()
	dayOfWeek := polishDaysOfWeek[t.Weekday()]

	return fmt.Sprintf("%s, %d %s %d", dayOfWeek, day, month, year)
}

func handleDay(s *discordgo.Session, m *discordgo.MessageCreate, scrap *scraper.Scraper) {
	day, _ := scrap.GetPage("day")

	entities := day.Entities
	if day.ExpiresAt.Before(time.Now()) {
		entities = day.Scraper.GetEntities(day)
		day.ExpiresAt = time.Now().Add(6 * time.Hour)
	}
	log.Println("[👍] Day cache hit")
	scrap.SavePage(day)

	display := ""
	for _, a := range entities {
		display = display + " _**" + a.Title + "**_  " + "\n" + a.Content + "\n" + "\n"
	}
	emb := &discordgo.MessageEmbed{
		Title:       formatDatePolish(time.Now()),
		Description: display,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, emb)

}
