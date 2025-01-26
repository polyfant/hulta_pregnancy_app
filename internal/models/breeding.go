package models

import "time"

type BreedingRecord struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    HorseID   uint      `json:"horse_id"`
    UserID    string    `json:"user_id"`
    Date      time.Time `json:"date"`
    Status    string    `json:"status"`
    Notes     string    `json:"notes,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
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


