package models

import "time"

type Check struct {
	CheckID   string    `gorm:"primaryKey;column:check_id"`
	UserID    string    `gorm:"column:user_id;not null"`
	Location  string    `gorm:"type:geometry(POINT,4326);column:location;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}
