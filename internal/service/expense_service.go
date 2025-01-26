package service

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type ExpenseService struct {
	expenseRepo         repository.ExpenseRepository
	recurringExpenseRepo repository.RecurringExpenseRepository
}

func NewExpenseService(
	expenseRepo repository.ExpenseRepository,
	recurringExpenseRepo repository.RecurringExpenseRepository,
) *ExpenseService {
	return &ExpenseService{
		expenseRepo:         expenseRepo,
		recurringExpenseRepo: recurringExpenseRepo,
	}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, expense *models.Expense) error {
	if expense.Amount <= 0 {
		return models.ErrInvalidAmount
	}

	return s.expenseRepo.Create(ctx, expense)
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	if expense.Amount <= 0 {
		return models.ErrInvalidAmount
	}

	// Validate expense type
	switch expense.ExpenseType {
	case models.ExpenseTypeFeed, models.ExpenseTypeVeterinary, models.ExpenseTypeFarrier,
		models.ExpenseTypeEquipment, models.ExpenseTypeTraining, models.ExpenseTypeCompetition,
		models.ExpenseTypeTransport, models.ExpenseTypeInsurance, models.ExpenseTypeBoarding,
		models.ExpenseTypeOther:
		// Valid type
	default:
		return models.ErrInvalidExpenseType
	}

	return s.expenseRepo.Update(ctx, expense)
}

func (s *ExpenseService) RecordExpense(ctx context.Context, expense *models.Expense) error {
	// Validate expense data
	if err := s.validateExpense(expense); err != nil {
		return err
	}

	// Set timestamps
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()

	return s.expenseRepo.Create(ctx, expense)
}

func (s *ExpenseService) GetHorseExpenses(ctx context.Context, horseID uint) ([]models.Expense, error) {
	return s.expenseRepo.GetByHorseID(ctx, horseID)
}

func (s *ExpenseService) GetUserTotalExpenses(ctx context.Context, userID string) (float64, error) {
	totalExpenses, err := s.expenseRepo.GetTotalExpensesByUser(ctx, userID)
	if err != nil {
		return 0, err
	}
	f, exact := totalExpenses.Float64()
	if !exact {
		// If the conversion is not exact, log a warning or handle it as needed
		fmt.Printf("Warning: Inexact conversion of total expenses for user %s\n", userID)
	}
	return f, nil
}

func (s *ExpenseService) GetExpensesByType(ctx context.Context, userID, expenseType string) ([]models.Expense, error) {
	return s.expenseRepo.GetExpensesByType(ctx, userID, expenseType)
}

func (s *ExpenseService) CreateRecurringExpense(ctx context.Context, recurringExpense *models.RecurringExpense) error {
	if recurringExpense.Amount <= 0 {
		return models.ErrInvalidAmount
	}

	// Validate expense type
	switch recurringExpense.ExpenseType {
	case models.ExpenseTypeFeed, models.ExpenseTypeVeterinary, models.ExpenseTypeFarrier,
		models.ExpenseTypeEquipment, models.ExpenseTypeTraining, models.ExpenseTypeCompetition,
		models.ExpenseTypeTransport, models.ExpenseTypeInsurance, models.ExpenseTypeBoarding,
		models.ExpenseTypeOther:
		// Valid type
	default:
		return models.ErrInvalidExpenseType
	}

	// Calculate next due date based on frequency
	switch recurringExpense.Frequency {
	case models.FrequencyDaily:
		recurringExpense.NextDueDate = recurringExpense.StartDate.AddDate(0, 0, 1)
	case models.FrequencyWeekly:
		recurringExpense.NextDueDate = recurringExpense.StartDate.AddDate(0, 0, 7)
	case models.FrequencyMonthly:
		recurringExpense.NextDueDate = recurringExpense.StartDate.AddDate(0, 1, 0)
	case models.FrequencyYearly:
		recurringExpense.NextDueDate = recurringExpense.StartDate.AddDate(1, 0, 0)
	default:
		return models.ErrInvalidFrequency
	}

	return s.recurringExpenseRepo.Create(ctx, recurringExpense)
}

func (s *ExpenseService) ProcessDueRecurringExpenses(ctx context.Context) error {
	dueExpenses, err := s.recurringExpenseRepo.GetDueRecurringExpenses(ctx)
	if err != nil {
		return err
	}

	for _, recurringExpense := range dueExpenses {
		// Create actual expense
		var horseID uint
		if recurringExpense.HorseID != nil {
			horseID = *recurringExpense.HorseID
		}
		expense := &models.Expense{
			UserID:      recurringExpense.UserID,
			HorseID:     horseID,
			Amount:      recurringExpense.Amount,
			Description: recurringExpense.Description,
			Date:        time.Now(),
			ExpenseType: recurringExpense.ExpenseType,
		}

		// Record the expense
		if err := s.RecordExpense(ctx, expense); err != nil {
			return err
		}

		// Update next due date based on frequency
		recurringExpense.NextDueDate = calculateNextDueDate(recurringExpense)
		// You might want to update the recurring expense in the repository
	}

	return nil
}

func (s *ExpenseService) validateExpense(expense *models.Expense) error {
	if expense.Amount < 0 {
		return models.ErrInvalidAmount
	}
	if expense.ExpenseType == "" {
		return models.ErrInvalidExpenseType
	}
	return nil
}

func (s *ExpenseService) validateRecurringExpense(recurringExpense *models.RecurringExpense) error {
	if recurringExpense.Amount < 0 {
		return models.ErrInvalidAmount
	}
	if recurringExpense.ExpenseType == "" {
		return models.ErrInvalidExpenseType
	}
	if recurringExpense.Frequency == "" {
		return models.ErrInvalidFrequency
	}
	if recurringExpense.StartDate.IsZero() {
		return fmt.Errorf("recurring expense start date is required")
	}
	return nil
}

// Helper function to calculate next due date
func calculateNextDueDate(recurringExpense models.RecurringExpense) time.Time {
	switch recurringExpense.Frequency {
	case models.FrequencyDaily:
		return recurringExpense.NextDueDate.AddDate(0, 0, 1)
	case models.FrequencyWeekly:
		return recurringExpense.NextDueDate.AddDate(0, 0, 7)
	case models.FrequencyMonthly:
		return recurringExpense.NextDueDate.AddDate(0, 1, 0)
	case models.FrequencyYearly:
		return recurringExpense.NextDueDate.AddDate(1, 0, 0)
	default:
		return recurringExpense.NextDueDate
	}
}
