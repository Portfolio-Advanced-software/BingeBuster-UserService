package messaging

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ProduceMessage(conn *amqp.Connection, payload interface{}, queueName string) {
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		queueName, // queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	FailOnError(err, "Failed to declare queue")

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	FailOnError(err, "Failed to marshal JSON payload")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Publish the message to the queue
	err = ch.PublishWithContext(ctx,
		"",     // exchange name
		q.Name, // routing key
		false,  // mandatory flag
		false,  // immediate flag
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonPayload,
		},
	)
	FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", jsonPayload)
}
