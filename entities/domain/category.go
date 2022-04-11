package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Title string
}