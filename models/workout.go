package models

type Workout struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null"`
	Time     int    `json:"target_time"`
	Calories int    `json:"estimated_calories"`
}
