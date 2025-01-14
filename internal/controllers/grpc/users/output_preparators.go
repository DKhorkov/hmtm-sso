package users

import (
	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/hmtm-sso/internal/entities"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func prepareUserOut(user entities.User) *sso.GetUserOut {
	return &sso.GetUserOut{
		ID:                user.ID,
		DisplayName:       user.DisplayName,
		Email:             user.Email,
		EmailConfirmed:    user.EmailConfirmed,
		Phone:             user.Phone,
		PhoneConfirmed:    user.PhoneConfirmed,
		Telegram:          user.Telegram,
		TelegramConfirmed: user.TelegramConfirmed,
		Avatar:            user.Avatar,
		CreatedAt:         timestamppb.New(user.CreatedAt),
		UpdatedAt:         timestamppb.New(user.UpdatedAt),
	}
}
