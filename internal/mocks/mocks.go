package mocks

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockHorseRepository struct {
	mock.Mock
}

func (m *MockHorseRepository) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Horse), args.Error(1)
}

func (m *MockHorseRepository) ListByUserID(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *MockHorseRepository) Create(ctx context.Context, horse *models.Horse) error {
	args := m.Called(ctx, horse)
	return args.Error(0)
}

func (m *MockHorseRepository) Update(ctx context.Context, horse *models.Horse) error {
	args := m.Called(ctx, horse)
	return args.Error(0)
}

func (m *MockHorseRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockHorseRepository) GetPregnant(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *MockHorseRepository) GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error) {
	args := m.Called(ctx, horseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FamilyTree), args.Error(1)
}

func (m *MockHorseRepository) ListByUser(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

type MockPregnancyRepository struct {
	mock.Mock
}

func (m *MockPregnancyRepository) GetByHorseID(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	args := m.Called(ctx, horseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Pregnancy), args.Error(1)
}

func (m *MockPregnancyRepository) Create(ctx context.Context, pregnancy *models.Pregnancy) error {
	args := m.Called(ctx, pregnancy)
	return args.Error(0)
}

func (m *MockPregnancyRepository) Update(ctx context.Context, pregnancy *models.Pregnancy) error {
	args := m.Called(ctx, pregnancy)
	return args.Error(0)
}

func (m *MockPregnancyRepository) GetEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
	args := m.Called(ctx, horseID)
	return args.Get(0).([]models.PregnancyEvent), args.Error(1)
}

func (m *MockPregnancyRepository) GetActive(ctx context.Context, userID string) ([]models.Pregnancy, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Pregnancy), args.Error(1)
}

func (m *MockPregnancyRepository) AddPreFoaling(ctx context.Context, sign *models.PreFoalingSign) error {
	args := m.Called(ctx, sign)
	return args.Error(0)
}

func (m *MockPregnancyRepository) GetPreFoaling(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error) {
	args := m.Called(ctx, horseID)
	return args.Get(0).([]models.PreFoalingSign), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockHealthRepository struct {
	mock.Mock
}

func (m *MockHealthRepository) CreateRecord(ctx context.Context, record *models.HealthRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockHealthRepository) GetRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error) {
	args := m.Called(ctx, horseID)
	return args.Get(0).([]models.HealthRecord), args.Error(1)
}

func (m *MockHealthRepository) UpdateRecord(ctx context.Context, record *models.HealthRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockHealthRepository) DeleteRecord(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockBreedingRepository struct {
	mock.Mock
}

func (m *MockBreedingRepository) CreateRecord(ctx context.Context, record *models.BreedingRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *MockBreedingRepository) GetRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error) {
	args := m.Called(ctx, horseID)
	return args.Get(0).([]models.BreedingRecord), args.Error(1)
}

func (m *MockBreedingRepository) UpdateRecord(ctx context.Context, record *models.BreedingRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

// Continue with other repositories... 