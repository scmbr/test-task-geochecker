package dto

type CreateIncidentInput struct {
	OperatorID string
	Latitude   float64
	Longitude  float64
	Radius     float64
}
