package web

import "time"

type EventResponse struct {
	ID uint `json:"id"`
	Title string `json:"title"`
	HostedBy string `json:"hosted_by"`
	Cover string `json:"cover"`
	DatetimeEvent string `json:"datetime_event"`
	Location string `json:"location"`
	Description string `json:"description"`
	CategoryID uint `json:"category_id"`
	Category CategoryResponse `json:"category"`
	UserID uint `json:"user_id"`
	User UserResponse `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}