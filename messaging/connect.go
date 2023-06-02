package messaging

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ(rabbitmqUrl string) (*amqp.Connection, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitmqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	return conn, nil
}
