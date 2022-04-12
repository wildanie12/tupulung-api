package utilities

import (
	"tupulung/entities"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&entities.User{},
	)
}