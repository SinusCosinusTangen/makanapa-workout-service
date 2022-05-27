package models

type Request struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	UserID string `json:"user_id" gorm:"not null"`
	Time   int    `json:"target_time" gorm:"not null"`
}
