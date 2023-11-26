package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func getIndividual(c *colly.Collector) {
	c.OnHTML(".text", func(h *colly.HTMLElement) {
		fmt.Println("Quote: ", h.Text)
	})

	c.OnHTML(".author", func(h *colly.HTMLElement) {
		fmt.Println("Author: ", h.Text)
	})
}

type Quote struct {
	Quote string
	Author string
}

func main() {
	quotes := []Quote{}

	c := colly.NewCollector(
		colly.AllowedDomains("quotes.toscrape.com"),
	)

	// c.Visit("https://en.wikipedia.org/wiki/Main_Page")

	c.OnRequest(func(req *colly.Request) {
		req.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
		fmt.Println("Visiting: ", req.URL)
	})

	c.OnError(func(_ *colly.Response, err error){
		fmt.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response code: ", r.StatusCode)
	})

	getIndividual(c)

	c.OnHTML(".quote", func(h *colly.HTMLElement) {
		div := h.DOM
		quote := div.Find(".text").Text()
		author := div.Find(".author").Text()
		
		q := Quote{
			Quote: quote,
			Author: author,
		}

		quotes = append(quotes, q)

		//fmt.Printf("Quote: %s\nBy: \\%s\\\n\n", quote, author)
	})

	c.Visit("http://quotes.toscrape.com/")

	fmt.Println(quotes)
}