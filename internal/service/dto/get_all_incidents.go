package dto

type GetAllIncidentOutput struct {
	Total     uint32
	Incidents []GetIncidentOutput
}
