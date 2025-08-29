package main

import (
	"context"
	"grpc-course-protobuf/pb/user"
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

	// bidirectional streaming
	userClient := user.NewUserServiceClient(clientConn)
	res, err := userClient.CreateUser(context.Background(), &user.User{
		Age: -1,
	})
	if err != nil {
		// st, ok := status.FromError(err)
		// if ok {
		// 	// error berasal dari grpc
		// 	if st.Code() == codes.InvalidArgument {
		// 		log.Println("there is validation error ", st.Message())
		// 	} else if st.Code() == codes.Unknown {
		// 		log.Println("there is Unknown error ", st.Message())
		// 	} else if st.Code() == codes.Internal {
		// 		log.Println("there is Internal error ", st.Message())
		// 	}
		// 	return
		// }

		log.Println("Failed to send message ", err)
		return
	}
	if !res.Base.IsSuccess {
		if res.Base.StatusCode == 400 {
			log.Println("There is validation error: ", res.Base.Message)
		} else if res.Base.StatusCode == 500 {
			log.Println("There is internal error: ", res.Base.Message)
		}
		return
	}
	log.Println("Response from server ", res.Base.Message)
}
