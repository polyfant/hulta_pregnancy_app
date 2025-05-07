package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// --- Mock ExpenseRepository ---
type mockExpenseRepo struct{ mock.Mock }

func (m *mockExpenseRepo) Create(ctx context.Context, expense *models.Expense) (*models.Expense, error) {
	args := m.Called(ctx, expense)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expense), args.Error(1)
}
func (m *mockExpenseRepo) Update(ctx context.Context, expense *models.Expense) (*models.Expense, error) {
	args := m.Called(ctx, expense)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expense), args.Error(1)
}
func (m *mockExpenseRepo) GetByID(ctx context.Context, expenseID uint) (*models.Expense, error) {
	args := m.Called(ctx, expenseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expense), args.Error(1)
}
func (m *mockExpenseRepo) Delete(ctx context.Context, expenseID uint) error {
	args := m.Called(ctx, expenseID)
	return args.Error(0)
}
func (m *mockExpenseRepo) GetByHorseID(ctx context.Context, horseID uint) ([]models.Expense, error) {
	args := m.Called(ctx, horseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Expense), args.Error(1)
}
func (m *mockExpenseRepo) GetExpensesByType(ctx context.Context, userID, expenseType string) ([]models.Expense, error) {
	args := m.Called(ctx, userID, expenseType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Expense), args.Error(1)
}
func (m *mockExpenseRepo) GetTotalExpensesByUser(ctx context.Context, userID string) (decimal.Decimal, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

var _ repository.ExpenseRepository = (*mockExpenseRepo)(nil)

// --- Mock RecurringExpenseRepository ---
type mockRecurringRepo struct{ mock.Mock }

func (m *mockRecurringRepo) Create(ctx context.Context, recurring *models.RecurringExpense) (*models.RecurringExpense, error) {
	args := m.Called(ctx, recurring)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RecurringExpense), args.Error(1)
}
func (m *mockRecurringRepo) Update(ctx context.Context, recurring *models.RecurringExpense) (*models.RecurringExpense, error) {
	args := m.Called(ctx, recurring)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RecurringExpense), args.Error(1)
}
func (m *mockRecurringRepo) GetByID(ctx context.Context, id uint) (*models.RecurringExpense, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RecurringExpense), args.Error(1)
}
func (m *mockRecurringRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockRecurringRepo) GetDueRecurringExpenses(ctx context.Context) ([]models.RecurringExpense, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.RecurringExpense), args.Error(1)
}
func (m *mockRecurringRepo) GetByUserID(ctx context.Context, userID string) ([]models.RecurringExpense, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.RecurringExpense), args.Error(1)
}

var _ repository.RecurringExpenseRepository = (*mockRecurringRepo)(nil)

// --- Mock HorseRepository ---
type mockHorseRepo struct{ mock.Mock }

func (m *mockHorseRepo) Create(ctx context.Context, horse *models.Horse) error {
	args := m.Called(ctx, horse); return args.Error(0)
}
func (m *mockHorseRepo) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Horse), args.Error(1)
}
func (m *mockHorseRepo) Update(ctx context.Context, horse *models.Horse) error {
	args := m.Called(ctx, horse); return args.Error(0)
}
func (m *mockHorseRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id); return args.Error(0)
}
func (m *mockHorseRepo) ListByUser(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID); if args.Get(0) == nil { return nil, args.Error(1) }; return args.Get(0).([]models.Horse), args.Error(1)
}
func (m *mockHorseRepo) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID); if args.Get(0) == nil { return nil, args.Error(1) }; return args.Get(0).([]models.Horse), args.Error(1)
}
func (m *mockHorseRepo) GetOffspring(ctx context.Context, horseID uint) ([]models.Horse, error) {
	args := m.Called(ctx, horseID); if args.Get(0) == nil { return nil, args.Error(1) }; return args.Get(0).([]models.Horse), args.Error(1)
}
func (m *mockHorseRepo) GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error) {
	args := m.Called(ctx, horseID); if args.Get(0) == nil { return nil, args.Error(1) }; return args.Get(0).(*models.FamilyTree), args.Error(1)
}
func (m *mockHorseRepo) GetPregnant(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID); if args.Get(0) == nil { return nil, args.Error(1) }; return args.Get(0).([]models.Horse), args.Error(1)
}

var _ repository.HorseRepository = (*mockHorseRepo)(nil)

func TestCreateRecurringExpense_Validation(t *testing.T) {
	ctx := context.Background()

	// Validation is now primarily in handlers and model tags.
	// This test should focus on service-level behavior, e.g., repository interaction errors.
	testCases := []struct {
		name          string
		input         *models.RecurringExpense
		setupMock     func(mrr *mockRecurringRepo, mer *mockExpenseRepo, mhr *mockHorseRepo)
		expectedErrorContains string
	}{
		{
			name: "Repository Create Fails",
			input: &models.RecurringExpense{
				UserID:      "user-1",
				Amount:      decimal.NewFromInt(100),
				ExpenseType: models.ExpenseTypeFeed,
				Frequency:   models.FrequencyMonthly,
				StartDate:   time.Now().Add(-24 * time.Hour), // Valid past date
			},
			setupMock: func(mrr *mockRecurringRepo, mer *mockExpenseRepo, mhr *mockHorseRepo) {
				// Service calculates NextDueDate before calling Create
				mrr.On("Create", ctx, mock.AnythingOfType("*models.RecurringExpense")).Return(nil, errors.New("repo create error")).Once()
			},
			expectedErrorContains: "repo create error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create new mocks for each sub-test to ensure isolation
			currentMockExpRepo := &mockExpenseRepo{}
			currentMockRecRepo := &mockRecurringRepo{}
			currentMockHrsRepo := &mockHorseRepo{}
			currentSvc := NewExpenseService(currentMockExpRepo, currentMockRecRepo, currentMockHrsRepo)

			tc.setupMock(currentMockRecRepo, currentMockExpRepo, currentMockHrsRepo)

			_, err := currentSvc.CreateRecurringExpense(ctx, tc.input)

			if tc.expectedErrorContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrorContains)
			} else {
				assert.NoError(t, err)
			}
			currentMockRecRepo.AssertExpectations(t)
			currentMockExpRepo.AssertExpectations(t)
			currentMockHrsRepo.AssertExpectations(t)
		})
	}
}

func TestCreateRecurringExpense_Success(t *testing.T) {
	mockRecRepo := &mockRecurringRepo{}
	mockExpRepo := &mockExpenseRepo{}
	mockHrsRepo := &mockHorseRepo{}
	svc := NewExpenseService(mockExpRepo, mockRecRepo, mockHrsRepo)
	ctx := context.Background()

	startDate := time.Now().Add(-48 * time.Hour)
	re := &models.RecurringExpense{
		UserID:      "user-1",
		Amount:      decimal.NewFromInt(100),
		ExpenseType: models.ExpenseTypeFeed,
		Frequency:   models.FrequencyMonthly,
		StartDate:   startDate,
	}
	expectedNextDueDate := calculateNextDueDate(*re)

	expectedReturn := &models.RecurringExpense{
		Model:       gorm.Model{ID: 1},
		UserID:      re.UserID,
		Amount:      re.Amount,
		ExpenseType: re.ExpenseType,
		Frequency:   re.Frequency,
		StartDate:   re.StartDate,
		NextDueDate: expectedNextDueDate,
	}

	mockRecRepo.On("Create", ctx, mock.MatchedBy(func(arg *models.RecurringExpense) bool {
		return arg.UserID == re.UserID && 
		       arg.Amount.Equal(re.Amount) &&
		       arg.ExpenseType == re.ExpenseType &&
		       arg.Frequency == re.Frequency &&
		       arg.StartDate.Equal(re.StartDate) &&
		       arg.NextDueDate.Equal(expectedNextDueDate)
	})).Return(expectedReturn, nil).Once()

	createdRe, err := svc.CreateRecurringExpense(ctx, re)
	assert.NoError(t, err)
	assert.NotNil(t, createdRe)
	assert.Equal(t, expectedReturn.ID, createdRe.ID)
	assert.True(t, createdRe.NextDueDate.Equal(expectedNextDueDate), "NextDueDate should be correctly calculated and returned")
	mockRecRepo.AssertExpectations(t)
}

func TestProcessDueRecurringExpenses(t *testing.T) {
	mockRecRepo := &mockRecurringRepo{}
	mockExpRepo := &mockExpenseRepo{}
	mockHrsRepo := &mockHorseRepo{}
	svc := NewExpenseService(mockExpRepo, mockRecRepo, mockHrsRepo)
	ctx := context.Background()

	now := time.Now()
	dueItem := models.RecurringExpense{
		Model:       gorm.Model{ID: 1},
		UserID:      "user-1",
		Amount:      decimal.NewFromInt(50),
		ExpenseType: models.ExpenseTypeFeed,
		Frequency:   models.FrequencyMonthly,
		StartDate:   now.AddDate(0, -2, 0),
		NextDueDate: now.Truncate(24 * time.Hour),
	}

	mockRecRepo.On("GetDueRecurringExpenses", ctx).Return([]models.RecurringExpense{dueItem}, nil).Once()

	expenseMatcher := mock.MatchedBy(func(exp *models.Expense) bool {
		return exp.UserID == dueItem.UserID &&
			exp.Amount.Equal(dueItem.Amount) &&
			exp.Date.Equal(dueItem.NextDueDate)
	})
	mockExpRepo.On("Create", ctx, expenseMatcher).Return(&models.Expense{ID: 101, UserID: dueItem.UserID}, nil).Once()

	expectedUpdatedNextDueDate := calculateNextDueDate(dueItem)
	recurringMatcher := mock.MatchedBy(func(rec *models.RecurringExpense) bool {
		return rec.ID == dueItem.ID && rec.NextDueDate.Equal(expectedUpdatedNextDueDate)
	})
	mockRecRepo.On("Update", ctx, recurringMatcher).Return(&models.RecurringExpense{Model: gorm.Model{ID: dueItem.ID}, NextDueDate: expectedUpdatedNextDueDate}, nil).Once()

	err := svc.ProcessDueRecurringExpenses(ctx)
	assert.NoError(t, err)

	mockRecRepo.AssertExpectations(t)
	mockExpRepo.AssertExpectations(t)
}

func TestProcessDueRecurringExpenses_Error(t *testing.T) {
	mockRecRepo := &mockRecurringRepo{}
	mockExpRepo := &mockExpenseRepo{}
	mockHrsRepo := &mockHorseRepo{}
	svc := NewExpenseService(mockExpRepo, mockRecRepo, mockHrsRepo)
	ctx := context.Background()

	mockRecRepo.On("GetDueRecurringExpenses", ctx).Return(nil, errors.New("db error")).Once()

	err := svc.ProcessDueRecurringExpenses(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	mockRecRepo.AssertExpectations(t)
}

// Further tests for other ExpenseService methods should be added here. 