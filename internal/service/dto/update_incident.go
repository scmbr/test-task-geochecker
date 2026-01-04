package dto

type UpdateIncidentInput struct {
	OperatorID string
	Latitude   float64
	Longitude  float64
	Radius     uint8
}
