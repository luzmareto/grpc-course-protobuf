package main

import (
	"context"
	"errors"
	"grp-course-protobuf/pb/chat"
	"grp-course-protobuf/pb/common"
	"grp-course-protobuf/pb/user"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// unary
type userService struct {
	user.UnimplementedUserServiceServer
}

func (us *userService) CreateUser(ctx context.Context, userRequest *user.User) (*user.CreateResponse, error) {
	// menggunakan wrapper response
	if userRequest.Age < 1 {
		return &user.CreateResponse{
			Base: &common.BaseResponse{
				StatusCode: 400,
				IsSuccess:  false,
				Message:    "validation error: age must be greater than 0",
			},
		}, nil
	}
	// membuat response server

	log.Println("User is created")
	return &user.CreateResponse{
		Base: &common.BaseResponse{
			StatusCode: 200,
			IsSuccess:  true,
			Message:    "User created",
		},
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

// bidirectional streaming
func (cs *chatService) Chat(stream grpc.BidiStreamingServer[chat.ChatMessage, chat.ChatMessage]) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return status.Errorf(codes.Unknown, "error receiving message")
		}

		log.Printf("Got message from %d content: %s", msg.UserId, msg.Content)

		time.Sleep(2 * time.Second)

		err = stream.Send(&chat.ChatMessage{
			UserId:  50,
			Content: "Reply from server",
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "error sending message")
		}
		err = stream.Send(&chat.ChatMessage{
			UserId:  50,
			Content: "Reply from server #2",
		})
		if err != nil {
			return status.Errorf(codes.Unknown, "error sending message")
		}
	}

	return nil
}

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
