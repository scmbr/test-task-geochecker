package domain

import "time"

type Operator struct {
	OperatorID string
	APIKeyHash string
	Name       string
	CreatedAt  time.Time
	RevokedAt  *time.Time
	UpdatedAt  *time.Time
}
