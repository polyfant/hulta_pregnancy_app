package models

import "time"

type DashboardStats struct {
	TotalHorses      int                `json:"totalHorses"`
	PregnantMares    int                `json:"pregnantMares"`
	UpcomingDueDates []PregnancySummary `json:"upcomingDueDates"`
	RecentHealth     []HealthSummary    `json:"recentHealth"`
	Costs            CostSummary        `json:"costs"`
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
