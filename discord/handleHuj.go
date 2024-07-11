package discord

import (
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func hujRandom() int {

	numbers := make([]int, 0)

	for i := 1; i <= 30; i++ {
		if i >= 10 && i <= 18 {
			for j := 0; j < 3; j++ {
				numbers = append(numbers, i)
			}
		} else {
			numbers = append(numbers, i)
		}
	}

	return numbers[rand.Intn(len(numbers))]
}

func handleHuj(s *discordgo.Session, m *discordgo.MessageCreate) {
	result := hujRandom()

	resultStr := strconv.Itoa(result)
	summary := "mikrus"

	if result > 25 {
		summary = "Potężna pała"
	} else if result > 16 {
		summary = "Kolega zadowolony"
	} else if result > 10 {
		summary = "Może być..."
	}

	var display string
	if len(m.Mentions) <= 0 {
		display = m.Author.Mention() + " ma " + resultStr + " cm huj.\n" + summary

		if m.Author.ID == pecetId || m.Author.ID == kszaqId {
			summary = "To jebany geniusz który osiągnął stan nirvany i oświecenia"
			resultStr = "999"
		}

	} else {
		userId := m.Mentions[0].Mention()[2 : len(m.Mentions[0].Mention())-1]
		if userId == pecetId || userId == kszaqId {
			summary = "To jebany geniusz który osiągnął stan nirvany i oświecenia"
			resultStr = "999"
		}
		display = m.Mentions[0].Mention() + " ma " + resultStr + " IQ.\n" + summary
	}
	s.ChannelMessageSend(m.ChannelID, display)
}
