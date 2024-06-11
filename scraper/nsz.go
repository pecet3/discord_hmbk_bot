package scraper

import (
	"log"
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

	baseUrl := "https://tygodnikszczytno.pl"
	route := baseUrl + "/news"

	articleChan := make(chan Article, 10) // Buffered channel for articles
	linkChan := make(chan string, 10)     // Buffered channel for links
	imageChan := make(chan string, 10)
	go func() {
		defer close(articleChan)
		defer close(linkChan)

		c.OnHTML("div.news", func(e *colly.HTMLElement) {
			var a Article
			src := e.Attr("src")
			log.Println(src)
			a.Title = e.ChildText("a")
			a.Content = e.ChildText("p.text")
			a.Date = e.ChildText("p.date")
			a.Image = e.ChildAttr("img", "src")
			log.Println(a.Image)
			findingStr := "CZYTAJ"
			if strings.Contains(a.Title, findingStr) {
				index := strings.Index(a.Title, findingStr)
				a.Title = a.Title[:index]
			}
			articleChan <- a
			imageChan <- a.Image
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
	var articles []Article
	var links []string
	var images []string
	for {
		select {
		case article, ok := <-articleChan:
			if ok {
				articles = append(articles, article)
			} else {
				articleChan = nil
			}
		case link, ok := <-linkChan:
			if ok {
				links = append(links, link)
			} else {
				linkChan = nil
			}
		case image, ok := <-imageChan:
			if ok {
				images = append(images, image)
			} else {
				imageChan = nil
			}
		}

		if articleChan == nil && linkChan == nil {
			break
		}
	}

	for i, art := range articles {
		if i < len(links) {
			art.Link = links[i]
		}

		art.Image = images[i]

		cn.Articles = append(cn.Articles, art)
	}
	log.Println(cn.Articles)

	cn.ExpiresAt = time.Now().Add(6 * time.Hour)
	return cn.Articles
}
