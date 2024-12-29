package entities

import "time"

type RefreshToken struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	TTL       time.Time `json:"ttl"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserDTO struct {
	Credentials LoginUserDTO
}

type TokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
