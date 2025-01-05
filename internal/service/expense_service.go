package service

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
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

func (s *ExpenseService) GetUserTotalExpenses(ctx context.Context, userID string) (decimal.Decimal, error) {
	return s.expenseRepo.GetTotalExpensesByUser(ctx, userID)
}

func (s *ExpenseService) GetExpensesByType(ctx context.Context, userID, expenseType string) ([]models.Expense, error) {
	return s.expenseRepo.GetExpensesByType(ctx, userID, expenseType)
}

func (s *ExpenseService) CreateRecurringExpense(ctx context.Context, recurringExpense *models.RecurringExpense) error {
	// Validate recurring expense
	if err := s.validateRecurringExpense(recurringExpense); err != nil {
		return err
	}

	// Set initial timestamps
	recurringExpense.CreatedAt = time.Now()
	recurringExpense.UpdatedAt = time.Now()

	return s.recurringExpenseRepo.Create(ctx, recurringExpense)
}

func (s *ExpenseService) ProcessDueRecurringExpenses(ctx context.Context) error {
	dueExpenses, err := s.recurringExpenseRepo.GetDueRecurringExpenses(ctx)
	if err != nil {
		return err
	}

	for _, recurringExpense := range dueExpenses {
		// Create actual expense
		expense := &models.Expense{
			UserID:      recurringExpense.UserID,
			HorseID:     recurringExpense.HorseID,
			Amount:      recurringExpense.Amount,
			Description: recurringExpense.Description,
			Date:        time.Now(),
			Type:        recurringExpense.Frequency, // Use frequency as type
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
	if expense.Amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("expense amount must be positive")
	}

	if expense.Date.IsZero() {
		return fmt.Errorf("expense date is required")
	}

	return nil
}

func (s *ExpenseService) validateRecurringExpense(recurringExpense *models.RecurringExpense) error {
	if recurringExpense.Amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("recurring expense amount must be positive")
	}

	if recurringExpense.Frequency == "" {
		return fmt.Errorf("recurring expense frequency is required")
	}

	return nil
}

// Helper function to calculate next due date
func calculateNextDueDate(recurringExpense models.RecurringExpense) time.Time {
	switch recurringExpense.Frequency {
	case "daily":
		return recurringExpense.NextDueDate.AddDate(0, 0, 1)
	case "weekly":
		return recurringExpense.NextDueDate.AddDate(0, 0, 7)
	case "monthly":
		return recurringExpense.NextDueDate.AddDate(0, 1, 0)
	case "yearly":
		return recurringExpense.NextDueDate.AddDate(1, 0, 0)
	default:
		return recurringExpense.NextDueDate
	}
}
