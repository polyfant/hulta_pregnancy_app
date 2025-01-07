package health

import (
	"context"
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHealthRepo struct {
	mock.Mock
}

func (m *mockHealthRepo) GetByHorseID(ctx context.Context, horseID uint) ([]models.HealthRecord, error) {
	args := m.Called(ctx, horseID)
	return args.Get(0).([]models.HealthRecord), args.Error(1)
}

func (m *mockHealthRepo) Create(ctx context.Context, record *models.HealthRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *mockHealthRepo) Update(ctx context.Context, record *models.HealthRecord) error {
	args := m.Called(ctx, record)
	return args.Error(0)
}

func (m *mockHealthRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockHorseRepo struct {
	mock.Mock
}

func (m *mockHorseRepo) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
	args := m.Called(ctx, id)
	if h, ok := args.Get(0).(*models.Horse); ok {
		return h, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockHorseRepo) ListByUser(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *mockHorseRepo) Create(ctx context.Context, horse *models.Horse) error {
	args := m.Called(ctx, horse)
	return args.Error(0)
}

func (m *mockHorseRepo) Update(ctx context.Context, horse *models.Horse) error {
	args := m.Called(ctx, horse)
	return args.Error(0)
}

func (m *mockHorseRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockHorseRepo) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *mockHorseRepo) GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error) {
	args := m.Called(ctx, horseID)
	return args.Get(0).(*models.FamilyTree), args.Error(1)
}

func (m *mockHorseRepo) GetOffspring(ctx context.Context, horseID uint) ([]models.Horse, error) {
	args := m.Called(ctx, horseID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

// ... implement other required methods

func TestCalculateDailyFeedRequirements(t *testing.T) {
	// Fix the season to Spring for consistent tests
	oldTimeNow := timeNow
	fixedTime := time.Date(2024, time.April, 15, 12, 0, 0, 0, time.UTC)
	timeNow = func() time.Time { return fixedTime }
	defer func() { timeNow = oldTimeNow }()

	tests := []struct {
		name     string
		horse    models.Horse
		activity ActivityLevel
		want     models.FeedRequirements
		wantErr  bool
	}{
		{
			name: "normal weight horse at maintenance",
			horse: models.Horse{
				Weight: 500,
			},
			activity: Maintenance,
			want: models.FeedRequirements{
				Hay:      10,    // 2% of body weight
				Grain:    2.5,   // 0.5% of body weight
				Minerals: 0.1,   // 100g
				Water:    25,    // 5% of body weight
			},
			wantErr: false,
		},
		{
			name: "pregnant mare in late gestation",
			horse: models.Horse{
				Weight:         500,
				IsPregnant:    true,
				ConceptionDate: timePtr(time.Now().AddDate(0, -9, 0)),
			},
			activity: Maintenance,
			want: models.FeedRequirements{
				Hay:      13,    // Base * 1.3
				Grain:    3,     // Base * 1.2
				Minerals: 0.15,  // Base * 1.5
				Water:    30,    // Base * 1.2
			},
			wantErr: false,
		},
		{
			name: "heavy work adjustments",
			horse: models.Horse{
				Weight: 500,
			},
			activity: HeavyWork,
			want: models.FeedRequirements{
				Hay:      13,    // Base * 1.3
				Grain:    5,     // Base * 2.0
				Minerals: 0.1,   // Unchanged
				Water:    40,    // Base * 1.6
			},
			wantErr: false,
		},
		{
			name: "invalid activity level",
			horse: models.Horse{
				Weight: 500,
			},
			activity: ActivityLevel(999),
			wantErr: true,
		},
		{
			name: "zero weight horse",
			horse: models.Horse{
				Weight: 0,
			},
			activity: Maintenance,
			want: models.FeedRequirements{
				Hay:      10,    // 2% of default weight (500)
				Grain:    2.5,   // 0.5% of default weight
				Minerals: 0.1,   // 100g
				Water:    25,    // 5% of default weight
			},
			wantErr: false,
		},
		{
			name: "extremely heavy horse",
			horse: models.Horse{
				Weight: 1500,
			},
			activity: Maintenance,
			wantErr: true, // Should fail validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			healthRepo := &mockHealthRepo{}
			horseRepo := &mockHorseRepo{}
			s := NewNutritionService(healthRepo, horseRepo)

			got, err := s.CalculateDailyFeedRequirements(tt.horse, tt.activity)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assertFeedRequirementsEqual(t, tt.want, got)
		})
	}
}

func assertFeedRequirementsEqual(t *testing.T, want, got models.FeedRequirements) {
	const delta = 0.01 // Allow small floating point differences
	assert.InDelta(t, want.Hay, got.Hay, delta)
	assert.InDelta(t, want.Grain, got.Grain, delta)
	assert.InDelta(t, want.Minerals, got.Minerals, delta)
	assert.InDelta(t, want.Water, got.Water, delta)
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestGetCurrentSeason(t *testing.T) {
	tests := []struct {
		name     string
		month    time.Month
		want     Season
	}{
		{"mid winter", time.January, Winter},
		{"early spring", time.March, Spring},
		{"mid spring", time.April, Spring},
		{"early summer", time.June, Summer},
		{"mid summer", time.July, Summer},
		{"early autumn", time.September, Autumn},
		{"mid autumn", time.October, Autumn},
		{"late winter", time.December, Winter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fixed time in the test month
			fixedTime := time.Date(2024, tt.month, 15, 12, 0, 0, 0, time.UTC)
			// Mock time.Now()
			oldTimeNow := timeNow
			timeNow = func() time.Time { return fixedTime }
			defer func() { timeNow = oldTimeNow }()

			got := getCurrentSeason()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAdjustForSeason(t *testing.T) {
	s := &NutritionService{}
	base := models.FeedRequirements{
		Hay:      10,
		Grain:    2.5,
		Minerals: 0.1,
		Water:    25,
	}

	tests := []struct {
		name   string
		season Season
		want   models.FeedRequirements
	}{
		{
			name:   "winter adjustments",
			season: Winter,
			want: models.FeedRequirements{
				Hay:      11.5,  // Base * 1.15
				Grain:    2.75,  // Base * 1.1
				Minerals: 0.1,   // Unchanged
				Water:    25,    // Unchanged
			},
		},
		{
			name:   "summer adjustments",
			season: Summer,
			want: models.FeedRequirements{
				Hay:      9,     // Base * 0.9
				Grain:    2.5,   // Unchanged
				Minerals: 0.1,   // Unchanged
				Water:    32.5,  // Base * 1.3
			},
		},
		{
			name:   "spring/autumn - no adjustments",
			season: Spring,
			want:   base,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.adjustForSeason(base, tt.season)
			assertFeedRequirementsEqual(t, tt.want, got)
		})
	}
} 