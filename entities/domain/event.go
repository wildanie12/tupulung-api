package domain

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title string
	HostedBy string
	Cover string
	UserID uint
	CategoryID uint
	DatetimeEvent time.Time
	Location string
	Description string
}