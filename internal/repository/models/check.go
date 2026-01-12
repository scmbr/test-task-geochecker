package models

import "time"

type Check struct {
	CheckID   string    `gorm:"primaryKey;column:check_id"`
	UserID    string    `gorm:"column:user_id;not null"`
	Longitude float64   `gorm:"column:longitude;not null"`
	Latitude  float64   `gorm:"column:latitude;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}
