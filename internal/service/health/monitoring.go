package health

import (
	"time"

	"github.com/polyfant/horse_tracking/internal/models"
)

// VitalSigns represents normal vital sign ranges for horses
type VitalSigns struct {
	Temperature    Range[float64]
	HeartRate     Range[int]
	RespirationRate Range[int]
	Description    string
	Category      string // Adult, Foal, Pregnant, etc.
}

type Range[T any] struct {
	Min T
	Max T
}

var VitalSignRanges = map[string]VitalSigns{
	"Adult": {
		Temperature:     Range[float64]{37.5, 38.5}, // Celsius
		HeartRate:      Range[int]{28, 44},         // Beats per minute
		RespirationRate: Range[int]{8, 16},          // Breaths per minute
		Description:     "Normal ranges for healthy adult horses",
		Category:       "Adult",
	},
	"Foal": {
		Temperature:     Range[float64]{37.8, 38.9},
		HeartRate:      Range[int]{80, 120},
		RespirationRate: Range[int]{20, 40},
		Description:     "Normal ranges for newborn foals (0-2 weeks)",
		Category:       "Foal",
	},
	"PregnantLate": {
		Temperature:     Range[float64]{37.5, 38.5},
		HeartRate:      Range[int]{34, 50},
		RespirationRate: Range[int]{10, 20},
		Description:     "Normal ranges for mares in late pregnancy",
		Category:       "Pregnant",
	},
}

// HealthAlert represents different types of health concerns
type HealthAlert struct {
	Category    string
	Symptoms    []string
	Urgency     int // 1-5, with 5 being most urgent
	Action      string
	Prevention  []string
}

var HealthAlerts = map[string]HealthAlert{
	"Colic": {
		Category: "Digestive",
		Symptoms: []string{
			"Pawing at ground",
			"Looking at flanks",
			"Rolling or attempting to roll",
			"Lack of appetite",
			"Reduced or no manure production",
			"Elevated heart rate",
		},
		Urgency: 5,
		Action:  "Contact veterinarian immediately. Remove feed. Walk horse if safe to do so.",
		Prevention: []string{
			"Regular dental care",
			"Consistent feeding schedule",
			"Fresh, clean water available",
			"Gradual feed changes",
			"Regular parasite control",
		},
	},
	"Laminitis": {
		Category: "Hoof",
		Symptoms: []string{
			"Reluctance to move",
			"Shifting weight",
			"Heat in hooves",
			"Strong digital pulse",
			"Characteristic stance (leaning back)",
		},
		Urgency: 4,
		Action:  "Contact veterinarian. Keep horse still. Apply cold therapy if acute.",
		Prevention: []string{
			"Proper diet management",
			"Limited access to rich grass",
			"Regular exercise",
			"Weight management",
			"Regular hoof care",
		},
	},
	"Respiratory": {
		Category: "Respiratory",
		Symptoms: []string{
			"Coughing",
			"Nasal discharge",
			"Increased breathing rate",
			"Exercise intolerance",
			"Abnormal breathing sounds",
		},
		Urgency: 3,
		Action:  "Monitor temperature and breathing rate. Contact vet if symptoms persist.",
		Prevention: []string{
			"Dust-free environment",
			"Good ventilation",
			"Clean bedding",
			"Regular vaccination",
			"Avoid dusty feed",
		},
	},
}

// VaccinationSchedule represents recommended vaccination timelines
type VaccinationSchedule struct {
	VaccineName string
	Frequency   string
	Important   bool
	Notes       string
}

var RecommendedVaccinations = []VaccinationSchedule{
	{
		VaccineName: "Tetanus",
		Frequency:   "Annually",
		Important:   true,
		Notes:       "Essential for all horses. Additional booster if injury occurs.",
	},
	{
		VaccineName: "Equine Influenza",
		Frequency:   "Every 6-12 months",
		Important:   true,
		Notes:       "More frequent for competition horses or those in contact with many other horses.",
	},
	{
		VaccineName: "Rhinopneumonitis (EHV)",
		Frequency:   "Every 6 months",
		Important:   true,
		Notes:       "Especially important for pregnant mares at 5, 7, and 9 months of gestation.",
	},
	{
		VaccineName: "West Nile Virus",
		Frequency:   "Annually",
		Important:   true,
		Notes:       "Before mosquito season. Additional booster may be needed in high-risk areas.",
	},
}

// DewormingSchedule represents recommended deworming protocols
type DewormingSchedule struct {
	Type     string
	Interval string
	Notes    []string
}

var DewormingProtocol = DewormingSchedule{
	Type:     "Strategic Deworming",
	Interval: "Based on fecal egg counts",
	Notes: []string{
		"Perform fecal egg count every 3-6 months",
		"Deworm based on results and veterinary recommendations",
		"Consider seasonal timing for specific parasites",
		"Keep records of products used and dates",
		"Monitor effectiveness through follow-up testing",
	},
}

// DentalCareSchedule represents recommended dental care intervals
type DentalCareSchedule struct {
	Age           string
	CheckInterval string
	Signs         []string
}

var DentalCareGuidelines = []DentalCareSchedule{
	{
		Age:           "Young Horse (2-5 years)",
		CheckInterval: "Every 6 months",
		Signs: []string{
			"Quidding (dropping feed)",
			"Weight loss",
			"Head tossing when ridden",
			"Difficulty maintaining bit contact",
			"Bad breath",
		},
	},
	{
		Age:           "Adult Horse (5-15 years)",
		CheckInterval: "Annually",
		Signs: []string{
			"Uneven wear patterns",
			"Difficulty chewing",
			"Feed packing in cheeks",
			"Resistance to bridle",
		},
	},
	{
		Age:           "Senior Horse (15+ years)",
		CheckInterval: "Every 6 months",
		Signs: []string{
			"Loss of condition",
			"Difficulty eating hay",
			"Slow eating",
			"Feed falling from mouth",
			"Excessive salivation",
		},
	},
}

type Service struct {
	db models.DataStore
}

// GetHealthAssessment provides a comprehensive health evaluation
func (s *Service) GetHealthAssessment(horse models.Horse) struct {
	VitalSignsCategory string
	Vaccinations      []VaccinationSchedule
	DentalCare        DentalCareSchedule
	Deworming         DewormingSchedule
	SpecialCare       []string
} {
	// Determine vital signs category
	category := "Adult"
	if horse.ConceptionDate != nil && time.Since(*horse.ConceptionDate) > 226*24*time.Hour {
		category = "PregnantLate"
	}

	// Determine dental care schedule based on age
	var dentalCare DentalCareSchedule
	age := time.Since(horse.DateOfBirth).Hours() / (24 * 365)
	switch {
	case age < 5:
		dentalCare = DentalCareGuidelines[0]
	case age < 15:
		dentalCare = DentalCareGuidelines[1]
	default:
		dentalCare = DentalCareGuidelines[2]
	}

	// Determine special care needs
	var specialCare []string
	if horse.ConceptionDate != nil {
		specialCare = append(specialCare, 
			"Regular pregnancy checks",
			"Modified exercise routine",
			"Increased nutritional requirements",
			"Monitor for pregnancy-related complications",
		)
	}

	return struct {
		VitalSignsCategory string
		Vaccinations      []VaccinationSchedule
		DentalCare        DentalCareSchedule
		Deworming         DewormingSchedule
		SpecialCare       []string
	}{
		VitalSignsCategory: category,
		Vaccinations:      RecommendedVaccinations,
		DentalCare:        dentalCare,
		Deworming:         DewormingProtocol,
		SpecialCare:       specialCare,
	}
}
