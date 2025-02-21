package main

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/requestid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
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

	// users, err := client.GetUsers(ctx, &emptypb.Empty{})
	// fmt.Println(users, err)

	userID, err := client.Register(ctx, &sso.RegisterIn{
		DisplayName: "test User",
		Email:       "alexqwerty35@yandex.ru",
		Password:    "Qwer1234@",
	})
	fmt.Println("Register: ", userID, err)

	tokens, err := client.Login(ctx, &sso.LoginIn{
		Email:    "alexqwerty35@yandex.ru",
		Password: "Qwer1234@",
	})
	fmt.Println(tokens, err)

	//
	// _, logoutErr := client.Logout(ctx, &sso.LogoutIn{
	//	AccessToken: tokens.GetAccessToken(),
	// })
	// fmt.Println(logoutErr)

	// _, err = client.VerifyEmail(ctx, &sso.VerifyEmailIn{VerifyEmailToken: "MjM="})
	// fmt.Println(err)

	// _, err = client.SendVerifyEmailMessage(ctx, &sso.SendVerifyEmailMessageIn{Email: "alexqwerty35@yandex.ru"})
	// fmt.Printf("%v\n", err)

	// _, err = client.ChangePassword(
	//	ctx,
	//	&sso.ChangePasswordIn{
	//		AccessToken: tokens.AccessToken,
	//		OldPassword: "K8NXoxwVE0vCEjJC",
	//		NewPassword: "Qwer1234@",
	//	},
	//)
	// fmt.Println(err)

	_, err = client.ForgetPassword(ctx, &sso.ForgetPasswordIn{AccessToken: tokens.GetAccessToken()})
	fmt.Println(err)
}
