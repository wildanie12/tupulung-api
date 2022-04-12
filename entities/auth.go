package entities

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type AuthRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}
