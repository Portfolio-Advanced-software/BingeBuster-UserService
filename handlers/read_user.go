package handlers

import (
	"context"
	"fmt"

	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/globals"
	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/models"
	userpb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) ReadUser(ctx context.Context, req *userpb.ReadUserReq) (*userpb.ReadUserRes, error) {
	// convert string id (from proto) to mongoDB ObjectId
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := globals.UserDb.FindOne(ctx, bson.M{"_id": oid})
	// Create an empty user to write our decode result to
	data := models.User{}
	// decode and write to data
	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find user with Object Id %s: %v", req.GetId(), err))
	}
	// Cast to ReadMovieRes type
	response := &userpb.ReadUserRes{
		User: &userpb.User{
			Id:               oid.Hex(),
			Email:            data.Email,
			Phone:            data.Phone,
			DateOfBirth:      data.DateOfBirth,
			FirstName:        data.FirstName,
			LastName:         data.LastName,
			CreditCardNumber: data.CreditCardNumber,
			ExpirationDate:   data.ExpirationDate,
			Cvc:              data.CVC,
		},
	}
	return response, nil
}
