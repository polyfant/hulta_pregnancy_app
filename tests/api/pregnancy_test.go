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

func TestGetPregnancyStatus(t *testing.T) {
    // Setup
    db := tests.SetupTestDB(t)
    router := setupTestRouter(db)

    tests := []struct {
        name       string
        horseID    string
        wantStatus int
    }{
        {
            name:       "Valid horse ID",
            horseID:    "1",
            wantStatus: http.StatusOK,
        },
        {
            name:       "Invalid horse ID",
            horseID:    "invalid",
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("GET", "/api/horses/"+tt.horseID+"/pregnancy", nil)
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.wantStatus, w.Code)
            
            if tt.wantStatus == http.StatusOK {
                var response models.Pregnancy
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err)
            }
        })
    }
}

func TestStartPregnancyTracking(t *testing.T) {
    db := tests.SetupTestDB(t)
    router := setupTestRouter(db)

    tests := []struct {
        name       string
        horseID    string
        payload    map[string]interface{}
        wantStatus int
    }{
        {
            name:    "Valid start tracking",
            horseID: "1",
            payload: map[string]interface{}{
                "conception_date": time.Now().AddDate(0, -1, 0),
            },
            wantStatus: http.StatusOK,
        },
        {
            name:    "Future conception date",
            horseID: "1",
            payload: map[string]interface{}{
                "conception_date": time.Now().AddDate(0, 1, 0),
            },
            wantStatus: http.StatusBadRequest,
        },
        {
            name:    "Missing conception date",
            horseID: "1",
            payload: map[string]interface{}{},
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.payload)
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("POST", "/api/horses/"+tt.horseID+"/pregnancy/start", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            router.ServeHTTP(w, req)
            assert.Equal(t, tt.wantStatus, w.Code)
        })
    }
}

func TestPreFoalingSigns(t *testing.T) {
    db := tests.SetupTestDB(t)
    router := setupTestRouter(db)

    tests := []struct {
        name       string
        horseID    string
        payload    map[string]interface{}
        wantStatus int
    }{
        {
            name:    "Record valid sign",
            horseID: "1",
            payload: map[string]interface{}{
                "description": "Udder changes",
                "date":       time.Now(),
                "notes":      "Significant changes observed",
            },
            wantStatus: http.StatusCreated,
        },
        {
            name:    "Missing description",
            horseID: "1",
            payload: map[string]interface{}{
                "date":  time.Now(),
                "notes": "Test notes",
            },
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.payload)
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("POST", "/api/horses/"+tt.horseID+"/pregnancy/foaling-signs", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            router.ServeHTTP(w, req)
            assert.Equal(t, tt.wantStatus, w.Code)
        })
    }
}

func TestEndPregnancyTracking(t *testing.T) {
	db := tests.SetupTestDB(t)
	router := setupTestRouter(db)

	tests := []struct {
		name       string
		horseID    string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name:    "Valid end tracking",
			horseID: "1",
			payload: map[string]interface{}{
				"end_date": time.Now(),
				"status":   "FOALED",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:    "Future end date",
			horseID: "1",
			payload: map[string]interface{}{
				"end_date": time.Now().AddDate(0, 1, 0),
				"status":   "FOALED",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:    "Invalid status",
			horseID: "1",
			payload: map[string]interface{}{
				"end_date": time.Now(),
				"status":   "INVALID_STATUS",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/horses/"+tt.horseID+"/pregnancy/end", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestGetPregnancyGuidelines(t *testing.T) {
	db := tests.SetupTestDB(t)
	router := setupTestRouter(db)

	tests := []struct {
		name       string
		horseID    string
		stage      string
		wantStatus int
	}{
		{
			name:       "Early gestation",
			horseID:    "1",
			stage:      "EARLY_GESTATION",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid stage",
			horseID:    "1",
			stage:      "INVALID_STAGE",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Non-existent horse",
			horseID:    "999",
			stage:      "EARLY_GESTATION",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/horses/%s/pregnancy/guidelines?stage=%s", tt.horseID, tt.stage), nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
} 