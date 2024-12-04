package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler func(delivery amqp091.Delivery)

func ConsumeMessages(queue string, handler MessageHandler) {
	ch := GetChannel()
	defer ch.Close()

	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to consume messages from queue: %v", err)
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			if err := d.Nack(false, true); err != nil {
				log.Printf("Error nacking message: %v", err)
				continue
			}
			handler(d)
			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message: %v", err)
			}
		}
	}()
	log.Printf("Listening for messages on queue: %s", queue)
	<-forever
}
