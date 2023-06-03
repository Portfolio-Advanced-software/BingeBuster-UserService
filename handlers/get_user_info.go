package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/globals"
	messaging "github.com/Portfolio-Advanced-software/BingeBuster-UserService/messaging"
	userpb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) RetrieveUserData(ctx context.Context, req *userpb.RetrieveUserDataReq) (*userpb.RetrieveUserDataRes, error) {
	// Create a message to request user data from other services
	message := map[string]interface{}{
		"user_id": req.GetUserId(),
		"action":  "retrieveUserData",
	}

	// Connect to RabbitMQ
	conn, err := messaging.ConnectToRabbitMQ(globals.RabbitMQUrl)
	if err != nil {
		log.Fatalf("Can't connect to RabbitMQ: %s", err)
	}

	// Produce a message to each relevant service's queue
	messaging.ProduceMessage(conn, message, "service1_queue")
	messaging.ProduceMessage(conn, message, "service2_queue")
	// Add more service queues as needed

	// Wait for responses from the other services
	responseCh := make(chan *userpb.OtherServiceUserData)
	errorCh := make(chan error)
	go consumeUserDataMessages(conn, responseCh, errorCh)

	// Process the received responses or errors
	var userData []*userpb.OtherServiceUserData
	for i := 0; i < numExpectedServices; i++ {
		select {
		case res := <-responseCh:
			userData = append(userData, res)
		case err := <-errorCh:
			log.Printf("Error while retrieving user data: %v", err)
		}
	}

	// Close the channels
	close(responseCh)
	close(errorCh)

	// If no user data is received from any service, return an error
	if len(userData) == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("No user data found for user ID %s", req.GetUserId()))
	}

	// Create the response
	response := &userpb.RetrieveUserDataRes{
		UserData: userData,
	}

	return response, nil
}

func consumeUserDataMessages(conn *messaging.RabbitMQConnection, responseCh chan<- *userpb.OtherServiceUserData, errorCh chan<- error) {
	// Create a channel for consuming messages
	msgs, err := messaging.ConsumeMessages(conn, "userData_queue")
	if err != nil {
		errorCh <- fmt.Errorf("Failed to consume messages: %v", err)
		return
	}

	// Process each received message
	for msg := range msgs {
		// Decode the message body
		var message map[string]interface{}
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			errorCh <- fmt.Errorf("Failed to decode message: %v", err)
			continue
		}

		// Extract the user data from the message
		// Adjust the data structure based on your specific requirements
		userData := &userpb.OtherServiceUserData{
			// Populate with data from the message
			// Example: Name: message["name"].(string),
			//          Age:  message["age"].(int32),
		}

		// Send the user data to the response channel
		responseCh <- userData

		// Acknowledge the message
		msg.Ack(false)
	}
}
