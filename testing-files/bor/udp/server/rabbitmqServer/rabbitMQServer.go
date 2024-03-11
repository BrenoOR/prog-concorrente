package rabbitmqserver

import (
	"context"
	"fmt"
	"log"
	"scrapeServer/commons"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RunRabbitMQ(port int, db *commons.DataBase) {

	database := *db

	rabbitMQServer, err := amqp.Dial(fmt.Sprint("amqp://guest:guest@localhost:", port, "/"))
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ at ", fmt.Sprint("amqp://guest:guest@localhost:", port, "/"), " => ", err)
	}
	defer rabbitMQServer.Close()

	rabbitMQChannel, err := rabbitMQServer.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel", err)
	}
	defer rabbitMQChannel.Close()

	err = rabbitMQChannel.ExchangeDeclare(
		"pages", // name
		"topic", // type
		false,   // durable
		true,    // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare an exchange => ", err)
	}

	queue, err := rabbitMQChannel.QueueDeclare(
		"requests", // name
		false,      // durable
		false,      // delete when unused
		true,       // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	err = rabbitMQChannel.QueueBind(
		queue.Name, // queue name
		"requests", // routing key
		"pages",    // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msg, err := rabbitMQChannel.Consume(
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

	holdConn := make(chan int)

	go func() {
		for d := range msg {
			go publishPageToRabbitMQ(rabbitMQChannel, string(d.Body), database.Pages[string(d.Body)])
		}
	}()

	<-holdConn
}

func publishPageToRabbitMQ(rabbitMQChannel *amqp.Channel, url string, page_content []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := rabbitMQChannel.PublishWithContext(ctx,
		"pages", // exchange
		url,     // routing key (url)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        page_content,
		})
	if err != nil {
		log.Fatal("Failed to publish a message => ", err)
	}

	fmt.Println("Published page", url, "to RabbitMQ.")
}
