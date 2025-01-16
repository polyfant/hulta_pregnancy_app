package models

// Priority represents the importance level
type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
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

// Vital signs normal ranges for horses
const (
	// Temperature ranges (in Celsius)
	MinNormalTemperature = 37.5
	MaxNormalTemperature = 38.5
	MinSafeTemperature  = 35.0 // Below this requires immediate attention
	MaxSafeTemperature  = 42.0 // Above this requires immediate attention

	// Heart rate ranges (beats per minute)
	MinNormalHeartRate = 28
	MaxNormalHeartRate = 44
	MinSafeHeartRate  = 20 // Below this requires immediate attention
	MaxSafeHeartRate  = 100 // Above this requires immediate attention

	// Respiration ranges (breaths per minute)
	MinNormalRespiration = 8
	MaxNormalRespiration = 16
	MinSafeRespiration  = 4  // Below this requires immediate attention
	MaxSafeRespiration  = 40 // Above this requires immediate attention
)

// Default values
const (
	DefaultGestationDays = 340 // Average horse pregnancy duration
)