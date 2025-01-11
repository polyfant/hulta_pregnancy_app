package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/tests"
	"github.com/stretchr/testify/assert"
)

func TestAddHealthRecord(t *testing.T) {
	db := tests.SetupTestDB(t)
	router := setupTestRouter(db)

	tests := []struct {
		name       string
		horseID    string
		record     models.HealthRecord
		wantStatus int
	}{
		{
			name:    "Valid health record",
			horseID: "1",
			record: models.HealthRecord{
				Type:        "VACCINATION",
				Date:        time.Now(),
				Description: "Annual vaccination",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:    "Missing type",
			horseID: "1",
			record: models.HealthRecord{
				Date:        time.Now(),
				Description: "Test record",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "Future date",
			horseID: "1",
			record: models.HealthRecord{
				Type:        "VACCINATION",
				Date:        time.Now().AddDate(0, 1, 0),
				Description: "Future record",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "Invalid horse ID",
			horseID: "invalid",
			record: models.HealthRecord{
				Type:        "VACCINATION",
				Date:        time.Now(),
				Description: "Test record",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.record)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", fmt.Sprintf("/api/horses/%s/health", tt.horseID), bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetHealthRecords(t *testing.T) {
	db := tests.SetupTestDB(t)
	router := setupTestRouter(db)

	tests := []struct {
		name       string
		horseID    string
		wantStatus int
		wantEmpty  bool
	}{
		{
			name:       "Existing horse",
			horseID:    "1",
			wantStatus: http.StatusOK,
			wantEmpty:  false,
		},
		{
			name:       "Non-existent horse",
			horseID:    "999",
			wantStatus: http.StatusNotFound,
			wantEmpty:  true,
		},
		{
			name:       "Invalid ID",
			horseID:    "abc",
			wantStatus: http.StatusBadRequest,
			wantEmpty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/horses/%s/health", tt.horseID), nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var response []models.HealthRecord
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				if tt.wantEmpty {
					assert.Empty(t, response)
				} else {
					assert.NotEmpty(t, response)
				}
			}
		})
	}
} 