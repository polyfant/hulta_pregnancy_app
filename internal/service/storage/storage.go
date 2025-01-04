package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) (*FileStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &FileStorage{basePath: basePath}, nil
}

func (fs *FileStorage) SaveHorsePhoto(horseID int64, file io.Reader, filename string) (string, error) {
	// Create horse-specific directory
	horsePath := filepath.Join(fs.basePath, fmt.Sprintf("horse_%d", horseID))
	if err := os.MkdirAll(horsePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create horse directory: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.TrimSuffix(filepath.Base(filename), ext), ext)
	fullPath := filepath.Join(horsePath, newFilename)

	// Create the file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy the file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	logger.Info("Saved horse photo", map[string]interface{}{
		"horseID":   horseID,
		"filename":  newFilename,
		"fullPath": fullPath,
	})

	return newFilename, nil
}

func (fs *FileStorage) GetHorsePhotoPath(horseID int64, filename string) string {
	return filepath.Join(fs.basePath, fmt.Sprintf("horse_%d", horseID), filename)
}

func (fs *FileStorage) DeleteHorsePhoto(horseID int64, filename string) error {
	fullPath := fs.GetHorsePhotoPath(horseID, filename)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	logger.Info("Deleted horse photo", map[string]interface{}{
		"horseID":   horseID,
		"filename":  filename,
		"fullPath": fullPath,
	})

	return nil
}
