package models

import "time"

type HealthRecord struct {
    ID          uint      `json:"id,omitempty" gorm:"primaryKey"`
    HorseID     uint      `json:"horse_id" gorm:"not null" validate:"required"`
    UserID      string    `json:"user_id" gorm:"index;not null" validate:"required"`
    Type        string    `json:"type" gorm:"size:50" validate:"required,max=50"`
    Date        time.Time `json:"date" validate:"required"`
    Description string    `json:"description" validate:"required,max=2000"`
    CreatedAt   time.Time `json:"created_at,omitempty"`
    UpdatedAt   time.Time `json:"updated_at,omitempty"`
} 