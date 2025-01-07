package repository

import (
	"context"
	"errors"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *GormUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormUserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *GormUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login", time.Now()).Error
}

func (r *GormUserRepository) GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error) {
	var stats models.DashboardStats
	
	// Get total horses
	var totalHorses int64
	if err := r.db.WithContext(ctx).Model(&models.Horse{}).Where("user_id = ?", userID).Count(&totalHorses).Error; err != nil {
		return nil, err
	}
	stats.TotalHorses = int(totalHorses)
	// Get pregnant horses
	var pregnantHorses int64
	if err := r.db.WithContext(ctx).Model(&models.Horse{}).Where("user_id = ? AND is_pregnant = ?", userID, true).Count(&pregnantHorses).Error; err != nil {
		return nil, err
	}
	

	return &stats, nil
}

// Custom errors
var (
	ErrUserNotFound = errors.New("user not found")
)
