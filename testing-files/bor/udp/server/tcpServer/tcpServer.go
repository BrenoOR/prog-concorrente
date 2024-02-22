package tcpserver

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"scrapeServer/commons"
	"strings"
	"time"
)

func RunTCP(port int, db *commons.DataBase) {
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

func getPageTCP(conn net.Conn, db *commons.DataBase) {
	defer conn.Close()

	buf := make([]byte, 50*1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Println(conn.LocalAddr(), err)
		return
	}

	page := bytes.Trim(buf, "\x00")
	//fmt.Println("[", time.Now().Format(time.RFC822), "] Getting page:", string(page), "to", conn.RemoteAddr())
	keys := make([]string, 0, len(db.Pages))
	for k := range db.Pages {
		//fmt.Println("Key:", k, "Page:", page)
		if strings.Contains(string(page), k) && len(k) == len(string(page)) {
			page_content := db.Pages[k]
			//fmt.Println("[", time.Now().Format(time.RFC822), "] Sending page:", k, "to", conn.RemoteAddr())
			conn.Write(page_content)
			return
		}
		keys = append(keys, k)
	}
	//fmt.Println("[", time.Now().Format(time.RFC822), "] Page:", string(page), "not found.")
	conn.Write([]byte(fmt.Sprint("Page not found. Try for:", keys[rand.Intn(len(keys))])))
}
