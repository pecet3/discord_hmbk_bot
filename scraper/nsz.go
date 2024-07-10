package scraper

import (
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

type NszScraper struct {
}

func (ns NszScraper) GetEntities(cn *Page) []Entity {

	c := colly.NewCollector()

	baseUrl := "https://tygodnikszczytno.pl"
	route := baseUrl + "/news"

	EntityChan := make(chan Entity, 10)
	linkChan := make(chan string, 10)
	go func() {
		defer close(EntityChan)
		defer close(linkChan)

		c.OnHTML("div.news", func(e *colly.HTMLElement) {
			var a Entity
			a.Title = e.ChildText("a")
			a.Content = e.ChildText("p.text")
			a.Date = e.ChildText("p.date")
			a.Image = e.ChildAttr("img", "src")
			findingStr := "CZYTAJ"
			if strings.Contains(a.Title, findingStr) {
				index := strings.Index(a.Title, findingStr)
				a.Title = a.Title[:index]
			}
			EntityChan <- a

		})

		c.OnHTML("a.readmore", func(e *colly.HTMLElement) {
			path := e.Attr("href")
			link := baseUrl + path
			linkChan <- link
		})

		c.Visit(route)
	}()
	c.OnRequest(func(r *colly.Request) {
		log.Println("News fetching from: ", r.URL)
	})
	var Entities []Entity
	var links []string

	for {
		select {
		case Entity, ok := <-EntityChan:
			if ok {
				Entities = append(Entities, Entity)
			} else {
				EntityChan = nil
			}
		case link, ok := <-linkChan:
			if ok {
				links = append(links, link)
			} else {
				linkChan = nil
			}

		}

		if EntityChan == nil && linkChan == nil {
			break
		}
	}
	cn.Entities = nil
	for i, art := range Entities {
		if i < len(links) {
			art.Link = links[i]
		}

		cn.Entities = append(cn.Entities, art)
	}
	return cn.Entities
}
