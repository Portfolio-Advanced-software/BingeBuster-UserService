package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/globals"
	messaging "github.com/Portfolio-Advanced-software/BingeBuster-UserService/messaging"
	userpb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserReq) (*userpb.DeleteUserRes, error) {
	// Create a filter for the userID field
	filter := bson.M{"userid": req.GetId()}

	// Delete the documents matching the filter
	result, err := globals.UserDb.DeleteMany(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete user(s) with id %s: %v", req.GetId(), err))
	}

	// Check if any documents were deleted
	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("No user(s)) with id %s found: %v", req.GetId(), err))
	}

	// Send a message to the watch history queue for deleting user's history
	message := map[string]interface{}{
		"user_id": req.GetId(),
		"action":  "deleteAllRecords",
	}

	conn, err := messaging.ConnectToRabbitMQ(globals.RabbitMQUrl)
	if err != nil {
		log.Fatalf("Can't connect to RabbitMQ: %s", err)
	}

	messaging.ProduceMessage(conn, message, "auth_queue")
	messaging.ProduceMessage(conn, message, "authz_queue")
	messaging.ProduceMessage(conn, message, "watch_history_queue")

	// Return response with success: true if no error is thrown (and thus document is removed)
	return &userpb.DeleteUserRes{
		Success: true,
	}, nil
}
