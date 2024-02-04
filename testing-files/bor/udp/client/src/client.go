package src

import (
	"fmt"
	"log"
	"net"
	"os"
)

func getPageUDP(page string, res *[]byte) {
	//fmt.Println("(UDP) Getting page: ", page, " (client)")
	connUDP := connectUDPServer(8081)
	defer connUDP.Close()

	_, err := connUDP.Write(([]byte)(page))
	if err != nil {
		log.Fatal(err)
	}

	_, err = connUDP.Read(*res)
	if err != nil {
		log.Fatal(err)
	}
}

func getPageTCP(page string, res *[]byte) {
	//fmt.Println("(TCP) Getting page: ", page, " (client)")
	connTCP := connectTCPServer(8082)
	defer connTCP.Close()

	_, err := connTCP.Write(([]byte)(page))
	if err != nil {
		log.Fatal(err)
	}

	_, err = connTCP.Read(*res)
	if err != nil {
		log.Fatal(err)
	}
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
