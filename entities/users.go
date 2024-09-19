package entities

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password" gorm:"not null"`
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
