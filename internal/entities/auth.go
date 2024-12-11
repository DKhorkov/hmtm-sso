package entities

import "time"

type RefreshToken struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"userID"`
	TTL       time.Time `json:"ttl"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserDTO struct {
	Credentials LoginUserDTO
}

type TokensDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
