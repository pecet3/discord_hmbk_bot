package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()

	c.OnHTML("div.news", func(e *colly.HTMLElement) {
		fmt.Println()
		fmt.Println(e.ChildText("a"))
		fmt.Println(e.ChildText("p.text"))
		fmt.Println()
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://tygodnikszczytno.pl/news/")

}
