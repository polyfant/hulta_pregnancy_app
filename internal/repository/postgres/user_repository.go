package postgres

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).
		Update("last_login", time.Now()).Error
}

func (r *UserRepository) GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}
	
	// Get total horses
	if err := r.db.WithContext(ctx).Model(&models.Horse{}).Where("user_id = ?", userID).Count(&stats.TotalHorses).Error; err != nil {
		return nil, err
	}

	// Get pregnant mares
	if err := r.db.WithContext(ctx).Model(&models.Horse{}).Where("user_id = ? AND is_pregnant = true", userID).Count(&stats.PregnantMares).Error; err != nil {
		return nil, err
	}

	return stats, nil
} 