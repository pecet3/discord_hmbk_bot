package scraper

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type NszScraper struct {
}

func (ns NszScraper) GetArticles(cn *CityNews) []Article {

	if !cn.ExpiresAt.Before(time.Now()) {
		return cn.Articles
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
		cn.Articles = append(cn.Articles, a)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("News fetching from: ", r.URL)
	})

	c.Visit("https://tygodnikszczytno.pl/news/")

	cn.ExpiresAt = time.Now().Add(6 * time.Hour)
	return cn.Articles
}
