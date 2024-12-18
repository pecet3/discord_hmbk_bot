package sessions

import (
	"sync"
	"time"
)

type Sessions struct {
	spamS map[string]*spamSession
	sMu   sync.Mutex
}

type spamSession struct {
	UserId    string
	ExpiresAt time.Time
}

func New() *Sessions {
	return &Sessions{
		spamS: make(map[string]*spamSession),
	}
}

func (bus *Sessions) AddSpamSession(userId string) {
	bus.sMu.Lock()
	defer bus.sMu.Unlock()
	bus.spamS[userId] = &spamSession{
		UserId:    userId,
		ExpiresAt: time.Now().Add(10 * time.Second),
	}
}

func (bus *Sessions) GetSpamSession(id string) (*spamSession, bool) {
	bus.sMu.Lock()
	defer bus.sMu.Unlock()
	session, exists := bus.spamS[id]
	return session, exists
}

func (bus *Sessions) RemoveSpamSession(id string) {
	bus.sMu.Lock()
	defer bus.sMu.Unlock()
	delete(bus.spamS, id)
}
