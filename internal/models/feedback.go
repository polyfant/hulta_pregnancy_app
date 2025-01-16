package models

import (
	"time"
)

// FeatureRequest represents a potential feature that users can vote on
type FeatureRequest struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"` // WEARABLES, MONITORING, ANALYTICS, etc.
	Status      string    `json:"status"`   // PROPOSED, PLANNED, IN_DEVELOPMENT, RELEASED
	VoteCount   int       `json:"vote_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ReleasedAt  *time.Time `json:"released_at,omitempty"`
}

// UserFeatureVote tracks user votes on features
type UserFeatureVote struct {
	UserID         string    `json:"user_id"`
	FeatureID      uint      `json:"feature_id"`
	VotedAt        time.Time `json:"voted_at"`
	Comment        string    `json:"comment,omitempty"`
	Priority       string    `json:"priority"` // MUST_HAVE, NICE_TO_HAVE, MEH
}

// FeatureSurveyResponse captures detailed user feedback
type FeatureSurveyResponse struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         string    `json:"user_id"`
	FeatureID      uint      `json:"feature_id"`
	WouldUse       bool      `json:"would_use"`
	WouldPay       bool      `json:"would_pay"`
	PricePoint     *float64  `json:"price_point,omitempty"`
	UseCase        string    `json:"use_case"`
	Feedback       string    `json:"feedback"`
	SubmittedAt    time.Time `json:"submitted_at"`
}

// InitialFeatures returns the initial set of votable features
func InitialFeatures() []FeatureRequest {
	now := time.Now()
	return []FeatureRequest{
		{
			Title:       "Smart Watch Integration",
			Description: "Receive real-time alerts and monitoring data on your smart watch. Get instant notifications about your horse's vital signs and important events.",
			Category:    "WEARABLES",
			Status:      "PROPOSED",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Title:       "Automated Camera Monitoring",
			Description: "24/7 video monitoring with AI-powered behavior analysis. Get alerts when unusual behavior is detected.",
			Category:    "MONITORING",
			Status:      "PROPOSED",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Title:       "Advanced Analytics Dashboard",
			Description: "Detailed analytics and insights about your horse's pregnancy journey with customizable charts and reports.",
			Category:    "ANALYTICS",
			Status:      "PROPOSED",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Title:       "Veterinarian Portal",
			Description: "Dedicated portal for vets to monitor their patients, with professional-grade analytics and reporting.",
			Category:    "PROFESSIONAL",
			Status:      "PROPOSED",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			Title:       "Multi-Horse Management",
			Description: "Tools for breeding facilities to manage multiple pregnancies simultaneously with comparative analytics.",
			Category:    "MANAGEMENT",
			Status:      "PROPOSED",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
}
