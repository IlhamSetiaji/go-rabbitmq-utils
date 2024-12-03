package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler func(delivery amqp091.Delivery)

func ConsumeMessages(queue string, handler MessageHandler) error {
	ch := GetChannel()
	defer ch.Close()

	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to consume messages from queue: %v", err)
		return err
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
	}()
	log.Printf("Listening for messages on queue: %s", queue)
	return nil
}
