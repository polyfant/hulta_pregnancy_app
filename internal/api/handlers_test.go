package api

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/mocks"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository/postgres"
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

// MockHorseService for testing
type MockHorseService struct {
	mock.Mock
}

func (m *MockHorseService) CreateHorse(horse *models.Horse) error {
	args := m.Called(horse)
	return args.Error(0)
}

func (m *MockHorseService) GetHorse(id uint) (*models.Horse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Horse), args.Error(1)
}

func (m *MockHorseService) ListHorsesByUser(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *MockHorseService) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

// MockUserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) UpdateLastLogin(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

// MockHealthService for testing
type MockHealthService struct {
	mock.Mock
}

func (m *MockHealthService) CreateHealthRecord(ctx context.Context, record *models.HealthRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockHealthService) UpdateHealthRecord(ctx context.Context, record *models.HealthRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockHealthService) GetHealthRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error) {
	args := m.Called(ctx, horseID)
	return args.Get(0).([]models.HealthRecord), args.Error(1)
}

func (m *MockHealthService) DeleteHealthRecord(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockHealthService) AddHealthRecord(ctx context.Context, record *models.HealthRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

type mockPregnancyRepo struct {
	repository.PregnancyRepository
	getEventsFn func(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error)
}

// Implement all required methods
func (m *mockPregnancyRepo) GetPregnancy(ctx context.Context, id uint) (*models.Pregnancy, error) {
	return nil, nil
}

func (m *mockPregnancyRepo) GetByHorseID(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	return nil, nil
}

func (m *mockPregnancyRepo) GetByUserID(ctx context.Context, userID string) ([]models.Pregnancy, error) {
	return nil, nil
}

func (m *mockPregnancyRepo) Create(ctx context.Context, pregnancy *models.Pregnancy) error {
	return nil
}

func (m *mockPregnancyRepo) Update(ctx context.Context, pregnancy *models.Pregnancy) error {
	return nil
}

func (m *mockPregnancyRepo) GetEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
	if m.getEventsFn != nil {
		return m.getEventsFn(ctx, horseID)
	}
	return []models.PregnancyEvent{}, nil
}

func (m *mockPregnancyRepo) AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error {
	return nil
}

func (m *mockPregnancyRepo) GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error) {
	return nil, nil
}

func (m *mockPregnancyRepo) GetPreFoalingChecklistItem(ctx context.Context, itemID uint) (*models.PreFoalingChecklistItem, error) {
	return nil, nil
}

func (m *mockPregnancyRepo) AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	return nil
}

func (m *mockPregnancyRepo) DeletePreFoalingChecklistItem(ctx context.Context, itemID uint) error {
	return nil
}

func (m *mockPregnancyRepo) InitializePreFoalingChecklist(ctx context.Context, horseID uint) error {
	return nil
}

func (m *mockPregnancyRepo) GetPreFoalingSigns(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error) {
	return nil, nil
}

func (m *mockPregnancyRepo) AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error {
	return nil
}

func (m *mockPregnancyRepo) GetCurrentPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	return nil, nil
}

func (m *mockPregnancyRepo) UpdatePregnancyStatus(ctx context.Context, horseID uint, isPregnant bool, conceptionDate *time.Time) error {
	return nil
}

func (m *mockPregnancyRepo) UpdatePreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	return nil
}

func setupTestHandler() (*Handler, *mocks.MockDB) {
	db := mocks.NewMockDB()
	mockDB := &database.PostgresDB{DB: db.GetDB()}
	
	// Create repositories with mock DB
	horseRepo := postgres.NewHorseRepository(db.GetDB())
	userRepo := postgres.NewUserRepository(db.GetDB())
	healthRepo := postgres.NewHealthRepository(db.GetDB())
	pregnancyRepo := &mockPregnancyRepo{
		getEventsFn: func(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
			return []models.PregnancyEvent{}, nil
		},
	}

	// Create services with repositories
	horseService := service.NewHorseService(horseRepo)
	userService := service.NewUserService(userRepo)
	healthService := service.NewHealthService(healthRepo)
	pregnancyService := service.NewPregnancyService(horseRepo, pregnancyRepo)

	config := HandlerConfig{
		Database:         mockDB,
		UserService:      userService,
		HorseService:     horseService,
		HealthService:    healthService,
		PregnancyService: pregnancyService,
		Cache:           cache.NewMemoryCache(),
	}

	return NewHandler(config), db
}
