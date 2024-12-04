package rabbitmq

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler func(delivery amqp091.Delivery)

// ConsumeMessages listens to a queue and processes messages using the provided handler.
// It ensures the consumer keeps running and automatically reconnects if the connection or channel fails.
func ConsumeMessages(queue string, handler MessageHandler, rabbitMQUrl string) {
	for {
		connection, err := amqp091.Dial(rabbitMQUrl)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v, retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer connection.Close()

		channel, err := connection.Channel()
		if err != nil {
			log.Printf("Failed to create channel: %v, retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer channel.Close()

		msgs, err := channel.Consume(queue, "", false, false, false, false, nil)
		if err != nil {
			log.Printf("Failed to consume messages from queue: %v, retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Listening for messages on queue: %s", queue)

		for d := range msgs {
			handler(d)

			// Acknowledge the message after successful processing
			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message: %v", err)
			}
		}

		// If the loop exits, reconnect
		log.Printf("Consumer disconnected. Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
