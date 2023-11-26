package main

import (
	"fmt"
	"net/http"
	"os"

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

	// Se o método Find encontrar um conjunto de elementos, o .Html() retorna o HTML do primeiro elemento desse conjunto
	river, err := document.Find("div.river").Html()
	checkError(err)

	//fmt.Println(river)
	writeFile(river, "index.html")
}