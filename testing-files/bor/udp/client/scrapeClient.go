package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	getPageUDP("http://quotes.toscrape.com")
	getPageTCP("http://quotes.toscrape.com/author/Albert-Einstein")
}

func getPageUDP(page string) {
	connUDP := connectUDPServer(8081)
	defer connUDP.Close()

	_, err := connUDP.Write(([]byte)("http://quotes.toscrape.com"))
	if err != nil {
		log.Fatal(err)
	}

	res1 := make([]byte, 500*1000) // buffer size 500KB
	_, err = connUDP.Read(res1)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res1))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".quote").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Find(".text").Text())
		fmt.Println(s.Find(".author").Text())
	})
}

func getPageTCP(page string) {
	connTCP := connectTCPServer(8082)
	defer connTCP.Close()

	_, err := connTCP.Write(([]byte)("http://quotes.toscrape.com/author/Albert-Einstein"))
	if err != nil {
		log.Fatal(err)
	}

	res := make([]byte, 500*1000) // buffer size 500KB
	_, err = connTCP.Read(res)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".product").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Find(".product_title").Text())
		fmt.Println(s.Find(".summary").Text())
	})

	doc.Find(".author-details").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Find(".author-title").Text())
		fmt.Println(s.Find(".author-born-date").Text())
		fmt.Println(s.Find(".author-born-location").Text())
	})
}

func connectUDPServer(port int) *net.UDPConn {
	scrapeServer, err := net.ResolveUDPAddr("udp", fmt.Sprint(":", port))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, scrapeServer)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return conn
}

func connectTCPServer(port int) *net.TCPConn {
	scrapeServer, err := net.ResolveTCPAddr("tcp", fmt.Sprint(":", port))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, scrapeServer)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return conn
}
