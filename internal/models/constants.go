package models

// Priority represents the importance level
type Priority string

const (
	PriorityHigh   Priority = "HIGH"
	PriorityMedium Priority = "MEDIUM"
	PriorityLow    Priority = "LOW"
)

// EventType represents different types of events
type EventType string

const (
	EventFoaling     EventType = "FOALING"
	EventVetCheck    EventType = "VET_CHECK"
	EventUltrasound  EventType = "ULTRASOUND"
	EventVaccination EventType = "VACCINATION"
	EventDeworming   EventType = "DEWORMING"
)

// BreedingStatus represents the status of breeding
type BreedingStatus string

const (
	BreedingStatusActive    BreedingStatus = "ACTIVE"
	BreedingStatusCompleted BreedingStatus = "COMPLETED"
	BreedingStatusFailed    BreedingStatus = "FAILED"
	BreedingStatusCancelled BreedingStatus = "CANCELLED"
)

// HealthRecordType represents the type of health record
type HealthRecordType string

const (
	HealthRecordTypeVetVisit    HealthRecordType = "VET_VISIT"
	HealthRecordTypeVaccination HealthRecordType = "VACCINATION"
	HealthRecordTypeDeworming   HealthRecordType = "DEWORMING"
	HealthRecordTypeDental      HealthRecordType = "DENTAL"
	HealthRecordTypeInjury      HealthRecordType = "INJURY"
	HealthRecordTypeOther       HealthRecordType = "OTHER"
)

// Default values
const (
	DefaultGestationDays = 340 // Average horse pregnancy duration
) 