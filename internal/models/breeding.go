package models

import "time"

type BreedingRecord struct {
    ID        uint      `json:"id,omitempty" gorm:"primaryKey"`
    HorseID   uint      `json:"horse_id" validate:"required"`
    UserID    string    `json:"user_id" validate:"required"`
    Date      time.Time `json:"date" validate:"required"`
    Status    string    `json:"status" validate:"required,max=50"`
    Notes     string    `json:"notes,omitempty" validate:"omitempty,max=5000"`
    CreatedAt time.Time `json:"created_at,omitempty"`
    UpdatedAt time.Time `json:"updated_at,omitempty"`
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


