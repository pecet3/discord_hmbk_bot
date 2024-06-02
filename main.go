package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("div.news", func(e *colly.HTMLElement) {
		// fmt.Println(e.ChildText("div"))
		fmt.Println(e.ChildText("a"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://tygodnikszczytno.pl/news/")
}
