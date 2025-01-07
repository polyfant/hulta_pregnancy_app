package health

import (
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type HealthService struct {
	healthRepo repository.HealthRepository
}

func NewHealthService(healthRepo repository.HealthRepository) *HealthService {
	return &HealthService{
		healthRepo: healthRepo,
	}
}

// ... rest of the file
