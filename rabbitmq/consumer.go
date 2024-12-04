package rabbitmq

import (
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler func(delivery amqp091.Delivery)

func ConsumeMessages(queue string, handler MessageHandler) {
	for {
		ch := GetChannel()
		if ch == nil {
			log.Printf("Failed to consume messages: channel is not open, retrying in 5 seconds...")
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
		if err != nil {
			log.Printf("Failed to consume messages from queue: %v, retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		log.Printf("Listening for messages on queue: %s", queue)

		// Process messages in a separate goroutine
		go func() {
			for d := range msgs {
				handler(d)

				// Acknowledge the message after successful processing
				if err := d.Ack(false); err != nil {
					log.Printf("Error acknowledging message: %v", err)
				}
			}
		}()

		// Block to keep the consumer active until an error occurs
		select {
		case <-ch.NotifyClose(make(chan *amqp091.Error)): // Channel is closed
			log.Printf("Channel closed, reconnecting...")
			time.Sleep(5 * time.Second) // Wait before reconnecting
		}
	}
}
