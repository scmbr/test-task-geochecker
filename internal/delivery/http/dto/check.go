package dto

type CreateCheckRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	Latitude  float64 `json:"latitude_id" binding:"required"`
	Longitude float64 `json:"longitude_id" binding:"required"`
}
