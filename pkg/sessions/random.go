package sessions

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type randomSession struct {
	User      *discordgo.User
	ExpiresAt time.Time
}

func (bus *Sessions) AddRandomSession(du *discordgo.User) {
	bus.rMu.Lock()
	defer bus.rMu.Unlock()
	bus.RandomS[du.ID] = &randomSession{
		User:      du,
		ExpiresAt: time.Now().Add(10 * time.Second),
	}
}

func (bus *Sessions) GetRandomSession(id string) (*randomSession, bool) {
	bus.rMu.Lock()
	defer bus.rMu.Unlock()
	session, exists := bus.RandomS[id]
	return session, exists
}

func (bus *Sessions) RemoveRandomSession(id string) {
	bus.rMu.Lock()
	defer bus.rMu.Unlock()
	delete(bus.RandomS, id)
}
