package repository

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/shopspring/decimal"
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

type BreedingRepository interface {
	GetCosts(ctx context.Context, horseID uint) ([]models.BreedingCost, error)
	CreateCost(ctx context.Context, cost *models.BreedingCost) error
	UpdateCost(ctx context.Context, cost *models.BreedingCost) error
	DeleteCost(ctx context.Context, id uint) error
	GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error)
	CreateRecord(ctx context.Context, record *models.BreedingRecord) error
	GetRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error)
}

type PregnancyRepository interface {
	GetPregnancy(ctx context.Context, id uint) (*models.Pregnancy, error)
	GetByHorseID(ctx context.Context, horseID uint) (*models.Pregnancy, error)
	Create(ctx context.Context, pregnancy *models.Pregnancy) error
	Update(ctx context.Context, pregnancy *models.Pregnancy) error
	GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error)
	AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error
	GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error)
	AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error
	GetPreFoalingSigns(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error)
	AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error
	GetCurrentPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error)
}
