package models

import "time"

type DashboardStats struct {
	TotalHorses      int64   `json:"totalHorses"`
	PregnantMares    int64   `json:"pregnantMares"`
	TotalExpenses    float64 `json:"totalExpenses"`
	UpcomingFoalings int64   `json:"upcomingFoalings"`
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

type UserDashboardStats struct {
    TotalHorses           int64 `json:"total_horses"`
    ActivePregnancies     int64 `json:"active_pregnancies"`
    RecentHealthRecords   int64 `json:"recent_health_records"`
    RecentBreedingRecords int64 `json:"recent_breeding_records"`
}
