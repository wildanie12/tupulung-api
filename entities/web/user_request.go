package web

type UserRequest struct {
	Name string `form:"name"`
	Email string `form:"email"`
	Password string `form:"password"`
	Gender string `form:"gender"`
	Address string `form:"address"`
	Avatar string `form:"avatar"`
	DOB string `form:"dob"`
}