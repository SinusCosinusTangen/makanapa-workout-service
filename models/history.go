package models

import (
	"time"
)

type History struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Timestamp time.Time `json:"timestamp" gorm:"not null;autoCreateTime"`
	Workouts  []Workout `json:"workouts" gorm:"many2many:history_workouts;"`
}
