package discord

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/pecet3/discord_hmbk_bot/paint"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func handlePaint(s *discordgo.Session, m *discordgo.MessageCreate, ps *paint.PaintSessions) {
	sessionId := uuid.NewString()

	url := os.Getenv("BASE_URL")
	session := paint.PaintSession{
		Id:         sessionId,
		ImgBytesCh: make(chan []byte),
		ExpiresAt:  time.Now().Add(2 * time.Hour),
	}

	ps.AddSession(sessionId, session)

	prvChannel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return
	}
	prvMsg := "Obyś nie skończył jak słynny akwarelista ( ͡° ͜ʖ ͡°)\n" +
		url + "paint/?session_id=" + sessionId
	s.ChannelMessageSend(prvChannel.ID, prvMsg)

	expiryCh := make(chan bool)
	go func(expiryCh chan bool) {
		for {
			s, isExists := ps.GetSession(sessionId)
			if !isExists {
				break
			}
			if s.ExpiresAt.Before(time.Now()) {
				expiryCh <- true
				break
			}
			time.Sleep(5 * time.Second)
		}
	}(expiryCh)

	for {
		select {
		case imgBytes := <-session.ImgBytesCh:
			imgName := m.Author.Username + "_" + sessionId + ".png"
			file := &discordgo.File{
				Name:        imgName,
				ContentType: "image/png",
				Reader:      bytes.NewReader(imgBytes),
			}

			msg := discordgo.MessageSend{
				File:    file,
				Content: "Oto piękny obraz autorstwa" + m.Author.Mention(),
			}
			log.Printf("> Painting %s has been sent", imgName)
			s.ChannelMessageSendComplex(m.ChannelID, &msg)
			ps.RemoveSession(sessionId)
			return
		case isExpired := <-expiryCh:
			if isExpired {
				ps.RemoveSession(sessionId)
				log.Printf("Paint session with id: %s, created by: %s has been removed (expired)", sessionId, m.Author.Username)
				return
			}
		}
	}
}
