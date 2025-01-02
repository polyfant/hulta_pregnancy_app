package models

import "time"

type Expense struct {
    ID          uint      `gorm:"primaryKey"`
    UserID      string    `gorm:"type:text;not null;index"`
    HorseID     uint      `gorm:"index"`
    Type        string    `gorm:"type:varchar(50);not null"`
    Amount      float64   `gorm:"type:decimal(10,2);not null;check:amount >= 0"`
    Date        time.Time `gorm:"not null"`
    Description string    `gorm:"type:text"`
    Receipt     string    `gorm:"type:text"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Constants for expense types
const (
    ExpenseTypeFeed        = "FEED"
    ExpenseTypeVet         = "VET"
    ExpenseTypeFarrier     = "FARRIER"
    ExpenseTypeEquipment   = "EQUIPMENT"
    ExpenseTypeOther       = "OTHER"
)

// Constants for expense frequencies
const (
    FrequencyDaily     = "DAILY"
    FrequencyWeekly    = "WEEKLY"
    FrequencyMonthly   = "MONTHLY"
    FrequencyQuarterly = "QUARTERLY"
    FrequencyYearly    = "YEARLY"
)

type RecurringExpense struct {
    ID          uint      `gorm:"primaryKey"`
    UserID      string    `gorm:"type:text;not null;index"`
    HorseID     *uint     `gorm:"index"`
    Type        string    `gorm:"type:varchar(50);not null"`
    Amount      float64   `gorm:"type:decimal(10,2);not null;check:amount >= 0"`
    Frequency   string    `gorm:"type:varchar(50);not null"`
    StartDate   time.Time `gorm:"not null"`
    EndDate     *time.Time
    Description string    `gorm:"type:text"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
} 