package repository

import (
	"context"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type GormExpenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) *GormExpenseRepository {
	return &GormExpenseRepository{db: db}
}

func (r *GormExpenseRepository) Create(ctx context.Context, expense *models.Expense) error {
	return r.db.WithContext(ctx).Create(expense).Error
}

func (r *GormExpenseRepository) GetByHorseID(ctx context.Context, horseID uint) ([]models.Expense, error) {
	var expenses []models.Expense
	result := r.db.WithContext(ctx).
		Where("horse_id = ?", horseID).
		Order("date DESC").
		Find(&expenses)
	return expenses, result.Error
}

func (r *GormExpenseRepository) GetTotalExpensesByUser(ctx context.Context, userID string) (decimal.Decimal, error) {
	var total decimal.Decimal
	result := r.db.WithContext(ctx).
		Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ?", userID).
		Scan(&total)
	
	return total, result.Error
}

func (r *GormExpenseRepository) GetExpensesByType(ctx context.Context, userID string, expenseType string) ([]models.Expense, error) {
	var expenses []models.Expense
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, expenseType).
		Order("date DESC").
		Find(&expenses)
	return expenses, result.Error
}
