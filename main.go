package main

import (
	"context"
	"grp-course-protobuf/pb/user"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type userService struct {
	user.UnimplementedUserServiceServer
}

func (us *userService) CreateUser(ctx context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	log.Println("user is created")
	return &user.CreateResponse{
		Message: "user created",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("There is error in your net listen ", err)
	}

	serv := grpc.NewServer()

	user.RegisterUserServiceServer(serv, &userService{})

	reflection.Register(serv)

	if err := serv.Serve(lis); err != nil {
		log.Fatal("Error running server ", err)
	}
}
