package service

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type UserService struct {
	userRepo     repository.UserRepository
	horseRepo    repository.HorseRepository
	expenseRepo  repository.ExpenseRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	horseRepo repository.HorseRepository,
	expenseRepo repository.ExpenseRepository,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		horseRepo:   horseRepo,
		expenseRepo: expenseRepo,
	}
}

func (s *UserService) CreateUserFromAuth0(ctx context.Context, user *models.User) error {
	// Validate user data
	if err := s.validateUserCreation(user); err != nil {
		return err
	}

	// Set initial timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastLogin = time.Now()
	user.IsActive = true

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) GetUserDashboard(ctx context.Context, userID string) (*models.UserDashboard, error) {
	// Fetch user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Fetch user's horses
	horses, err := s.horseRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Fetch total expenses
	totalExpenses, err := s.expenseRepo.GetTotalExpensesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Construct dashboard
	return &models.UserDashboard{
		User:          *user,
		Horses:        horses,
		TotalExpenses: totalExpenses,
	}, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, user *models.User) error {
	// Validate update
	if err := s.validateUserUpdate(user); err != nil {
		return err
	}

	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) RecordUserLogin(ctx context.Context, userID string) error {
	return s.userRepo.UpdateLastLogin(ctx, userID)
}

func (s *UserService) validateUserCreation(user *models.User) error {
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}

	if user.ID == "" {
		return fmt.Errorf("user ID is required")
	}

	return nil
}

func (s *UserService) validateUserUpdate(user *models.User) error {
	if user.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	// Additional validation logic
	return nil
}
