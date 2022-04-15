package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `gorm:"unique"`
	Password  string
	Name      string
	Gender    string
	Address   string
	Avatar    string
	DOB       time.Time
	DarkTheme string
	Events	  []Event `gorm:"many2many:participants;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:EventID"`
}

type UserRequest struct {
	Name     string `form:"name" validate:"required"`
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
	Gender   string `form:"gender" validate:"required"`
	Address  string `form:"address"`
	Avatar   string `form:"avatar"`
	DOB      string `form:"dob" validate:"required"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	Address   string    `json:"address"`
	Avatar    string    `json:"avatar"`
	DOB       time.Time `json:"dob"`
	DarkTheme bool      `json:"dark_theme"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
