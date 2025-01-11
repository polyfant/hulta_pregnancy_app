package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/tests"
	"github.com/stretchr/testify/assert"
)

func TestAddBreedingRecord(t *testing.T) {
	db := tests.SetupTestDB(t)
	router := setupTestRouter(db)

	tests := []struct {
		name       string
		horseID    string
		record     models.BreedingRecord
		wantStatus int
	}{
		{
			name:    "Valid breeding record",
			horseID: "1",
			record: models.BreedingRecord{
				Date:   time.Now(),
				Status: string(models.BreedingStatusActive),
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:    "Future date",
			horseID: "1",
			record: models.BreedingRecord{
				Date:   time.Now().AddDate(0, 1, 0),
				Status: string(models.BreedingStatusActive),
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "Invalid status",
			horseID: "1",
			record: models.BreedingRecord{
				Date:   time.Now(),
				Status: "INVALID_STATUS",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.record)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/horses/"+tt.horseID+"/breeding", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
} 