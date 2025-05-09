package models

import "time"

type Gender string

const (
	GenderMare     Gender = "MARE"
	GenderStallion Gender = "STALLION"
	GenderGelding  Gender = "GELDING"
)

type Horse struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	UserID         string     `json:"user_id"`
	Name           string     `json:"name"`
	Breed          string     `gorm:"type:varchar(100)"`
	Gender         Gender     `json:"gender"`
	BirthDate      time.Time  `json:"birth_date"`
	Weight         float64
	Height         float64
	Color          string     `gorm:"size:100"`
	IsPregnant     bool       `gorm:"default:false"`
	ConceptionDate *time.Time
	MotherId       *uint
	FatherId       *uint
	ExternalMother string
	ExternalFather string
	
	// Owner information
	OwnerName    string `json:"owner_name" gorm:"type:varchar(100)"`
	OwnerContact string `json:"owner_contact" gorm:"type:varchar(100)"`
	OwnerEmail   string `json:"owner_email" gorm:"type:varchar(100)"`
	OwnerPhone   string `json:"owner_phone" gorm:"type:varchar(20)"`
	
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
	
	DamID          *uint      `json:"dam_id"`
	SireID         *uint      `json:"sire_id"`
}

type HorseDetails struct {
	Horse     Horse       `json:"horse"`
	Expenses  []Expense   `json:"expenses"`
}

type HorseFamilyTree struct {
    Horse     *Horse   `json:"horse"`
    Parents   []Horse  `json:"parents"`
    Offspring []Horse  `json:"offspring"`
    Siblings  []Horse  `json:"siblings"`
}

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
	if h.BirthDate.IsZero() {
		return false
	}
	minBreedingAge := h.BirthDate.AddDate(3, 0, 0) // 3 years for both mares and stallions
	return time.Now().After(minBreedingAge)
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

type FamilyTree struct {
	Horse     *Horse   `json:"horse"`
	Mother    *Horse   `json:"mother,omitempty"`
	Father    *Horse   `json:"father,omitempty"`
	Offspring []*Horse `json:"offspring,omitempty"`
}
