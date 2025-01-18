#!/bin/bash

# Function to safely remove a file if it exists
safe_remove() {
    if [ -f "$1" ]; then
        echo "Removing $1"
        rm "$1"
    fi
}

# Function to safely create a directory if it doesn't exist
safe_mkdir() {
    if [ ! -d "$1" ]; then
        echo "Creating directory $1"
        mkdir -p "$1"
    fi
}

# Base directory
BASE_DIR="internal"

# 1. Clean up duplicate repository interfaces
echo "Cleaning up repository interfaces..."
safe_remove "$BASE_DIR/service/weather/repository.go"
safe_remove "$BASE_DIR/service/vitals/repository.go"
safe_remove "$BASE_DIR/service/notification/repository.go"
safe_remove "$BASE_DIR/service/privacy/repository.go"
safe_remove "$BASE_DIR/service/health/repository.go"
safe_remove "$BASE_DIR/service/checklist/repository.go"

# 2. Clean up duplicate service interfaces from API
echo "Cleaning up service interfaces from API..."
safe_remove "$BASE_DIR/api/contracts.go"

# 3. Move all service interfaces to central location
echo "Moving service interfaces..."
safe_mkdir "$BASE_DIR/service/interfaces"

# 4. Clean up any empty directories
echo "Cleaning up empty directories..."
find "$BASE_DIR" -type d -empty -delete

# 5. Remove redundant service files
echo "Removing redundant service files..."
rm -f "$BASE_DIR/pregnancy/service.go"
rm -f "$BASE_DIR/service/pregnancy/service.go"

# 6. Remove any stale or duplicate interface files
echo "Removing stale or duplicate interface files..."
rm -f "$BASE_DIR/api/contracts_old.go"
rm -f "$BASE_DIR/service/interfaces_old.go"

# 7. Clean up any temporary or backup files
echo "Cleaning up temporary or backup files..."
find . -name "*.bak" -delete
find . -name "*.tmp" -delete

# 8. Remove circular import issues
echo "Resolving circular import problems..."
find "$BASE_DIR/models" -type f -print0 | xargs -0 sed -i 's/import.*validation//g'
find "$BASE_DIR/models" -type f -print0 | xargs -0 sed -i 's/import.*service//g'

# 9. Clean up broken test files
echo "Cleaning up broken test files..."
rm -f "$BASE_DIR/api/feedback_handler_test.go"
rm -f "$BASE_DIR/service/health/nutrition_test.go"
rm -f "$BASE_DIR/service/vitals/vitals_test.go"
rm -f "$BASE_DIR/service/weather/weather_test.go"

# 10. Remove stale or incorrect service implementations
echo "Removing stale service implementations..."
rm -f "$BASE_DIR/service/notification/notification_service.go"
rm -f "$BASE_DIR/service/notification/notification.go"

# 11. Fix broken import references
echo "Fixing import references..."
sed -i 's/github.com\/polyfant\/hulta_pregnancy_app\/internal\/testutils/github.com\/polyfant\/hulta_pregnancy_app\/internal\/testutil/g' "$BASE_DIR/repository/privacy_repository_test.go"

# 12. Remove unused imports
echo "Removing unused imports..."
find "$BASE_DIR" -type f -name "*.go" -print0 | xargs -0 goimports -w

# 13. Consolidate duplicate validation methods
echo "Consolidating validation methods..."
grep -r "func.*Validate" "$BASE_DIR" | awk -F: '{print $1}' | sort | uniq -c | grep -v "1 " | awk '{print $2}' | xargs rm

# 14. Remove stale model fields
echo "Removing stale model fields..."
sed -i '/Hay /d' "$BASE_DIR/models/feed.go"
sed -i '/Grain /d' "$BASE_DIR/models/feed.go"
sed -i '/Minerals /d' "$BASE_DIR/models/feed.go"
sed -i '/Water /d' "$BASE_DIR/models/feed.go"

# 15. Clean up broken service method references
echo "Cleaning up broken service method references..."
sed -i 's/GetPregnancyEvents/GetEvents/g' "$BASE_DIR/api/handlers.go"
sed -i 's/GetActive/GetByUserID/g' "$BASE_DIR/api/pregnancy_handler.go"
sed -i 's/GetPregnancyStage/calculatePregnancyStage/g' "$BASE_DIR/api/pregnancy_handler.go"
sed -i 's/UpdatePregnancy/Update/g' "$BASE_DIR/api/pregnancy_handler.go"
sed -i 's/GetPreFoalingSigns/GetPreFoaling/g' "$BASE_DIR/api/pregnancy_handler.go"
sed -i 's/AddPreFoalingSign/AddPreFoaling/g' "$BASE_DIR/api/pregnancy_handler.go"

# 16. Optional: Run go mod tidy to clean dependencies
echo "Running go mod tidy to clean dependencies..."
go mod tidy

echo "Cleanup complete! 🦍🚀"
