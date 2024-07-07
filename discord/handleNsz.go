package discord

import (
	"strconv"
	"time"

	"github.com/pecet3/discord_hmbk_bot/scraper"

	"github.com/bwmarrin/discordgo"
)

func handleNsz(s *discordgo.Session, m *discordgo.MessageCreate, scrap *scraper.Scraper) {
	nsz, _ := scrap.GetPage("szczytno")
	entities := nsz.Entities
	if entities == nil {
		entities = nsz.Scraper.GetEntities(nsz)
		nsz.ExpiresAt = time.Now().Add(6 * time.Hour)
		scrap.SavePage(nsz)
	} else {
		if !nsz.ExpiresAt.Before(time.Now()) {
		} else {
			entities = nsz.Scraper.GetEntities(nsz)
			nsz.ExpiresAt = time.Now().Add(6 * time.Hour)
			scrap.SavePage(nsz)
		}
	}

	if len(NSZ)+2 > len(m.Content) {
		display := ""
		for i, a := range entities[:14] {
			display = display + " _**[" + strconv.Itoa(i+1) + "]**_  " + a.Title + "\n"
			i++
		}
		// scrap.SavePage(nsz)
		img := &discordgo.MessageEmbedImage{
			URL:    "https://i.ibb.co/kJdN894/nsz.jpg",
			Height: 100,
		}

		emb := &discordgo.MessageEmbed{
			Title:       "Gorące Newsy ze Szczytna \n !nsz [numer]",
			Description: display,
			Image:       img,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, emb)
		return
	}
	arg := m.Content[len(NSZ)+2:]
	page, _ := strconv.Atoi(arg)
	pageLeft := page - 1

	if page == 0 {
		pageLeft = 0
		page = 1
	}
	for _, a := range entities[pageLeft:page] {

		img := &discordgo.MessageEmbedImage{
			URL:    a.Image,
			Height: 100,
		}

		emb := &discordgo.MessageEmbed{
			URL:         a.Link,
			Title:       a.Title,
			Description: a.Date + "\n " + a.Content,
			Image:       img,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, emb)
	}

}
