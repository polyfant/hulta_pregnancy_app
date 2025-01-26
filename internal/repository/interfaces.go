package repository

import (
	"context"
	"time"

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
	GetOffspring(ctx context.Context, horseID uint) ([]models.Horse, error)
	GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error)
	GetPregnant(ctx context.Context, userID string) ([]models.Horse, error)
}

type ExpenseRepository interface {
	Create(ctx context.Context, expense *models.Expense) error
	Update(ctx context.Context, expense *models.Expense) error
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
	GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error)
}

type RecurringExpenseRepository interface {
	Create(ctx context.Context, recurringExpense *models.RecurringExpense) error
	GetByUserID(ctx context.Context, userID string) ([]models.RecurringExpense, error)
	GetDueRecurringExpenses(ctx context.Context) ([]models.RecurringExpense, error)
}

type BreedingRepository interface {
	GetCosts(ctx context.Context, horseID uint) ([]models.BreedingCost, error)
	Create(ctx context.Context, cost *models.BreedingCost) error
	GetRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error)
	CreateRecord(ctx context.Context, record *models.BreedingRecord) error
	UpdateRecord(ctx context.Context, record *models.BreedingRecord) error
	DeleteRecord(ctx context.Context, id uint) error
}

type PregnancyRepository interface {
	GetPregnancy(ctx context.Context, id uint) (*models.Pregnancy, error)
	GetByHorseID(ctx context.Context, horseID uint) (*models.Pregnancy, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Pregnancy, error)
	Create(ctx context.Context, pregnancy *models.Pregnancy) error
	Update(ctx context.Context, pregnancy *models.Pregnancy) error
	GetEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error)
	AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error
	GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error)
	GetPreFoalingChecklistItem(ctx context.Context, itemID uint) (*models.PreFoalingChecklistItem, error)
	AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error
	DeletePreFoalingChecklistItem(ctx context.Context, itemID uint) error
	InitializePreFoalingChecklist(ctx context.Context, horseID uint) error
	GetPreFoalingSigns(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error)
	AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error
	GetCurrentPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error)
	UpdatePregnancyStatus(ctx context.Context, horseID uint, isPregnant bool, conceptionDate *time.Time) error
	UpdatePreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error
	GetActive(ctx context.Context, userID string) ([]models.Pregnancy, error)
	AddPreFoaling(ctx context.Context, sign *models.PreFoalingSign) error
	GetPreFoaling(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error)
}

type HealthRepository interface {
	CreateRecord(ctx context.Context, record *models.HealthRecord) error
	GetRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error)
	UpdateRecord(ctx context.Context, record *models.HealthRecord) error
	DeleteRecord(ctx context.Context, id uint) error
}
