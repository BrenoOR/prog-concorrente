package main

import (
	"fmt"
	"os"
	"sync"
	"time"
	"web-scraper/scraper"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func scrapTrial_v2(trial int) int64 {
	//wg := sync.WaitGroup{}
	ch := make(chan int, 10)
	URLsToVisit_v2 := []string{}
	quotes := scraper.QuoteSlc{}
	authors := scraper.AuthorSlc{}
	pokemons := scraper.PokemonSlc{}
	times := scraper.RequestTimeSlc{}

	//wg.Add(1)
	//URLsToVisit_v2 = append(URLsToVisit_v2, "https://scrapeme.live/shop/")
	URLsToVisit_v2 = append(URLsToVisit_v2, "http://quotes.toscrape.com")

	//for i := 2; i <= 48; i++ {
	//	//wg.Add(1)
	//	URLsToVisit_v2 = append(URLsToVisit_v2, "https://scrapeme.live/shop/page/"+fmt.Sprint(i)+"/")
	//}

	start := time.Now()
	for _, page := range URLsToVisit_v2 {
		ch <- 1
		go scraper.Scrape_v2(page, &quotes, &authors, &pokemons, &times, &ch)
		//time.Sleep(3 * time.Second)
		//fmt.Println(len(ch))
	}
	//wg.Wait()
	for len(ch) > 0 {
		time.Sleep(1 * time.Second)
	}

	elapsed := time.Since(start)

	//fmt.Println("Pokemons: ", len(pokemons.Get()))
	fmt.Println("Quotes: ", len(quotes.Get()))
	fmt.Println("Authors: ", len(authors.Get()))

	total := 0
	sum := int64(0)
	for _, t := range times.Get() {
		sum += t
		total++
	}

	fmt.Println("Average Request Time: ", sum/int64(total), "ms")

	return elapsed.Milliseconds()
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
	f, err := os.Create("../../udp/server/pages/pagelist.txt")
	checkError(err)
	defer f.Close()

	_, err = f.WriteString("https://scrapeme.live/shop/\n")
	checkError(err)
	f.Sync()

	_, err = f.WriteString("https://scrapeme.live/shop/page/1/\n")
	checkError(err)
	f.Sync()

	_, err = f.WriteString("http://quotes.toscrape.com\n")
	checkError(err)
	f.Sync()

	start := time.Now()
	wg.Add(2)
	go scraper.Scrape("https://scrapeme.live/shop/", &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg)
	go scraper.Scrape("http://quotes.toscrape.com", &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg)

	wg.Wait()
	for !finished {
		for len(URLsToVisit.Slc) > 0 {
			nextURL, toVisit := URLsToVisit.Pop()
			if toVisit {
				fmt.Println("Visiting: ", nextURL)
				_, err = f.WriteString(nextURL + "\n")
				checkError(err)

				wg.Add(1)
				go scraper.Scrape(nextURL, &URLsToVisit, &URLsVisited, &quotes, &authors, &pokemons, &wg)
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

	//fmt.Println("Pokemons: ", len(pokemons.Get()))
	fmt.Println("Quotes: ", len(quotes.Get()))
	fmt.Println("Authors: ", len(authors.Get()))

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
	//wc_v2 := scrapTrial_v2(1)
	//nc := scrapTrial_NC(1)

	fmt.Println("Concurrency: Scraped in ", wc, "ms")
	//fmt.Println("No Concurrency: Scraped in ", wc_v2, "ms")
}
