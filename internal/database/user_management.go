package database

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DatabaseUserConfig struct {
	ReadOnlyUser   string
	WriteUser      string
	AdminUser      string
}

func generateSecurePassword() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func CreateDatabaseUsers(db *gorm.DB) (*DatabaseUserConfig, error) {
	users := &DatabaseUserConfig{
		ReadOnlyUser: fmt.Sprintf("horse_readonly_%d", time.Now().UnixNano()),
		WriteUser:    fmt.Sprintf("horse_write_%d", time.Now().UnixNano()),
		AdminUser:    fmt.Sprintf("horse_admin_%d", time.Now().UnixNano()),
	}

	// Read-only user
	readOnlyPass := generateSecurePassword()
	if err := db.Exec(fmt.Sprintf(`
		CREATE USER %s WITH PASSWORD '%s';
		GRANT CONNECT ON DATABASE horse_tracker TO %s;
		GRANT USAGE ON SCHEMA public TO %s;
		GRANT SELECT ON ALL TABLES IN SCHEMA public TO %s;
		GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO %s;
	`, 
		users.ReadOnlyUser, readOnlyPass, 
		users.ReadOnlyUser, 
		users.ReadOnlyUser,
		users.ReadOnlyUser,
		users.ReadOnlyUser,
	)).Error; err != nil {
		return nil, fmt.Errorf("failed to create read-only user: %v", err)
	}

	// Write user
	writePass := generateSecurePassword()
	if err := db.Exec(fmt.Sprintf(`
		CREATE USER %s WITH PASSWORD '%s';
		GRANT CONNECT ON DATABASE horse_tracker TO %s;
		GRANT USAGE ON SCHEMA public TO %s;
		GRANT INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO %s;
		GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO %s;
	`, 
		users.WriteUser, writePass, 
		users.WriteUser,
		users.WriteUser,
		users.WriteUser,
		users.WriteUser,
	)).Error; err != nil {
		return nil, fmt.Errorf("failed to create write user: %v", err)
	}

	// Admin user
	adminPass := generateSecurePassword()
	if err := db.Exec(fmt.Sprintf(`
		CREATE USER %s WITH PASSWORD '%s' SUPERUSER;
		GRANT ALL PRIVILEGES ON DATABASE horse_tracker TO %s;
	`, 
		users.AdminUser, adminPass, 
		users.AdminUser,
	)).Error; err != nil {
		return nil, fmt.Errorf("failed to create admin user: %v", err)
	}

	return users, nil
}

func (u *DatabaseUserConfig) GetConnectionStrings(host, port, dbName string) map[string]string {
	return map[string]string{
		"readonly": fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", 
			host, port, u.ReadOnlyUser, dbName),
		"write": fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", 
			host, port, u.WriteUser, dbName),
		"admin": fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", 
			host, port, u.AdminUser, dbName),
	}
}
