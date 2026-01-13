package dto

import "time"

type CreateIncidentRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Radius    uint16  `json:"radius" binding:"required"`
}

type GetAllIncidentsRequest struct {
	Offset int `json:"offset" binding:"required"`
	Limit  int `json:"limit" binding:"required"`
}
type GetAllIncidentsResponse struct {
	Total     int32                 `json:"total"`
	Incidents []GetIncidentResponse `json:"incidents"`
}
type GetIncidentResponse struct {
	ID         string    `json:"id"`
	OperatorID string    `json:"operator_id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Radius     uint16    `json:"radius"`
	CreatedAt  time.Time `json:"created_at"`
}
type UpdateIncidentRequest struct {
	OperatorID *string  `json:"operator_id"`
	Latitude   *float64 `json:"latitude"`
	Longitude  *float64 `json:"longitude" `
	Radius     *uint16  `json:"radius"`
}
