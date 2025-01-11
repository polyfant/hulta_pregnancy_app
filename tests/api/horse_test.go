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

func TestListHorses(t *testing.T) {
    // Setup
    db := tests.SetupTestDB(t)
    router := setupTestRouter(db)

    // Test
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/horses", nil)
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response []models.Horse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
}

func TestAddHorse(t *testing.T) {
    // Setup
    db := tests.SetupTestDB(t)
    router := setupTestRouter(db)

    newHorse := models.Horse{
        Name:   "Test Horse",
        Breed:  "Test Breed",
        Gender: models.GenderMare,
    }

    body, _ := json.Marshal(newHorse)

    // Test
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/horses", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)

    var response models.Horse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, newHorse.Name, response.Name)
}

func TestGetHorse(t *testing.T) {
    db := tests.SetupTestDB(t)
    router := setupTestRouter(db)

    tests := []struct {
        name       string
        horseID    string
        wantStatus int
    }{
        {
            name:       "Existing horse",
            horseID:    "1",
            wantStatus: http.StatusOK,
        },
        {
            name:       "Non-existent horse",
            horseID:    "999",
            wantStatus: http.StatusNotFound,
        },
        {
            name:       "Invalid ID format",
            horseID:    "abc",
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("GET", "/api/horses/"+tt.horseID, nil)
            router.ServeHTTP(w, req)
            assert.Equal(t, tt.wantStatus, w.Code)
        })
    }
}

func TestAddHorseValidation(t *testing.T) {
    db := tests.SetupTestDB(t)
    router := setupTestRouter(db)

    tests := []struct {
        name       string
        horse      models.Horse
        wantStatus int
    }{
        {
            name: "Missing name",
            horse: models.Horse{
                Breed:  "Test Breed",
                Gender: models.GenderMare,
            },
            wantStatus: http.StatusBadRequest,
        },
        {
            name: "Invalid gender",
            horse: models.Horse{
                Name:   "Test Horse",
                Breed:  "Test Breed",
                Gender: "INVALID",
            },
            wantStatus: http.StatusBadRequest,
        },
        {
            name: "Future birth date",
            horse: models.Horse{
                Name:      "Test Horse",
                Breed:     "Test Breed",
                Gender:    models.GenderMare,
                BirthDate: time.Now().AddDate(1, 0, 0),
            },
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            body, _ := json.Marshal(tt.horse)
            w := httptest.NewRecorder()
            req, _ := http.NewRequest("POST", "/api/horses", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            router.ServeHTTP(w, req)
            assert.Equal(t, tt.wantStatus, w.Code)
        })
    }
} 