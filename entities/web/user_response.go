package web

import "time"

type UserResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Gender string `json:"gender"`
	Address string `json:"address"`
	Avatar string `json:"avatar"`
	DOB string `json:"dob"`
	DarkTheme bool `json:"dark_theme"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}