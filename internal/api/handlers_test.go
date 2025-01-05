package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
)

// MockDatabase implements the Database interface for testing
type MockDatabase struct {
	mock.Mock
}

// MockHorseRepository for testing
type MockHorseRepository struct {
	mock.Mock
}

func (m *MockHorseRepository) Create(ctx context.Context, horse *models.Horse) error {
	args := m.Called(ctx, horse)
	return args.Error(0)
}

func (m *MockHorseRepository) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Horse), args.Error(1)
}

func (m *MockHorseRepository) Update(ctx context.Context, horse *models.Horse) error {
	args := m.Called(ctx, horse)
	return args.Error(0)
}

func (m *MockHorseRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockHorseRepository) ListByUser(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *MockHorseRepository) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

// MockExpenseRepository for testing
type MockExpenseRepository struct {
	mock.Mock
}

func (m *MockExpenseRepository) Create(ctx context.Context, expense *models.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) GetByHorseID(ctx context.Context, horseID uint) ([]models.Expense, error) {
	args := m.Called(ctx, horseID)
	return args.Get(0).([]models.Expense), args.Error(1)
}

func (m *MockExpenseRepository) GetTotalExpensesByUser(ctx context.Context, userID string) (decimal.Decimal, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockExpenseRepository) GetExpensesByType(ctx context.Context, userID, expenseType string) ([]models.Expense, error) {
	args := m.Called(ctx, userID, expenseType)
	return args.Get(0).([]models.Expense), args.Error(1)
}

// Test AddHorse Handler
func TestAddHorseHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock repositories
	mockHorseRepo := new(MockHorseRepository)
	mockExpenseRepo := new(MockExpenseRepository)
	mockUserRepo := new(MockUserRepository)
	mockRecurringExpenseRepo := new(MockRecurringExpenseRepository)

	// Create handler with mock repositories
	handler := NewHandler(
		nil, 
		mockHorseRepo, 
		mockExpenseRepo, 
		mockRecurringExpenseRepo, 
		mockUserRepo,
	)

	// Prepare test horse
	testHorse := &models.Horse{
		Name:      "Test Horse",
		BirthDate: time.Now().AddDate(-5, 0, 0),
		UserID:    "test-user-id",
	}

	// Expect horse creation to succeed
	mockHorseRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	// Create Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Prepare request body
	jsonBody, _ := json.Marshal(testHorse)
	c.Request, _ = http.NewRequest(http.MethodPost, "/horses", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")
	
	// Set user context
	c.Set("user_id", "test-user-id")

	// Call handler
	handler.AddHorse(c)

	// Assert response
	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify mock expectations
	mockHorseRepo.AssertExpectations(t)
}

// Test ListHorses Handler
func TestListHorsesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock repositories
	mockHorseRepo := new(MockHorseRepository)
	mockExpenseRepo := new(MockExpenseRepository)
	mockUserRepo := new(MockUserRepository)
	mockRecurringExpenseRepo := new(MockRecurringExpenseRepository)

	// Create handler with mock repositories
	handler := NewHandler(
		nil, 
		mockHorseRepo, 
		mockExpenseRepo, 
		mockRecurringExpenseRepo, 
		mockUserRepo,
	)

	// Prepare test horses
	testHorses := []models.Horse{
		{
			Name:      "Horse 1",
			BirthDate: time.Now().AddDate(-5, 0, 0),
			UserID:    "test-user-id",
		},
		{
			Name:      "Horse 2",
			BirthDate: time.Now().AddDate(-3, 0, 0),
			UserID:    "test-user-id",
		},
	}

	// Expect horse list to be retrieved
	mockHorseRepo.On("ListByUser", mock.Anything, "test-user-id").Return(testHorses, nil)

	// Create Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Prepare request
	c.Request, _ = http.NewRequest(http.MethodGet, "/horses", nil)
	
	// Set user context
	c.Set("user_id", "test-user-id")

	// Call handler
	handler.ListHorses(c)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Decode response
	var horses []models.Horse
	err := json.Unmarshal(w.Body.Bytes(), &horses)
	assert.NoError(t, err)
	assert.Len(t, horses, 2)

	// Verify mock expectations
	mockHorseRepo.AssertExpectations(t)
}

// Test GetPregnantHorses Handler
func TestGetPregnantHorsesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock repositories
	mockHorseRepo := new(MockHorseRepository)
	mockExpenseRepo := new(MockExpenseRepository)
	mockUserRepo := new(MockUserRepository)
	mockRecurringExpenseRepo := new(MockRecurringExpenseRepository)

	// Create handler with mock repositories
	handler := NewHandler(
		nil, 
		mockHorseRepo, 
		mockExpenseRepo, 
		mockRecurringExpenseRepo, 
		mockUserRepo,
	)

	// Prepare pregnant horses
	pregnantHorses := []models.Horse{
		{
			Name:        "Pregnant Horse 1",
			IsPregnant:  true,
			UserID:      "test-user-id",
		},
	}

	// Expect pregnant horses to be retrieved
	mockHorseRepo.On("GetPregnantHorses", mock.Anything, "test-user-id").Return(pregnantHorses, nil)

	// Create Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Prepare request
	c.Request, _ = http.NewRequest(http.MethodGet, "/horses/pregnant", nil)
	
	// Set user context
	c.Set("user_id", "test-user-id")

	// Call handler
	handler.GetPregnantHorses(c)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Decode response
	var horses []models.Horse
	err := json.Unmarshal(w.Body.Bytes(), &horses)
	assert.NoError(t, err)
	assert.Len(t, horses, 1)
	assert.True(t, horses[0].IsPregnant)

	// Verify mock expectations
	mockHorseRepo.AssertExpectations(t)
}

// MockUserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockRecurringExpenseRepository for testing
type MockRecurringExpenseRepository struct {
	mock.Mock
}

func (m *MockRecurringExpenseRepository) Create(ctx context.Context, recurringExpense *models.RecurringExpense) error {
	args := m.Called(ctx, recurringExpense)
	return args.Error(0)
}

func (m *MockRecurringExpenseRepository) GetByUserID(ctx context.Context, userID string) ([]models.RecurringExpense, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.RecurringExpense), args.Error(1)
}

func (m *MockRecurringExpenseRepository) GetDueRecurringExpenses(ctx context.Context) ([]models.RecurringExpense, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.RecurringExpense), args.Error(1)
}
