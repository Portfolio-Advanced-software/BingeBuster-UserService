package handlers

import (
	"context"
	"fmt"

	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/globals"
	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/models"
	userpb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserReq) (*userpb.CreateUserRes, error) {
	// Essentially doing req.User to access the struct with a nil check
	user := req.GetUser()
	if user == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid user")
	}
	// Now we have to convert this into a User type to convert into BSON
	data := models.User{
		// ID:    Empty, so it gets omitted and MongoDB generates a unique Object ID upon insertion.
		Email:            user.GetEmail(),
		Phone:            user.GetPhone(),
		DateOfBirth:      user.GetDateOfBirth(),
		FirstName:        user.GetFirstName(),
		LastName:         user.GetLastName(),
		CreditCardNumber: user.GetCreditCardNumber(),
		ExpirationDate:   user.GetExpirationDate(),
		CVC:              user.GetCvc(),
	}

	// Insert the data into the database, result contains the newly generated Object ID for the new document
	result, err := globals.UserDb.InsertOne(globals.MongoCtx, data)
	// check for potential errors
	if err != nil {
		// return internal gRPC error to be handled later
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	// add the id to movie, first cast the "generic type" (go doesn't have real generics yet) to an Object ID.
	oid := result.InsertedID.(primitive.ObjectID)
	// Convert the object id to it's string counterpart
	user.Id = oid.Hex()
	// return the blog in a CreateMovieRes type
	return &userpb.CreateUserRes{User: user}, nil
}
