package discord

import (
	"bytes"
	"log"
	"os"
	"time"
	"webscraping/paint"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func handlePaint(s *discordgo.Session, m *discordgo.MessageCreate, ps *paint.PaintSessions) {
	sessionId := uuid.NewString()

	url := os.Getenv("BASE_URL")
	session := paint.PaintSession{
		Id:         sessionId,
		ImgBytesCh: make(chan []byte),
		ExpiresAt:  time.Now().Add(4 * time.Hour),
	}

	ps.AddSession(sessionId, session)
	s.ChannelMessageSend(m.ChannelID, url+"/?session_id="+sessionId)

	imgBytes := <-session.ImgBytesCh
	imgName := m.Author.Username + "_" + sessionId + ".png"
	file := &discordgo.File{
		Name:        imgName,
		ContentType: "image/png",
		Reader:      bytes.NewReader(imgBytes),
	}

	msg := discordgo.MessageSend{
		File:    file,
		Content: "Oto piękny obraz " + m.Author.Mention(),
	}
	log.Printf("> Painting %s has been sent", imgName)
	s.ChannelMessageSendComplex(m.ChannelID, &msg)

	ps.RemoveSession(sessionId)

}
