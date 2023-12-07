package main

import (
	"fmt"
	"sync"
	"time"
	"web-scraper/scraper"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func scrapTrial(trial int) int64 {
	wg := sync.WaitGroup{}
	URLsVisited := scraper.Slice_CS{}
	URLsToVisit := scraper.Slice_CS{}
	quotes := scraper.QuoteSlc{}
	authors := scraper.AuthorSlc{}
	pokemons := scraper.PokemonSlc{}
	finished := false

	URLsVisited.Append("https://scrapeme.live/shop/page/1/")

	start := time.Now()
	wg.Add(1)
	go scraper.Scrape("https://scrapeme.live/shop/", &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg)
	//go scraper.Scrape("http://quotes.toscrape.com", &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg)

	wg.Wait()
	for !finished {
		if len(URLsToVisit.Slc) < 50 {
			for len(URLsToVisit.Slc) > 0 {
				nextURL, toVisit := URLsToVisit.Pop()
				if toVisit {
					wg.Add(1)
					go scraper.Scrape(nextURL, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg)
				} else {
					finished = true
				}
			}
			wg.Wait()
		}
	}
	wg.Wait()

	elapsed := time.Since(start)

	fmt.Println("Pokemons: ", len(pokemons.Get()))

	return elapsed.Milliseconds()
}

func scrapTrial_NC(trial int) int64 {
	URLsVisited := scraper.Slice_CS{}
	URLsToVisit := scraper.Slice_CS{}
	quotes := scraper.QuoteSlc{}
	authors := scraper.AuthorSlc{}
	pokemons := scraper.PokemonSlc{}

	URLsVisited.Append("https://scrapeme.live/shop/page/1/")

	start := time.Now()
	scraper.Scrape_NC("https://scrapeme.live/shop/", &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons)
	//scraper.Scrape_NC("http://quotes.toscrape.com", &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons)

	for len(URLsToVisit.Slc) > 0 {
		nextURL, toVisit := URLsToVisit.Pop()
		if toVisit {
			fmt.Println("URL: ", nextURL)
			scraper.Scrape_NC(nextURL, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons)
		}
		fmt.Println("URLsToVisit: ", len(URLsToVisit.Slc))
	}

	elapsed := time.Since(start)

	fmt.Println("Pokemons: ", len(pokemons.Get()))

	return elapsed.Milliseconds()
}

func main() {

	wc := scrapTrial(1)
	nc := scrapTrial_NC(1)

	fmt.Println("Concurrency: Scraped in ", wc, "ms")
	fmt.Println("No Concurrency: Scraped in ", nc, "ms")
}
