package pregnancy

import (
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPregnancyStage(t *testing.T) {
    calc := NewCalculator()
    
    tests := []struct {
        name              string
        daysSinceConception int
        expectedStage     models.PregnancyStage
    }{
        {"Early stage", 50, models.PregnancyStageEarlyGestation},
        {"Mid stage", 150, models.PregnancyStageMidGestation},
        {"Late stage", 300, models.PregnancyStageLateGestation},
        {"Overdue", 345, models.PregnancyStageOverdue},
        {"High risk", 380, models.PregnancyStageHighRisk},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            stage := calc.GetPregnancyStage(tt.daysSinceConception)
            assert.Equal(t, tt.expectedStage, stage)
        })
    }
}

func TestGetStageInfo(t *testing.T) {
    calc := NewCalculator()
    now := time.Now()
    
    tests := []struct {
        name           string
        conceptionDate time.Time
        expectOverdue  bool
        expectDaysOverdue int
    }{
        {
            name: "Normal pregnancy",
            conceptionDate: now.AddDate(0, 0, -300),
            expectOverdue: false,
            expectDaysOverdue: 0,
        },
        {
            name: "Overdue pregnancy",
            conceptionDate: now.AddDate(0, 0, -350),
            expectOverdue: true,
            expectDaysOverdue: 10,
        },
        {
            name: "High risk overdue",
            conceptionDate: now.AddDate(0, 0, -380),
            expectOverdue: true,
            expectDaysOverdue: 40,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            info := calc.GetStageInfo(tt.conceptionDate)
            assert.Equal(t, tt.expectOverdue, info.IsOverdue)
            assert.Equal(t, tt.expectDaysOverdue, info.DaysOverdue)
        })
    }
}

func TestCalculator(t *testing.T) {
    calc := NewCalculator()
    now := time.Now()

    // Test all pregnancy stages with more edge cases
    t.Run("GetPregnancyStage edge cases", func(t *testing.T) {
        tests := []struct {
            name              string
            daysSinceConception int
            expectedStage     models.PregnancyStage
        }{
            {"Day 0", 0, models.PregnancyStageEarlyGestation},
            {"Last day of early stage", 98, models.PregnancyStageEarlyGestation},
            {"First day of mid stage", 99, models.PregnancyStageMidGestation},
            {"Last day of mid stage", 196, models.PregnancyStageMidGestation},
            {"First day of late stage", 197, models.PregnancyStageLateGestation},
            {"Last normal day", 340, models.PregnancyStageLateGestation},
            {"First overdue day", 341, models.PregnancyStageOverdue},
            {"Last overdue day", 370, models.PregnancyStageOverdue},
            {"First high risk day", 371, models.PregnancyStageHighRisk},
            {"Extremely late", 400, models.PregnancyStageHighRisk},
        }

        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                stage := calc.GetPregnancyStage(tt.daysSinceConception)
                assert.Equal(t, tt.expectedStage, stage)
            })
        }
    })

    // Test GetStageInfo with more detailed checks
    t.Run("GetStageInfo comprehensive checks", func(t *testing.T) {
        tests := []struct {
            name           string
            conceptionDate time.Time
            expect        struct {
                Stage          models.PregnancyStage
                DaysSoFar     int
                WeeksSoFar    int
                DaysRemaining int
                WeeksRemaining int
                Progress      float64
                IsOverdue     bool
                DaysOverdue   int
            }
        }{
            {
                name: "Early pregnancy",
                conceptionDate: now.AddDate(0, 0, -50),
                expect: struct {
                    Stage          models.PregnancyStage
                    DaysSoFar     int
                    WeeksSoFar    int
                    DaysRemaining int
                    WeeksRemaining int
                    Progress      float64
                    IsOverdue     bool
                    DaysOverdue   int
                }{
                    Stage:          models.PregnancyStageEarlyGestation,
                    DaysSoFar:     50,
                    WeeksSoFar:    7,
                    DaysRemaining: 290,
                    WeeksRemaining: 41,
                    Progress:      14.7, // 50/340 * 100
                    IsOverdue:     false,
                    DaysOverdue:   0,
                },
            },
            {
                name: "Just overdue",
                conceptionDate: now.AddDate(0, 0, -341),
                expect: struct {
                    Stage          models.PregnancyStage
                    DaysSoFar     int
                    WeeksSoFar    int
                    DaysRemaining int
                    WeeksRemaining int
                    Progress      float64
                    IsOverdue     bool
                    DaysOverdue   int
                }{
                    Stage:          models.PregnancyStageOverdue,
                    DaysSoFar:     341,
                    WeeksSoFar:    48,
                    DaysRemaining: -1,
                    WeeksRemaining: 0,
                    Progress:      100.0,
                    IsOverdue:     true,
                    DaysOverdue:   1,
                },
            },
        }

        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                info := calc.GetStageInfo(tt.conceptionDate)
                assert.Equal(t, tt.expect.Stage, info.Stage)
                assert.Equal(t, tt.expect.DaysSoFar, info.DaysSoFar)
                assert.Equal(t, tt.expect.WeeksSoFar, info.WeeksSoFar)
                assert.Equal(t, tt.expect.DaysRemaining, info.DaysRemaining)
                assert.Equal(t, tt.expect.WeeksRemaining, info.WeeksRemaining)
                assert.InDelta(t, tt.expect.Progress, info.Progress, 0.1)
                assert.Equal(t, tt.expect.IsOverdue, info.IsOverdue)
                assert.Equal(t, tt.expect.DaysOverdue, info.DaysOverdue)
            })
        }
    })

    // Test edge cases for CalculateProgress
    t.Run("CalculateProgress edge cases", func(t *testing.T) {
        tests := []struct {
            name           string
            conceptionDate time.Time
            expectProgress float64
            expectDaysRemaining int
        }{
            {
                name: "Just started",
                conceptionDate: now,
                expectProgress: 0,
                expectDaysRemaining: 340,
            },
            {
                name: "Half way",
                conceptionDate: now.AddDate(0, 0, -170),
                expectProgress: 50,
                expectDaysRemaining: 170,
            },
            {
                name: "Completed",
                conceptionDate: now.AddDate(0, 0, -340),
                expectProgress: 100,
                expectDaysRemaining: 0,
            },
            {
                name: "Over completed",
                conceptionDate: now.AddDate(0, 0, -400),
                expectProgress: 100,
                expectDaysRemaining: 0,
            },
        }

        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                progress, daysRemaining := calc.CalculateProgress(tt.conceptionDate)
                assert.InDelta(t, tt.expectProgress, progress, 0.1)
                assert.Equal(t, tt.expectDaysRemaining, daysRemaining)
            })
        }
    })
}

func TestCalculateDueDateInfo(t *testing.T) {
    calc := NewCalculator()
    now := time.Now()

    tests := []struct {
        name           string
        conceptionDate time.Time
        expect        struct {
            daysUntilDue    int
            weeksUntilDue   int
            isInDueWindow   bool
            daysSpread      int  // Days between earliest and latest
        }
    }{
        {
            name: "Normal pregnancy - mid term",
            conceptionDate: now.AddDate(0, 0, -170), // Half way through
            expect: struct {
                daysUntilDue    int
                weeksUntilDue   int
                isInDueWindow   bool
                daysSpread      int
            }{
                daysUntilDue:  170,
                weeksUntilDue: 24,
                isInDueWindow: false,
                daysSpread:    50, // 370 - 320 days
            },
        },
        {
            name: "Almost due",
            conceptionDate: now.AddDate(0, 0, -319), // One day before earliest
            expect: struct {
                daysUntilDue    int
                weeksUntilDue   int
                isInDueWindow   bool
                daysSpread      int
            }{
                daysUntilDue:  21,
                weeksUntilDue: 3,
                isInDueWindow: false,
                daysSpread:    50,
            },
        },
        {
            name: "In due window",
            conceptionDate: now.AddDate(0, 0, -330), // Between min and max
            expect: struct {
                daysUntilDue    int
                weeksUntilDue   int
                isInDueWindow   bool
                daysSpread      int
            }{
                daysUntilDue:  10,
                weeksUntilDue: 1,
                isInDueWindow: true,
                daysSpread:    50,
            },
        },
        {
            name: "Overdue",
            conceptionDate: now.AddDate(0, 0, -380), // Past latest
            expect: struct {
                daysUntilDue    int
                weeksUntilDue   int
                isInDueWindow   bool
                daysSpread      int
            }{
                daysUntilDue:  -40,
                weeksUntilDue: -5,
                isInDueWindow: false,
                daysSpread:    50,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            info := calc.CalculateDueDateInfo(tt.conceptionDate)

            // Test days and weeks until due
            assert.Equal(t, tt.expect.daysUntilDue, info.DaysUntilDue)
            assert.Equal(t, tt.expect.weeksUntilDue, info.WeeksUntilDue)
            
            // Test due window status
            assert.Equal(t, tt.expect.isInDueWindow, info.IsInDueWindow)
            
            // Test due date ranges
            daysSpread := int(info.LatestDueDate.Sub(info.EarliestDueDate).Hours() / 24)
            assert.Equal(t, tt.expect.daysSpread, daysSpread)
            
            // Test date ordering
            assert.True(t, info.EarliestDueDate.Before(info.ExpectedDueDate))
            assert.True(t, info.ExpectedDueDate.Before(info.LatestDueDate))
            
            // Validate due dates are calculated correctly
            assert.Equal(t, tt.conceptionDate.AddDate(0, 0, calc.defaultGestationDays), info.ExpectedDueDate)
            assert.Equal(t, tt.conceptionDate.AddDate(0, 0, calc.minGestationDays), info.EarliestDueDate)
            assert.Equal(t, tt.conceptionDate.AddDate(0, 0, calc.maxGestationDays), info.LatestDueDate)
        })
    }
} 