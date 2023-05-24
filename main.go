package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	models "github.com/Portfolio-Advanced-software/BingeBuster-UserService/models"
	mongodb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/mongodb"
	userpb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
}

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
	result, err := userdb.InsertOne(mongoCtx, data)
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

func (s *UserServiceServer) ReadUser(ctx context.Context, req *userpb.ReadUserReq) (*userpb.ReadUserRes, error) {
	// convert string id (from proto) to mongoDB ObjectId
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := userdb.FindOne(ctx, bson.M{"_id": oid})
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

func (s *UserServiceServer) ListUsers(req *userpb.ListUsersReq, stream userpb.UserService_ListUsersServer) error {
	// Initiate a movie type to write decoded data to
	data := &models.User{}
	// collection.Find returns a cursor for our (empty) query
	cursor, err := userdb.Find(context.Background(), bson.M{})
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
	result := userdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

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

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserReq) (*userpb.DeleteUserRes, error) {
	// Get the ID (string) from the request message and convert it to an Object ID
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	// Check for errors
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	// DeleteOne returns DeleteResult which is a struct containing the amount of deleted docs (in this case only 1 always)
	// So we return a boolean instead
	_, err = userdb.DeleteOne(ctx, bson.M{"_id": oid})
	// Check for errors
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete user with id %s: %v", req.GetId(), err))
	}
	// Return response with success: true if no error is thrown (and thus document is removed)
	return &userpb.DeleteUserRes{
		Success: true,
	}, nil
}

const (
	port = ":50053"
)

var db *mongo.Client
var userdb *mongo.Collection
var mongoCtx context.Context

var mongoUsername = "user-service"
var mongoPwd = "vLxxhmS0eJFwmteF"
var connUri = "mongodb+srv://" + mongoUsername + ":" + mongoPwd + "@cluster0.fpedw5d.mongodb.net/"

var dbName = "UserService"
var collectionName = "Users"

func main() {
	// Configure 'log' package to give file name and line number on eg. log.Fatal
	// Pipe flags to one another (log.LstdFLags = log.Ldate | log.Ltime)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Starting server on port :50053...")

	// Set listener to start server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Unable to listen on port %p: %v", lis.Addr(), err)
	}

	// Set options, here we can configure things like TLS support
	opts := []grpc.ServerOption{}
	// Create new gRPC server with (blank) options
	s := grpc.NewServer(opts...)
	// Create UserService type
	srv := &UserServiceServer{}

	// Register the service with the server
	userpb.RegisterUserServiceServer(s, srv)

	// Initialize MongoDb client
	fmt.Println("Connecting to MongoDB...")
	db = mongodb.ConnectToMongoDB(connUri)

	// Bind our collection to our global variable for use in other methods
	userdb = db.Database(dbName).Collection(collectionName)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :50053")

	// Right way to stop the server using a SHUTDOWN HOOK
	// Create a channel to receive OS signals
	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	// As long as user doesn't press CTRL+C a message is not passed and our main routine keeps running
	<-c

	// After receiving CTRL+C Properly stop the server
	fmt.Println("\nStopping the server...")
	s.Stop()
	lis.Close()
	fmt.Println("Closing MongoDB connection")
	db.Disconnect(mongoCtx)
	fmt.Println("Done.")

}
