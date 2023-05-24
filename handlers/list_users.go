package handlers

import (
	"context"
	"fmt"

	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/globals"
	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/models"
	userpb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) ListUsers(req *userpb.ListUsersReq, stream userpb.UserService_ListUsersServer) error {
	// Initiate a movie type to write decoded data to
	data := &models.User{}
	// collection.Find returns a cursor for our (empty) query
	cursor, err := globals.UserDb.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	// An expression with defer will be called at the end of the function
	defer cursor.Close(context.Background())
	// cursor.Next() returns a boolean, if false there are no more items and loop will break
	for cursor.Next(context.Background()) {
		// Decode the data at the current pointer and write it to data
		err := cursor.Decode(data)
		// check error
		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}
		// If no error is found send blog over stream
		stream.Send(&userpb.ListUsersRes{
			User: &userpb.User{
				Id:               data.ID.Hex(),
				Email:            data.Email,
				Phone:            data.Phone,
				DateOfBirth:      data.DateOfBirth,
				FirstName:        data.FirstName,
				LastName:         data.LastName,
				CreditCardNumber: data.CreditCardNumber,
				ExpirationDate:   data.ExpirationDate,
				Cvc:              data.CVC,
			},
		})
	}
	// Check if the cursor has any errors
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}
	return nil
}
