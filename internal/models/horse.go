package models

import (
	"fmt"
	"time"
)

type Horse struct {
	ID             int64      `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	Breed          string     `json:"breed" db:"breed"`
	Gender         string     `json:"gender" db:"gender"`
	DateOfBirth    time.Time  `json:"dateOfBirth" db:"date_of_birth"`
	Weight         float64    `json:"weight" db:"weight"`
	IsPregnant     bool       `json:"isPregnant" db:"is_pregnant"`
	ConceptionDate *time.Time `json:"conceptionDate,omitempty" db:"conception_date"`
	MotherID       *int64     `json:"motherId,omitempty" db:"mother_id"`
	FatherID       *int64     `json:"fatherId,omitempty" db:"father_id"`
	ExternalMother string     `json:"externalMother,omitempty" db:"external_mother"`
	ExternalFather string     `json:"externalFather,omitempty" db:"external_father"`
	Age            string     `json:"age,omitempty" db:"-"` // Calculated field
}

// Gender constants
const (
	GenderMare     = "MARE"
	GenderStallion = "STALLION"
	GenderGelding  = "GELDING"
)

// CalculateAge returns the horse's age in years and months
func (h *Horse) CalculateAge(now time.Time) string {
	years := now.Year() - h.DateOfBirth.Year()
	months := int(now.Month() - h.DateOfBirth.Month())

	if months < 0 {
		years--
		months += 12
	}

	if years > 0 {
		if months > 0 {
			return fmt.Sprintf("%d years, %d months", years, months)
		}
		return fmt.Sprintf("%d years", years)
	}
	return fmt.Sprintf("%d months", months)
}

type HealthRecord struct {
	ID      int64     `json:"id" db:"id"`
	HorseID int64     `json:"horseId" db:"horse_id"`
	Date    time.Time `json:"date" db:"date"`
	Type    string    `json:"type" db:"type"`
	Notes   string    `json:"notes" db:"notes"`
}

type FamilyTree struct {
	Horse          Horse           `json:"horse"`
	Mother         *FamilyMember   `json:"mother,omitempty"`
	Father         *FamilyMember   `json:"father,omitempty"`
	Offspring      []FamilyMember  `json:"offspring,omitempty"`
	Siblings       []FamilyMember  `json:"siblings,omitempty"`
}

type FamilyMember struct {
	ID             int64     `json:"id,omitempty"`
	Name           string    `json:"name"`
	Breed          string    `json:"breed,omitempty"`
	Gender         string    `json:"gender,omitempty"`
	DateOfBirth    time.Time `json:"dateOfBirth,omitempty"`
	Age            string    `json:"age,omitempty"`
	IsExternal     bool      `json:"isExternal"`
	ExternalSource string    `json:"externalSource,omitempty"`
}
