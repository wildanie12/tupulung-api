package entities

type Like struct {
	ID      uint `gorm:"primary_key;auto_increment;not_null"`
	EventID uint
	UserID  uint
}

type LikeRequest struct {
	EventID uint `form:"event_id"`
}
