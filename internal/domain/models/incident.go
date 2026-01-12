package models

import "time"

type Incident struct {
	IncidentID string     `gorm:"primaryKey;column:incident_id"`
	OperatorID string     `gorm:"column:operator_id;not null"`
	Location   string     `gorm:"type:geometry(POINT,4326);column:location;not null"`
	Radius     uint32     `gorm:"column:radius;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;default:null"`
	UpdatedAt  *time.Time `gorm:"column:updated_at;default:null"`
}
type UpdateIncidentInput struct {
	OperatorID *string  `json:"operator_id"`
	Latitude   *float64 `json:"latitude"`
	Longitude  *float64 `json:"longitude"`
	Radius     *uint32  `json:"radius"`
}
