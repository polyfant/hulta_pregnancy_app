package postgres

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func NewExpenseRepository(db *gorm.DB) repository.ExpenseRepository {
	return &expenseRepository{db: db}
}

type expenseRepository struct {
	db *gorm.DB
}

func (r *expenseRepository) Create(ctx context.Context, expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *expenseRepository) GetByHorseID(ctx context.Context, horseID uint) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Where("horse_id = ?", horseID).Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) GetExpensesByType(ctx context.Context, userID string, expenseType string) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Where("user_id = ? AND type = ?", userID, expenseType).Find(&expenses).Error
	return expenses, err
}

func (r *expenseRepository) GetTotalExpensesByUser(ctx context.Context, userID string) (decimal.Decimal, error) {
	var total decimal.Decimal
	err := r.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ?", userID).
		Scan(&total).Error
	return total, err
} 