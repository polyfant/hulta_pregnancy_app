package database

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
)

type DatabaseBackup struct {
	dsn          string
	backupDir    string
	backupPrefix string
}

func NewDatabaseBackup(dsn, backupDir string) *DatabaseBackup {
	return &DatabaseBackup{
		dsn:          dsn,
		backupDir:    backupDir,
		backupPrefix: "horse_tracking_db_backup_",
	}
}

func (db *DatabaseBackup) Backup() error {
	// Ensure backup directory exists
	if err := os.MkdirAll(db.backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupFilename := fmt.Sprintf("%s%s.sql", db.backupPrefix, timestamp)
	backupPath := filepath.Join(db.backupDir, backupFilename)

	// PostgreSQL pg_dump command
	cmd := fmt.Sprintf("pg_dump %s > %s", db.dsn, backupPath)
	
	logger.Info("Starting database backup", 
		"backup_path", backupPath)

	// Execute backup
	output, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("backup failed: %w, output: %s", err, string(output))
	}

	logger.Info("Database backup completed successfully", 
		"backup_path", backupPath)

	return nil
}

func (db *DatabaseBackup) ScheduleBackups(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := db.Backup(); err != nil {
				logger.Error(err, "Scheduled backup failed")
			}
		}
	}()
}

// Retention management: Delete old backups
func (db *DatabaseBackup) ManageBackupRetention(maxBackups int) error {
	files, err := os.ReadDir(db.backupDir)
	if err != nil {
		return err
	}

	var backupFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), db.backupPrefix) {
			backupFiles = append(backupFiles, file)
		}
	}

	// Sort files by modification time, oldest first
	sort.Slice(backupFiles, func(i, j int) bool {
		infoI, _ := backupFiles[i].Info()
		infoJ, _ := backupFiles[j].Info()
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// Remove old backups
	for i := 0; i < len(backupFiles)-maxBackups; i++ {
		os.Remove(filepath.Join(db.backupDir, backupFiles[i].Name()))
	}

	return nil
}
