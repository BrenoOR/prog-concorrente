package udpserver

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

func RunUDP(port int, db *commons.DataBase) {
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

func getPageUDP(db *commons.DataBase, udpServer net.PacketConn, addr net.Addr, page string) {
	//fmt.Println("[", time.Now().Format(time.RFC822), "] Getting page:", page, "to", addr)

	keys := make([]string, 0, len(db.Pages))
	for k := range db.Pages {
		//fmt.Println("Key size:", len(k), "Page size:", len(page))
		if strings.Contains(page, k) && len(k) == len(page) {
			//fmt.Println("[", time.Now().Format(time.RFC822), "] Sending page:", k, "to", addr)
			udpServer.WriteTo(db.Pages[k], addr)
			return
		}
		keys = append(keys, k)
	}
	//fmt.Println("[", time.Now().Format(time.RFC822), "] Page:", page, "not found.")
	udpServer.WriteTo([]byte(fmt.Sprint("Page not found. Try for:", keys[rand.Intn(len(keys))])), addr)
}
