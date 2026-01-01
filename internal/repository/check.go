package repository

import (
	"context"

	"github.com/scmbr/test-task-geochecker/internal/domain/models"
	"gorm.io/gorm"
)

type CheckRepo struct {
	db *gorm.DB
}

func NewCheckRepository(db *gorm.DB) *CheckRepo {
	return &CheckRepo{db: db}
}
func (r *CheckRepo) Create(ctx context.Context, check models.Check) error {
	if err := r.db.WithContext(ctx).Create(check).Error; err != nil {
		return err
	}
	return nil
}
func (r *CheckRepo) GetById(ctx context.Context, id string) (*models.Check, error) {
	var check models.Check
	res := r.db.WithContext(ctx).Where("check_id = ?", id).First(&check)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &check, nil
}
