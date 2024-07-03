package discord

import (
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func weightedRandom() int {

	numbers := make([]int, 0)

	for i := 1; i <= 200; i++ {
		if i >= 80 && i <= 120 {
			for j := 0; j < 5; j++ {
				numbers = append(numbers, i)
			}
		} else {
			numbers = append(numbers, i)
		}
	}

	return numbers[rand.Intn(len(numbers))]
}

var pecetId = strconv.Itoa(282817551401091072)
var kszaqId = strconv.Itoa(377032854179282944)

func handleIq(s *discordgo.Session, m *discordgo.MessageCreate) {
	result := weightedRandom()

	resultStr := strconv.Itoa(result)
	summary := "To jebany debil."

	if result > 140 {
		summary = "To człowiek mądry, szlachetny i bardzo inteligentny."
	} else if result > 120 {
		summary = "To bardzo mądry człowiek."
	} else if result > 100 {
		summary = "To człowiek niemądry niegłupi."
	}

	var display string
	if len(m.Mentions) <= 0 {
		display = m.Author.Mention() + " ma " + resultStr + " IQ.\n" + summary

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
