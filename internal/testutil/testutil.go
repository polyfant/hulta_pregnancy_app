package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB creates a test database connection
func SetupTestDB() (*gorm.DB, error) {
	// Load test environment variables
	if err := godotenv.Load("../../.env.test"); err != nil {
		return nil, fmt.Errorf("error loading .env.test file: %w", err)
	}

	// Get database connection parameters from environment variables
	host := getEnvOrDefault("TEST_DB_HOST", "localhost")
	port := getEnvOrDefault("TEST_DB_PORT", "5432")
	user := getEnvOrDefault("TEST_DB_USER", "postgres")
	password := getEnvOrDefault("TEST_DB_PASSWORD", "postgres")
	dbname := getEnvOrDefault("TEST_DB_NAME", "horse_tracking_test")

	// Create database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Auto-migrate all required tables
	err = db.AutoMigrate(
		&models.VitalSigns{},
		&models.VitalSignsPrediction{},
		&models.VitalSignsAlert{},
		&models.FeatureRequest{},
		&models.UserFeatureVote{},
		&models.FeatureSurveyResponse{},
		&models.Pregnancy{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate tables: %w", err)
	}

	return db, nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// CreateTestServer creates a test HTTP server with the given handler
func CreateTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

// CreateTestWebSocketServer creates a test WebSocket server
func CreateTestWebSocketServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

// CreateTestWebSocketClient creates a test WebSocket client
func CreateTestWebSocketClient(server *httptest.Server, path string) (*websocket.Conn, error) {
	url := "ws" + server.URL[4:] + path
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	return conn, err
}

// MockTime represents a fixed time for testing
var MockTime = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

// MockTimeNow returns a fixed time for testing
func MockTimeNow() time.Time {
	return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
}

// CreateTestContext creates a test Gin context
func CreateTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

// ParseResponse parses the JSON response from a test request
func ParseResponse(w *httptest.ResponseRecorder, v interface{}) error {
	return json.NewDecoder(w.Body).Decode(v)
}

// CreateTestVitalSigns creates test vital signs data
func CreateTestVitalSigns(horseID uint) *models.VitalSigns {
	return &models.VitalSigns{
		HorseID:     horseID,
		Temperature: 38.0,
		HeartRate:   40,
		Respiration: 12,
		RecordedAt:  MockTimeNow(),
	}
}

// CreateTestFeatureRequest creates a test feature request
func CreateTestFeatureRequest(db *gorm.DB, title string) (*models.FeatureRequest, error) {
	feature := &models.FeatureRequest{
		Title:       title,
		Description: "Test Description",
		Status:      "PLANNED",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := db.Create(feature).Error
	if err != nil {
		return nil, err
	}

	return feature, nil
}

// CreateTestVitalSignsPrediction creates a test vital signs prediction
func CreateTestVitalSignsPrediction(db *gorm.DB, horseID uint) (*models.VitalSignsPrediction, error) {
	prediction := &models.VitalSignsPrediction{
		HorseID:            horseID,
		PredictedFoaling:   time.Now().Add(24 * time.Hour),
		FoalingProbability: 0.85,
		RiskLevel:          "LOW",
		Alerts:             pq.StringArray{"Alert 1", "Alert 2"},
		CreatedAt:          time.Now(),
	}

	err := db.Create(prediction).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create test vital signs prediction: %w", err)
	}

	return prediction, nil
}

// AssertWebSocketMessage asserts that a WebSocket message matches expected data
func AssertWebSocketMessage(t testing.TB, conn *websocket.Conn, timeout time.Duration, expected interface{}) error {
	t.Helper()
	
	done := make(chan error, 1)
	go func() {
		_, message, err := conn.ReadMessage()
		if err != nil {
			done <- err
			return
		}

		var actual interface{}
		if err := json.Unmarshal(message, &actual); err != nil {
			done <- err
			return
		}

		expectedJSON, _ := json.Marshal(expected)
		actualJSON, _ := json.Marshal(actual)
		if string(expectedJSON) != string(actualJSON) {
			done <- fmt.Errorf("expected %s but got %s", string(expectedJSON), string(actualJSON))
			return
		}

		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for websocket message")
	}
}
