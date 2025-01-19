package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type MigrationConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Schema   string
}

func getMigrationsPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "migrations")
}

// cleanSQL removes comments and splits into statements while preserving PL/pgSQL functions
func cleanSQL(content string) []string {
	var statements []string
	var currentStmt strings.Builder
	lines := strings.Split(content, "\n")
	inFunction := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// Skip empty lines and goose comments
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "-- +goose") {
			continue
		}

		// Check if we're entering a function definition
		if strings.Contains(strings.ToUpper(trimmedLine), "CREATE OR REPLACE FUNCTION") {
			inFunction = true
		}

		currentStmt.WriteString(line)
		currentStmt.WriteString("\n")

		// If we're not in a function, look for statement terminators
		if !inFunction {
			if strings.HasSuffix(trimmedLine, ";") {
				if stmt := strings.TrimSpace(currentStmt.String()); stmt != "" {
					statements = append(statements, stmt)
				}
				currentStmt.Reset()
			}
		} else {
			// For functions, look for the function terminator
			if strings.Contains(trimmedLine, "$func$;") {
				inFunction = false
				if stmt := strings.TrimSpace(currentStmt.String()); stmt != "" {
					statements = append(statements, stmt)
				}
				currentStmt.Reset()
			}
		}
	}

	// Add any remaining statement
	if stmt := strings.TrimSpace(currentStmt.String()); stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}

func RunMigrations(db *gorm.DB) error {
	migrationsPath := getMigrationsPath()
	log.Printf("Looking for migrations in: %s", migrationsPath)

	// Read migration files
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("could not read migrations directory: %v", err)
	}

	// Get only .sql files and sort them
	var migrations []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			migrations = append(migrations, f.Name())
		}
	}
	sort.Strings(migrations)

	// Execute each migration file
	for _, migration := range migrations {
		log.Printf("Executing migration: %s", migration)
		
		content, err := os.ReadFile(filepath.Join(migrationsPath, migration))
		if err != nil {
			return fmt.Errorf("could not read migration file %s: %v", migration, err)
		}

		// Clean and split the SQL content
		statements := cleanSQL(string(content))

		// Execute each statement
		for _, stmt := range statements {
			if err := db.Exec(stmt).Error; err != nil {
				return fmt.Errorf("error executing statement in %s: %v\nStatement: %s", migration, err, stmt)
			}
		}

		log.Printf("Successfully executed migration: %s", migration)
	}

	log.Printf("All migrations completed successfully")
	return nil
}
