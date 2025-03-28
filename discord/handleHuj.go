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
		display = m.Author.Mention() + " ma " + resultStr + " cm huja.\n" + summary
	} else {
		display = m.Mentions[0].Mention() + " ma " + resultStr + " cm huja.\n" + summary
	}
	s.ChannelMessageSend(m.ChannelID, display)
}
