package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/6.0"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
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
