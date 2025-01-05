package models

import "time"

type Horse struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         string    `gorm:"type:text;not null;index"`
	Name           string    `gorm:"type:varchar(100);not null"`
	Breed           string    `gorm:"type:varchar(100)"`
	Gender         string    `gorm:"type:varchar(50)"`
	BirthDate      time.Time
	Weight         float64
	Height         float64
	Color          string    `gorm:"size:100"`
	IsPregnant     bool      `gorm:"default:false"`
	ConceptionDate *time.Time
	MotherId       *uint
	FatherId       *uint
	ExternalMother string
	ExternalFather string
	
	LastVetCheck   *time.Time
	LastHeatDate   *time.Time
	CycleLength    int
	
	Mother         *Horse
	Father         *Horse
	Pregnancies    []Pregnancy
	HealthRecords  []HealthRecord
	BreedingCosts  []BreedingCost
	
	Notes          string    `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`
}

type HorseDetails struct {
	Horse     Horse       `json:"horse"`
	Expenses  []Expense   `json:"expenses"`
}

const (
	GenderMare     = "MARE"
	GenderStallion = "STALLION"
	GenderGelding  = "GELDING"
)

func (h *Horse) Age() int {
	if h.BirthDate.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - h.BirthDate.Year()
	
	if now.Month() < h.BirthDate.Month() || 
	   (now.Month() == h.BirthDate.Month() && now.Day() < h.BirthDate.Day()) {
		age--
	}
	return age
}

func (h *Horse) IsBreedingAge() bool {
	age := h.Age()
	switch h.Gender {
	case GenderMare:
		return age >= 3 && age <= 20
	case GenderStallion:
		return age >= 3 && age <= 25
	default:
		return false
	}
}

func (h *Horse) CanBreed() bool {
	return h.IsBreedingAge() && (h.Gender == GenderMare || h.Gender == GenderStallion)
}

func (h *Horse) DaysPregnant() int {
	if !h.IsPregnant || h.ConceptionDate == nil {
		return 0
	}
	return int(time.Since(*h.ConceptionDate).Hours() / 24)
}

func (h *Horse) ExpectedFoalingDate() *time.Time {
	if !h.IsPregnant || h.ConceptionDate == nil {
		return nil
	}
	foalingDate := h.ConceptionDate.AddDate(0, 0, 340)
	return &foalingDate
}

func (h *Horse) ValidateGender() bool {
	return h.Gender == GenderMare || h.Gender == GenderStallion || h.Gender == GenderGelding
}

func (h *Horse) ValidatePregnancy() bool {
	if h.IsPregnant {
		return h.Gender == GenderMare && h.ConceptionDate != nil
	}
	return true
}
