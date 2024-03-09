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

	for {
		for k, v := range database.Pages {
			publishPageToRabbitMQ(port, k, v)
		}
	}
}

func publishPageToRabbitMQ(port int, url string, page_content []byte) {
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
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare an exchange", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = rabbitMQChannel.PublishWithContext(ctx,
		"pages", // exchange
		url,     // routing key (url)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        page_content,
		})
	if err != nil {
		log.Fatal("Failed to publish a message", err)
	}

	time.Sleep(5 * time.Second)
}
