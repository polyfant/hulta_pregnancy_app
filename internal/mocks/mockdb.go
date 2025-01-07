package mocks

import "gorm.io/gorm"

type MockDB struct {
    *gorm.DB
}

func NewMockDB() *MockDB {
    return &MockDB{}
} 