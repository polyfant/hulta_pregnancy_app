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
	UserID         string     `json:"user_id" validate:"required"`
	Name           string     `json:"name" validate:"required,min=1,max=100"`
	Breed          string     `json:"breed" gorm:"type:varchar(100)" validate:"omitempty,max=100"`
	Gender         Gender     `json:"gender" validate:"required,oneof=MARE STALLION GELDING"`
	BirthDate      time.Time  `json:"birth_date" validate:"required"`
	Weight         float64    `json:"weight" validate:"omitempty,gt=0"`
	Height         float64    `json:"height" validate:"omitempty,gt=0"`
	Color          string     `json:"color" gorm:"size:100" validate:"omitempty,max=100"`
	IsPregnant     bool       `json:"is_pregnant" gorm:"default:false"`
	ConceptionDate *time.Time `json:"conception_date" validate:"omitempty,required_if=IsPregnant true"`
	MotherId       *uint      `json:"mother_id,omitempty"`
	FatherId       *uint      `json:"father_id,omitempty"`
	ExternalMother string     `json:"external_mother,omitempty" validate:"omitempty,max=100"`
	ExternalFather string     `json:"external_father,omitempty" validate:"omitempty,max=100"`
	
	// Owner information
	OwnerName    string `json:"owner_name" gorm:"type:varchar(100)" validate:"omitempty,max=100"`
	OwnerContact string `json:"owner_contact" gorm:"type:varchar(100)" validate:"omitempty,max=100"`
	OwnerEmail   string `json:"owner_email" gorm:"type:varchar(100)" validate:"omitempty,email,max=100"`
	OwnerPhone   string `json:"owner_phone" gorm:"type:varchar(20)" validate:"omitempty,max=20"`
	
	LastVetCheck   *time.Time `json:"last_vet_check,omitempty"`
	LastHeatDate   *time.Time `json:"last_heat_date,omitempty"`
	CycleLength    int        `json:"cycle_length,omitempty" validate:"omitempty,gt=0"`
	
	Mother         *Horse `json:"-"`
	Father         *Horse `json:"-"`
	Pregnancies    []Pregnancy `json:"-"`
	HealthRecords  []HealthRecord `json:"-"`
	BreedingCosts  []BreedingCost `json:"-"`
	
	Notes          string     `json:"notes,omitempty" gorm:"type:text" validate:"omitempty,max=5000"`
	CreatedAt      time.Time  `json:"created_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	
	DamID          *uint      `json:"dam_id,omitempty"`
	SireID         *uint      `json:"sire_id,omitempty"`
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
		return h.Gender == GenderMare && h.ConceptionDate != nil && !h.ConceptionDate.IsZero() && h.ConceptionDate.Before(time.Now())
	}
	return true
}

type FamilyTree struct {
	Horse     *Horse   `json:"horse"`
	Mother    *Horse   `json:"mother,omitempty"`
	Father    *Horse   `json:"father,omitempty"`
	Offspring []*Horse `json:"offspring,omitempty"`
}
