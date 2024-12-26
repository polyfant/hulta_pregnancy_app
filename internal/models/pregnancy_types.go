package models

import (
	"time"
)

// PregnancyStage represents different stages of pregnancy
type PregnancyStage string

const (
	EarlyGestation    PregnancyStage = "EARLY_GESTATION"    // 0-114 days
	MidGestation      PregnancyStage = "MID_GESTATION"      // 115-225 days
	LateGestation     PregnancyStage = "LATE_GESTATION"     // 226-310 days
	PreFoaling        PregnancyStage = "PRE_FOALING"        // 311-330 days
	Foaling           PregnancyStage = "FOALING"            // 331+ days
)

// PregnancyEvent represents a significant event during pregnancy
type PregnancyEvent struct {
	ID          int64     `json:"id" db:"id"`
	HorseID     int64     `json:"horseId" db:"horse_id"`
	Date        time.Time `json:"date" db:"date"`
	Type        string    `json:"type" db:"type"`
	Description string    `json:"description" db:"description"`
	Notes       string    `json:"notes,omitempty" db:"notes"`
}

// PregnancyStatus represents the current pregnancy status of a horse
type PregnancyStatus struct {
	IsPregnant       bool            `json:"isPregnant"`
	ConceptionDate   time.Time       `json:"conceptionDate,omitempty"`
	CurrentStage     string          `json:"currentStage,omitempty"`
	DaysInPregnancy  int             `json:"daysInPregnancy,omitempty"`
	ExpectedDueDate  time.Time       `json:"expectedDueDate,omitempty"`
	LastEvent        *PregnancyEvent `json:"lastEvent,omitempty"`
	NextMilestones   []string        `json:"nextMilestones,omitempty"`
	UpcomingChecks   []string        `json:"upcomingChecks,omitempty"`
	NutritionGuide   []string        `json:"nutritionGuide,omitempty"`
	ExerciseGuide    []string        `json:"exerciseGuide,omitempty"`
	WarningSignsList []string        `json:"warningSignsList,omitempty"`
}

// PreFoalingSign represents a sign of impending foaling
type PreFoalingSign struct {
	ID           int64     `json:"id" db:"id"`
	HorseID      int64     `json:"horseId" db:"horse_id"`
	Name         string    `json:"name" db:"name"`
	Observed     bool      `json:"observed" db:"observed"`
	DateObserved time.Time `json:"dateObserved,omitempty" db:"date_observed"`
	Notes        string    `json:"notes,omitempty" db:"notes"`
}

// PregnancyGuideline contains comprehensive guidelines for each stage
type PregnancyGuideline struct {
	Stage           string   `json:"stage"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Recommendations []string `json:"recommendations"`
	Warnings        []string `json:"warnings"`
	Checkpoints     []string `json:"checkpoints"`
}

// PregnancyGuidelines contains guidelines for a specific stage
type PregnancyGuidelines struct {
	Stage            string   `json:"stage"`
	NutritionGuide   []string `json:"nutritionGuide"`
	ExerciseGuide    []string `json:"exerciseGuide"`
	HealthChecks     []string `json:"healthChecks"`
	WarningSignsList []string `json:"warningSignsList"`
	Preparations     []string `json:"preparations"`
}

// Common pregnancy event types
const (
	EventConception     = "CONCEPTION"
	EventVetCheck      = "VET_CHECK"
	EventVaccination   = "VACCINATION"
	EventComplication  = "COMPLICATION"
	EventPreFoaling    = "PRE_FOALING"
	EventFoaling       = "FOALING"
	EventMiscarriage   = "MISCARRIAGE"
	EventAbortion      = "ABORTION"
)

// StandardMilestones defines standard pregnancy milestones with their descriptions
var StandardMilestones = []PregnancyGuideline{
	{
		Stage:           string(EarlyGestation),
		Title:           "Early Pregnancy",
		Description:     "Early Pregnancy",
		Recommendations: []string{},
		Warnings:        []string{},
		Checkpoints:     []string{},
	},
	{
		Stage:           string(MidGestation),
		Title:           "Middle Pregnancy",
		Description:     "Middle Pregnancy",
		Recommendations: []string{},
		Warnings:        []string{},
		Checkpoints:     []string{},
	},
	{
		Stage:           string(LateGestation),
		Title:           "Late Pregnancy",
		Description:     "Late Pregnancy",
		Recommendations: []string{},
		Warnings:        []string{},
		Checkpoints:     []string{},
	},
	{
		Stage:           string(PreFoaling),
		Title:           "Near Term",
		Description:     "Near Term",
		Recommendations: []string{},
		Warnings:        []string{},
		Checkpoints:     []string{},
	},
	{
		Stage:           string(Foaling),
		Title:           "Imminent Foaling",
		Description:     "Imminent Foaling",
		Recommendations: []string{},
		Warnings:        []string{},
		Checkpoints:     []string{},
	},
}
