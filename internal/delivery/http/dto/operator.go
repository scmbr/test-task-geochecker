package dto

type CreateOperatorRequest struct {
	Name   string `json:"name"  binding:"required,min=1,max=100"`
	APIKey string `json:"api_key" binding:"required"`
}
