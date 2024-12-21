package models

import "time"

type User struct {
	ID             int64     `json:"id" db:"id"`
	Email          string    `json:"email" db:"email"`
	HashedPassword string    `json:"-" db:"hashed_password"`
	LastSync       time.Time `json:"last_sync" db:"last_sync"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type SyncData struct {
	UserID    int64           `json:"user_id"`
	Timestamp time.Time       `json:"timestamp"`
	Horses    []Horse         `json:"horses"`
	Health    []HealthRecord  `json:"health"`
	Events    []PregnancyEvent `json:"events"`
}
