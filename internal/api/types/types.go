package types

// ErrorResponse represents an error response
type ErrorResponse struct {
    Error string `json:"error"`
}

// DashboardStats represents the statistics shown on the dashboard
type DashboardStats struct {
    TotalHorses       int     `json:"total_horses"`
    PregnantHorses    int     `json:"pregnant_horses"`
    TotalExpenses     float64 `json:"total_expenses"`
    UpcomingEvents    int     `json:"upcoming_events"`
    ActivePregnancies int     `json:"active_pregnancies"`
}

// APIResponse represents a standard API response
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string     `json:"error,omitempty"`
} 