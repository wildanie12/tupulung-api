package entities

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Title string
}

type CategoryRequest struct {
	Title string `form:"title"`
}

type CategoryResponse struct {
	ID uint `json:"id"`
	Title string `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}