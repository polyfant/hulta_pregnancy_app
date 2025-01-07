package service

import "time"

const (
    DefaultGestationDays = 340
)

func CalculateDueDate(conceptionDate time.Time) time.Time {
    return conceptionDate.AddDate(0, 0, DefaultGestationDays)
} 