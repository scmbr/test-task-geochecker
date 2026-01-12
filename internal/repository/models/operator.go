package models

import "time"

type Operator struct {
	OperatorID string     `gorm:"primaryKey;column:operator_id"`
	APIKeyHash string     `gorm:"column:api_key_hash;unique;not null"`
	Name       string     `gorm:"column:name;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	RevokedAt  *time.Time `gorm:"column:revoked_at;default:null"`
	UpdatedAt  *time.Time `gorm:"column:updated_at;default:null"`
}
