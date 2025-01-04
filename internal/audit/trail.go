package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AuditTrail struct {
	mu       sync.RWMutex
	logPath  string
	logFile  *os.File
}

type AuditEvent struct {
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	Details   map[string]interface{} `json:"details"`
}

func NewAuditTrail(logDir string) (*AuditTrail, error) {
	logPath := filepath.Join(logDir, "audit.log")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file: %w", err)
	}

	return &AuditTrail{
		logPath:  logPath,
		logFile:  logFile,
	}, nil
}

func (a *AuditTrail) LogEvent(userID, action string, details map[string]interface{}) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	event := AuditEvent{
		ID:        uuid.New(),
		Timestamp: time.Now(),
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal audit event: %w", err)
	}

	_, err = a.logFile.Write(append(jsonEvent, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write audit log: %w", err)
	}

	return nil
}

func (a *AuditTrail) Close() error {
	return a.logFile.Close()
}
