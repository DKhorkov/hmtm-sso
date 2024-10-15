package entities

import "time"

type RefreshToken struct {
	ID        int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserID    int       `json:"user_id" gorm:"not null"`
	TTL       time.Time `json:"TTL" gorm:"not null"`
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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
