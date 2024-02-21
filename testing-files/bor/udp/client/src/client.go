package src

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type Args struct {
	Url string
}

func getPageUDP(page string, res *[]byte, rttMutex *sync.Mutex, rttMean *int64) {
	//fmt.Println("(UDP) Getting page: ", page, " (client)")
	connUDP := connectUDPServer(8081)
	defer connUDP.Close()

	start := time.Now()

	_, err := connUDP.Write(([]byte)(page))
	if err != nil {
		log.Fatal(err)
	}

	_, err = connUDP.Read(*res)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()
	rttMutex.Lock()
	*rttMean += end.Sub(start).Microseconds()
	rttMutex.Unlock()
}

func getPageTCP(page string, res *[]byte, rttMutex *sync.Mutex, rttMean *int64) {
	//fmt.Println("(TCP) Getting page: ", page, " (client)")
	connTCP := connectTCPServer(8082)
	defer connTCP.Close()

	start := time.Now()

	_, err := connTCP.Write(([]byte)(page))
	if err != nil {
		log.Fatal(err)
	}

	_, err = connTCP.Read(*res)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()
	rttMutex.Lock()
	*rttMean += end.Sub(start).Microseconds()
	rttMutex.Unlock()
}

func getPageGoRPC(page string, res *[]byte, rttMutex *sync.Mutex, rttMean *int64) {
	args := Args{}
	args.Url = page

	client := connectGoRPCServer(8083)

	start := time.Now()

	err := client.Call("PageServer.GetPage", args, res)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()
	rttMutex.Lock()
	*rttMean += end.Sub(start).Microseconds()
	rttMutex.Unlock()
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

func connectGoRPCServer(port int) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", fmt.Sprint(":", port))
	if err != nil {
		log.Fatal(err)
	}

	return client
}
