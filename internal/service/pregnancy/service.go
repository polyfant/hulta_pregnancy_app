package pregnancy

import (
	"context"
	"github.com/polyfant/hulta_pregnancy_app/internal/middleware"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetPregnancies(ctx context.Context) ([]models.Pregnancy, error) {
	userID := ctx.Value(middleware.UserIDKey).(string)
	var pregnancies []models.Pregnancy
	err := s.db.Where("user_id = ?", userID).Find(&pregnancies).Error
	return pregnancies, err
}

// Add other service methods as needed
