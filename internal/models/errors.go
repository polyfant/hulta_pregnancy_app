package models

import "errors"

var (
	// Expense errors
	ErrInvalidAmount      = errors.New("expense amount must be non-negative")
	ErrInvalidExpenseType = errors.New("invalid expense type")
	ErrInvalidFrequency   = errors.New("invalid frequency type")
)
