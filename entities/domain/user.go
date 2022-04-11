package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email string
	Password string
	Name string
	Gender string
	Address string
	Avatar string
	DOB string
	DarkTheme string
}