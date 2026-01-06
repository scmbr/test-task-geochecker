package service

import "errors"

var (
	ErrIncidentNotFound      = errors.New("incident not found")
	ErrIncidentAlreadyExists = errors.New("incident already exists")
	ErrCheckNotFound         = errors.New("check not found")
	ErrCheckAlreadyExists    = errors.New("check already exists")
	ErrOperatorNotFound      = errors.New("operator not found")
	ErrOperatorAlreadyExists = errors.New("operator already exists")
	ErrInvalidAPIKey         = errors.New("invalid api key")
)
