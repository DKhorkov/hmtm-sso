package entities

import "time"

type RefreshToken struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"userId"`
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
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type TokensDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
