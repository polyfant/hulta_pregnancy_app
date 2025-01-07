package mocks

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func NewMockDB() *MockDB {
	return &MockDB{}
}

// Implement methods to mock gorm.DB interactions
func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	mockedArgs := m.Called(query, args)
	return mockedArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Error() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) GetDB() *gorm.DB {
	return &gorm.DB{}  // Return a mock gorm.DB instance
}