package dto

type CheckInput struct {
	UserID    string
	Latitude  float64
	Longitude float64
}
type GetCheckOutput struct {
	ID        string
	UserID    string
	Latitude  float64
	Longitude float64
}
