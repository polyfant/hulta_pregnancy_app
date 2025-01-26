package service

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type HorseServiceImpl struct {
	repo repository.HorseRepository
}

func NewHorseService(repo repository.HorseRepository) *HorseServiceImpl {
	return &HorseServiceImpl{repo: repo}
}

func (s *HorseServiceImpl) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *HorseServiceImpl) ListByUserID(ctx context.Context, userID string) ([]models.Horse, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *HorseServiceImpl) Create(ctx context.Context, horse *models.Horse) error {
	return s.repo.Create(context.Background(), horse)
}

func (s *HorseServiceImpl) Update(ctx context.Context, horse *models.Horse) error {
	return s.repo.Update(context.Background(), horse)
}

func (s *HorseServiceImpl) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(context.Background(), id)
}

func (s *HorseServiceImpl) GetPregnant(ctx context.Context, userID string) ([]models.Horse, error) {
	return s.repo.GetPregnant(ctx, userID)
}

func (s *HorseServiceImpl) GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error) {
	return s.repo.GetFamilyTree(ctx, horseID)
}

func (s *HorseServiceImpl) GetOffspring(ctx context.Context, horseID uint) ([]models.Horse, error) {
	return s.repo.GetOffspring(ctx, horseID)
}

func (s *HorseServiceImpl) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	return s.repo.GetPregnantHorses(ctx, userID)
}
