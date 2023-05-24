package handlers

import (
	"context"
	"fmt"

	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/globals"
	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/models"
	userpb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserReq) (*userpb.UpdateUserRes, error) {
	// Get the user data from the request
	user := req.GetUser()

	// Convert the Id string to a MongoDB ObjectId
	oid, err := primitive.ObjectIDFromHex(user.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied user id to a MongoDB ObjectId: %v", err),
		)
	}

	// Convert the data to be updated into an unordered Bson document
	update := bson.M{
		"email":            user.GetEmail(),
		"phone":            user.GetPhone(),
		"dateofbirth":      user.GetDateOfBirth(),
		"firstname":        user.GetFirstName(),
		"lastname":         user.GetLastName(),
		"creditcardnumber": user.GetCreditCardNumber(),
		"expirationdate":   user.GetExpirationDate(),
		"cvc":              user.GetCvc(),
	}

	// Convert the oid into an unordered bson document to search by id
	filter := bson.M{"_id": oid}

	// Result is the BSON encoded result
	// To return the updated document instead of original we have to add options.
	result := globals.UserDb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	// Decode result and write it to 'decoded'
	decoded := models.User{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find user with supplied ID: %v", err),
		)
	}
	return &userpb.UpdateUserRes{
		User: &userpb.User{
			Id:               decoded.ID.Hex(),
			Email:            decoded.Email,
			Phone:            decoded.Phone,
			DateOfBirth:      decoded.DateOfBirth,
			FirstName:        decoded.FirstName,
			LastName:         decoded.LastName,
			CreditCardNumber: decoded.CreditCardNumber,
			ExpirationDate:   decoded.ExpirationDate,
			Cvc:              decoded.CVC,
		},
	}, nil
}
