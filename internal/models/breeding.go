package models

import "time"

type BreedingRecord struct {
    ID          uint      `gorm:"primaryKey"`
    UserID      string    `gorm:"index;not null"`
    MareID      uint      `gorm:"not null"`
    StallionID  *uint     
    Date        time.Time `gorm:"not null"`
    Notes       string
    Successful  *bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type BreedingCost struct {
    ID          uint      `gorm:"primaryKey"`
    UserID      string    `gorm:"index;not null"`
    HorseID     uint      `gorm:"not null"`
    Type        string    `gorm:"size:50"`
    Amount      float64
    Date        time.Time
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Constants for breeding status
const (
    BreedingStatusPending    = "PENDING"
    BreedingStatusSuccessful = "SUCCESSFUL"
    BreedingStatusFailed     = "FAILED"
) 