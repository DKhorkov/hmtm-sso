package entities

import "time"

type User struct {
	ID             uint64    `json:"id"`
	DisplayName    string    `json:"display_name"`
	Email          string    `json:"email"`
	EmailConfirmed bool      `json:"email_confirmed"`
	Password       string    `json:"password"`
	Phone          *string   `json:"phone,omitempty"`
	PhoneConfirmed bool      `json:"phone_confirmed"`
	Telegram       *string   `json:"telegram,omitempty"`
	Avatar         *string   `json:"avatar,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
