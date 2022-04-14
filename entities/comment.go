package entities

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	EventID uint
	UserID uint
	Comment string
	Event Event `gorm:"foreignKey:EventID;references:ID"`
	User User `gorm:"foreignKey:UserID;references:ID"`
}

type CommentRequest struct {
	Comment string `form:"comment"`
}

type CommentResponse struct {
	ID uint `json:"id"`
	EventID uint `json:"event_id"`
	UserID uint `json:"user_id"`
	User UserResponse `json:"user"`
	Comment string `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}