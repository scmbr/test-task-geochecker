package dto

import "time"

type CreateIncidentRequest struct {
	Latitude  float64 `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude float64 `json:"longitude" binding:"required,gte=-180,lte=180"`
	Radius    uint16  `json:"radius" binding:"required,gt=0"`
}

type GetAllIncidentsRequest struct {
	Offset int `form:"offset" binding:"omitempty,min=0"`
	Limit  int `form:"limit" binding:"omitempty,min=1"`
}
type GetAllIncidentsResponse struct {
	Total     int32                 `json:"total"`
	Incidents []GetIncidentResponse `json:"incidents"`
}
type GetIncidentResponse struct {
	IncidentID string    `json:"incident_id"`
	OperatorID string    `json:"operator_id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Radius     uint16    `json:"radius"`
	CreatedAt  time.Time `json:"created_at"`
}
type UpdateIncidentRequest struct {
	OperatorID *string  `json:"operator_id" binding:"omitempty,uuid4"`
	Latitude   *float64 `json:"latitude" binding:"omitempty,gte=-90,lte=90"`
	Longitude  *float64 `json:"longitude" binding:"omitempty,gte=-180,lte=180"`
	Radius     *uint16  `json:"radius" binding:"omitempty,gt=0"`
}
type GetIncidentStatsByIdResponse struct {
	IncidentID   string    `json:"incident_id"`
	UserCount    int       `json:"user_count"`
	SinceMinutes time.Time `json:"since_minutes"`
}
