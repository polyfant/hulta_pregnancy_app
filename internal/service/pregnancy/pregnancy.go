package pregnancy

import (
	"fmt"
	"time"

	"github.com/polyfant/horse_tracking/internal/logger"
	"github.com/polyfant/horse_tracking/internal/models"
)

// PregnancyGuidelines contains detailed information about each stage of pregnancy
type PregnancyGuidelines struct {
	Stage       models.PregnancyStage
	DayRange    string
	Monitoring  []string
	Nutrition   []string
	Exercise    []string
	Risks       []string
	Preparation []string
	Progress    string
	DaysRemaining string
	Priority    string
}

var pregnancyStageGuidelines = map[models.PregnancyStage]PregnancyGuidelines{
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
	models.FinalGestation: {
		Stage:    models.FinalGestation,
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

type PreFoalingSign struct {
	Name        string
	Description string
	Urgency     int // 1-5, with 5 being most urgent
	TimeToFoal  string
}

var PreFoalingSigns = []PreFoalingSign{
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

func NewService(db models.DataStore) *Service {
	return &Service{db: db}
}

func (s *Service) GetPregnancyGuidelines(horse models.Horse) (*PregnancyGuidelines, error) {
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

func (s *Service) GetPregnancyStage(horse models.Horse) models.PregnancyStage {
	if horse.ConceptionDate == nil {
		return ""
	}

	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)

	switch {
	case daysPregnant < 114:
		return models.EarlyGestation
	case daysPregnant < 226:
		return models.MidGestation
	case daysPregnant < 310:
		return models.LateGestation
	default:
		return models.FinalGestation
	}
}

func (s *Service) CheckPreFoalingSigns(horse models.Horse) []PreFoalingSign {
	if horse.ConceptionDate == nil {
		return nil
	}

	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)
	if daysPregnant < 300 {
		return nil
	}

	// Return all signs to watch for in late pregnancy
	return PreFoalingSigns
}

func (s *Service) RecordPreFoalingSign(horseID int64, signName string) error {
	// Record the observation in pregnancy events
	event := &models.PregnancyEvent{
		HorseID:     horseID,
		Date:        time.Now(),
		Description: fmt.Sprintf("Pre-foaling sign observed: %s", signName),
	}

	if err := s.db.AddPregnancyEvent(event); err != nil {
		return fmt.Errorf("failed to record pre-foaling sign: %w", err)
	}

	// Log the observation
	logger.Info("Recorded pre-foaling sign", map[string]interface{}{
		"horseID": horseID,
		"sign":    signName,
	})

	return nil
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
