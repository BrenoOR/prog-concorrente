package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"encoding/csv"

	"github.com/PuerkitoBio/goquery"
)


func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func writeFile(data, filename string) {
	file, err := os.Create(filename)
	defer file.Close()
	checkError(err)

	file.WriteString(data)
}

func main() {
	url := "https://techcrunch.com/"

	response, err := http.Get(url)
	defer response.Body.Close()
	checkError(err)
	// Se o código de status for acima de 400, erro no cliente
	if response.StatusCode > 400 {
		fmt.Println("Status code: ", response.StatusCode)
	}
	
	// Converte a página html para structs
	document, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)
	
	file, err := os.Create("posts.csv")
	checkError(err)
	writer := csv.NewWriter(file)
	// Se o método Find encontrar um conjunto de elementos, o .Each() itera sobre os objetos selecionados e executa um função para cada elemento
	document.Find("div.river").Find("div.post-block").Each(func(index int, item *goquery.Selection){
		h2 := item.Find("h2")
		title := strings.TrimSpace(h2.Text()) // Eliminante espaços em branco
		url, _ := h2.Find("a").Attr("href")

		excerpt := strings.TrimSpace(item.Find("div.post-block__content").Text())
		
		//fmt.Printf("Title: %s\nURL: %s\nExcerpt: %s\n\n", title, url, excerpt)

		posts := []string{title, url, excerpt}

		// Escrevendo os elementos em um csv
		writer.Write(posts)
	})

	writer.Flush()

	//fmt.Println(river)
	//writeFile(river, "index.html")
}