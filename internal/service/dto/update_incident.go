package dto

type UpdateIncidentInput struct {
	ID         string
	OperatorID string
	Latitude   float64
	Longitude  float64
	Radius     float64
	CreatedAt  string
}
