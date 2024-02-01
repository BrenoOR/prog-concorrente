package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn := connectUDPServer(8081)
	defer conn.Close()

	_, err := conn.Write(([]byte)("https://scrapeme.live/shop/"))
	if err != nil {
		log.Fatal(err)
	}

	res := make([]byte, 500*1000) // buffer size 500KB
	_, err = conn.Read(res)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(res))
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
