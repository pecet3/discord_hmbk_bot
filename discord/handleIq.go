package discord

import (
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
	summary := ""

	if result > 120 {
		summary = "To bardzo mądry człowiek."
	} else if result > 100 {
		summary = "To człowiek niemądry niegłupi."
	} else if result < 80 {
		summary = "To debil."
	}
	var display string
	if len(m.Mentions) <= 0 {
		display = m.Author.Mention() + " ma " + resultStr + " IQ.\n" + summary

	} else {
		display = m.Mentions[0].Mention() + " ma " + resultStr + " IQ.\n" + summary
	}
	s.ChannelMessageSend(m.ChannelID, display)
}
