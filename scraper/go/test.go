package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("User-Agent", "Mozilla/6.0")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(r.StatusCode)
		fmt.Println(r.Request.Headers)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
		fmt.Println(r.Request.Headers)
	})

	c.Visit("https://www.albumoftheyear.org/")
}