package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"scrape-client/src"
	"strings"
	"sync"
	"time"
)

func printProgressBar(it int, total int, trialMean int64, rttMean int64, trialNCMean int64, rttNCMean int64) {
	percentage := float64(it) / float64(total) * 100
	filledLength := int(50 * it / total)
	end := ">"

	if it == total {
		end = "="
	}

	bar := strings.Repeat("=", filledLength) + end + strings.Repeat(" ", 50-filledLength)
	fmt.Printf("\r[%s] %.2f%% [Means (ms): %d | %d | %d | %d]", bar, percentage, trialMean, rttMean, trialNCMean, rttNCMean)
	if it == total {
		fmt.Println()
	}
}

func main() {
	totalTrials := 1000
	trialTotal := int64(0)
	rttTotal := int64(0)
	trialNCTotal := int64(0)
	rttNCTotal := int64(0)

	connType := ""
	connTypes := make([]string, 0)
	connTypes = append(connTypes, "udp", "tcp", "rpc")

	args := os.Args[1:]
	if len(args) == 1 {
		switch args[0] {
		case "udp":
			connType = "udp"
		case "tcp":
			connType = "tcp"
		case "rpc":
			connType = "rpc"
		case "help":
			fmt.Println("Connection types available are:", connTypes)
		default:
			log.Fatal("Connection type not defined.")
		}
	} else {
		log.Fatal("Provide exactly one connection type or 'help' for more info.")
	}

	if connType == "" {
		log.Fatal("Connection type not provided.")
	}

	data, err := os.Create(connType + ".csv")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer data.Close()

	writer := csv.NewWriter(data)
	defer writer.Flush()

	headers := []string{"Trial", strings.ToUpper(connType), "Mean RTT (" + strings.ToUpper(connType) + ")", strings.ToUpper(connType) + " (No Concurrency)", "Mean RTT (" + strings.ToUpper(connType) + "NC)"}
	writer.Write(headers)

	for i := 0; i < totalTrials; i++ {
		trialID := i + 1
		trialNC, rttNC := scrapeTrialNC(trialID, connType)
		trial, rtt := scrapeTrial(trialID, connType)

		row := []int64{int64(trialID), trial, rtt, trialNC, rttNC}

		strRow := []string{}
		for _, cell := range row {
			strRow = append(strRow, fmt.Sprint(cell))
		}
		writer.Write(strRow)

		trialTotal += trial
		rttTotal += rtt
		trialNCTotal += trialNC
		rttNCTotal += rttNC

		printProgressBar(i+1, totalTrials, trialTotal/int64(trialID), rttTotal/int64(trialID), trialNCTotal/int64(trialID), rttNCTotal/int64(trialID))
	}
}

func scrapeTrial(trial int, connType string) (int64, int64) {
	wg := sync.WaitGroup{}
	URLsVisited := src.Slice_CS{}
	URLsToVisit := src.Slice_CS{}
	quotes := src.QuoteSlc{}
	authors := src.AuthorSlc{}
	pokemons := src.PokemonSlc{}
	finished := false

	rttMean := int64(0)
	rttMutex := sync.Mutex{}

	client := src.ConnectGoRPCServer(8083)
	clientMutex := sync.Mutex{}
	defer client.Close()

	URLsVisited.Append("https://scrapeme.live/shop/page/1/")

	start := time.Now()

	wg.Add(1)
	go src.Scrape("https://scrapeme.live/shop/", connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg, &rttMutex, &rttMean, client, &clientMutex)
	//go src.Scrape("http://quotes.toscrape.com", connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg)

	wg.Wait()
	for !finished {
		for len(URLsToVisit.Slc) > 0 {
			nextURL, toVisit := URLsToVisit.Pop()
			if toVisit {
				//fmt.Println("Visiting: ", nextURL)

				wg.Add(1)
				go src.Scrape(nextURL, connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg, &rttMutex, &rttMean, client, &clientMutex)
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
	rttMean /= int64(len(URLsVisited.Slc) - 1)

	//fmt.Println("Pokemons: ", len(pokemons.Get()))
	//fmt.Println("Pages Visited: ", len(URLsVisited.Get()))
	//fmt.Println("Pages To Visit: ", len(URLsToVisit.Get()))
	//fmt.Println("Quotes: ", len(quotes.Get()))
	//fmt.Println("Authors: ", len(authors.Get()))

	return elapsed.Microseconds(), rttMean
}

func scrapeTrialNC(trial int, connType string) (int64, int64) {
	URLsVisited := src.Slice_CS{}
	URLsToVisit := src.Slice_CS{}
	quotes := src.QuoteSlc{}
	authors := src.AuthorSlc{}
	pokemons := src.PokemonSlc{}

	rttMean := int64(0)
	rttMutex := sync.Mutex{}

	client := src.ConnectGoRPCServer(8083)
	clientMutex := sync.Mutex{}
	defer client.Close()

	URLsVisited.Append("https://scrapeme.live/shop/page/1/")

	start := time.Now()

	src.ScrapeNC("https://scrapeme.live/shop/", connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &rttMutex, &rttMean, client, &clientMutex)
	//src.ScrapeNC("http://quotes.toscrape.com", connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons)

	//fmt.Println("Pages Visited: ", len(URLsVisited.Get()))
	//fmt.Println("Pages To Visit: ", len(URLsToVisit.Get()))

	//time.Sleep(2 * time.Second)

	for len(URLsToVisit.Slc) > 0 {
		nextURL, toVisit := URLsToVisit.Pop()
		if toVisit {
			//fmt.Println("Visiting: ", nextURL)

			src.ScrapeNC(nextURL, connType, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &rttMutex, &rttMean, client, &clientMutex)

			//fmt.Println("Pages Visited: ", len(URLsVisited.Get()))
			//fmt.Println("Pages To Visit: ", len(URLsToVisit.Get()))

			//time.Sleep(2 * time.Second)
		}
	}

	elapsed := time.Since(start)
	rttMean /= int64(len(URLsVisited.Slc) - 1)

	//fmt.Println("Pokemons: ", len(pokemons.Get()))
	//fmt.Println("Pages Visited: ", len(URLsVisited.Get()))
	//fmt.Println("Pages To Visit: ", len(URLsToVisit.Get()))
	//fmt.Println("Quotes: ", len(quotes.Get()))
	//fmt.Println("Authors: ", len(authors.Get()))

	return elapsed.Microseconds(), rttMean
}
