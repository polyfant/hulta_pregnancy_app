package models

import "time"

type HealthRecord struct {
    ID          uint      `gorm:"primaryKey"`
    HorseID     uint      `gorm:"not null"`
    UserID      string    `gorm:"index;not null"`
    Type        string    `gorm:"size:50"`
    Date        time.Time
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
} 