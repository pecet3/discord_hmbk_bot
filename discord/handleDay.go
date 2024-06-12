package discord

import (
	"log"
	"webscraping/scraper"

	"github.com/bwmarrin/discordgo"
)

func handleDay(s *discordgo.Session, m *discordgo.MessageCreate, scrap *scraper.Scraper) {
	day, _ := scrap.GetPage("day")
	log.Println(2)
	entities := day.Scraper.GetEntities(day)
	log.Println(entities)

	// display := ""
	// for _, a := range entities {
	// 	display = display + " _**" + a.Title + "**_  " + a.Content + "\n"
	// }
	// emb := &discordgo.MessageEmbed{
	// 	Title:       time.Now().Format("20060102150405"),
	// 	Description: display,
	// }
	// s.ChannelMessageSendEmbed(m.ChannelID, emb)
	// scrap.SavePage(day)

}
