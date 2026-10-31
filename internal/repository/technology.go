package repository

import (
	"context"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"gorm.io/gorm"
)

type TechnologyRepository struct {
	DB *gorm.DB
}

func NewTechnologyRepository(db *gorm.DB) *TechnologyRepository {
	return &TechnologyRepository{DB: db}
}

func (repo *TechnologyRepository) List(ctx context.Context) ([]models.ProjectTechnology, error) {
	technologies := make([]models.ProjectTechnology, 0)
	err := repo.DB.WithContext(ctx).Order("id ASC").Find(&technologies).Error
	return technologies, err
}

func (repo *TechnologyRepository) Create(ctx context.Context, technology *models.ProjectTechnology) error {
	return repo.DB.WithContext(ctx).Create(technology).Error
}

func (repo *TechnologyRepository) Delete(ctx context.Context, id int64) error {
	result := repo.DB.WithContext(ctx).Delete(&models.ProjectTechnology{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
