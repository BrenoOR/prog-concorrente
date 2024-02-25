package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"scrapeServer/commons"
	gorpcserver "scrapeServer/gorpcServer"
	tcpserver "scrapeServer/tcpServer"
	udpserver "scrapeServer/udpServer"
)

func printProgressBar(it int, total int) {
	percentage := float64(it) / float64(total) * 100
	filledLength := int(50 * it / total)
	end := ">"

	if it == total {
		end = "="
	}

	bar := strings.Repeat("=", filledLength) + end + strings.Repeat(" ", 50-filledLength)
	fmt.Printf("\r[%s] %.2f%%", bar, percentage)
	if it == total {
		fmt.Println()
	}
}

func main() {
	checkDirs()

	db := commons.DataBase{}
	db.Pages = make(map[string][]byte)
	prepareDataBase(&db)

	connTypes := make([]string, 0)
	connTypes = append(connTypes, "udp", "tcp", "rpc")

	args := os.Args[1:]
	if len(args) == 1 {
		switch args[0] {
		case "udp":
			udpserver.RunUDP(8081, &db)
		case "tcp":
			tcpserver.RunTCP(8082, &db)
		case "rpc":
			gorpcserver.RunGoRPC(8083, &db)
		case "help":
			fmt.Println("Connection types available are:", connTypes)
		default:
			log.Fatal("Connection type not defined.")
		}
	} else {
		log.Fatal("Provide exactly one connection type or 'help' for more info.")
	}

}

func checkDirs() {
	if _, err := os.Stat("pages/quotes.toscrape.com/"); os.IsNotExist(err) {
		os.MkdirAll("pages/quotes.toscrape.com/", 0755)
	}
	if _, err := os.Stat("pages/quotes.toscrape.com/author"); os.IsNotExist(err) {
		os.MkdirAll("pages/quotes.toscrape.com/author", 0755)
	}
	if _, err := os.Stat("pages/scrapeme.live/"); os.IsNotExist(err) {
		os.MkdirAll("pages/scrapeme.live/", 0755)
	}
	if _, err := os.Stat("pages/scrapeme.live/shop"); os.IsNotExist(err) {
		os.MkdirAll("pages/scrapeme.live/shop", 0755)
	}
	if _, err := os.Stat("pages/scrapeme.live/shop/page"); os.IsNotExist(err) {
		os.MkdirAll("pages/scrapeme.live/shop/page", 0755)
	}

	if _, err := os.Stat("pages/quotes.toscrape.com/index.html"); os.IsNotExist(err) {
		downloadDataBase()
	}
	if _, err := os.Stat("pages/scrapeme.live/shop/index.html"); os.IsNotExist(err) {
		downloadDataBase()
	}
}

func downloadDataBase() {
	file_content, err := os.ReadFile("pages/pagelist.txt")
	if err != nil {
		log.Fatal(err)
	}

	pages := strings.Split(string(file_content), "\n")
	wg := sync.WaitGroup{}

	it := 0
	mutexIT := sync.Mutex{}
	mutexMap := sync.Mutex{}

	for _, page := range pages {
		//fmt.Println(page)
		wg.Add(1)
		go downloadPage(page, &it, len(pages), &mutexIT, &mutexMap, &wg)
	}
	wg.Wait()
}

func prepareDataBase(db *commons.DataBase) {
	file_content, err := os.ReadFile("pages/pagelist.txt")
	if err != nil {
		log.Fatal(err)
	}

	pages := strings.Split(string(file_content), "\n")
	wg := sync.WaitGroup{}

	it := 0
	mutexIT := sync.Mutex{}
	mutexMap := sync.Mutex{}

	for _, page := range pages {
		//fmt.Println(page)
		wg.Add(1)
		go loadPage(db, page, &it, len(pages), &mutexIT, &mutexMap, &wg)
	}
	wg.Wait()

	if len(db.Pages) == 0 {
		log.Fatal("No pages loaded.")
	} else {
		fmt.Println("Database loaded with", len(db.Pages), "pages.")
	}
}

func downloadPage(page string, it *int, totalRange int, mutexIT *sync.Mutex, mutexMap *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		mutexIT.Lock()
		*it++
		mutexIT.Unlock()
		printProgressBar(*it, totalRange)
	}()

	res, err := http.Get(page)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal("Error: ", res.Status)
	}

	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	pageName := strings.Split(page, "://")[1]

	if pageName[len(pageName)-1] == '/' {
		pageName = pageName[:len(pageName)-1]
	}

	if len(strings.Split(pageName, "/")) == 1 {
		pageName += "/index"
	} else if strings.Split(pageName, "/")[len(strings.Split(pageName, "/"))-1] == "shop" {
		pageName += "/index"
	}

	err = os.WriteFile("pages/"+pageName+".html", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func loadPage(db *commons.DataBase, page string, it *int, totalRange int, mutexIT *sync.Mutex, mutexMap *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		mutexIT.Lock()
		*it++
		mutexIT.Unlock()
		printProgressBar(*it, totalRange)
	}()

	pageName := strings.Split(page, "://")[1]

	if pageName[len(pageName)-1] == '/' {
		pageName = pageName[:len(pageName)-1]
	}

	if len(strings.Split(pageName, "/")) == 1 {
		pageName += "/index"
	} else if strings.Split(pageName, "/")[len(strings.Split(pageName, "/"))-1] == "shop" {
		pageName += "/index"
	}

	content, err := os.ReadFile("pages/" + pageName + ".html")
	if err != nil {
		log.Fatal(err)
	}

	mutexMap.Lock()
	db.Pages[page] = content
	mutexMap.Unlock()
}
