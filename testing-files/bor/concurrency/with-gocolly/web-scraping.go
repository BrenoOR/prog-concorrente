package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"
	"web-scraper-colly/scraper"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func scrapTrial(trial int) int64 {
	wg := sync.WaitGroup{}

	quotes := []scraper.Quote{}
	authors := []scraper.Author{}
	pokemons := []scraper.Pokemon{}
	pagesToVisit := scraper.SlcURL{MutexURL: sync.Mutex{}, URLs: []string{}}
	visitedPages := scraper.SlcURL{MutexURL: sync.Mutex{}, URLs: []string{}}

	visitedPages.URLs = append(visitedPages.URLs, "https://scrapeme.live/shop/page/1/")

	start := time.Now()

	wg.Add(1)
	go scraper.Scraper("https://scrapeme.live/shop/", &quotes, &authors, &pokemons, &pagesToVisit, &visitedPages, &wg)

	wg.Wait()

	for len(pagesToVisit.URLs) > 0 {
		for i, page := range pagesToVisit.URLs {
			if i < 50 {
				//fmt.Println(i, page)
				wg.Add(1)
				go scraper.Scraper(page, &quotes, &authors, &pokemons, &pagesToVisit, &visitedPages, &wg)
			} else {
				break
			}
		}
		wg.Wait()
	}

	elapsed := time.Since(start)

	//strBuilder := strings.Builder{}
	//strBuilder.WriteString("pokemons-")
	//strBuilder.WriteString(fmt.Sprint(trial))
	//strBuilder.WriteString(".csv")
	//
	//file, err := os.Create(strBuilder.String())
	//checkError(err)
	//
	//writer := csv.NewWriter(file)
	//writer.Write([]string{"Name", "Description", "Weight", "Dimensions", "Categories", "Price", "Quantity", "SKU", "Tags"})
	//for _, p := range pokemons {
	//	writer.Write(p.CsvEntry())
	//}
	//
	//writer.Flush()

	return elapsed.Milliseconds()
}

func scrapTrialNoConcurrence(trial int) int64 {
	quotes := []scraper.Quote{}
	authors := []scraper.Author{}
	pokemons := []scraper.Pokemon{}
	visitedPages := []scraper.URL{}

	start := time.Now()

	scraper.ScraperNoConcurrence("https://scrapeme.live/shop/", &quotes, &authors, &pokemons, &visitedPages)

	elapsed := time.Since(start)

	//strBuilder := strings.Builder{}
	//strBuilder.WriteString("pokemons-nc-")
	//strBuilder.WriteString(fmt.Sprint(trial))
	//strBuilder.WriteString(".csv")
	//
	//file, err := os.Create(strBuilder.String())
	//checkError(err)
	//
	//writer := csv.NewWriter(file)
	//writer.Write([]string{"Name", "Description", "Weight", "Dimensions", "Categories", "Price", "Quantity", "SKU", "Tags"})
	//for _, p := range pokemons {
	//	writer.Write(p.CsvEntry())
	//}
	//
	//writer.Flush()

	return elapsed.Milliseconds()
}

func main() {
	trials := []int64{}
	for i := 0; i < 1; i++ {
		fmt.Println("Trial: ", i)
		trials = append(trials, scrapTrial(i))
		fmt.Println("Trial ", i, " finished in: ", trials[i], "ms")
	}
	var sum int64 = 0
	for _, trialET := range trials {
		sum += trialET
	}

	file, err := os.Create("trial-timings.csv")
	checkError(err)

	writer := csv.NewWriter(file)
	writer.Write([]string{"Trial", "ET"})
	for i, trial := range trials {
		writer.Write([]string{fmt.Sprint(i), fmt.Sprint(trial)})
	}

	writer.Flush()

	fmt.Println("Average time: ", sum/int64(len(trials)), "ms")

	trialsNC := []int64{}
	for i := 0; i < 1; i++ {
		fmt.Println("Trial: ", i)
		trialsNC = append(trialsNC, scrapTrialNoConcurrence(i))
		fmt.Println("Trial ", i, " finished in: ", trialsNC[i], "ms")
	}
	var sumNC int64 = 0
	for _, trialET := range trialsNC {
		sumNC += trialET
	}

	fileNC, errNC := os.Create("trial-timings-nc.csv")
	checkError(errNC)

	writerNC := csv.NewWriter(fileNC)
	writerNC.Write([]string{"Trial", "ET"})
	for i, trial := range trialsNC {
		writerNC.Write([]string{fmt.Sprint(i), fmt.Sprint(trial)})
	}

	writerNC.Flush()

	fmt.Println("Average time: ", sumNC/int64(len(trialsNC)), "ms")
}
