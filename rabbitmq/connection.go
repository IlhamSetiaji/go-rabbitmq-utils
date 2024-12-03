package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

var connection *amqp091.Connection
var channel *amqp091.Channel

// InitializeConnection establishes a RabbitMQ connection and channel.
func InitializeConnection(rabbitMQURL string) error {
	var err error

	// Establish a connection to RabbitMQ
	connection, err = amqp091.Dial(rabbitMQURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return err
	}

	// Open a channel on the connection
	channel, err = connection.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return err
	}

	log.Println("RabbitMQ connection established successfully.")
	return nil
}

// GetChannel returns the shared RabbitMQ channel for publishing/consuming.
func GetChannel() *amqp091.Channel {
	return channel
}

// CloseConnection cleans up the connection and channel.
func CloseConnection() {
	if channel != nil {
		err := channel.Close()
		if err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		} else {
			log.Println("RabbitMQ channel closed successfully.")
		}
	}

	if connection != nil {
		err := connection.Close()
		if err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		} else {
			log.Println("RabbitMQ connection closed successfully.")
		}
	}
}
