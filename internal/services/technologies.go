package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/uaaaalim/alimmun-restapi/internal/models"
	"github.com/uaaaalim/alimmun-restapi/internal/repository"
	"gorm.io/gorm"
)

var (
	TechnologyInvalid       = errors.New("name must contain 1 to 100 characters; slug must contain 1 to 100 lowercase letters, digits or single hyphens between words")
	TechnologyAlreadyExists = errors.New("technology name or slug already exists")
	TechnologyNotFound      = errors.New("technology not found")
	TechnologyInvalidImage  = errors.New("image must be an absolute HTTP or HTTPS URL of at most 2048 characters")
	technologySlugPattern   = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
)

type TechnologyService struct {
	TechnologyRepo *repository.TechnologyRepository
}

func NewTechnologyService(repo *repository.TechnologyRepository) *TechnologyService {
	return &TechnologyService{TechnologyRepo: repo}
}

func (s *TechnologyService) List(ctx context.Context) ([]models.ProjectTechnology, error) {
	technologies, err := s.TechnologyRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list technologies: %w", err)
	}
	return technologies, nil
}

func (s *TechnologyService) Create(ctx context.Context, name, slug string, image *string) (*models.ProjectTechnology, error) {
	name = strings.TrimSpace(name)
	slug = strings.TrimSpace(slug)
	if name == "" || utf8.RuneCountInString(name) > 100 || len(slug) > 100 || !technologySlugPattern.MatchString(slug) {
		return nil, TechnologyInvalid
	}

	if image != nil {
		value := strings.TrimSpace(*image)
		parsed, err := url.Parse(value)
		if err != nil || len(value) > 2048 || parsed.Hostname() == "" ||
			(parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.User != nil {
			return nil, TechnologyInvalidImage
		}
		image = &value
	}

	technology := &models.ProjectTechnology{Name: name, Slug: slug, Image: image}
	if err := s.TechnologyRepo.Create(ctx, technology); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, TechnologyAlreadyExists
		}
		return nil, fmt.Errorf("create technology: %w", err)
	}
	return technology, nil
}

func (s *TechnologyService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return TechnologyInvalid
	}
	if err := s.TechnologyRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return TechnologyNotFound
		}
		return fmt.Errorf("delete technology: %w", err)
	}
	return nil
}
