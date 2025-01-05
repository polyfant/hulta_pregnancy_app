package repository

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type HorseRepository interface {
	Create(ctx context.Context, horse *models.Horse) error
	GetByID(ctx context.Context, id uint) (*models.Horse, error)
	Update(ctx context.Context, horse *models.Horse) error
	Delete(ctx context.Context, id uint) error
	ListByUser(ctx context.Context, userID string) ([]models.Horse, error)
	GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error)
}

type ExpenseRepository interface {
	Create(ctx context.Context, expense *models.Expense) error
	GetByHorseID(ctx context.Context, horseID uint) ([]models.Expense, error)
	GetTotalExpensesByUser(ctx context.Context, userID string) (decimal.Decimal, error)
	GetExpensesByType(ctx context.Context, userID string, expenseType string) ([]models.Expense, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdateLastLogin(ctx context.Context, userID string) error
}

type RecurringExpenseRepository interface {
	Create(ctx context.Context, recurringExpense *models.RecurringExpense) error
	GetByUserID(ctx context.Context, userID string) ([]models.RecurringExpense, error)
	GetDueRecurringExpenses(ctx context.Context) ([]models.RecurringExpense, error)
}

type PregnancyRepository interface {
	Create(pregnancy *models.Pregnancy) error
	GetActivePregnancyByHorseID(horseID uint) (*models.Pregnancy, error)
	Update(pregnancy *models.Pregnancy) error
	GetPregnantHorses() ([]models.Pregnancy, error)
	GetPregnancy(pregnancyID int64) (*models.Pregnancy, error)
	GetPregnancies(userID string) ([]models.Pregnancy, error)
	AddPregnancy(pregnancy *models.Pregnancy) error
	UpdatePregnancy(pregnancy *models.Pregnancy) error
	AddPregnancyEvent(event *models.PregnancyEvent) error
}
