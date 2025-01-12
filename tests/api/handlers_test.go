package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/api"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupTestHandler(t)
	router := gin.New()
	api.SetupRouter(router, handler)

	t.Run("GetHorse", func(t *testing.T) {
		horseID := uint(1)
		expectedHorse := &models.Horse{
			ID:     horseID,
			Name:   "Test Horse",
			UserID: "test_user",
		}

		mockHorseRepo.On("GetByID", mock.Anything, horseID).
			Return(expectedHorse, nil).Once()

		w := performRequest(router, "GET", "/api/horses/1", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Horse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedHorse.ID, response.ID)
		assert.Equal(t, expectedHorse.Name, response.Name)

		mockHorseRepo.AssertExpectations(t)
	})
}

func performRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	w := httptest.NewRecorder()

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	r.ServeHTTP(w, req)
	return w
} 