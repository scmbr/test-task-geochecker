package dto

import "time"

type ValidateOperatorOutput struct {
	OperatorID string
	Name       string
	CreatedAt  time.Time
	RevokedAt  *time.Time
}
