package pregnancy

import (
	"fmt"
	"time"


	"github.com/polyfant/horse_tracking/internal/models"
)

// internalPregnancyGuidelines contains detailed information about each stage of pregnancy
// This is an internal type used by the service to provide more detailed guidelines
type internalPregnancyGuidelines struct {
	Stage         models.PregnancyStage
	DayRange      string
	Monitoring    []string
	Nutrition     []string
	Exercise      []string
	Risks         []string
	Preparation   []string
	Progress      string
	DaysRemaining string
	Priority      string
}

var pregnancyStageGuidelines = map[models.PregnancyStage]internalPregnancyGuidelines{
	models.EarlyGestation: {
		Stage:    models.EarlyGestation,
		DayRange: "0-113 days",
		Monitoring: []string{
			"Watch for signs of continued estrus or return to heat",
			"Monitor general health and appetite",
			"Check for any vaginal discharge",
			"Observe for signs of colic or discomfort",
		},
		Nutrition: []string{
			"Maintain regular diet if mare is in good condition",
			"Ensure access to clean, fresh water",
			"Provide high-quality forage",
			"Continue regular mineral supplementation",
		},
		Exercise: []string{
			"Continue normal exercise routine",
			"Avoid strenuous activities",
			"Allow regular turnout",
		},
		Risks: []string{
			"Early embryonic death (highest risk first 42 days)",
			"Twin pregnancies",
			"Infections",
		},
		Preparation: []string{
			"Schedule early pregnancy check (14-16 days)",
			"Plan for follow-up ultrasound (28-30 days)",
			"Consider checking for twins",
		},
	},
	models.MidGestation: {
		Stage:    models.MidGestation,
		DayRange: "114-226 days",
		Monitoring: []string{
			"Regular body condition scoring",
			"Monitor for any signs of illness",
			"Watch for changes in behavior",
			"Check udder development (shouldn't be significant yet)",
		},
		Nutrition: []string{
			"Gradually increase feed quality",
			"Ensure adequate protein intake",
			"Maintain proper calcium:phosphorus ratio",
			"Consider adding omega-3 supplements",
		},
		Exercise: []string{
			"Maintain moderate exercise",
			"Continue daily turnout",
			"Avoid high-intensity work",
		},
		Risks: []string{
			"Placental issues",
			"Nutritional deficiencies",
			"Stress-related complications",
		},
		Preparation: []string{
			"Plan vaccination schedule",
			"Consider deworming program",
			"Begin planning for foaling arrangements",
		},
	},
	models.LateGestation: {
		Stage:    models.LateGestation,
		DayRange: "227-310 days",
		Monitoring: []string{
			"Watch for udder development",
			"Monitor body condition closely",
			"Check for edema",
			"Observe fetal movement",
			"Watch for any signs of premature labor",
		},
		Nutrition: []string{
			"Increase feed quantity (by 30-50%)",
			"Provide high-quality protein",
			"Ensure adequate mineral intake",
			"Consider adding vitamin E supplement",
		},
		Exercise: []string{
			"Reduce exercise intensity",
			"Maintain light activity",
			"Ensure safe turnout conditions",
		},
		Risks: []string{
			"Premature labor",
			"Placentitis",
			"Colic",
			"Nutritional imbalances",
		},
		Preparation: []string{
			"Prepare foaling area",
			"Assemble foaling kit",
			"Review foaling procedures",
			"Have veterinarian contacts ready",
		},
	},
	models.PreFoaling: {
		Stage:    models.PreFoaling,
		DayRange: "311-340 days",
		Monitoring: []string{
			"Check udder development twice daily",
			"Monitor for waxing of teats",
			"Watch for relaxation of pelvic ligaments",
			"Observe for restlessness or nesting behavior",
			"Check vulvar area for changes",
			"Monitor temperature twice daily",
		},
		Nutrition: []string{
			"Maintain increased feed levels",
			"Ensure easy access to fresh water",
			"Consider adding electrolytes",
			"Feed smaller meals more frequently",
		},
		Exercise: []string{
			"Light hand walking only",
			"Provide safe turnout in small area",
			"Avoid any strenuous activity",
		},
		Risks: []string{
			"Dystocia (difficult birth)",
			"Red bag delivery",
			"Premature placental separation",
			"Post-foaling complications",
		},
		Preparation: []string{
			"Have foaling kit ready and accessible",
			"Keep watch schedule organized",
			"Install monitoring cameras if using",
			"Have veterinarian on standby",
			"Prepare mare and foal IDs",
		},
	},
}

// internalPreFoalingSign is an internal type used by the service
type internalPreFoalingSign struct {
	Name        string
	Description string
	Urgency     int
	TimeToFoal  string
}

var preFoalingSigns = []internalPreFoalingSign{
	{
		Name:        "Udder Development",
		Description: "Gradual filling of the udder, usually begins 4-6 weeks before foaling",
		Urgency:     1,
		TimeToFoal:  "4-6 weeks",
	},
	{
		Name:        "Vulvar Relaxation",
		Description: "Relaxation and lengthening of the vulva",
		Urgency:     2,
		TimeToFoal:  "1-2 weeks",
	},
	{
		Name:        "Waxing",
		Description: "Waxy secretions on teat ends",
		Urgency:     4,
		TimeToFoal:  "12-72 hours",
	},
	{
		Name:        "Milk Dripping",
		Description: "Active dripping of milk",
		Urgency:     5,
		TimeToFoal:  "12-24 hours",
	},
	{
		Name:        "Restlessness",
		Description: "Frequent lying down and getting up, pawing, looking at sides",
		Urgency:     5,
		TimeToFoal:  "1-4 hours",
	},
	{
		Name:        "Sweating",
		Description: "Patchy sweating, particularly on neck and flanks",
		Urgency:     5,
		TimeToFoal:  "1-4 hours",
	},
}

type Service struct {
	db models.DataStore
}

func (s *Service) GetPregnancyGuidelines(horse models.Horse) (any, error) {
	if !horse.IsPregnant {
		return nil, fmt.Errorf("horse is not pregnant")
	}

	guidelines, err := s.GetPregnancyGuidelinesForHorse(horse)
	if err != nil {
		return nil, err
	}

	return guidelines, nil
}

func (s *Service) StartPregnancy(horseID int64, conceptionDate time.Time) error {
	// Check if horse exists and is a mare
	horse, err := s.db.GetHorse(horseID)
	if err != nil {
		return fmt.Errorf("failed to get horse: %w", err)
	}

	if horse.Gender != models.GenderMare {
		return fmt.Errorf("only mares can be pregnant")
	}

	// Add conception event
	event := &models.PregnancyEvent{
		HorseID:     horseID,
		Type:        models.EventConception,
		Description: "Pregnancy started",
		Date:        conceptionDate,
	}

	if err := s.db.AddPregnancyEvent(event); err != nil {
		return fmt.Errorf("failed to add conception event: %w", err)
	}

	// Update horse pregnancy status
	if err := s.db.UpdateHorsePregnancyStatus(horseID, true, conceptionDate); err != nil {
		return fmt.Errorf("failed to update horse pregnancy status: %w", err)
	}

	return nil
}

func (s *Service) EndPregnancy(horseID int64, outcome string, notes string) error {
	// Check if horse is pregnant
	horse, err := s.db.GetHorse(horseID)
	if err != nil {
		return fmt.Errorf("failed to get horse: %w", err)
	}

	if !horse.IsPregnant {
		return fmt.Errorf("horse is not pregnant")
	}

	// Add foaling event
	event := &models.PregnancyEvent{
		HorseID:     horseID,
		Type:        models.EventFoaling,
		Description: outcome,
		Notes:       notes,
		Date:        time.Now(),
	}

	if err := s.db.AddPregnancyEvent(event); err != nil {
		return fmt.Errorf("failed to add foaling event: %w", err)
	}

	// Update horse pregnancy status
	if err := s.db.UpdateHorsePregnancyStatus(horseID, false, time.Time{}); err != nil {
		return fmt.Errorf("failed to update horse pregnancy status: %w", err)
	}

	return nil
}

func (s *Service) GetPregnancyEvents(horseID int64) ([]models.PregnancyEvent, error) {
	return s.db.GetPregnancyEvents(horseID)
}

func (s *Service) AddPregnancyEvent(horseID int64, eventType string, description string, notes string) error {
	event := &models.PregnancyEvent{
		HorseID:     horseID,
		Type:        eventType,
		Description: description,
		Notes:       notes,
		Date:        time.Now(),
	}

	return s.db.AddPregnancyEvent(event)
}

func (s *Service) RecordPreFoalingSign(horseID int64, signName string) error {
	// Find the sign in our predefined list
	var foundSign *internalPreFoalingSign
	for _, sign := range preFoalingSigns {
		if sign.Name == signName {
			foundSign = &sign
			break
		}
	}

	if foundSign == nil {
		return fmt.Errorf("invalid pre-foaling sign: %s", signName)
	}

	// Record the sign
	signRecord := models.PreFoalingSign{
		HorseID:      horseID,
		Name:         signName,
		Observed:     true,
		DateObserved: time.Now(),
	}

	// Save the pre-foaling sign
	err := s.db.AddPreFoalingSign(&signRecord)
	if err != nil {
		return fmt.Errorf("error recording pre-foaling sign: %v", err)
	}

	// Add the sign to pregnancy events for tracking
	event := &models.PregnancyEvent{
		HorseID:     horseID,
		Type:        "PRE_FOALING_SIGN",
		Description: fmt.Sprintf("Pre-foaling sign observed: %s", signName),
		Date:        time.Now(),
	}

	if err := s.db.AddPregnancyEvent(event); err != nil {
		return fmt.Errorf("failed to record pre-foaling sign event: %w", err)
	}

	return nil
}

func (s *Service) GetPregnancyStage(horse models.Horse) models.PregnancyStage {
	if horse.ConceptionDate == nil {
		return models.EarlyGestation
	}

	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)

	switch {
	case daysPregnant <= 113:
		return models.EarlyGestation
	case daysPregnant <= 225:
		return models.MidGestation
	case daysPregnant <= 310:
		return models.LateGestation
	case daysPregnant <= 330:
		return models.PreFoaling
	default:
		return models.Foaling
	}
}

func (s *Service) CheckPreFoalingSigns(horse models.Horse) []internalPreFoalingSign {
	if horse.ConceptionDate == nil {
		return nil
	}

	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)
	if daysPregnant < 300 {
		return nil
	}

	// Return all signs to watch for in late pregnancy
	return preFoalingSigns
}

func (s *Service) GetFoalingChecklist() []string {
	return []string{
		"Clean, bedded foaling stall or paddock",
		"Working phone and emergency contact numbers",
		"Flashlight with fresh batteries",
		"Clean towels and sheets",
		"Tail wrap",
		"Sterile scissors",
		"Clean string or umbilical clamps",
		"Iodine for umbilical stump",
		"Clean bucket and warm water",
		"Watch or clock for timing stages",
		"Veterinarian's phone number",
		"Camera for documentation",
		"Notebook and pen",
		"Clean halter and lead rope for mare",
		"Small foal halter",
		"Obstetrical lubricant",
		"Disposable gloves",
		"Clean receptacle for placenta",
	}
}

func (s *Service) GetPostFoalingChecklist() []string {
	return []string{
		"Mare appears bright and alert",
		"Foal is breathing normally",
		"Foal attempts to rise within 30 minutes",
		"Foal stands within 1-2 hours",
		"Foal nurses within 2-3 hours",
		"Mare has passed complete placenta within 3 hours",
		"No excessive bleeding from mare",
		"Foal has passed meconium",
		"Mare's temperature is normal",
		"Both mare and foal are bonding well",
		"Umbilical stump has been dipped in iodine",
		"Colostrum quality has been checked",
		"IgG levels checked in foal (12-18 hours after birth)",
	}
}

func (s *Service) StartPregnancyTracking(horseID int64, conceptionDate time.Time) error {
	// Update horse's conception date
	err := s.db.UpdateHorseConceptionDate(horseID, conceptionDate)
	if err != nil {
		return fmt.Errorf("error updating conception date: %v", err)
	}

	// Add initial pregnancy event
	event := &models.PregnancyEvent{
		HorseID:     horseID,
		Date:        time.Now(),
		Type:        models.EventConception,
		Description: "Pregnancy tracking started",
		Notes:       fmt.Sprintf("Conception date recorded as %s", conceptionDate.Format("2006-01-02")),
	}

	return s.db.AddPregnancyEvent(event)
}

func (s *Service) EndPregnancyTracking(horseID int64, outcome string) error {
	// Clear conception date
	err := s.db.UpdateHorseConceptionDate(horseID, time.Time{})
	if err != nil {
		return fmt.Errorf("error clearing conception date: %v", err)
	}

	// Add final pregnancy event
	event := &models.PregnancyEvent{
		HorseID:     horseID,
		Date:        time.Now(),
		Type:        models.EventFoaling,
		Description: "Pregnancy tracking ended",
		Notes:       fmt.Sprintf("Outcome: %s", outcome),
	}

	return s.db.AddPregnancyEvent(event)
}

func (s *Service) GetPregnancyGuidelinesByStage(stage models.PregnancyStage) ([]models.PregnancyGuideline, error) {
	guidelines := []models.PregnancyGuideline{
		{
			Stage:       string(models.EarlyGestation),
			Title:       "Early Gestation (0-113 days)",
			Description: "Critical period for embryonic development.",
			Recommendations: []string{
				"Maintain regular exercise routine",
				"Continue normal feeding schedule",
				"Schedule initial pregnancy check",
			},
			Warnings: []string{
				"Watch for signs of pregnancy loss",
				"Monitor for twin pregnancies",
			},
			Checkpoints: []string{
				"14-16 day pregnancy check",
				"45 day ultrasound",
				"Monitor appetite and behavior",
			},
		},
		{
			Stage:       string(models.MidGestation),
			Title:       "Mid Gestation (114-225 days)",
			Description: "Period of steady fetal growth.",
			Recommendations: []string{
				"Maintain balanced nutrition",
				"Continue moderate exercise",
				"Schedule routine checkups",
			},
			Warnings: []string{
				"Watch for signs of discomfort",
				"Monitor body condition",
			},
			Checkpoints: []string{
				"Monthly health checks",
				"Vaccination updates",
				"Deworming schedule",
			},
		},
		{
			Stage:       string(models.LateGestation),
			Title:       "Late Gestation (226-310 days)",
			Description: "Period of rapid fetal growth.",
			Recommendations: []string{
				"Increase feed quality",
				"Reduce strenuous exercise",
				"Prepare foaling area",
			},
			Warnings: []string{
				"Watch for signs of premature labor",
				"Monitor udder development",
			},
			Checkpoints: []string{
				"Weekly health checks",
				"Prepare foaling kit",
				"Monitor vital signs",
			},
		},
		{
			Stage:       string(models.PreFoaling),
			Title:       "Pre-Foaling (311-330 days)",
			Description: "Final preparation for foaling.",
			Recommendations: []string{
				"Monitor closely",
				"Prepare foaling area",
				"Have vet contact ready",
			},
			Warnings: []string{
				"Watch for early labor signs",
				"Monitor udder changes",
			},
			Checkpoints: []string{
				"Daily health checks",
				"Monitor temperature",
				"Check udder development",
			},
		},
		{
			Stage:       string(models.Foaling),
			Title:       "Foaling (331+ days)",
			Description: "Active preparation for imminent birth.",
			Recommendations: []string{
				"Monitor 24/7",
				"Have veterinarian on call",
				"Prepare for immediate post-foaling care",
			},
			Warnings: []string{
				"Watch for prolonged labor",
				"Monitor mare and foal bonding",
			},
			Checkpoints: []string{
				"Time contractions",
				"Ensure proper presentation",
				"Check placenta after birth",
			},
		},
	}

	if stage == models.EarlyGestation {
		return []models.PregnancyGuideline{guidelines[0]}, nil
	} else if stage == models.MidGestation {
		return []models.PregnancyGuideline{guidelines[1]}, nil
	} else if stage == models.LateGestation {
		return []models.PregnancyGuideline{guidelines[2]}, nil
	} else if stage == models.PreFoaling {
		return []models.PregnancyGuideline{guidelines[3]}, nil
	} else if stage == models.Foaling {
		return []models.PregnancyGuideline{guidelines[4]}, nil
	}

	return nil, fmt.Errorf("invalid pregnancy stage: %s", stage)
}

func (s *Service) GetPregnancyGuidelinesForHorse(horse models.Horse) (*internalPregnancyGuidelines, error) {
	if horse.ConceptionDate == nil {
		return nil, fmt.Errorf("horse is not pregnant")
	}

	calculator := NewPregnancyCalculator(*horse.ConceptionDate)
	stage := calculator.GetStage()
	schedule := calculator.GetMonitoringSchedule()
	currentDay := calculator.GetCurrentDay()
	daysRemaining := calculator.GetRemainingDays()
	progress := calculator.GetProgressPercentage()

	// Get base guidelines for the stage
	guidelines := pregnancyStageGuidelines[stage]

	// Add dynamic monitoring instructions based on the schedule
	if schedule.TemperatureCheck {
		guidelines.Monitoring = append(guidelines.Monitoring,
			fmt.Sprintf("Check temperature every %d hours", schedule.CheckFrequency))
	}
	if schedule.UdderCheck {
		guidelines.Monitoring = append(guidelines.Monitoring,
			"Monitor udder development for signs of waxing")
	}
	if schedule.VulvaCheck {
		guidelines.Monitoring = append(guidelines.Monitoring,
			"Check vulva for relaxation and color changes")
	}

	// Add progress information
	guidelines.Progress = fmt.Sprintf("Day %d of pregnancy (%.1f%% complete)",
		currentDay, progress)
	guidelines.DaysRemaining = fmt.Sprintf("%d days until estimated due date",
		daysRemaining)
	guidelines.Priority = schedule.Priority

	return &guidelines, nil
}

func (s *Service) getGuidelinesForStage(stage models.PregnancyStage) models.PregnancyGuidelines {
	baseGuidelines := pregnancyStageGuidelines[stage]
	return models.PregnancyGuidelines{
		Stage:            string(stage),
		NutritionGuide:   baseGuidelines.Nutrition,
		ExerciseGuide:    baseGuidelines.Exercise,
		HealthChecks:     baseGuidelines.Monitoring,
		WarningSignsList: baseGuidelines.Risks,
		Preparations:     baseGuidelines.Preparation,
	}
}

func (s *Service) getNextMilestones(currentStage models.PregnancyStage) []string {
	milestones := map[models.PregnancyStage][]string{
		models.EarlyGestation: {
			"14-16 day pregnancy check",
			"45 day ultrasound",
			"First trimester completion",
		},
		models.MidGestation: {
			"Vaccination updates",
			"Deworming schedule",
			"Second trimester completion",
		},
		models.LateGestation: {
			"Begin foaling preparation",
			"Weekly health checks",
			"Prepare foaling kit",
		},
		models.PreFoaling: {
			"Daily health monitoring",
			"Watch for udder development",
			"Monitor temperature changes",
		},
		models.Foaling: {
			"Watch for imminent labor signs",
			"Monitor contractions",
			"Post-foaling care",
		},
	}

	return milestones[currentStage]
}

func NewService(db models.DataStore) *Service {
	return &Service{db: db}
}

func (s *Service) GetPregnancyStatus(horseID int64) (*models.PregnancyStatus, error) {
	horse, err := s.db.GetHorse(horseID)
	if err != nil {
		return nil, fmt.Errorf("error getting horse: %v", err)
	}

	if horse.ConceptionDate == nil {
		return &models.PregnancyStatus{
			IsPregnant: false,
		}, nil
	}

	now := time.Now()
	daysInPregnancy := int(now.Sub(*horse.ConceptionDate).Hours() / 24)
	expectedDueDate := horse.ConceptionDate.AddDate(0, 0, 340) // Average horse pregnancy is 340 days

	// Get the current stage based on days in pregnancy
	var currentStage models.PregnancyStage
	switch {
	case daysInPregnancy <= 113:
		currentStage = models.EarlyGestation
	case daysInPregnancy <= 225:
		currentStage = models.MidGestation
	case daysInPregnancy <= 310:
		currentStage = models.LateGestation
	case daysInPregnancy <= 330:
		currentStage = models.PreFoaling
	default:
		currentStage = models.Foaling
	}

	// Get the latest event
	events, err := s.db.GetPregnancyEvents(horse.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting pregnancy events: %v", err)
	}

	var lastEvent *models.PregnancyEvent
	if len(events) > 0 {
		lastEvent = &events[len(events)-1]
	}

	// Get guidelines for the current stage
	guidelines := s.getGuidelinesForStage(currentStage)

	return &models.PregnancyStatus{
		IsPregnant:       true,
		ConceptionDate:   *horse.ConceptionDate,
		CurrentStage:     string(currentStage),
		DaysInPregnancy:  daysInPregnancy,
		ExpectedDueDate:  expectedDueDate,
		LastEvent:        lastEvent,
		NextMilestones:   s.getNextMilestones(currentStage),
		UpcomingChecks:   guidelines.HealthChecks,
		NutritionGuide:   guidelines.NutritionGuide,
		ExerciseGuide:    guidelines.ExerciseGuide,
		WarningSignsList: guidelines.WarningSignsList,
	}, nil
}
