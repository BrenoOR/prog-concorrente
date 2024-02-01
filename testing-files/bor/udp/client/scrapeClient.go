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
	conn := connectUDPServer(8081)
	defer conn.Close()

	_, err := conn.Write(([]byte)("http://quotes.toscrape.com"))
	if err != nil {
		log.Fatal(err)
	}

	res := make([]byte, 500*1000) // buffer size 500KB
	_, err = conn.Read(res)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".quote").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Find(".text").Text())
		fmt.Println(s.Find(".author").Text())
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
