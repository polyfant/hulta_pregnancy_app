package service

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type GrowthService interface {
	RecordGrowthMeasurement(ctx context.Context, foalID uint, weight, height float64) error
	GetFoalGrowthData(ctx context.Context, foalID uint) ([]models.GrowthData, error)
	AnalyzeGrowthTrends(ctx context.Context, foalID uint) (*GrowthAnalysis, error)
}

type GrowthServiceImpl struct {
	growthRepo repository.GrowthRepository
	horseRepo  repository.HorseRepository
}

type GrowthAnalysis struct {
	AverageWeightGain float64 `json:"averageWeightGain"`
	AverageHeightGain float64 `json:"averageHeightGain"`
	ProjectedWeight   float64 `json:"projectedWeight"`
	ProjectedHeight   float64 `json:"projectedHeight"`
	GrowthStatus      string  `json:"growthStatus"`
}

func NewGrowthService(growthRepo repository.GrowthRepository, horseRepo repository.HorseRepository) *GrowthServiceImpl {
	return &GrowthServiceImpl{
		growthRepo: growthRepo,
		horseRepo:  horseRepo,
	}
}

func (s *GrowthServiceImpl) RecordGrowthMeasurement(ctx context.Context, foalID uint, weight, height float64) error {
	// Validate horse exists
	horse, err := s.horseRepo.GetByID(ctx, foalID)
	if err != nil {
		return fmt.Errorf("invalid foal ID: %w", err)
	}

	// Get existing growth data to calculate age
	existingData, _ := s.growthRepo.GetGrowthDataByFoalID(ctx, foalID)
	age := len(existingData)

	// Determine expected values based on breed (simplified)
	expectedWeight := calculateExpectedWeight(horse.Breed, age)
	expectedHeight := calculateExpectedHeight(horse.Breed, age)

	growthData := &models.GrowthData{
		FoalID:          foalID,
		Age:             age,
		Weight:          weight,
		Height:          height,
		ExpectedWeight:  expectedWeight,
		ExpectedHeight:  expectedHeight,
		MeasurementDate: time.Now(),
	}

	return s.growthRepo.CreateGrowthData(ctx, growthData)
}

func (s *GrowthServiceImpl) GetFoalGrowthData(ctx context.Context, foalID uint) ([]models.GrowthData, error) {
	return s.growthRepo.GetGrowthDataByFoalID(ctx, foalID)
}

func (s *GrowthServiceImpl) AnalyzeGrowthTrends(ctx context.Context, foalID uint) (*GrowthAnalysis, error) {
	growthData, err := s.growthRepo.GetGrowthDataByFoalID(ctx, foalID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve growth data: %w", err)
	}

	if len(growthData) < 2 {
		return nil, fmt.Errorf("insufficient growth data for analysis")
	}

	// Calculate average gains
	totalWeightGain := growthData[len(growthData)-1].Weight - growthData[0].Weight
	totalHeightGain := growthData[len(growthData)-1].Height - growthData[0].Height
	averageWeightGain := totalWeightGain / float64(len(growthData)-1)
	averageHeightGain := totalHeightGain / float64(len(growthData)-1)

	// Simple growth status determination
	growthStatus := "Normal"
	if averageWeightGain < 0.5 {
		growthStatus = "Slow"
	} else if averageWeightGain > 1.5 {
		growthStatus = "Rapid"
	}

	// Naive projection (linear extrapolation)
	lastData := growthData[len(growthData)-1]
	projectedWeight := lastData.Weight + (averageWeightGain * 2)
	projectedHeight := lastData.Height + (averageHeightGain * 2)

	return &GrowthAnalysis{
		AverageWeightGain: averageWeightGain,
		AverageHeightGain: averageHeightGain,
		ProjectedWeight:   projectedWeight,
		ProjectedHeight:   projectedHeight,
		GrowthStatus:      growthStatus,
	}, nil
}

// Simplified breed-specific growth expectation calculations
func calculateExpectedWeight(breed string, age int) float64 {
	switch breed {
	case "Thoroughbred":
		return 50 + (float64(age) * 1.2)
	case "Warmblood":
		return 55 + (float64(age) * 1.5)
	case "Arabian":
		return 45 + (float64(age) * 1.0)
	default:
		return 50 + (float64(age) * 1.2)
	}
}

func calculateExpectedHeight(breed string, age int) float64 {
	switch breed {
	case "Thoroughbred":
		return 100 + (float64(age) * 0.8)
	case "Warmblood":
		return 110 + (float64(age) * 1.0)
	case "Arabian":
		return 95 + (float64(age) * 0.7)
	default:
		return 100 + (float64(age) * 0.8)
	}
}
