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

func (s *HorseService) GetHorse(ctx context.Context, id uint) (*models.Horse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *HorseService) ListHorsesByUser(ctx context.Context, userID string) ([]models.Horse, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *HorseService) CreateHorse(horse *models.Horse) error {
	return s.repo.Create(context.Background(), horse)
}

func (s *HorseService) UpdateHorse(horse *models.Horse) error {
	return s.repo.Update(context.Background(), horse)
}

func (s *HorseService) DeleteHorse(id uint) error {
	return s.repo.Delete(context.Background(), id)
}

func (s *HorseService) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	return s.repo.GetPregnantHorses(ctx, userID)
}

func (s *HorseService) GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error) {
	return s.repo.GetFamilyTree(ctx, horseID)
}
