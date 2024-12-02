package entities

import "time"

type RefreshToken struct {
	ID        uint64    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserID    uint64    `json:"userID" gorm:"not null"`
	TTL       time.Time `json:"ttl" gorm:"not null"`
	Value     string    `json:"value" gorm:"unique; not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"not null"`
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
