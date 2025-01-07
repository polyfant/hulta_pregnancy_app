package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type GormRecurringExpenseRepository struct {
	db *gorm.DB
}

func NewRecurringExpenseRepository(db *gorm.DB) *GormRecurringExpenseRepository {
	return &GormRecurringExpenseRepository{db: db}
}

func (r *GormRecurringExpenseRepository) Create(ctx context.Context, recurringExpense *models.RecurringExpense) error {
	return r.db.WithContext(ctx).Create(recurringExpense).Error
}

func (r *GormRecurringExpenseRepository) GetByUserID(ctx context.Context, userID string) ([]models.RecurringExpense, error) {
	var recurringExpenses []models.RecurringExpense
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("next_due_date ASC").
		Find(&recurringExpenses)
	return recurringExpenses, result.Error
}

func (r *GormRecurringExpenseRepository) GetDueRecurringExpenses(ctx context.Context) ([]models.RecurringExpense, error) {
	var recurringExpenses []models.RecurringExpense
	result := r.db.WithContext(ctx).
		Where("next_due_date <= ?", time.Now()).
		Order("next_due_date ASC").
		Find(&recurringExpenses)
	return recurringExpenses, result.Error
}
