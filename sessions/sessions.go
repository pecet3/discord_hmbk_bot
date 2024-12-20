package sessions

import (
	"sync"
)

type Sessions struct {
	spamS   map[string]*spamSession
	sMu     sync.RWMutex
	RandomS map[string]*randomSession
	rMu     sync.RWMutex
}

func New() *Sessions {
	return &Sessions{
		spamS:   make(map[string]*spamSession),
		RandomS: make(map[string]*randomSession),
	}
}
