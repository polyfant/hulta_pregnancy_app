package models

import "time"

// DataStore interface
type DataStore interface {
	// Horse operations
	GetHorse(id int64) (Horse, error)
	GetAllHorses() ([]Horse, error)
	AddHorse(horse *Horse) error
	UpdateHorse(horse *Horse) error
	DeleteHorse(id int64) error
	GetHorsesByUser(userID string) ([]Horse, error)

	// Health record operations
	GetHealthRecords(horseID int64) ([]HealthRecord, error)
	AddHealthRecord(record *HealthRecord) error
	UpdateHealthRecord(record *HealthRecord) error

	// Pregnancy operations
	GetPregnancies(userID string) ([]Pregnancy, error)
	GetPregnancy(id int64) (Pregnancy, error)
	AddPregnancy(pregnancy *Pregnancy) error
	UpdatePregnancy(pregnancy *Pregnancy) error
	GetPregnancyEvents(pregnancyID int64) ([]PregnancyEvent, error)
	AddPregnancyEvent(event *PregnancyEvent) error

	// Pre-foaling operations
	GetPreFoalingSigns(horseID int64) ([]PreFoalingSign, error)
	AddPreFoalingSign(sign *PreFoalingSign) error

	// Pre-foaling checklist operations
	GetPreFoalingChecklist(horseID int64) ([]PreFoalingChecklistItem, error)
	AddPreFoalingChecklistItem(item *PreFoalingChecklistItem) error
	UpdatePreFoalingChecklistItem(item *PreFoalingChecklistItem) error
	DeletePreFoalingChecklistItem(id int64) error

	// Expense operations
	GetExpenses(userID string, from, to time.Time) ([]Expense, error)
	GetHorseExpenses(horseID int64, from, to time.Time) ([]Expense, error)
	AddExpense(expense *Expense) error
	UpdateExpense(expense *Expense) error
	DeleteExpense(id int64) error
	GetExpenseSummary(userID string, from, to time.Time) (map[string]float64, error)

	// Recurring expense operations
	GetRecurringExpenses(userID string) ([]RecurringExpense, error)
	AddRecurringExpense(expense *RecurringExpense) error
	UpdateRecurringExpense(expense *RecurringExpense) error
	DeleteRecurringExpense(id int64) error

	// Add these methods
	GetBreedingCosts(horseID uint) ([]BreedingCost, error)
	GetBreedingRecords(horseID int64) ([]BreedingRecord, error)
	AddBreedingRecord(record *BreedingRecord) error
}
