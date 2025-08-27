package main

import (
	"context"
	"errors"
	"grp-course-protobuf/pb/chat"
	"io"
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

	// client streaming
	chatClient := chat.NewChatServiceClient(clientConn)
	stream, err := chatClient.ReceiveMessage(context.Background(), &chat.ReceiveMessageRequest{
		UserId: 30})
	if err != nil {
		log.Fatal("failed to send message ", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal("failed to receive message ", err)
		}
		log.Printf("got  message to %d content %s", msg.UserId, msg.Content)
	}

}
