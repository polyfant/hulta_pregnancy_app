package postgres

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"gorm.io/gorm"
)

func NewRecurringExpenseRepository(db *gorm.DB) repository.RecurringExpenseRepository {
	return &recurringExpenseRepository{db: db}
}

type recurringExpenseRepository struct {
	db *gorm.DB
}

func (r *recurringExpenseRepository) Create(ctx context.Context, expense *models.RecurringExpense) error {
	return r.db.Create(expense).Error
}

func (r *recurringExpenseRepository) GetByUserID(ctx context.Context, userID string) ([]models.RecurringExpense, error) {
	var expenses []models.RecurringExpense
	err := r.db.Where("user_id = ?", userID).Find(&expenses).Error
	return expenses, err
}

func (r *recurringExpenseRepository) GetDueRecurringExpenses(ctx context.Context) ([]models.RecurringExpense, error) {
	var expenses []models.RecurringExpense
	err := r.db.Where("next_due_date <= ?", time.Now()).Find(&expenses).Error
	return expenses, err
} 