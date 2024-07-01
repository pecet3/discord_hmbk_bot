package paint

import (
	"errors"
	"sync"
)

type PaintSession struct {
	Id       string
	FinishCh chan bool
}

type PaintSessions struct {
	Map   map[string]PaintSession
	Mutex sync.Mutex
}

func NewPaintSessions() *PaintSessions {
	return &PaintSessions{
		Map: make(map[string]PaintSession),
	}
}

func (ps *PaintSessions) AddSession(id string, session PaintSession) {
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()
	ps.Map[id] = session
}

func (ps *PaintSessions) GetSession(id string) (PaintSession, bool) {
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()
	session, exists := ps.Map[id]
	return session, exists
}

func (ps *PaintSessions) RemoveSession(id string) {
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()
	delete(ps.Map, id)
}

func (ps *PaintSessions) UpdateSession(id string, session PaintSession) error {
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()
	if _, exists := ps.Map[id]; !exists {
		return errors.New("sesja nie istnieje")
	}
	ps.Map[id] = session
	return nil
}

func (ps *PaintSessions) ListSessions() []string {
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()
	ids := make([]string, 0, len(ps.Map))
	for id := range ps.Map {
		ids = append(ids, id)
	}
	return ids
}

func (ps *PaintSessions) SessionCount() int {
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()
	return len(ps.Map)
}
