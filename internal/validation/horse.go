package validation

import (
	"errors"
	"time"
	"github.com/polyfant/horse_tracking/internal/models"
)

func ValidateHorse(horse *models.Horse) error {
	if horse.Name == "" {
		return errors.New("horse name is required")
	}
	if horse.Breed == "" {
		return errors.New("horse breed is required")
	}
	if !horse.DateOfBirth.IsZero() && horse.DateOfBirth.After(time.Now()) {
		return errors.New("date of birth cannot be in the future")
	}
	if horse.ConceptionDate != nil {
		if horse.ConceptionDate.After(time.Now()) {
			return errors.New("conception date cannot be in the future")
		}
		// Validate pregnancy duration
		if time.Since(*horse.ConceptionDate) > 365*24*time.Hour {
			return errors.New("invalid conception date: pregnancy duration exceeds one year")
		}
	}
	return nil
}

func ValidateHealthRecord(record *models.HealthRecord) error {
	if record.HorseID <= 0 {
		return errors.New("invalid horse ID")
	}
	if record.Type == "" {
		return errors.New("health record type is required")
	}
	if record.Date.IsZero() {
		return errors.New("health record date is required")
	}
	if record.Date.After(time.Now()) {
		return errors.New("health record date cannot be in the future")
	}
	return nil
}

func ValidatePregnancyEvent(event *models.PregnancyEvent) error {
	if event.HorseID <= 0 {
		return errors.New("invalid horse ID")
	}
	if event.Description == "" {
		return errors.New("event description is required")
	}
	if event.Date.IsZero() {
		return errors.New("event date is required")
	}
	if event.Date.After(time.Now()) {
		return errors.New("event date cannot be in the future")
	}
	return nil
}
