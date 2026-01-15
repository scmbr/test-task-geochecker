package repository

import (
	"context"
	"time"

	"github.com/scmbr/test-task-geochecker/internal/domain"
	"github.com/scmbr/test-task-geochecker/internal/repository/models"
	"gorm.io/gorm"
)

type OperatorRepo struct {
	db *gorm.DB
}

func NewOperatorRepository(db *gorm.DB) *OperatorRepo {
	return &OperatorRepo{db: db}
}
func (r *OperatorRepo) GetActiveByAPIKeyHash(ctx context.Context, apiKeyHash string) (*domain.Operator, error) {
	var operator domain.Operator

	res := r.db.WithContext(ctx).
		Where("api_key_hash = ?", apiKeyHash).
		Where("revoked_at IS NULL").
		First(&operator)

	if res.Error == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return &operator, nil
}
func (r *OperatorRepo) Create(ctx context.Context, operator *domain.Operator) error {
	if err := r.db.WithContext(ctx).Create(&operator).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}
func (r *OperatorRepo) Revoke(ctx context.Context, operatorID string) error {
	res := r.db.WithContext(ctx).
		Model(&models.Operator{}).
		Where("operator_id = ?", operatorID).
		Where("revoked_at IS NULL").
		Update("revoked_at", time.Now())

	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	if res.Error != nil {
		return res.Error
	}

	return nil
}
