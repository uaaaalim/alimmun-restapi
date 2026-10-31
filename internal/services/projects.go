package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"github.com/uaaaalim/alimmun-restapi/internal/repository"
	"gorm.io/gorm"
)

var (
	ProjectTechnologiesInvalid = errors.New("one or more technologies do not exist or IDs are repeated")
	ProjectNotFound            = errors.New("project not found")
)

type ProjectService struct {
	ProjectRepo *repository.ProjectRepository
}

func NewProjectService(projectRepo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{ProjectRepo: projectRepo}
}

func (s *ProjectService) GetByPage(ctx context.Context, page, pageSize int) ([]models.Project, error) {
	return s.ProjectRepo.GetByPage(ctx, page, pageSize)
}

func (s *ProjectService) Create(ctx context.Context, project *models.Project, technologyIDs []int64) error {
	technologies, err := s.ProjectRepo.GetTechnologiesByIDs(ctx, technologyIDs)
	if err != nil {
		return fmt.Errorf("get project technologies: %w", err)
	}
	if len(technologies) != len(technologyIDs) {
		return ProjectTechnologiesInvalid
	}
	project.Technologies = technologies
	if err := s.ProjectRepo.Create(ctx, project); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return ProjectTechnologiesInvalid
		}
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

func (s *ProjectService) Delete(ctx context.Context, id int64) error {
	if err := s.ProjectRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ProjectNotFound
		}
		return fmt.Errorf("delete project: %w", err)
	}
	return nil
}
