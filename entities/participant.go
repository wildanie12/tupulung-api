package entities

type Participant struct {
	ID      uint `gorm:"primary_key;auto_increment;not_null"`
	EventID uint
	UserID  uint
}

type ParticipantRequest struct {
	EventID uint `form:"event_id"`
}
