package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

type DataBase struct {
	pages map[string][]byte
}

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
	db := DataBase{}
	db.pages = make(map[string][]byte)
	prepareDataBase(&db)

	wg := sync.WaitGroup{}
	defer wg.Wait()

	wg.Add(1)
	go runUDP(8081, &db)
}

func prepareDataBase(db *DataBase) {

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
}

func loadPage(db *DataBase, page string, it *int, totalRange int, mutexIT *sync.Mutex, mutexMap *sync.Mutex, wg *sync.WaitGroup) {
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

	mutexMap.Lock()
	db.pages[page] = content
	mutexMap.Unlock()

	//fmt.Println("Page", page, "loaded with", len(db.pages[page]), "bytes.")
}

func runUDP(port int, db *DataBase) {
	if len(db.pages) == 0 {
		log.Fatal("No pages loaded.")
	} else {
		fmt.Println("Database loaded with", len(db.pages), "pages.")
	}

	udpServer, err := net.ListenPacket("udp", fmt.Sprint(":", port))
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()
	defer log.Println("UDP Server address:", udpServer.LocalAddr(), "closed.")
	log.Println("UDP Server address:", udpServer.LocalAddr(), "open.")

	for {
		buf := make([]byte, 1024)
		_, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println("Received request from", addr.String(), "for page:", string(buf))

		go getPageUDP(db, udpServer, addr, string(buf))
	}
}

func getPageUDP(db *DataBase, udpServer net.PacketConn, addr net.Addr, page string) {
	page_content, onMap := db.pages[page]
	if !onMap {
		fmt.Println("Page", page, "not found.")
		keys := make([]string, 0, len(db.pages))
		for k := range db.pages {
			//fmt.Println("Key:", k, "Page:", page)
			if strings.Contains(page, k) {
				fmt.Println("Opsie. Page", page, "found.")
				page_content = db.pages[k]
				fmt.Println("Sending page:", page, "to", addr.String())
				fmt.Println("Page size:", len(page_content), "bytes")
				udpServer.WriteTo(page_content, addr)
				return
			}
			keys = append(keys, k)
		}
		udpServer.WriteTo([]byte(fmt.Sprint("Page not found. Try for:", keys[rand.Intn(len(keys))])), addr)
		//udpServer.WriteTo([]byte(fmt.Sprint("Page not found. Try for: ", keys[0])), addr)
	} else {
		fmt.Println("Sending page:", page, "to", addr.String())
		fmt.Println("Page size:", len(db.pages[page]), "bytes")
		udpServer.WriteTo(page_content, addr)
	}
}
