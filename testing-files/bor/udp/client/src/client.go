package src

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Args struct {
	Url string
}

func getPageUDP(page string, res *[]byte, rttMutex *sync.Mutex, rttMean *int64) {
	//fmt.Println("(UDP) Getting page: ", page, " (client)")
	args := Args{}
	args.Url = page

	req, err := json.Marshal(args)
	if err != nil {
		log.Fatal(err)
	}

	connUDP := ConnectUDPServer(8081)
	defer connUDP.Close()

	start := time.Now()

	_, err = connUDP.Write(req)
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
	args := Args{}
	args.Url = page

	req, err := json.Marshal(args)
	if err != nil {
		log.Fatal(err)
	}

	connTCP := ConnectTCPServer(8082)
	defer connTCP.Close()

	start := time.Now()

	_, err = connTCP.Write(bytes.Trim(req, "\x00"))
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

func getPageGoRPC(page string, res *[]byte, rttMutex *sync.Mutex, rttMean *int64, client *rpc.Client, clientMutex *sync.Mutex) {
	args := Args{}
	args.Url = page

	clientMutex.Lock()

	start := time.Now()

	err := client.Call("PageServer.GetPage", args, res)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()

	clientMutex.Unlock()

	rttMutex.Lock()
	*rttMean += end.Sub(start).Microseconds()
	rttMutex.Unlock()
}

func getPageRabbitMQ(page string, res *[]byte, rttMutex *sync.Mutex, rttMean *int64, connRabbitMQ *amqp.Connection) {
	chRabbitMQ, err := connRabbitMQ.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer chRabbitMQ.Close()

	err = chRabbitMQ.ExchangeDeclare(
		"pages", // name
		"topic", // type
		false,   // durable
		true,    // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	queue, err := chRabbitMQ.QueueDeclare(
		page,  // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	err = chRabbitMQ.QueueBind(
		queue.Name, // queue name
		page,       // routing key
		"pages",    // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	err = chRabbitMQ.PublishWithContext(ctx,
		"pages",    // exchange
		"requests", // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(page),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	cancel()

	msg, err := chRabbitMQ.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("Getting page:", page, "=> onqueue:", queue.Name)

	holdConn := make(chan int)

	start := time.Now()

	go func(st *time.Time, ch *chan int) {
		for d := range msg {
			*res = d.Body

			end := time.Now()

			fmt.Println("Received page:", page, "in", end.Sub(*st).Microseconds(), "microseconds.")

			rttMutex.Lock()
			*rttMean += end.Sub(*st).Microseconds()
			rttMutex.Unlock()
			*ch <- 1
			break
		}
	}(&start, &holdConn)

	<-holdConn
}

func ConnectUDPServer(port int) *net.UDPConn {
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

func ConnectTCPServer(port int) *net.TCPConn {
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

func ConnectGoRPCServer(port int) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", fmt.Sprint(":", port))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func ConnectRabbitMQServer(port int) *amqp.Connection {
	client, err := amqp.Dial(fmt.Sprint("amqp://guest:guest@localhost:", port, "/"))
	if err != nil {
		log.Fatal(err)
	}

	return client
}
