package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type ExportService struct {
	exportPath string
}

func NewExportService(exportPath string) (*ExportService, error) {
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create export directory: %w", err)
	}
	return &ExportService{exportPath: exportPath}, nil
}

func (es *ExportService) ExportHorsesToCSV(horses []models.Horse) (string, error) {
	filename := fmt.Sprintf("horses_export_%s.csv", time.Now().Format("2006-01-02_15-04-05"))
	filepath := fmt.Sprintf("%s/%s", es.exportPath, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Name", "Breed", "Date of Birth", "Conception Date", "Mother ID", "Father ID"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, horse := range horses {
		motherID := ""
		if horse.MotherId != nil {
			motherID = strconv.FormatInt(int64(*horse.MotherId), 10)
		}
		fatherID := ""
		if horse.FatherId != nil {
			fatherID = strconv.FormatInt(int64(*horse.FatherId), 10)
		}
		conceptionDate := ""
		if horse.ConceptionDate != nil {
			conceptionDate = horse.ConceptionDate.Format("2006-01-02")
		}

		record := []string{
			strconv.FormatInt(int64(horse.ID), 10),
			horse.Name,
			horse.Breed,
			horse.BirthDate.Format("2006-01-02"),
			conceptionDate,
			motherID,
			fatherID,
		}

		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write record: %w", err)
		}
	}

	logger.Info("Exported horses to CSV", map[string]interface{}{
		"filename": filename,
		"count":    len(horses),
	})

	return filename, nil
}

func (es *ExportService) ExportHealthRecordsToCSV(records []models.HealthRecord) (string, error) {
	filename := fmt.Sprintf("health_records_export_%s.csv", time.Now().Format("2006-01-02_15-04-05"))
	filepath := fmt.Sprintf("%s/%s", es.exportPath, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Horse ID", "Date", "Type", "Notes"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, record := range records {
		row := []string{
			strconv.FormatInt(int64(record.ID), 10),
			strconv.FormatInt(int64(record.HorseID), 10),
			record.Date.Format("2006-01-02"),
			record.Type,
			record.Description,
		}

		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write record: %w", err)
		}
	}

	logger.Info("Exported health records to CSV", map[string]interface{}{
		"filename": filename,
		"count":    len(records),
	})

	return filename, nil
}
