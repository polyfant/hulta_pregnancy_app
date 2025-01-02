package models

import "time"

type DashboardStats struct {
	TotalHorses      int     `json:"totalHorses"`
	PregnantMares    int     `json:"pregnantMares"`
	TotalExpenses    float64 `json:"totalExpenses"`
	UpcomingFoalings int     `json:"upcomingFoalings"`
}

type PregnancySummary struct {
	HorseID      int64     `json:"horseId"`
	HorseName    string    `json:"horseName"`
	DueDate      time.Time `json:"dueDate"`
	DaysRemaining int      `json:"daysRemaining"`
}

type HealthSummary struct {
	HorseID   int64     `json:"horseId"`
	HorseName string    `json:"horseName"`
	Type      string    `json:"type"`
	Date      time.Time `json:"date"`
	Notes     string    `json:"notes"`
}

type CostSummary struct {
	TotalCosts     float64 `json:"totalCosts"`
	MonthlyAverage float64 `json:"monthlyAverage"`
	RecentCosts    float64 `json:"recentCosts"` // Last 30 days
}
