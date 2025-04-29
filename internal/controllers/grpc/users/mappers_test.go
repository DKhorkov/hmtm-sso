package users

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"
	"github.com/DKhorkov/hmtm-sso/internal/entities"
)

func TestMapUserToOut(t *testing.T) {
	testCases := []struct {
		name     string
		user     entities.User
		expected *sso.GetUserOut
	}{
		{
			name: "full user",
			user: entities.User{
				ID:                1,
				DisplayName:       "John Doe",
				Email:             "john@example.com",
				EmailConfirmed:    true,
				Password:          "hashedpassword",
				Phone:             pointers.New("1234567890"),
				PhoneConfirmed:    true,
				Telegram:          pointers.New("@johndoe"),
				TelegramConfirmed: true,
				Avatar:            pointers.New("avatar.jpg"),
				CreatedAt:         time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:         time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			expected: &sso.GetUserOut{
				ID:                1,
				DisplayName:       "John Doe",
				Email:             "john@example.com",
				EmailConfirmed:    true,
				Phone:             pointers.New("1234567890"),
				PhoneConfirmed:    true,
				Telegram:          pointers.New("@johndoe"),
				TelegramConfirmed: true,
				Avatar:            pointers.New("avatar.jpg"),
				CreatedAt:         timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt:         timestamppb.New(time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "minimal user",
			user: entities.User{
				ID:             2,
				DisplayName:    "Jane Doe",
				Email:          "jane@example.com",
				EmailConfirmed: false,
				Password:       "hashedpassword",
				CreatedAt:      time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2023, 2, 2, 0, 0, 0, 0, time.UTC),
			},
			expected: &sso.GetUserOut{
				ID:                2,
				DisplayName:       "Jane Doe",
				Email:             "jane@example.com",
				EmailConfirmed:    false,
				PhoneConfirmed:    false,
				TelegramConfirmed: false,
				CreatedAt:         timestamppb.New(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt:         timestamppb.New(time.Date(2023, 2, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "user with partial optional fields",
			user: entities.User{
				ID:                3,
				DisplayName:       "Bob Smith",
				Email:             "bob@example.com",
				EmailConfirmed:    true,
				Password:          "hashedpassword",
				Phone:             nil,
				PhoneConfirmed:    false,
				Telegram:          pointers.New("@bobsmith"),
				TelegramConfirmed: false,
				Avatar:            nil,
				CreatedAt:         time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:         time.Date(2023, 3, 2, 0, 0, 0, 0, time.UTC),
			},
			expected: &sso.GetUserOut{
				ID:                3,
				DisplayName:       "Bob Smith",
				Email:             "bob@example.com",
				EmailConfirmed:    true,
				PhoneConfirmed:    false,
				Telegram:          pointers.New("@bobsmith"),
				TelegramConfirmed: false,
				CreatedAt:         timestamppb.New(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt:         timestamppb.New(time.Date(2023, 3, 2, 0, 0, 0, 0, time.UTC)),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := mapUserToOut(tc.user)

			// Проверка полей
			require.Equal(t, tc.expected.ID, result.ID)
			require.Equal(t, tc.expected.DisplayName, result.DisplayName)
			require.Equal(t, tc.expected.Email, result.Email)
			require.Equal(t, tc.expected.EmailConfirmed, result.EmailConfirmed)
			require.Equal(t, tc.expected.Phone, result.Phone)
			require.Equal(t, tc.expected.PhoneConfirmed, result.PhoneConfirmed)
			require.Equal(t, tc.expected.Telegram, result.Telegram)
			require.Equal(t, tc.expected.TelegramConfirmed, result.TelegramConfirmed)
			require.Equal(t, tc.expected.Avatar, result.Avatar)

			// Проверка временных меток
			require.Equal(t, tc.expected.CreatedAt.AsTime(), result.CreatedAt.AsTime())
			require.Equal(t, tc.expected.UpdatedAt.AsTime(), result.UpdatedAt.AsTime())
		})
	}
}
