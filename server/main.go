package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"strconv"

	pb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/protos"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var users []*pb.UserInfo

type userServer struct {
	pb.UnimplementedUserServer
}

func main() {
	initUsers()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterUserServer(s, &userServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initUsers() {
	user1 := &pb.UserInfo{Id: "1", Watchedmovies: "4"}
	user2 := &pb.UserInfo{Id: "2", Watchedmovies: "16"}

	users = append(users, user1)
	users = append(users, user2)
}

func (s *userServer) GetUsers(in *pb.Empty,
	stream pb.User_GetUsersServer) error {
	log.Printf("Received: %v", in)
	for _, user := range users {
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	return nil
}

func (s *userServer) GetUser(ctx context.Context,
	in *pb.Id) (*pb.UserInfo, error) {
	log.Printf("Received: %v", in)

	res := &pb.UserInfo{}

	for _, user := range users {
		if user.GetId() == in.GetValue() {
			res = user
			break
		}
	}
	return res, nil
}

func (s *userServer) CreateUser(ctx context.Context,
	in *pb.UserInfo) (*pb.Id, error) {
	log.Printf("Received: &v", in)
	res := pb.Id{}
	res.Value = strconv.Itoa(rand.Intn(100000000))
	in.Id = res.GetValue()
	users = append(users, in)
	return &res, nil
}

func (s *userServer) UpdateUser(ctx context.Context,
	in *pb.UserInfo) (*pb.Status, error) {
	log.Printf("Received: &v", in)

	res := pb.Status{}
	for index, user := range users {
		if user.GetId() == in.GetId() {
			users = append(users[:index], users[index+1:]...)
			in.Id = user.GetId()
			users = append(users, in)
			res.Value = 1
			break
		}
	}
	return &res, nil
}

func (s *userServer) DeleteUser(ctx context.Context,
	in *pb.Id) (*pb.Status, error) {
	log.Printf("Received: &v", in)

	res := pb.Status{}
	for index, user := range users {
		if user.GetId() == in.GetValue() {
			users = append(users[:index], users[index+1:]...)
			res.Value = 1
			break
		}
	}
	return &res, nil
}
