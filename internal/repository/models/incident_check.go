package models

type IncidentsCheck struct {
	IncidentID string `gorm:"column:incident_id;primaryKey"`
	CheckID    string `gorm:"column:check_id;primaryKey"`
}
