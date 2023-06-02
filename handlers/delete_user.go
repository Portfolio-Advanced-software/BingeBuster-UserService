package handlers

import (
	"context"
	"fmt"

	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/globals"
	messaging "github.com/Portfolio-Advanced-software/BingeBuster-UserService/messaging"
	userpb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserReq) (*userpb.DeleteUserRes, error) {
	// Get the ID (string) from the request message and convert it to an Object ID
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	// Check for errors
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	// DeleteOne returns DeleteResult which is a struct containing the amount of deleted docs (in this case only 1 always)
	// So we return a boolean instead
	_, err = globals.UserDb.DeleteOne(ctx, bson.M{"_id": oid})
	// Check for errors
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete user with id %s: %v", req.GetId(), err))
	}

	// Send a message to the watch history queue for deleting user's history
	message := map[string]interface{}{
		"userId": req.GetId(),
		"action": "deleteAllRecords",
	}

	conn, err := messaging.ConnectToRabbitMQ(globals.RabbitMQUrl)

	queueName := "watch_history_queue"
	messaging.ProduceMessage(conn, message, queueName)

	// Return response with success: true if no error is thrown (and thus document is removed)
	return &userpb.DeleteUserRes{
		Success: true,
	}, nil
}
