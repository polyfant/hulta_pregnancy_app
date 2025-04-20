package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockExpenseRepo struct{ mock.Mock }
func (m *mockExpenseRepo) Create(ctx context.Context, expense *models.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}
func (m *mockExpenseRepo) Update(ctx context.Context, expense *models.Expense) error { return nil }
func (m *mockExpenseRepo) GetByHorseID(ctx context.Context, horseID uint) ([]models.Expense, error) { return nil, nil }
func (m *mockExpenseRepo) GetExpensesByType(ctx context.Context, userID, expenseType string) ([]models.Expense, error) { return nil, nil }
func (m *mockExpenseRepo) GetTotalExpensesByUser(ctx context.Context, userID string) (decimal.Decimal, error) { return decimal.NewFromInt(0), nil }

// RecurringExpenseRepository mock
type mockRecurringRepo struct{ mock.Mock }
func (m *mockRecurringRepo) Create(ctx context.Context, recurring *models.RecurringExpense) error {
	args := m.Called(ctx, recurring)
	return args.Error(0)
}
func (m *mockRecurringRepo) GetDueRecurringExpenses(ctx context.Context) ([]models.RecurringExpense, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.RecurringExpense), args.Error(1)
}
func (m *mockRecurringRepo) GetByUserID(ctx context.Context, userID string) ([]models.RecurringExpense, error) {
	return nil, nil
}

func TestCreateRecurringExpense_Validation(t *testing.T) {
	svc := &ExpenseService{
		expenseRepo:         &mockExpenseRepo{},
		recurringExpenseRepo: &mockRecurringRepo{},
	}
	ctx := context.Background()

	re := &models.RecurringExpense{
		UserID:      "user-1",
		Amount:      -10,
		ExpenseType: "",
		Frequency:   "",
		StartDate:   time.Now(),
	}
	err := svc.CreateRecurringExpense(ctx, re)
	assert.ErrorIs(t, err, models.ErrInvalidAmount)

	re.Amount = 100
	err = svc.CreateRecurringExpense(ctx, re)
	assert.ErrorIs(t, err, models.ErrInvalidExpenseType)

	re.ExpenseType = models.ExpenseTypeFeed
	err = svc.CreateRecurringExpense(ctx, re)
	assert.ErrorIs(t, err, models.ErrInvalidFrequency)
}

func TestCreateRecurringExpense_Success(t *testing.T) {
	repo := &mockRecurringRepo{}
	svc := &ExpenseService{
		expenseRepo:         &mockExpenseRepo{},
		recurringExpenseRepo: repo,
	}
	ctx := context.Background()

	re := &models.RecurringExpense{
		UserID:      "user-1",
		Amount:      100,
		ExpenseType: models.ExpenseTypeFeed,
		Frequency:   models.FrequencyMonthly,
		StartDate:   time.Now(),
	}

	repo.On("Create", ctx, re).Return(nil).Once()

	err := svc.CreateRecurringExpense(ctx, re)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProcessDueRecurringExpenses(t *testing.T) {
	repo := &mockRecurringRepo{}
	expRepo := &mockExpenseRepo{}
	svc := &ExpenseService{
		expenseRepo:         expRepo,
		recurringExpenseRepo: repo,
	}
	ctx := context.Background()

	due := models.RecurringExpense{
		UserID:      "user-1",
		Amount:      50,
		ExpenseType: models.ExpenseTypeFeed,
		Frequency:   models.FrequencyMonthly,
		StartDate:   time.Now().AddDate(0, -1, 0),
		NextDueDate: time.Now(),
	}

	repo.On("GetDueRecurringExpenses", ctx).Return([]models.RecurringExpense{due}, nil).Once()
	expRepo.On("Create", ctx, mock.AnythingOfType("*models.Expense")).Return(nil).Once()

	err := svc.ProcessDueRecurringExpenses(ctx)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	expRepo.AssertExpectations(t)
}

func TestProcessDueRecurringExpenses_Error(t *testing.T) {
	repo := &mockRecurringRepo{}
	svc := &ExpenseService{
		expenseRepo:         &mockExpenseRepo{},
		recurringExpenseRepo: repo,
	}
	ctx := context.Background()

	repo.On("GetDueRecurringExpenses", ctx).Return([]models.RecurringExpense{}, errors.New("db error")).Once()

	err := svc.ProcessDueRecurringExpenses(ctx)
	assert.Error(t, err)
} 