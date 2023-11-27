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

func main() {
	wg := sync.WaitGroup{}

	quotes := []scraper.Quote{}
	authors := []scraper.Author{}
	pokemons := []scraper.Pokemon{}
	visitedPages := []scraper.URL{}

	start := time.Now()

	wg.Add(1)
	go scraper.Scraper("https://scrapeme.live/shop/", &quotes, &authors, &pokemons, &visitedPages, &wg)

	wg.Wait()
	//for _, q := range quotes {
	//	q.Print()
	//}
	//
	//fmt.Println()
	//
	//for _, a := range authors {
	//	a.Print()
	//}
	//
	//for _, p := range pokemons {
	//	p.Print()
	//}
	elapsed := time.Since(start)

	fmt.Println("Visited pages:")
	for i, page := range visitedPages {
		fmt.Println(i, ": ", page.Url)
	}

	file, err := os.Create("pokemons.csv")
	checkError(err)

	writer := csv.NewWriter(file)
	writer.Write([]string{"Name", "Description", "Weight", "Dimensions", "Categories", "Price", "Quantity", "SKU", "Tags"})
	for _, p := range pokemons {
		writer.Write(p.CsvEntry())
	}

	writer.Flush()

	fmt.Println(elapsed.Milliseconds())
}
