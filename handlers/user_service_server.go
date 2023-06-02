package handlers

import (
	userpb "github.com/Portfolio-Advanced-software/BingeBuster-UserService/proto"
)

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
}
