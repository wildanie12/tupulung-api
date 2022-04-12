package entities

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title         string
	HostedBy      string
	Cover         string
	UserID        uint
	CategoryID    uint
	DatetimeEvent time.Time
	Location      string
	Description   string
	User          User     `gorm:"foreignKey:UserID;references:ID"`
	Category      Category `gorm:"foreignKey:CategoryID;references:ID"`
}

type EventRequest struct {
	Title         string `form:"title"`
	HostedBy      string `form:"hosted_by"`
	Cover         string `form:"cover"`
	DatetimeEvent string `form:"datetime_event"`
	CategoryID    uint   `form:"category_id"`
	Location      string `form:"location"`
	Description   string `form:"description"`
}

type EventResponse struct {
	ID            uint             `json:"id"`
	Title         string           `json:"title"`
	HostedBy      string           `json:"hosted_by"`
	Cover         string           `json:"cover"`
	DatetimeEvent time.Time        `json:"datetime_event"`
	Location      string           `json:"location"`
	Description   string           `json:"description"`
	CategoryID    uint             `json:"category_id"`
	Category      CategoryResponse `json:"category"`
	UserID        uint             `json:"user_id"`
	User          UserResponse     `json:"user"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}
