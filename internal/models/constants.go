package models

// Pregnancy-related constants
type PregnancyStage string

const (
    PregnancyStageUnknown       PregnancyStage = "UNKNOWN"
    PregnancyStageEarlyGestation PregnancyStage = "EARLY_GESTATION"
    PregnancyStageMidGestation   PregnancyStage = "MID_GESTATION"
    PregnancyStageLateGestation  PregnancyStage = "LATE_GESTATION"
    PregnancyStageOverdue       PregnancyStage = "OVERDUE"
    PregnancyStageHighRisk      PregnancyStage = "HIGH_RISK"

    PregnancyStatusActive    string = "ACTIVE"
    PregnancyStatusCompleted string = "COMPLETED"
    PregnancyStatusAborted   string = "ABORTED"

    DefaultGestationDays = 340 // Average horse pregnancy duration
)

// Priority levels
type Priority string

const (
    PriorityHigh   Priority = "HIGH"
    PriorityMedium Priority = "MEDIUM"
    PriorityLow    Priority = "LOW"
)

// Horse-related constants
type Gender string

const (
    GenderMare     Gender = "MARE"
    GenderStallion Gender = "STALLION"
)

// Event types
type EventType string

const (
    EventFoaling     EventType = "FOALING"
    EventVetCheck    EventType = "VET_CHECK"
    EventUltrasound  EventType = "ULTRASOUND"
    EventVaccination EventType = "VACCINATION"
    EventDeworming   EventType = "DEWORMING"
)

type BreedingStatusType string

const (
    BreedingStatusActive     BreedingStatusType = "ACTIVE"
    BreedingStatusCompleted  BreedingStatusType = "COMPLETED"
    BreedingStatusCancelled  BreedingStatusType = "CANCELLED"
) 