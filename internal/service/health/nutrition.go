package health

import (
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// NutritionGuidelines represents feeding recommendations
type NutritionGuidelines struct {
	Category           string
	DailyFeedPercent  float64 // Percent of body weight
	HayPercent        float64 // Percent of total feed
	GrainPercent      float64 // Percent of total feed
	Supplements       []string
	FeedingFrequency  int // Times per day
	SpecialNotes      []string
}

var NutritionRecommendations = map[string]NutritionGuidelines{
	"Maintenance": {
		Category:          "Adult Maintenance",
		DailyFeedPercent: 1.5,
		HayPercent:       100,
		GrainPercent:     0,
		Supplements: []string{
			"Salt block",
			"Mineral supplement if needed based on forage analysis",
		},
		FeedingFrequency: 2,
		SpecialNotes: []string{
			"Adjust based on work level",
			"Monitor body condition",
			"Ensure constant access to fresh water",
			"Provide salt block",
		},
	},
	"Performance": {
		Category:          "Working Horse",
		DailyFeedPercent: 2.0,
		HayPercent:       70,
		GrainPercent:     30,
		Supplements: []string{
			"Electrolytes",
			"Joint supplements if needed",
			"Vitamin E and Selenium",
		},
		FeedingFrequency: 3,
		SpecialNotes: []string{
			"Feed according to work intensity",
			"Provide hay before grain",
			"Consider adding electrolytes after heavy work",
			"Monitor hydration",
		},
	},
	"EarlyPregnancy": {
		Category:          "Early Pregnancy",
		DailyFeedPercent: 1.5,
		HayPercent:       100,
		GrainPercent:     0,
		Supplements: []string{
			"Prenatal vitamin supplement",
			"Calcium supplement",
			"Omega-3 fatty acids",
		},
		FeedingFrequency: 2,
		SpecialNotes: []string{
			"Maintain good body condition",
			"Focus on quality forage",
			"Monitor for any feed aversions",
			"Ensure balanced mineral intake",
		},
	},
	"LatePregnancy": {
		Category:          "Late Pregnancy",
		DailyFeedPercent: 2.0,
		HayPercent:       80,
		GrainPercent:     20,
		Supplements: []string{
			"Enhanced prenatal supplement",
			"Additional calcium",
			"Vitamin E",
			"Omega-3 fatty acids",
		},
		FeedingFrequency: 4,
		SpecialNotes: []string{
			"Increase feed gradually",
			"Monitor body condition closely",
			"Feed smaller meals more frequently",
			"Ensure adequate protein intake",
		},
	},
	"Lactating": {
		Category:          "Lactating Mare",
		DailyFeedPercent: 2.5,
		HayPercent:       60,
		GrainPercent:     40,
		Supplements: []string{
			"Calcium supplement",
			"Protein supplement",
			"Vitamin and mineral mix",
			"Omega-3 fatty acids",
		},
		FeedingFrequency: 4,
		SpecialNotes: []string{
			"Highest nutritional demands",
			"Monitor milk production",
			"Ensure adequate water intake",
			"Watch for weight loss",
		},
	},
}

// BodyConditionScore represents the Henneke body condition scoring system
type BodyConditionScore struct {
	Score       int
	Description string
	Indicators  []string
	Action      string
}

var BodyConditionScores = []BodyConditionScore{
	{
		Score:       1,
		Description: "Poor",
		Indicators: []string{
			"Prominent spinous processes, ribs, tailhead, and hip bones",
			"No fatty tissue can be felt",
			"Severe concavity under tail",
		},
		Action: "Immediate veterinary attention needed. Gradually increase caloric intake.",
	},
	{
		Score:       3,
		Description: "Thin",
		Indicators: []string{
			"Slight fat covering over ribs",
			"Spinous processes easily visible",
			"Tailhead prominent but individual vertebrae not visible",
		},
		Action: "Increase feed quality and quantity. Monitor for underlying health issues.",
	},
	{
		Score:       5,
		Description: "Moderate (Ideal)",
		Indicators: []string{
			"Ribs not visible but easily felt",
			"Back is level",
			"Fat around tailhead beginning to feel spongy",
		},
		Action: "Maintain current feeding program. Regular monitoring.",
	},
	{
		Score:       7,
		Description: "Fleshy",
		Indicators: []string{
			"Ribs felt with firm pressure",
			"Fat filling between ribs",
			"Noticeable fat deposits along withers and behind shoulders",
		},
		Action: "Reduce caloric intake. Increase exercise if possible.",
	},
	{
		Score:       9,
		Description: "Extremely Fat",
		Indicators: []string{
			"Ribs not palpable",
			"Massive fat deposits over spinous processes",
			"Deep positive crease along back",
		},
		Action: "Immediate diet restriction needed. Risk of metabolic issues.",
	},
}

// WaterRequirements calculates daily water needs
type WaterRequirements struct {
	BaseAmount     float64 // Liters per day
	ActivityAdd    float64 // Additional liters based on activity
	WeatherAdd     float64 // Additional liters based on weather
	PregnancyAdd   float64 // Additional liters if pregnant
	LactationAdd   float64 // Additional liters if lactating
}

func calculateWaterNeeds(horse models.Horse, temperature float64, isActive bool) WaterRequirements {
	// Base calculation: 5-7% of body weight
	baseAmount := horse.Weight * 0.06 // Using 6% as average

	var req WaterRequirements
	req.BaseAmount = baseAmount

	// Activity adjustment
	if isActive {
		req.ActivityAdd = baseAmount * 0.2 // 20% increase for active horses
	}

	// Weather adjustment (hot weather)
	if temperature > 25 { // Celsius
		req.WeatherAdd = baseAmount * 0.3 // 30% increase in hot weather
	}

	// Pregnancy adjustment
	if horse.ConceptionDate != nil {
		daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)
		if daysPregnant > 226 {
			req.PregnancyAdd = baseAmount * 0.15 // 15% increase in late pregnancy
		}
	}

	return req
}

func (s *Service) GetNutritionPlan(horse models.Horse) struct {
	Guidelines    NutritionGuidelines
	WaterNeeds    WaterRequirements
	Supplements   []string
	SpecialNotes  []string
} {
	var guidelines NutritionGuidelines

	// Determine appropriate nutrition category
	if horse.ConceptionDate != nil {
		daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)
		if daysPregnant > 226 {
			guidelines = NutritionRecommendations["LatePregnancy"]
		} else {
			guidelines = NutritionRecommendations["EarlyPregnancy"]
		}
	} else {
		guidelines = NutritionRecommendations["Maintenance"]
	}

	// Calculate water requirements (assuming moderate temperature and activity)
	waterNeeds := calculateWaterNeeds(horse, 20, false)

	// Additional supplements based on individual needs
	var supplements []string
	supplements = append(supplements, guidelines.Supplements...)

	// Add season-specific notes
	currentMonth := time.Now().Month()
	var specialNotes []string
	specialNotes = append(specialNotes, guidelines.SpecialNotes...)

	switch {
	case currentMonth >= time.June && currentMonth <= time.August:
		specialNotes = append(specialNotes,
			"Monitor for heat stress",
			"Consider electrolyte supplementation",
			"Provide access to shade",
			"Check water temperature",
		)
	case currentMonth >= time.December && currentMonth <= time.February:
		specialNotes = append(specialNotes,
			"Increase hay for warmth",
			"Ensure water isn't frozen",
			"Monitor body condition more frequently",
			"Consider adding warm mashes",
		)
	}

	return struct {
		Guidelines    NutritionGuidelines
		WaterNeeds    WaterRequirements
		Supplements   []string
		SpecialNotes  []string
	}{
		Guidelines:   guidelines,
		WaterNeeds:   waterNeeds,
		Supplements:  supplements,
		SpecialNotes: specialNotes,
	}
}
