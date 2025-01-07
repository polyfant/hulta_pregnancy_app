package api

import "time"

// Common types used across handlers
type ErrorResponse struct {
    Error string `json:"error"`
}

type PregnancyProgress struct {
    DueDate       time.Time `json:"due_date"`
    Progress      float64   `json:"progress"`
    DaysRemaining int      `json:"days_remaining"`
    Stage         string    `json:"stage"`
} 