package models

import "time"

type User struct {
	ID             int64      `json:"id" db:"id"`
	Email          string     `json:"email" db:"email"`
	HashedPassword string     `json:"-" db:"hashed_password"`
	LastSync       time.Time  `json:"last_sync" db:"last_sync"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin      *time.Time `json:"last_login" db:"last_login"`
	IsActive       bool       `json:"is_active" db:"is_active"`
}

type SyncData struct {
	UserID    int64           `json:"user_id"`
	Timestamp time.Time       `json:"timestamp"`
	Horses    []Horse         `json:"horses"`
	Health    []HealthRecord  `json:"health"`
	Events    []PregnancyEvent `json:"events"`
}

type UserDashboard struct {
	TotalHorses          int    `json:"total_horses"`
	PregnantHorses       int    `json:"pregnant_horses"`
	TotalExpenses        string `json:"total_expenses"`
	UpcomingEvents       int    `json:"upcoming_events"`
	ActivePregnancies    int    `json:"active_pregnancies"`
}
