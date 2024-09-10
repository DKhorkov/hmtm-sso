package main

import (
	"fmt"
	"github.com/DKhorkov/hmtm-sso/internal/config"
	"github.com/DKhorkov/hmtm-sso/protobuf/generated/go/sso"
)

func main() {
	settings := config.GetConfig()
	fmt.Println(settings)
	fmt.Println(sso.GetUserRequest{})
}
