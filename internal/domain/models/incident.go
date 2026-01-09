package models

import "time"

type Incident struct {
	IncidentID string     `gorm:"primaryKey;column:incident_id"`
	OperatorID string     `gorm:"column:operator_id;not null"`
	Longitude  float64    `gorm:"column:longitude;not null"`
	Latitude   float64    `gorm:"column:latitude;not null"`
	Radius     uint8      `gorm:"column:radius;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;default:null"`
	UpdatedAt  *time.Time `gorm:"column:updated_at;default:null"`
}
type UpdateIncidentInput struct {
	OperatorID *string  `json:"operator_id"`
	Latitude   *float64 `json:"latitude"`
	Longitude  *float64 `json:"longitude"`
	Radius     *uint8   `json:"radius"`
}
