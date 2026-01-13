package dto

import "time"

type CreateIncidentInput struct {
	OperatorID string
	Latitude   float64
	Longitude  float64
	Radius     uint16
}
type GetAllIncidentsInput struct {
	Limit  int
	Offset int
}
type GetAllIncidentsOutput struct {
	Total     uint32
	Incidents []GetIncidentOutput
}

type GetIncidentOutput struct {
	ID         string
	OperatorID string
	Latitude   float64
	Longitude  float64
	Radius     uint16
	CreatedAt  time.Time
}
type UpdateIncidentInput struct {
	OperatorID *string
	Latitude   *float64
	Longitude  *float64
	Radius     *uint16
}
