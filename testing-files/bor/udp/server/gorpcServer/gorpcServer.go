package gorpcserver

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"scrapeServer/commons"
	"strings"
	"sync"
	"time"
)

type PageServer int64

var database *commons.DataBase
var dbMutex *sync.Mutex

func RunGoRPC(port int, db *commons.DataBase) {
	dbMutex = &sync.Mutex{}

	dbMutex.Lock()
	database = db
	dbMutex.Unlock()

	pageServer := new(PageServer)
	rpc.Register(pageServer)

	rpc.HandleHTTP()

	gorpcServer, err := net.Listen("tcp", fmt.Sprint(":", port))
	if err != nil {
		log.Fatal(gorpcServer.Addr(), err)
	}
	defer gorpcServer.Close()
	defer fmt.Println("[", time.Now().Format(time.RFC822), "] TCP Server address:", gorpcServer.Addr(), "closed.")
	fmt.Println("[", time.Now().Format(time.RFC822), "] TCP Server address:", gorpcServer.Addr(), "open.")

	http.Serve(gorpcServer, nil)
}

func (s *PageServer) GetPage(args *commons.Args, reply *[]byte) error {
	if args.Url == "" {
		return errors.New("invalid page provided")
	}

	dbMutex.Lock()
	db := database
	dbMutex.Unlock()

	keys := make([]string, 0, len(db.Pages))
	for k := range db.Pages {
		//fmt.Println("Key:", k, "Page:", page)
		if strings.Contains(args.Url, k) && len(k) == len(args.Url) {
			page_content := db.Pages[k]
			//fmt.Println("[", time.Now().Format(time.RFC822), "] Sending page:", k)
			*reply = page_content
			return nil
		}
		keys = append(keys, k)
	}

	//fmt.Println("[", time.Now().Format(time.RFC822), "] Page:", args.Url, "not found.")
	*reply = []byte(fmt.Sprint("Page not found. Try for:", keys[rand.Intn(len(keys))]))

	return nil
}
