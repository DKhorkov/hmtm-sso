package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/libs/requestid"
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

	ctx := metadata.AppendToOutgoingContext(context.Background(), requestid.Key, requestid.New())
	// tokens, err := client.Login(
	//	ctx,
	//	&sso.LoginIn{
	//		Email:     "test@mail.test",
	//		Password:  "qwer1234",
	//	},
	//)
	// fmt.Println(tokens, err)

	users, err := client.GetUsers(ctx, &emptypb.Empty{})
	fmt.Println(users, err)

	// userID, err := client.Register(ctx, &sso.RegisterIn{
	//	DisplayName: "test User",
	//	Email:       "sometestemail2@yandex.ru",
	//	Password:    "test@Password2",
	// })
	// fmt.Println("Register: ", userID, err)
}
