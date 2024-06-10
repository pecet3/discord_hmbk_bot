package news

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Article struct {
	Title   string
	Content string
}

type CityNews struct {
	Name      string
	Articles  []Article
	ExpiresAt time.Time
}

func (n *CityNews) GetArticles() []Article {
	if !n.ExpiresAt.Before(time.Now()) {
		return n.Articles
	}
	c := colly.NewCollector()

	var a Article

	c.OnHTML("div.news", func(e *colly.HTMLElement) {
		a.Title = e.ChildText("a")
		a.Content = e.ChildText("p.text")

		findingStr := "CZYTAJ"

		if strings.Contains(a.Title, findingStr) {
			index := strings.Index(a.Title, findingStr)
			a.Title = a.Title[:index]
		}
		n.Articles = append(n.Articles, a)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://tygodnikszczytno.pl/news/")

	n.ExpiresAt = time.Now().Add(6 * time.Hour)
	return n.Articles
}
