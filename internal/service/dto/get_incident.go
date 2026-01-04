package dto

import "time"

type GetIncidentOutput struct {
	ID         string
	OperatorID string
	Latitude   float64
	Longitude  float64
	Radius     uint8
	CreatedAt  time.Time
}
