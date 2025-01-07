package models

import "time"

type BreedingRecord struct {
    ID        uint           `gorm:"primaryKey"`
    HorseID   uint          `gorm:"not null"`
    Date      time.Time
    Status    BreedingStatus `gorm:"size:50"`
    CreatedAt time.Time
    UpdatedAt time.Time
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