package repository

import (
	"context"

	"github.com/scmbr/test-task-geochecker/internal/domain"
	"gorm.io/gorm"
)

type CheckRepo struct {
	db *gorm.DB
}

func NewCheckRepository(db *gorm.DB) *CheckRepo {
	return &CheckRepo{db: db}
}
func (r *CheckRepo) Create(ctx context.Context, check *domain.Check) error {
	if err := r.db.WithContext(ctx).Create(check).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}
func (r *CheckRepo) GetById(ctx context.Context, id string) (*domain.Check, error) {
	var check domain.Check
	res := r.db.WithContext(ctx).Where("check_id = ?", id).First(&check)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &check, nil
}
