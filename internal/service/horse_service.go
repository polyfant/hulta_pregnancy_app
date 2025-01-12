package service

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type HorseService struct {
	repo repository.HorseRepository
}

func NewHorseService(repo repository.HorseRepository) *HorseService {
	return &HorseService{repo: repo}
}

func (s *HorseService) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *HorseService) ListByUserID(ctx context.Context, userID string) ([]models.Horse, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *HorseService) Create(ctx context.Context, horse *models.Horse) error {
	return s.repo.Create(context.Background(), horse)
}

func (s *HorseService) Update(ctx context.Context, horse *models.Horse) error {
	return s.repo.Update(context.Background(), horse)
}

func (s *HorseService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(context.Background(), id)
}

func (s *HorseService) GetPregnant(ctx context.Context, userID string) ([]models.Horse, error) {
	return s.repo.GetPregnant(ctx, userID)
}

func (s *HorseService) GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error) {
	return s.repo.GetFamilyTree(ctx, horseID)
}
