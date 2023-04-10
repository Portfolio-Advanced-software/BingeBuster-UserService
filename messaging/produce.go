package messaging

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ProduceMessage(message string, queue string) {
	conn, err := amqp.Dial("amqps://tnhdeowx:tInXH7wKtKdyn-v97fZ_HGM5XmHsDTNl@rattlesnake.rmq.cloudamqp.com/tnhdeowx")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		queue, // queue name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := []byte(message)

	// Publish the message to the queue
	err = ch.PublishWithContext(ctx,
		"",     // exchange name
		q.Name, // routing key
		false,  // mandatory flag
		false,  // immediate flag
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
