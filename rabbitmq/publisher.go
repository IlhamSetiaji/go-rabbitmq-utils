package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func PublishMessage(exchange, routingKey string, body []byte) error {
	ch := GetChannel()
	defer ch.Close()

	err := ch.Publish(exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}
	log.Printf("Message published to routing key: %s", routingKey)
	return nil
}
