package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
	"web-scraper-nc/scrapernc"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	startTime := time.Now()

	quotes := []scrapernc.Quote{}
	authors := []scrapernc.Author{}
	pokemons := []scrapernc.Pokemon{}
	visitedPages := []scrapernc.URL{}

	scrapernc.Scraper("https://scrapeme.live/shop/", &quotes, &authors, &pokemons, &visitedPages)

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

	fmt.Println("Visited pages:")
	for i, page := range visitedPages {
		fmt.Println(i, ": ", page.Url)
	}

	file, err := os.Create("pokemonsnc.csv")
	checkError(err)
	fmt.Println("testando")

	writer := csv.NewWriter(file)
	writer.Write([]string{"Name", "Description", "Weight", "Dimensions", "Categories", "Price", "Quantity", "SKU", "Tags"})
	for _, p := range pokemons {
		writer.Write(p.CsvEntry())
	}

	writer.Flush()

	endTime := time.Now()

	executionTime := endTime.Sub(startTime).Seconds()

	fmt.Printf("\n\nTempo de execução: %fs\n", executionTime)
}
