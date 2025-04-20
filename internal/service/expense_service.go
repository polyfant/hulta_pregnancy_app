package service

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
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
		err := models.ErrInvalidAmount
		logger.Error(err, "invalid amount in CreateExpense", "userID", expense.UserID, "amount", expense.Amount)
		return err
	}

	err := s.expenseRepo.Create(ctx, expense)
	if err != nil {
		logger.Error(err, "failed to create expense", "userID", expense.UserID)
	}
	return err
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	if expense.Amount <= 0 {
		err := models.ErrInvalidAmount
		logger.Error(err, "invalid amount in UpdateExpense", "userID", expense.UserID, "amount", expense.Amount)
		return err
	}

	// Validate expense type
	switch expense.ExpenseType {
	case models.ExpenseTypeFeed, models.ExpenseTypeVeterinary, models.ExpenseTypeFarrier,
		models.ExpenseTypeEquipment, models.ExpenseTypeTraining, models.ExpenseTypeCompetition,
		models.ExpenseTypeTransport, models.ExpenseTypeInsurance, models.ExpenseTypeBoarding,
		models.ExpenseTypeOther:
		// Valid type
	default:
		err := models.ErrInvalidExpenseType
		logger.Error(err, "invalid expense type in UpdateExpense", "userID", expense.UserID, "type", expense.ExpenseType)
		return err
	}

	err := s.expenseRepo.Update(ctx, expense)
	if err != nil {
		logger.Error(err, "failed to update expense", "userID", expense.UserID)
	}
	return err
}

func (s *ExpenseService) RecordExpense(ctx context.Context, expense *models.Expense) error {
	if err := s.validateExpense(expense); err != nil {
		logger.Error(err, "invalid expense in RecordExpense", "userID", expense.UserID)
		return err
	}

	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()

	err := s.expenseRepo.Create(ctx, expense)
	if err != nil {
		logger.Error(err, "failed to record expense", "userID", expense.UserID)
	}
	return err
}

func (s *ExpenseService) GetHorseExpenses(ctx context.Context, horseID uint) ([]models.Expense, error) {
	expenses, err := s.expenseRepo.GetByHorseID(ctx, horseID)
	if err != nil {
		logger.Error(err, "failed to get horse expenses", "horseID", horseID)
	}
	return expenses, err
}

func (s *ExpenseService) GetUserTotalExpenses(ctx context.Context, userID string) (float64, error) {
	totalExpenses, err := s.expenseRepo.GetTotalExpensesByUser(ctx, userID)
	if err != nil {
		logger.Error(err, "failed to get total expenses by user", "userID", userID)
		return 0, err
	}
	f, exact := totalExpenses.Float64()
	if !exact {
		logger.Warn("inexact conversion of total expenses", "userID", userID)
	}
	return f, nil
}

func (s *ExpenseService) GetExpensesByType(ctx context.Context, userID, expenseType string) ([]models.Expense, error) {
	expenses, err := s.expenseRepo.GetExpensesByType(ctx, userID, expenseType)
	if err != nil {
		logger.Error(err, "failed to get expenses by type", "userID", userID, "type", expenseType)
	}
	return expenses, err
}

func (s *ExpenseService) CreateRecurringExpense(ctx context.Context, recurringExpense *models.RecurringExpense) error {
	if recurringExpense.Amount <= 0 {
		err := models.ErrInvalidAmount
		logger.Error(err, "invalid amount in CreateRecurringExpense", "userID", recurringExpense.UserID, "amount", recurringExpense.Amount)
		return err
	}

	// Validate expense type
	switch recurringExpense.ExpenseType {
	case models.ExpenseTypeFeed, models.ExpenseTypeVeterinary, models.ExpenseTypeFarrier,
		models.ExpenseTypeEquipment, models.ExpenseTypeTraining, models.ExpenseTypeCompetition,
		models.ExpenseTypeTransport, models.ExpenseTypeInsurance, models.ExpenseTypeBoarding,
		models.ExpenseTypeOther:
		// Valid type
	default:
		err := models.ErrInvalidExpenseType
		logger.Error(err, "invalid expense type in CreateRecurringExpense", "userID", recurringExpense.UserID, "type", recurringExpense.ExpenseType)
		return err
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
		err := models.ErrInvalidFrequency
		logger.Error(err, "invalid frequency in CreateRecurringExpense", "userID", recurringExpense.UserID, "frequency", recurringExpense.Frequency)
		return err
	}

	err := s.recurringExpenseRepo.Create(ctx, recurringExpense)
	if err != nil {
		logger.Error(err, "failed to create recurring expense", "userID", recurringExpense.UserID)
	}
	return err
}

func (s *ExpenseService) ProcessDueRecurringExpenses(ctx context.Context) error {
	dueExpenses, err := s.recurringExpenseRepo.GetDueRecurringExpenses(ctx)
	if err != nil {
		logger.Error(err, "failed to get due recurring expenses")
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
			logger.Error(err, "failed to record recurring expense", "userID", recurringExpense.UserID)
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
