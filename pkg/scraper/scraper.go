package scraper

import (
	"errors"
	"sync"
	"time"
)

type Entity struct {
	Title   string
	Content string
	Date    string
	Link    string
	Image   string
}

type Page struct {
	Name      string
	Entities  []Entity
	ExpiresAt time.Time
	Scraper   IScraper
}

type Scraper struct {
	PagesMap map[string]*Page
	mu       sync.RWMutex
}

func New() *Scraper {
	return &Scraper{
		PagesMap: make(map[string]*Page),
	}
}

func (s *Scraper) GetPage(name string) (*Page, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cn, exists := s.PagesMap[name]

	if !exists {
		return nil, errors.New("there is no such city")
	}
	return cn, nil
}

func (s *Scraper) SavePage(p *Page) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.PagesMap[p.Name] = p
}

type IScraper interface {
	GetEntities(p *Page) []Entity
}
