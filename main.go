package main

import (
	"context"
	"errors"
	"grp-course-protobuf/pb/chat"
	"grp-course-protobuf/pb/user"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

// client streaming
type chatService struct {
	chat.UnimplementedChatServiceServer
}

// client streaming
func (cs *chatService) SendMessage(stream grpc.ClientStreamingServer[chat.ChatMessage, chat.ChatResponse]) error {
	// trik infinite loop agar server bisa menerima banyak pesan
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return status.Errorf(codes.Unknown, "error receiving message: %v", err)
		}
		log.Printf("Receive message: %s, to %d", req.Content, req.UserId)
	}

	return stream.SendAndClose(&chat.ChatResponse{
		Message: "Thanks for the message",
	})
}

// server streaming
func (cs *chatService) ReceiveMessage(req *chat.ReceiveMessageRequest, stream grpc.ServerStreamingServer[chat.ChatMessage]) error {
	log.Printf("Got connection request from %d\n", req.UserId)

	for i := 0; i < 10; i++ {
		err := stream.Send(&chat.ChatMessage{
			UserId:  123,
			Content: "Hi",
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "error sending message to client %v", err)
		}
	}
	return nil
}

// func (UnimplementedChatServiceServer) Chat(grpc.BidiStreamingServer[ChatMessage, ChatMessage]) error {
// 	return status.Errorf(codes.Unimplemented, "method Chat not implemented")
// }

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("There is error in your net listen ", err)
	}

	serv := grpc.NewServer()

	user.RegisterUserServiceServer(serv, &userService{})
	chat.RegisterChatServiceServer(serv, &chatService{})

	reflection.Register(serv)

	if err := serv.Serve(lis); err != nil {
		log.Fatal("Error running server ", err)
	}
}
