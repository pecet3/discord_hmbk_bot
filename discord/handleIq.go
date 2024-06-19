package discord

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func weightedRandom() int {

	numbers := make([]int, 0)

	for i := 1; i <= 200; i++ {
		if i >= 100 && i <= 140 {
			for j := 0; j < 5; j++ {
				numbers = append(numbers, i)
			}
		} else {
			numbers = append(numbers, i)
		}
	}

	return numbers[rand.Intn(len(numbers))]
}

func handleIq(s *discordgo.Session, m *discordgo.MessageCreate) {
	result := weightedRandom()
	resultStr := strconv.Itoa(result)
	log.Println(m.Content)
	summary := ""
	log.Println(m.Mentions[0].Mention())
	if result > 120 {
		summary = "To bardzo mądry człowiek."
	} else {
		summary = "Inteligencją nie grzeszy..."
	}

	display := m.Author.Mention() + " ma " + resultStr + " IQ.\n" + summary + m.Mentions[0].Mention()

	s.ChannelMessageSend(m.ChannelID, display)

}
