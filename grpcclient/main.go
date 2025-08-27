package main

import (
	"context"
	"grp-course-protobuf/pb/user"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// client server dengan melakukan bypass security karena masih di sisi development local
	clientConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to create client ", err)
	}

	userClient := user.NewUserServiceClient(clientConn)

	response, err := userClient.CreateUser(context.Background(), &user.User{
		Id:      1,
		Age:     13,
		Balance: 130000,
		Address: &user.Address{
			Id:          123,
			FullAddress: "Jalan Merdeka No 123",
			Province:    "DKI Jakarta",
			City:        "Jakarta Selatan",
		},
	})
	if err != nil {
		log.Fatal("error calling user client", err)
	}
	log.Println("response from server", response.Message)
}
