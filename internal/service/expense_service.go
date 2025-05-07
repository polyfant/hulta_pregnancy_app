package service

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

// Helper function to calculate next due date
func calculateNextDueDate(recurringExpense models.RecurringExpense) time.Time {
	currentDueDate := recurringExpense.NextDueDate
	// If NextDueDate is zero (e.g., for a new record), use StartDate.
	// This also handles cases where StartDate might be explicitly set as the base for recalculation.
	if currentDueDate.IsZero() {
		currentDueDate = recurringExpense.StartDate
	}

	switch recurringExpense.Frequency {
	case models.FrequencyDaily:
		return currentDueDate.AddDate(0, 0, 1)
	case models.FrequencyWeekly:
		return currentDueDate.AddDate(0, 0, 7)
	case models.FrequencyMonthly:
		return currentDueDate.AddDate(0, 1, 0)
	case models.FrequencyYearly:
		return currentDueDate.AddDate(1, 0, 0)
	default:
		logger.Warn("Unknown or unsupported frequency for recurring expense", 
			"frequency", recurringExpense.Frequency, 
			"recurringExpenseID", recurringExpense.ID, 
			"userID", recurringExpense.UserID)
		return currentDueDate // Return current to avoid breaking, but signals an issue.
	}
}

// --- Interface Definition ---
// It's good practice to define an interface for the service
type ExpenseService interface {
	// Expense methods
	CreateExpense(ctx context.Context, expense *models.Expense) (*models.Expense, error)
	UpdateExpense(ctx context.Context, expense *models.Expense) (*models.Expense, error)
	GetExpenseByID(ctx context.Context, userID string, expenseID uint) (*models.Expense, error)
	DeleteExpense(ctx context.Context, userID string, expenseID uint) error
	GetHorseExpenses(ctx context.Context, userID string, horseID uint) ([]models.Expense, error)
	GetUserTotalExpenses(ctx context.Context, userID string) (float64, error)
	GetExpensesByType(ctx context.Context, userID, expenseType string) ([]models.Expense, error)

	// Recurring Expense methods
	CreateRecurringExpense(ctx context.Context, recurringExpense *models.RecurringExpense) (*models.RecurringExpense, error)
	UpdateRecurringExpense(ctx context.Context, recurringExpense *models.RecurringExpense) (*models.RecurringExpense, error)
	GetRecurringExpenseByID(ctx context.Context, userID string, recurringExpenseID uint) (*models.RecurringExpense, error)
	DeleteRecurringExpense(ctx context.Context, userID string, recurringExpenseID uint) error
	GetUserRecurringExpenses(ctx context.Context, userID string) ([]models.RecurringExpense, error)
	ProcessDueRecurringExpenses(ctx context.Context) error
}

// --- Implementation ---

type expenseServiceImpl struct { // Renamed struct for clarity with interface
	expenseRepo         repository.ExpenseRepository
	recurringExpenseRepo repository.RecurringExpenseRepository
	horseRepo           repository.HorseRepository // Added horseRepo for ownership checks
}

// Ensure implementation satisfies the interface
var _ ExpenseService = (*expenseServiceImpl)(nil)

func NewExpenseService(
	expenseRepo repository.ExpenseRepository,
	recurringExpenseRepo repository.RecurringExpenseRepository,
	horseRepo repository.HorseRepository, // Added horseRepo parameter
) ExpenseService { // Return the interface type
	return &expenseServiceImpl{
		expenseRepo:         expenseRepo,
		recurringExpenseRepo: recurringExpenseRepo,
		horseRepo:           horseRepo, // Initialize horseRepo
	}
}

// --- Expense Methods ---

func (s *expenseServiceImpl) CreateExpense(ctx context.Context, expense *models.Expense) (*models.Expense, error) {
	// Validation is now done in the handler
	createdExpense, err := s.expenseRepo.Create(ctx, expense)
	if err != nil {
		logger.Error(err, "failed to create expense", "userID", expense.UserID)
		return nil, fmt.Errorf("repository failed to create expense: %w", err)
	}
	logger.Info("Successfully created expense", "expenseID", createdExpense.ID, "userID", expense.UserID)
	return createdExpense, nil
}

func (s *expenseServiceImpl) UpdateExpense(ctx context.Context, expense *models.Expense) (*models.Expense, error) {
	// Validation is now done in the handler

	existingExpense, err := s.expenseRepo.GetByID(ctx, expense.ID)
	if err != nil {
		logger.Error(err, "failed to get expense for update", "expenseID", expense.ID, "userID", expense.UserID)
		return nil, fmt.Errorf("repository failed to get expense for update: %w", err)
	}

	if existingExpense.UserID != expense.UserID {
		logger.Warn("user attempting to update expense they do not own", "targetUserID", existingExpense.UserID, "callerUserID", expense.UserID, "expenseID", expense.ID)
		return nil, fmt.Errorf("permission denied: user does not own this expense")
	}

	updatedExpense, err := s.expenseRepo.Update(ctx, expense)
	if err != nil {
		logger.Error(err, "failed to update expense", "userID", expense.UserID, "expenseID", expense.ID)
		return nil, fmt.Errorf("repository failed to update expense: %w", err)
	}
	logger.Info("Successfully updated expense", "expenseID", updatedExpense.ID, "userID", expense.UserID)
	return updatedExpense, nil
}

func (s *expenseServiceImpl) GetExpenseByID(ctx context.Context, userID string, expenseID uint) (*models.Expense, error) {
	expense, err := s.expenseRepo.GetByID(ctx, expenseID)
	if err != nil {
		logger.Error(err, "failed to get expense by ID from repository", "expenseID", expenseID, "userID", userID)
		return nil, fmt.Errorf("repository failed to get expense by ID: %w", err)
	}

	if expense.UserID != userID {
		logger.Warn("user attempting to access expense they do not own", "targetUserID", expense.UserID, "callerUserID", userID, "expenseID", expenseID)
		return nil, fmt.Errorf("permission denied: user does not own this expense")
	}

	logger.Info("Successfully fetched expense by ID", "expenseID", expenseID, "userID", userID)
	return expense, nil
}

func (s *expenseServiceImpl) DeleteExpense(ctx context.Context, userID string, expenseID uint) error {
	// GetExpenseByID includes ownership check
	expense, err := s.GetExpenseByID(ctx, userID, expenseID)
	if err != nil {
		// GetExpenseByID already logs
		return err // Error already formatted and logged by GetExpenseByID
	}

	err = s.expenseRepo.Delete(ctx, expense.ID) // Use expense.ID from the fetched record
	if err != nil {
		logger.Error(err, "failed to delete expense from repository", "expenseID", expense.ID, "userID", userID)
		return fmt.Errorf("repository failed to delete expense: %w", err)
	}

	logger.Info("Successfully deleted expense", "expenseID", expense.ID, "userID", userID)
	return nil
}

func (s *expenseServiceImpl) GetHorseExpenses(ctx context.Context, userID string, horseID uint) ([]models.Expense, error) {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		logger.Error(err, "failed to get horse for expense listing", "horseID", horseID, "userID", userID)
		return nil, fmt.Errorf("failed to get horse: %w", err)
	}
	if horse.UserID != userID {
		logger.Warn("user attempting to list expenses for a horse they do not own", "targetHorseUserID", horse.UserID, "callerUserID", userID, "horseID", horseID)
		return nil, fmt.Errorf("permission denied: user does not own this horse")
	}

	expenses, err := s.expenseRepo.GetByHorseID(ctx, horseID)
	if err != nil {
		logger.Error(err, "failed to get horse expenses from repository", "horseID", horseID, "userID", userID)
		return nil, fmt.Errorf("repository failed to get horse expenses: %w", err)
	}
	logger.Info("Successfully fetched expenses for horse", "horseID", horseID, "userID", userID, "count", len(expenses))
	return expenses, nil
}

func (s *expenseServiceImpl) GetUserTotalExpenses(ctx context.Context, userID string) (float64, error) {
	totalExpensesDecimal, err := s.expenseRepo.GetTotalExpensesByUser(ctx, userID)
	if err != nil {
		logger.Error(err, "failed to get total expenses by user from repository", "userID", userID)
		return 0, fmt.Errorf("repository failed to get total expenses for user: %w", err)
	}
	f, exact := totalExpensesDecimal.Float64()
	if !exact {
		// This scenario (inexact conversion) is unlikely with typical currency values but good to log.
		logger.Warn("Inexact conversion from decimal to float64 for total expenses", "userID", userID, "decimalValue", totalExpensesDecimal.String())
	}
	logger.Info("Successfully fetched total expenses for user", "userID", userID, "total", f)
	return f, nil
}

func (s *expenseServiceImpl) GetExpensesByType(ctx context.Context, userID, expenseType string) ([]models.Expense, error) {
	expenses, err := s.expenseRepo.GetExpensesByType(ctx, userID, expenseType)
	if err != nil {
		logger.Error(err, "failed to get expenses by type from repository", "userID", userID, "type", expenseType)
		return nil, fmt.Errorf("repository failed to get expenses by type: %w", err)
	}
	logger.Info("Successfully fetched expenses by type", "userID", userID, "type", expenseType, "count", len(expenses))
	return expenses, nil
}

// --- Recurring Expense Methods ---

func (s *expenseServiceImpl) CreateRecurringExpense(ctx context.Context, recurringExpense *models.RecurringExpense) (*models.RecurringExpense, error) {
	// Validation is now done in handler

	// Calculate initial next due date based on frequency and start date
	recurringExpense.NextDueDate = calculateNextDueDate(*recurringExpense) // Ensure this uses StartDate correctly

	createdRecExpense, err := s.recurringExpenseRepo.Create(ctx, recurringExpense)
	if err != nil {
		logger.Error(err, "failed to create recurring expense in repository", "userID", recurringExpense.UserID)
		return nil, fmt.Errorf("repository failed to create recurring expense: %w", err)
	}
	logger.Info("Successfully created recurring expense", "recurringExpenseID", createdRecExpense.ID, "userID", recurringExpense.UserID)
	return createdRecExpense, nil
}

func (s *expenseServiceImpl) UpdateRecurringExpense(ctx context.Context, recurringExpense *models.RecurringExpense) (*models.RecurringExpense, error) {
	// Validation is now done in handler

	existingRecExpense, err := s.recurringExpenseRepo.GetByID(ctx, recurringExpense.ID)
	if err != nil {
		logger.Error(err, "failed to get recurring expense for update from repository", "recurringExpenseID", recurringExpense.ID, "userID", recurringExpense.UserID)
		return nil, fmt.Errorf("repository failed to get recurring expense for update: %w", err)
	}

	if existingRecExpense.UserID != recurringExpense.UserID {
		logger.Warn("user attempting to update recurring expense they do not own", 
			"targetUserID", existingRecExpense.UserID, 
			"callerUserID", recurringExpense.UserID, 
			"recurringExpenseID", recurringExpense.ID)
		return nil, fmt.Errorf("permission denied: user does not own this recurring expense")
	}

	// Recalculate NextDueDate if relevant fields (StartDate, Frequency, Interval) might have changed.
	// The handler should pass the updated model. calculateNextDueDate uses these fields.
	recurringExpense.NextDueDate = calculateNextDueDate(*recurringExpense)

	updatedRecExpense, err := s.recurringExpenseRepo.Update(ctx, recurringExpense)
	if err != nil {
		logger.Error(err, "failed to update recurring expense in repository", "userID", recurringExpense.UserID, "recurringExpenseID", recurringExpense.ID)
		return nil, fmt.Errorf("repository failed to update recurring expense: %w", err)
	}
	logger.Info("Successfully updated recurring expense", "recurringExpenseID", updatedRecExpense.ID, "userID", recurringExpense.UserID)
	return updatedRecExpense, nil
}

func (s *expenseServiceImpl) GetRecurringExpenseByID(ctx context.Context, userID string, recurringExpenseID uint) (*models.RecurringExpense, error) {
	record, err := s.recurringExpenseRepo.GetByID(ctx, recurringExpenseID)
	if err != nil {
		logger.Error(err, "failed to get recurring expense by ID from repository", "recurringExpenseID", recurringExpenseID, "userID", userID)
		return nil, fmt.Errorf("repository failed to get recurring expense by ID: %w", err)
	}

	if record.UserID != userID {
		logger.Warn("user attempting to access recurring expense they do not own", "targetUserID", record.UserID, "callerUserID", userID, "recurringExpenseID", recurringExpenseID)
		return nil, fmt.Errorf("permission denied: user does not own this recurring expense")
	}

	logger.Info("Successfully fetched recurring expense by ID", "recurringExpenseID", recurringExpenseID, "userID", userID)
	return record, nil
}

func (s *expenseServiceImpl) DeleteRecurringExpense(ctx context.Context, userID string, recurringExpenseID uint) error {
	// GetRecurringExpenseByID includes ownership check and logging for fetch errors
	record, err := s.GetRecurringExpenseByID(ctx, userID, recurringExpenseID)
	if err != nil {
		return err // Error already formatted and logged by GetRecurringExpenseByID
	}

	err = s.recurringExpenseRepo.Delete(ctx, record.ID) // Use record.ID from the fetched record
	if err != nil {
		logger.Error(err, "failed to delete recurring expense from repository", "recurringExpenseID", record.ID, "userID", userID)
		return fmt.Errorf("repository failed to delete recurring expense: %w", err)
	}

	logger.Info("Successfully deleted recurring expense", "recurringExpenseID", record.ID, "userID", userID)
	return nil
}

func (s *expenseServiceImpl) GetUserRecurringExpenses(ctx context.Context, userID string) ([]models.RecurringExpense, error) {
	expenses, err := s.recurringExpenseRepo.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error(err, "failed to get user recurring expenses from repository", "userID", userID)
		return nil, fmt.Errorf("repository failed to get user recurring expenses: %w", err)
	}
	logger.Info("Successfully fetched recurring expenses for user", "userID", userID, "count", len(expenses))
	return expenses, nil
}

func (s *expenseServiceImpl) ProcessDueRecurringExpenses(ctx context.Context) error {
	dueExpenses, err := s.recurringExpenseRepo.GetDueRecurringExpenses(ctx)
	if err != nil {
		logger.Error(err, "failed to get due recurring expenses from repository")
		return fmt.Errorf("repository failed to get due recurring expenses: %w", err)
	}

	if len(dueExpenses) == 0 {
		logger.Info("No due recurring expenses to process.")
		return nil
	}

	logger.Info("Processing due recurring expenses", "count", len(dueExpenses))
	var processedCount, errorCount int

	for i := range dueExpenses { // Iterate by index to modify the original slice item for NextDueDate update
		recurringExpense := dueExpenses[i] // Get a mutable copy for modifications if necessary, or work with dueExpenses[i] directly

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
			Date:        recurringExpense.NextDueDate, // Use the current NextDueDate for the generated expense
			ExpenseType: recurringExpense.ExpenseType,
			CreatedAt:   time.Now(), // Set creation time for the new expense record
			UpdatedAt:   time.Now(), // Set update time for the new expense record
		}

		createdExpense, err := s.expenseRepo.Create(ctx, expense)
		if err != nil {
			logger.Error(err, "failed to create expense from recurring expense", 
				"recurringExpenseID", recurringExpense.ID, 
				"userID", recurringExpense.UserID)
			errorCount++
			continue // Skip updating this recurring expense if creating the actual expense failed
		}
		logger.Info("Successfully created expense from recurring item", "expenseID", createdExpense.ID, "recurringExpenseID", recurringExpense.ID)

		// Update NextDueDate for the original recurring expense in the slice/database
		dueExpenses[i].NextDueDate = calculateNextDueDate(dueExpenses[i])
		// Also update UpdatedAt for the recurring expense
		dueExpenses[i].UpdatedAt = time.Now()

		_, err = s.recurringExpenseRepo.Update(ctx, &dueExpenses[i]) // Pass pointer to the item in the slice
		if err != nil {
			logger.Error(err, "failed to update NextDueDate for recurring expense in repository", 
				"recurringExpenseID", dueExpenses[i].ID,
				"userID", dueExpenses[i].UserID)
			errorCount++
			// If updating NextDueDate fails, the expense might be processed again on the next run.
			// This is a critical error and might need manual intervention or a more robust retry/dead-letter queue.
		} else {
			processedCount++
			logger.Info("Successfully updated NextDueDate for recurring expense", "recurringExpenseID", dueExpenses[i].ID, "newNextDueDate", dueExpenses[i].NextDueDate)
		}
	}

	logger.Info("Finished processing due recurring expenses", "totalDue", len(dueExpenses), "successfullyProcessed", processedCount, "errorsEncountered", errorCount)
	if errorCount > 0 {
		return fmt.Errorf("encountered %d errors while processing %d due recurring expenses", errorCount, len(dueExpenses))
	}
	return nil
}