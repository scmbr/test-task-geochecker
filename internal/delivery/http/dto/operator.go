package dto

type CreateOperatorRequest struct {
	Name   string `json:"name"  binding:"required"`
	APIKey string `json:"api_key" binding:"required"`
}
