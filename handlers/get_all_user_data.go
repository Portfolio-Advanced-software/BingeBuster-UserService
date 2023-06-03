package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/globals"
	messaging "github.com/Portfolio-Advanced-software/BingeBuster-UserService/messaging"
	userpb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/proto"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *UserServiceServer) GetAllUserData(ctx context.Context, req *userpb.GetAllUserDataReq) (*userpb.GetAllUserDataRes, error) {
	// Connect to RabbitMQ
	conn, err := messaging.ConnectToRabbitMQ(globals.RabbitMQUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	var responseStrings []string

	// Define the callback function for consuming messages
	consumeCallback := func(body []byte) error {
		// Append the received string to the response strings slice
		responseStrings = append(responseStrings, string(body))

		println("the received message" + string(body))
		// Notify the wait group that a response has been received
		wg.Done()
		return nil
	}

	// Start consuming messages from the queue
	go messaging.ConsumeMessage(conn, "user_data", consumeCallback)

	// Prepare the message
	message := map[string]interface{}{
		"user_id": req.GetId(),
		"action":  "getAllRecords",
	}

	// Publish the message to the exchange
	messaging.ProduceMessage(conn, message, "auth_queue")
	messaging.ProduceMessage(conn, message, "authz_queue")
	messaging.ProduceMessage(conn, message, "watch_history_queue")

	// Wait for three responses to arrive
	wg.Add(3)
	wg.Wait()

	// Combine the response strings into a single string
	combinedString := strings.Join(responseStrings, "")

	id := req.GetId()

	var result bson.M
	if err := globals.UserDb.FindOne(ctx, bson.M{"userid": id}).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to find record: %v", err)
	}

	// Convert the result to a JSON string
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal record to JSON: %w", err)
	}

	// Convert the JSON byte array to a string
	dataString := string(jsonData)

	// Add something to the combinedString
	combinedString += dataString

	// Return response with success: true if no error is thrown (and thus document is removed)
	return &userpb.GetAllUserDataRes{
		Data: combinedString,
	}, nil
}
