package domain

import (
	"errors"
	"time"
)

type Operator struct {
	OperatorID string
	APIKeyHash string
	Name       string
	CreatedAt  time.Time
	RevokedAt  *time.Time
	UpdatedAt  *time.Time
}

func NewOperator(operatorID, name, apiKeyHash string) (*Operator, error) {
	if operatorID == "" {
		return nil, errors.New("operatorID cannot be empty")
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if apiKeyHash == "" {
		return nil, errors.New("apiKeyHash cannot be empty")
	}

	now := time.Now().UTC()
	return &Operator{
		OperatorID: operatorID,
		Name:       name,
		APIKeyHash: apiKeyHash,
		CreatedAt:  now,
		UpdatedAt:  &now,
	}, nil
}
