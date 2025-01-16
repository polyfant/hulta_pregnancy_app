package vitals

import (
	"context"
	"math"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/websocket"
)

// Service defines the interface for vital signs operations
type Service interface {
	RecordVitalSigns(ctx context.Context, vitals *models.VitalSigns) error
	GetVitalSigns(ctx context.Context, horseID uint, from, to time.Time) ([]*models.VitalSigns, error)
	GetLatestVitalSigns(ctx context.Context, horseID uint) (*models.VitalSigns, error)
	GetAlerts(ctx context.Context, horseID uint, includeAcknowledged bool) ([]*models.VitalSignsAlert, error)
	GetAlert(ctx context.Context, alertID uint) (*models.VitalSignsAlert, error)
	AcknowledgeAlert(ctx context.Context, alertID uint) error
	GetTrends(ctx context.Context, horseID uint, from, to time.Time) ([]*models.VitalSignsTrend, error)
}

// Repository defines the interface for data access operations
type Repository interface {
	SaveVitalSigns(ctx context.Context, vitals *models.VitalSigns) error
	GetVitalSigns(ctx context.Context, horseID uint, from, to time.Time) ([]*models.VitalSigns, error)
	GetLatestVitalSigns(ctx context.Context, horseID uint) (*models.VitalSigns, error)
	GetAlerts(ctx context.Context, horseID uint, includeAcknowledged bool) ([]*models.VitalSignsAlert, error)
	GetAlert(ctx context.Context, alertID uint) (*models.VitalSignsAlert, error)
	AcknowledgeAlert(ctx context.Context, alertID uint) error
	SaveAlert(ctx context.Context, alert *models.VitalSignsAlert) error
	SaveTrend(ctx context.Context, trend *models.VitalSignsTrend) error
}

// service implements the Service interface
type service struct {
	repo Repository
	hub  *websocket.Hub
}

// NewService creates a new vital signs service
func NewService(repo Repository, hub *websocket.Hub) Service {
	return &service{
		repo: repo,
		hub:  hub,
	}
}

// RecordVitalSigns records new vital signs and checks for alerts
func (s *service) RecordVitalSigns(ctx context.Context, vitals *models.VitalSigns) error {
	if err := s.repo.SaveVitalSigns(ctx, vitals); err != nil {
		return err
	}

	if err := s.checkForAlerts(ctx, vitals); err != nil {
		return err
	}

	s.broadcastUpdate(vitals)
	return nil
}

// GetVitalSigns retrieves vital signs history for a horse
func (s *service) GetVitalSigns(ctx context.Context, horseID uint, from, to time.Time) ([]*models.VitalSigns, error) {
	return s.repo.GetVitalSigns(ctx, horseID, from, to)
}

// GetLatestVitalSigns retrieves the most recent vital signs for a horse
func (s *service) GetLatestVitalSigns(ctx context.Context, horseID uint) (*models.VitalSigns, error) {
	return s.repo.GetLatestVitalSigns(ctx, horseID)
}

// GetAlerts retrieves alerts for a horse
func (s *service) GetAlerts(ctx context.Context, horseID uint, includeAcknowledged bool) ([]*models.VitalSignsAlert, error) {
	return s.repo.GetAlerts(ctx, horseID, includeAcknowledged)
}

// GetAlert retrieves a specific alert
func (s *service) GetAlert(ctx context.Context, alertID uint) (*models.VitalSignsAlert, error) {
	return s.repo.GetAlert(ctx, alertID)
}

// AcknowledgeAlert marks an alert as acknowledged
func (s *service) AcknowledgeAlert(ctx context.Context, alertID uint) error {
	return s.repo.AcknowledgeAlert(ctx, alertID)
}

// GetTrends retrieves and analyzes vital signs trends
func (s *service) GetTrends(ctx context.Context, horseID uint, from, to time.Time) ([]*models.VitalSignsTrend, error) {
	// Retrieve vital signs data
	vitals, err := s.repo.GetVitalSigns(ctx, horseID, from, to)
	if err != nil {
		return nil, err
	}

	// Prepare data for analysis
	var (
		temperatures []float64
		heartRates  []float64
		respirations []float64
		timestamps  []time.Time
	)

	for _, v := range vitals {
		temperatures = append(temperatures, v.Temperature)
		heartRates = append(heartRates, float64(v.HeartRate))
		respirations = append(respirations, float64(v.Respiration))
		timestamps = append(timestamps, v.RecordedAt)
	}

	// Analyze trends for each metric
	trends := []*models.VitalSignsTrend{
		s.analyzeTrend("temperature", temperatures, timestamps),
		s.analyzeTrend("heart_rate", heartRates, timestamps),
		s.analyzeTrend("respiratory_rate", respirations, timestamps),
	}

	// Save trends
	for _, trend := range trends {
		trend.HorseID = horseID
		trend.CreatedAt = time.Now()
		if err := s.repo.SaveTrend(ctx, trend); err != nil {
			return nil, err
		}
	}

	return trends, nil
}

// analyzeTrend analyzes a single metric for trends
func (s *service) analyzeTrend(metricType string, values []float64, timestamps []time.Time) *models.VitalSignsTrend {
	if len(values) < 2 {
		return nil
	}

	// Calculate trend direction and magnitude
	startValue := values[0]
	endValue := values[len(values)-1]
	change := endValue - startValue
	magnitude := math.Abs(change)

	var direction string
	switch {
	case change > 0:
		direction = "increasing"
	case change < 0:
		direction = "decreasing"
	default:
		direction = "stable"
	}

	return &models.VitalSignsTrend{
		MetricType: metricType,
		Direction:  direction,
		Magnitude:  magnitude,
		StartTime:  timestamps[0],
		EndTime:    timestamps[len(timestamps)-1],
	}
}

// checkForAlerts checks vital signs against thresholds and creates alerts if needed
func (s *service) checkForAlerts(ctx context.Context, v *models.VitalSigns) error {
	thresholds := models.DefaultThresholds()

	// Check temperature
	if v.Temperature < thresholds.TemperatureMin {
		alert := &models.VitalSignsAlert{
			HorseID:   v.HorseID,
			AlertType: "temperature_low",
			Type:      "vital_signs",
			Message:   "Temperature is below normal range",
			Severity:  "warning",
			Parameter: "temperature",
			Value:     v.Temperature,
			CreatedAt: time.Now(),
		}
		if err := s.repo.SaveAlert(ctx, alert); err != nil {
			return err
		}
	} else if v.Temperature > thresholds.TemperatureMax {
		alert := &models.VitalSignsAlert{
			HorseID:   v.HorseID,
			AlertType: "temperature_high",
			Type:      "vital_signs",
			Message:   "Temperature is above normal range",
			Severity:  "warning",
			Parameter: "temperature",
			Value:     v.Temperature,
			CreatedAt: time.Now(),
		}
		if err := s.repo.SaveAlert(ctx, alert); err != nil {
			return err
		}
	}

	// Check heart rate
	if int(v.HeartRate) < thresholds.HeartRateMin {
		alert := &models.VitalSignsAlert{
			HorseID:   v.HorseID,
			AlertType: "heart_rate_low",
			Type:      "vital_signs",
			Message:   "Heart rate is below normal range",
			Severity:  "warning",
			Parameter: "heart_rate",
			Value:     float64(v.HeartRate),
			CreatedAt: time.Now(),
		}
		if err := s.repo.SaveAlert(ctx, alert); err != nil {
			return err
		}
	} else if int(v.HeartRate) > thresholds.HeartRateMax {
		alert := &models.VitalSignsAlert{
			HorseID:   v.HorseID,
			AlertType: "heart_rate_high",
			Type:      "vital_signs",
			Message:   "Heart rate is above normal range",
			Severity:  "warning",
			Parameter: "heart_rate",
			Value:     float64(v.HeartRate),
			CreatedAt: time.Now(),
		}
		if err := s.repo.SaveAlert(ctx, alert); err != nil {
			return err
		}
	}

	// Check respiration
	if int(v.Respiration) < thresholds.RespirationMin {
		alert := &models.VitalSignsAlert{
			HorseID:   v.HorseID,
			AlertType: "respiration_low",
			Type:      "vital_signs",
			Message:   "Respiration rate is below normal range",
			Severity:  "warning",
			Parameter: "respiration",
			Value:     float64(v.Respiration),
			CreatedAt: time.Now(),
		}
		if err := s.repo.SaveAlert(ctx, alert); err != nil {
			return err
		}
	} else if int(v.Respiration) > thresholds.RespirationMax {
		alert := &models.VitalSignsAlert{
			HorseID:   v.HorseID,
			AlertType: "respiration_high",
			Type:      "vital_signs",
			Message:   "Respiration rate is above normal range",
			Severity:  "warning",
			Parameter: "respiration",
			Value:     float64(v.Respiration),
			CreatedAt: time.Now(),
		}
		if err := s.repo.SaveAlert(ctx, alert); err != nil {
			return err
		}
	}

	return nil
}

// broadcastUpdate broadcasts vital signs updates to connected clients
func (s *service) broadcastUpdate(vitals *models.VitalSigns) {
	if s.hub != nil {
		s.hub.Broadcast(vitals)
	}
}
