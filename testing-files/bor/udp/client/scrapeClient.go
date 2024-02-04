package main

import (
	"fmt"
	"scrape-client/src"
	"sync"
	"time"
)

func main() {
	udpTrial := scrapeTrial(1, "udp")
	fmt.Println("UDP Trial 1: ", udpTrial, "ms")
	udpTrialNC := scrapeTrialNC(1, "udp")
	fmt.Println("UDP Trial 1 (No Concurrency): ", udpTrialNC, "ms")
	tcpTrial := scrapeTrial(1, "tcp")
	fmt.Println("TCP Trial 1: ", tcpTrial, "ms")
	tcpTrialNC := scrapeTrialNC(1, "tcp")
	fmt.Println("TCP Trial 1 (No Concurrency): ", tcpTrialNC, "ms")
}

func scrapeTrial(trial int, connType string) int64 {
	wg := sync.WaitGroup{}
	URLsVisited := src.Slice_CS{}
	URLsToVisit := src.Slice_CS{}
	quotes := src.QuoteSlc{}
	authors := src.AuthorSlc{}
	pokemons := src.PokemonSlc{}
	finished := false

	URLsVisited.Append("https://scrapeme.live/shop/page/1/")

	start := time.Now()

	wg.Add(1)
	go src.Scrape("https://scrapeme.live/shop/", connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg)
	//go src.Scrape("http://quotes.toscrape.com", connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg)

	wg.Wait()
	for !finished {
		for len(URLsToVisit.Slc) > 0 {
			nextURL, toVisit := URLsToVisit.Pop()
			if toVisit {
				// fmt.Println("Visiting: ", nextURL)

				wg.Add(1)
				go src.Scrape(nextURL, connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg)
			}
		}
		wg.Wait()
		if len(URLsToVisit.Slc) > 0 {
			finished = false
		} else if len(URLsToVisit.Slc) == 0 {
			finished = true
		}
	}
	wg.Wait()

	elapsed := time.Since(start)

	fmt.Println("Pokemons: ", len(pokemons.Get()))
	fmt.Println("Pages Visited: ", len(URLsVisited.Get()))
	fmt.Println("Pages To Visit: ", len(URLsToVisit.Get()))
	//fmt.Println("Quotes: ", len(quotes.Get()))
	//fmt.Println("Authors: ", len(authors.Get()))

	return elapsed.Milliseconds()
}

func scrapeTrialNC(trial int, connType string) int64 {
	URLsVisited := src.Slice_CS{}
	URLsToVisit := src.Slice_CS{}
	quotes := src.QuoteSlc{}
	authors := src.AuthorSlc{}
	pokemons := src.PokemonSlc{}

	URLsVisited.Append("https://scrapeme.live/shop/page/1/")

	start := time.Now()

	src.ScrapeNC("https://scrapeme.live/shop/", connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons)
	//src.ScrapeNC("http://quotes.toscrape.com", connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons)

	//fmt.Println("Pages Visited: ", len(URLsVisited.Get()))
	//fmt.Println("Pages To Visit: ", len(URLsToVisit.Get()))

	//time.Sleep(2 * time.Second)

	for len(URLsToVisit.Slc) > 0 {
		nextURL, toVisit := URLsToVisit.Pop()
		if toVisit {
			//fmt.Println("Visiting: ", nextURL)

			src.ScrapeNC(nextURL, connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons)

			//fmt.Println("Pages Visited: ", len(URLsVisited.Get()))
			//fmt.Println("Pages To Visit: ", len(URLsToVisit.Get()))

			//time.Sleep(2 * time.Second)
		}
	}

	elapsed := time.Since(start)

	fmt.Println("Pokemons: ", len(pokemons.Get()))
	fmt.Println("Pages Visited: ", len(URLsVisited.Get()))
	fmt.Println("Pages To Visit: ", len(URLsToVisit.Get()))
	//fmt.Println("Quotes: ", len(quotes.Get()))
	//fmt.Println("Authors: ", len(authors.Get()))

	return elapsed.Milliseconds()
}
