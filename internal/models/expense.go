package models

import (
	"time"

	"gorm.io/gorm"
)

type ExpenseType string
type Frequency string

const (
	ExpenseTypeFeed        ExpenseType = "feed"
	ExpenseTypeVeterinary  ExpenseType = "veterinary"
	ExpenseTypeFarrier     ExpenseType = "farrier"
	ExpenseTypeEquipment   ExpenseType = "equipment"
	ExpenseTypeTraining    ExpenseType = "training"
	ExpenseTypeCompetition ExpenseType = "competition"
	ExpenseTypeTransport   ExpenseType = "transport"
	ExpenseTypeInsurance   ExpenseType = "insurance"
	ExpenseTypeBoarding    ExpenseType = "boarding"
	ExpenseTypeOther       ExpenseType = "other"
)

const (
	FrequencyDaily   Frequency = "DAILY"
	FrequencyWeekly  Frequency = "WEEKLY"
	FrequencyMonthly Frequency = "MONTHLY"
	FrequencyYearly  Frequency = "YEARLY"
)

type Expense struct {
	ID          uint        `gorm:"primaryKey"`
	UserID      string      `gorm:"type:text;not null;index"`
	HorseID     uint        `gorm:"index"`
	ExpenseType ExpenseType `gorm:"type:expense_type;not null"`
	Amount      float64     `gorm:"type:decimal(10,2);not null;check:amount >= 0"`
	Date        time.Time   `gorm:"not null"`
	Description string      `gorm:"type:text"`
	Receipt     string      `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RecurringExpense struct {
	gorm.Model
	UserID      string      `json:"user_id"`
	HorseID     *uint       `json:"horse_id"`
	ExpenseType ExpenseType `gorm:"type:expense_type;not null"`
	Amount      float64     `gorm:"type:decimal(10,2);not null;check:amount >= 0"`
	Frequency   Frequency   `gorm:"type:varchar(50);not null"`
	StartDate   time.Time   `gorm:"not null"`
	EndDate     *time.Time
	Description string      `gorm:"type:text"`
	NextDueDate time.Time   `json:"next_due_date"`
}