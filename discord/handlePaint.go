package discord

import (
	"log"
	"webscraping/paint"

	"github.com/bwmarrin/discordgo"
)

const ADR = "http://localhost:8080/"

func handlePaint(s *discordgo.Session, m *discordgo.MessageCreate, ps *paint.PaintSessions) {
	sessionId := "test"

	session := paint.PaintSession{
		Id:       sessionId,
		FinishCh: make(chan bool),
	}

	ps.AddSession(sessionId, session)
	s.ChannelMessageSend(m.ChannelID, ADR+sessionId)

	isFinish := <-session.FinishCh

	if isFinish {
		log.Println("finish")
	}
	s.ChannelMessageSend(m.ChannelID, "finish")

}
