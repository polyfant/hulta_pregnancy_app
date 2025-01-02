#!/bin/bash

# Configuration
PGUSER="postgres"
PGPASSWORD="your_password_here"
BACKUP_DIR="/var/lib/postgresql/backups"
DB_NAME="HE_horse_db"
RETENTION_DAYS=7

# Export password for pg_dump
export PGPASSWORD

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Generate backup filename with timestamp
TIMESTAMP=$(date +%Y_%m_%d_%H_%M_%S)
BACKUP_FILE="$BACKUP_DIR/${DB_NAME}_${TIMESTAMP}.backup"

# Perform backup
pg_dump -h localhost -U "$PGUSER" -F c -b -v -f "$BACKUP_FILE" "$DB_NAME"

# Remove backups older than retention period
find "$BACKUP_DIR" -name "${DB_NAME}_*.backup" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: $BACKUP_FILE" 