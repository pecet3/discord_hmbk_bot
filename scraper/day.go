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

type ReceiverEnt struct {
	Titles   []string
	Contents []string
}

func (ns DayScraper) GetEntities(cn *Page) []Entity {

	c := colly.NewCollector()

	baseUrl := "https://www.kalbi.pl/kalendarz-swiat-nietypowych"
	route := baseUrl

	receiverCh := make(chan ReceiverEnt)
	go func() {
		defer close(receiverCh)
		c.OnHTML("div.descritoptions-of-holiday", func(e *colly.HTMLElement) {
			titlesStr := e.ChildText("a")
			contentStr := e.ChildText("p")
			titles := splitByUpperCase(titlesStr)
			contents := strings.Split(contentStr, "\n")

			for len(titles) > 0 {
				en := ReceiverEnt{
					Titles:   titles,
					Contents: contents,
				}
				receiverCh <- en
			}
		})
		c.OnRequest(func(r *colly.Request) {
			log.Println("Days fetching from: ", r.URL)
		})

		c.Visit(route)
	}()

	recv := <-receiverCh

	var entities []Entity
	if len(recv.Contents) != len(recv.Titles) {
		return []Entity{}
	}
	for i, t := range recv.Titles {
		ent := Entity{
			Title:   t,
			Content: recv.Contents[i],
		}
		entities = append(entities, ent)
	}
	if len(entities) <= 0 {
		return []Entity{}
	}
	cn.Entities = entities
	return cn.Entities
}
