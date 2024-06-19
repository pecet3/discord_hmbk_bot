package scraper

import (
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type DayScraper struct {
}

func splitByUpperCase(s string) []string {
	re := regexp.MustCompile(`([a-ząćęłńóśźż])([A-ZĄĆĘŁŃÓŚŹŻ])`)
	withSpaces := re.ReplaceAllString(s, "$1,$2")
	parts := strings.Split(withSpaces, ",")
	return parts
}

func (ns DayScraper) GetEntities(cn *Page) []Entity {

	c := colly.NewCollector()

	baseUrl := "https://www.kalbi.pl/kalendarz-swiat-nietypowych"
	route := baseUrl

	var entities []Entity
	ch := make(chan int)
	go func() {
		c.OnHTML("div.descritoptions-of-holiday", func(e *colly.HTMLElement) {
			titlesStr := e.ChildText("a")
			contentStr := e.ChildText("p")
			titles := splitByUpperCase(titlesStr)
			contents := strings.Split(contentStr, "\n")
			for i, title := range titles {
				var ent Entity
				ent.Title = title
				ent.Content = contents[i]
				entities = append(entities, ent)
				if i == len(titles)-1 {
					ch <- len(titles)
				}
			}
		})
		c.OnRequest(func(r *colly.Request) {
			log.Println("Days fetching from: ", r.URL)
		})

		c.Visit(route)
	}()
	i := <-ch
	cn.Entities = entities[:i]
	return cn.Entities
}
