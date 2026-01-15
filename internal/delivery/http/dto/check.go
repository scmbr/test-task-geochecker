package dto

type CreateCheckRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}
type CreateCheckResponse struct {
	Count     uint16                `json:"count" binding:"required"`
	Incidents []GetIncidentResponse `json:"incidents" binding:"required"`
}
