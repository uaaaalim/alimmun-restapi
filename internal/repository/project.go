package repository

import (
	"context"
	"fmt"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	DB *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{DB: db}
}

func (repo *ProjectRepository) GetTechnologiesByIDs(ctx context.Context, ids []int64) ([]models.ProjectTechnology, error) {
	technologies := make([]models.ProjectTechnology, 0)
	if len(ids) == 0 {
		return technologies, nil
	}
	err := repo.DB.WithContext(ctx).Where("id IN ?", ids).Order("id ASC").Find(&technologies).Error
	return technologies, err
}

func (repo *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	// Save links to existing technologies in GORM's default transaction.
	return repo.DB.WithContext(ctx).Omit("Technologies.*").Create(project).Error
}

func (repo *ProjectRepository) Delete(ctx context.Context, id int64) error {
	result := repo.DB.WithContext(ctx).Delete(&models.Project{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (repo *ProjectRepository) GetByPage(ctx context.Context, page, pageSize int) ([]models.Project, error) {
	if page < 1 || pageSize < 1 || pageSize > 100 {
		return nil, fmt.Errorf("invalid pagination parameters")
	}

	projects := make([]models.Project, 0)
	offset := (page - 1) * pageSize

	err := repo.DB.WithContext(ctx).
		Preload("Technologies", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&projects).Error
	if err != nil {
		return nil, fmt.Errorf("get projects by page: %w", err)
	}

	return projects, nil
}
