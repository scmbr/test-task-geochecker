package repository

import (
	"context"

	"github.com/scmbr/test-task-geochecker/internal/domain"
	"github.com/scmbr/test-task-geochecker/internal/repository/models"
	"gorm.io/gorm"
)

type CheckRepo struct {
	db *gorm.DB
}

func NewCheckRepository(db *gorm.DB) *CheckRepo {
	return &CheckRepo{db: db}
}
func (r *CheckRepo) Create(ctx context.Context, check *domain.Check) error {
	checkModel := models.CheckDomainToModel(check)
	if err := r.db.WithContext(ctx).Create(checkModel).Error; err != nil {
		return err
	}
	return nil
}
func (r *CheckRepo) GetById(ctx context.Context, id string) (*domain.Check, error) {
	var check models.Check
	res := r.db.WithContext(ctx).Where("check_id = ?", id).First(&check)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	checkDomain, err := models.CheckModelToDomain(&check)
	if err != nil {
		return nil, err
	}
	return checkDomain, nil
}
