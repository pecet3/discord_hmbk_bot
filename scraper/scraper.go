package scraper

import (
	"errors"
	"sync"
	"time"
)

type Article struct {
	Title   string
	Content string
	Date    string
	Link    string
	Image   string
}

type CityNews struct {
	Name      string
	Articles  []Article
	ExpiresAt time.Time
	Scraper   IScraper
}

type Scraper struct {
	CitiesMap map[string]*CityNews
	mu        sync.RWMutex
}

func New() *Scraper {
	return &Scraper{
		CitiesMap: make(map[string]*CityNews),
	}
}

func (s *Scraper) GetCity(cName string) (*CityNews, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cn, exists := s.CitiesMap[cName]

	if !exists {
		return nil, errors.New("there is no such city")
	}
	return cn, nil
}

func (s *Scraper) SaveCity(cn *CityNews) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.CitiesMap[cn.Name] = cn
}

type IScraper interface {
	GetArticles(cn *CityNews) []Article
}
