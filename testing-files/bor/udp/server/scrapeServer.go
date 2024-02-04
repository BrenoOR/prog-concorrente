package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
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

	wg.Add(2)
	go runUDP(8081, &db)
	go runTCP(8082, &db)
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

	if len(db.pages) == 0 {
		log.Fatal("No pages loaded.")
	} else {
		fmt.Println("Database loaded with", len(db.pages), "pages.")
	}
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
	database := *db
	udpServer, err := net.ListenPacket("udp", fmt.Sprint(":", port))
	if err != nil {
		log.Fatal(udpServer.LocalAddr(), err)
	}
	defer udpServer.Close()
	defer fmt.Println("[", time.Now().Format(time.RFC822), "] UDP Server address:", udpServer.LocalAddr(), "closed.")
	fmt.Println("[", time.Now().Format(time.RFC822), "] UDP Server address:", udpServer.LocalAddr(), "open.")

	for {
		buf := make([]byte, 50*1024)
		_, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			log.Println(udpServer.LocalAddr(), err)
			continue
		}

		page := bytes.Trim(buf, "\x00")

		go getPageUDP(&database, udpServer, addr, string(page))
	}
}

func getPageUDP(db *DataBase, udpServer net.PacketConn, addr net.Addr, page string) {
	fmt.Println("[", time.Now().Format(time.RFC822), "] Getting page:", page, "to", addr)

	keys := make([]string, 0, len(db.pages))
	for k := range db.pages {
		//fmt.Println("Key size:", len(k), "Page size:", len(page))
		if strings.Contains(page, k) && len(k) == len(page) {
			fmt.Println("[", time.Now().Format(time.RFC822), "] Sending page:", k, "to", addr)
			udpServer.WriteTo(db.pages[k], addr)
			return
		}
		keys = append(keys, k)
	}
	fmt.Println("[", time.Now().Format(time.RFC822), "] Page:", page, "not found.")
	udpServer.WriteTo([]byte(fmt.Sprint("Page not found. Try for:", keys[rand.Intn(len(keys))])), addr)
}

func runTCP(port int, db *DataBase) {
	database := *db
	tcpServer, err := net.Listen("tcp", fmt.Sprint(":", port))
	if err != nil {
		log.Fatal(tcpServer.Addr(), err)
	}
	defer tcpServer.Close()
	defer fmt.Println("[", time.Now().Format(time.RFC822), "] TCP Server address:", tcpServer.Addr(), "closed.")
	fmt.Println("[", time.Now().Format(time.RFC822), "] TCP Server address:", tcpServer.Addr(), "open.")

	for {
		conn, err := tcpServer.Accept()
		if err != nil {
			log.Println(tcpServer.Addr(), err)
			continue
		}

		go getPageTCP(conn, &database)
	}
}

func getPageTCP(conn net.Conn, db *DataBase) {
	defer conn.Close()

	buf := make([]byte, 50*1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Println(conn.LocalAddr(), err)
		return
	}

	page := bytes.Trim(buf, "\x00")
	fmt.Println("[", time.Now().Format(time.RFC822), "] Getting page:", string(page), "to", conn.RemoteAddr())
	keys := make([]string, 0, len(db.pages))
	for k := range db.pages {
		//fmt.Println("Key:", k, "Page:", page)
		if strings.Contains(string(page), k) && len(k) == len(string(page)) {
			page_content := db.pages[k]
			fmt.Println("[", time.Now().Format(time.RFC822), "] Sending page:", k, "to", conn.RemoteAddr())
			conn.Write(page_content)
			return
		}
		keys = append(keys, k)
	}
	fmt.Println("[", time.Now().Format(time.RFC822), "] Page:", string(page), "not found.")
	conn.Write([]byte(fmt.Sprint("Page not found. Try for:", keys[rand.Intn(len(keys))])))
}
