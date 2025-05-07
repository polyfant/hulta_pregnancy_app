package models

import (
	"time"

	"github.com/shopspring/decimal"
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
	ID          uint        `json:"id,omitempty" gorm:"primaryKey"`
	UserID      string      `json:"user_id" gorm:"type:text;not null;index" validate:"required"`
	HorseID     uint        `json:"horse_id" gorm:"index" validate:"required"`
	ExpenseType ExpenseType `json:"expense_type" gorm:"type:expense_type;not null" validate:"required,oneof=feed veterinary farrier equipment training competition transport insurance boarding other"`
	Amount      decimal.Decimal `json:"amount" gorm:"type:decimal(10,2);not null;check:amount >= 0" validate:"required,gt=0"`
	Date        time.Time   `json:"date" gorm:"not null" validate:"required"`
	Description string      `json:"description,omitempty" gorm:"type:text" validate:"omitempty,max=2000"`
	Receipt     string      `json:"receipt,omitempty" gorm:"type:text" validate:"omitempty,uri"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
}

type RecurringExpense struct {
	gorm.Model
	UserID      string      `json:"user_id" validate:"required"`
	HorseID     *uint       `json:"horse_id" validate:"omitempty"`
	ExpenseType ExpenseType `json:"expense_type" gorm:"type:expense_type;not null" validate:"required,oneof=feed veterinary farrier equipment training competition transport insurance boarding other"`
	Amount      decimal.Decimal `json:"amount" gorm:"type:decimal(10,2);not null;check:amount >= 0" validate:"required,gt=0"`
	Frequency   Frequency   `json:"frequency" gorm:"type:varchar(50);not null" validate:"required,oneof=DAILY WEEKLY MONTHLY YEARLY"`
	StartDate   time.Time   `json:"start_date" gorm:"not null" validate:"required"`
	EndDate     *time.Time  `json:"end_date,omitempty" gorm:"default:null" validate:"omitempty,gtfield=StartDate"`
	Description string      `json:"description,omitempty" gorm:"type:text" validate:"omitempty,max=2000"`
	NextDueDate time.Time   `json:"next_due_date,omitempty"`
}