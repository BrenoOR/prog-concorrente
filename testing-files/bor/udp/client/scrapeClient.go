package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"scrape-client/src"
	"sync"
	"time"
)

func main() {
	totalTrials := 1000
	info := [][]int64{}

	for i := 0; i < totalTrials; i++ {
		trialID := i + 1
		udpTrial := scrapeTrial(trialID, "udp")
		fmt.Println("UDP Trial ", trialID, ": ", udpTrial, "ms")
		udpTrialNC := scrapeTrialNC(trialID, "udp")
		fmt.Println("UDP Trial ", trialID, " (No Concurrency): ", udpTrialNC, "ms")
		tcpTrial := scrapeTrial(trialID, "tcp")
		fmt.Println("TCP Trial ", trialID, ": ", tcpTrial, "ms")
		tcpTrialNC := scrapeTrialNC(trialID, "tcp")
		fmt.Println("TCP Trial ", trialID, " (No Concurrency): ", tcpTrialNC, "ms")

		row := []int64{int64(trialID), udpTrial, udpTrialNC, tcpTrial, tcpTrialNC}
		info = append(info, row)
	}

	data, err := os.Create("data.csv")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer data.Close()

	writer := csv.NewWriter(data)
	defer writer.Flush()
	headers := []string{"Trial", "UDP", "UDP (No Concurrency)", "TCP", "TCP (No Concurrency)"}
	writer.Write(headers)
	for _, row := range info {
		strRow := []string{}
		for _, cell := range row {
			strRow = append(strRow, fmt.Sprint(cell))
		}
		writer.Write(strRow)
	}
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
