package service

import (
	"context"
	"fmt"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

// UserService handles user-related business logic
type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

// Core user management methods
func (s *UserServiceImpl) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *UserServiceImpl) Create(ctx context.Context, user *models.User) error {
	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *UserServiceImpl) Update(ctx context.Context, user *models.User) error {
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *UserServiceImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (s *UserServiceImpl) GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx, userID)
}

func (s *UserServiceImpl) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, user *models.User) error {
	return s.repo.Update(ctx, user)
}
