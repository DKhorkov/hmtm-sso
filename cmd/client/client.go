package main

import (
	"context"
	"fmt"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/libs/requestid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	sso.AuthServiceClient
	sso.UsersServiceClient
}

func main() {
	clientConnection, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", "0.0.0.0", 8070),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		panic(err)
	}

	client := &Client{
		UsersServiceClient: sso.NewUsersServiceClient(clientConnection),
		AuthServiceClient:  sso.NewAuthServiceClient(clientConnection),
	}

	requestID := requestid.New()
	// tokens, err := client.Login(
	//	context.Background(),
	//	&sso.LoginIn{
	//		RequestID: requestID,
	//		Email:     "test@mail.test",
	//		Password:  "qwer1234",
	//	},
	//)
	// fmt.Println(tokens, err)

	// users, err := client.GetUsers(context.Background(), &sso.GetUsersIn{RequestID: requestID})
	// fmt.Println(users, err)

	userID, err := client.Register(context.Background(), &sso.RegisterIn{
		RequestID: requestID,
		Email:     "sometestemail@yandex.ru",
		Password:  "test@Password2",
	})
	fmt.Println("Register: ", userID, err)
}
