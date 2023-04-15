package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/protos"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserClient(conn)

	runGetUsers(client)

}

func runGetUsers(client pb.UserClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Empty{}
	stream, err := client.GetUsers(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetUsers(_) = _, %v", client, err)
	}
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetUsers(_) = _, %v", client, err)
		}
		log.Printf("UserInfo: %v", row)
	}
}

func runGetUser(client pb.UserClient, userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: userid}
	res, err := client.GetUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetUser(_) = _, %v", client, err)
	}
	log.Printf("UserInfo: %v", res)
}

func runCreateUser(client pb.UserClient, watchedmovies string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.UserInfo{Watchedmovies: watchedmovies}
	res, err := client.CreateUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.CreateUser(_) = _, %v", client, err)
	}
	if res.GetValue() != "" {
		log.Printf("CreateUser Id: %v", res)
	} else {
		log.Printf("CreateUser Failed")
	}

}

func runUpdateUser(client pb.UserClient, userid string, watchedmovies string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.UserInfo{Id: userid, Watchedmovies: watchedmovies}
	res, err := client.UpdateUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.UpdateUser(_) = _, %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("UpdateUser Success")
	} else {
		log.Printf("UpdateUser Failed")
	}
}

func runDeleteUser(client pb.UserClient, userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: userid}
	res, err := client.DeleteUser(ctx, req)
	if err != nil {
		log.Fatalf("%v.DeleteUser(_) = _, %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("DeleteUser Success")
	} else {
		log.Printf("DeleteUser Failed")
	}
}
