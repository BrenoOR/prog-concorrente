package rabbitmqserver

import (
	"context"
	"fmt"
	"log"
	"scrapeServer/commons"
	"sync"
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

	wg := sync.WaitGroup{}
	for {
		for k, v := range database.Pages {
			if rabbitMQChannel.IsClosed() {
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
			}
			wg.Add(1)
			go publishPageToRabbitMQ(rabbitMQChannel, k, v, &wg)
			time.Sleep(time.Millisecond * 10)
		}
		wg.Wait()
		fmt.Println("[", time.Now().Format(time.RFC822), "] All pages published to RabbitMQ.")
		time.Sleep(time.Millisecond * 100)
	}
}

func publishPageToRabbitMQ(rabbitMQChannel *amqp.Channel, url string, page_content []byte, wg *sync.WaitGroup) {
	defer wg.Done()

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
}
