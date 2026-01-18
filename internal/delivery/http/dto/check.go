package dto

type CreateCheckRequest struct {
	UserID    string  `json:"user_id" binding:"required,uuid4"`
	Latitude  float64 `json:"latitude" binding:"required,gte=-90,lte=90"`
	Longitude float64 `json:"longitude" binding:"required,gte=-180,lte=180"`
}

type CreateCheckResponse struct {
	Count     uint16                `json:"count" binding:"required"`
	Incidents []GetIncidentResponse `json:"incidents" binding:"required"`
}
