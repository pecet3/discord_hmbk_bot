package discord

import (
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pecet3/discord_hmbk_bot/pkg/sessions"
)

const intervalHour = 6

func getRandomText() string {
	texts := [5]string{
		" to totalny kozak",
		" ma największego penisa na serwerze",
		" to największy ruchacz matek graczy",
		" podkrada palonko",
		" pożycza palonko i nie oddaje",
	}

	index := rand.Intn(5)
	return texts[index]
}

func loopRandom(sessions *sessions.Sessions, discord *discordgo.Session) {
	for {
		time.Sleep(time.Hour * intervalHour)
		i := 0
		lenRandomS := len(sessions.RandomS)
		if lenRandomS == 0 {
			continue
		}
		winnerIndex := rand.Intn(lenRandomS)
		log.Println("Winner index is: ", winnerIndex)
		for uuid, u := range sessions.RandomS {
			if i == winnerIndex {
				discord.ChannelMessageSend(FONTANNA_ID, u.User.Username+getRandomText())
			}
			sessions.RemoveRandomSession(uuid)
			i++
		}
	}
}
