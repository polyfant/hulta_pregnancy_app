#!/bin/bash
# Backup Script for Hulta Pregnancy App

BACKUP_DIR="/opt/hulta-pregnancy-app/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Create backup directory if not exists
mkdir -p $BACKUP_DIR

# Backup PostgreSQL Database
docker exec postgres pg_dump -U horsetracker horse_tracking_db > $BACKUP_DIR/db_backup_$TIMESTAMP.sql

# Backup Docker volumes
docker run --rm -v /opt/hulta-pregnancy-app:/backup ubuntu tar cvf /backup/volumes_backup_$TIMESTAMP.tar /var/lib/docker/volumes

# Cleanup old backups (keep last 7 days)
find $BACKUP_DIR -type f -mtime +7 -delete

echo "Backup completed: $TIMESTAMP"
