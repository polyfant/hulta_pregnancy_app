# Changelog

## [Unreleased]

### Added

-   ML-powered growth predictions 🧠
-   Environmental impact monitoring 🌍
-   Privacy Dashboard with real-time scoring 🔒
-   Automated privacy audits
-   Granular privacy controls
-   Data masking capabilities
-   Anonymous usage analytics
-   Scheduled data cleanup
-   Export/purge functionality
-   Created internal/api/types/types.go for common API types
-   Added PostgreSQL database connection handling
-   Added basic configuration management
-   Added repository constructors
-   Added pregnancy service implementation
-   Added breeding service interface

### Changed

-   Migrated frontend from vanilla JS to React
-   Refactored API calls to use React Query
-   Replaced custom styling with Mantine UI
-   Updated project documentation
-   Moved ErrorResponse to types package
-   Updated handler error responses to use types.ErrorResponse
-   Improved main.go error handling and structure
-   Updated PregnancyEvent model with UserID field
-   Refactored service initialization in main.go

### Fixed

-   Improved data encryption
-   Added comprehensive audit logging
-   Enhanced privacy controls
-   Added environmental monitoring
-   Fixed service pointer vs value handling
-   Fixed database configuration structure
-   Fixed repository constructor implementations

## [0.3.1] - 2025-01-19

### Fixed
- Fixed database migration handling using goose 🔧
- Standardized API paths to use /api/v1 prefix consistently
- Updated frontend API client to use consistent versioned endpoints
- Fixed AddHorse component to use proper API client with auth headers

## [0.3.0] - 2025-01-15

### Added
- Privacy change logging system 🔒
- Data retention controls for weather and health data
- Automated data cleanup based on retention settings
- Transaction support for privacy preference updates

### Changed
- Updated privacy preferences to use upsert functionality
- Improved test database setup and cleanup
- Enhanced error handling in privacy repository
- Optimized weather data retention period to 30 days

### Fixed
- Fixed privacy change log table creation
- Improved transaction rollback handling
- Fixed weather data deletion test
- Updated test assertions for better reliability

## [0.2.0] - 2025-01-01

### Added

-   React frontend with TypeScript
-   Mantine UI component library
-   React Query for state management
-   Comprehensive frontend validation
-   Enhanced user experience with loading states

### Changed

-   Migrated frontend from vanilla JS to React
-   Refactored API calls to use React Query
-   Replaced custom styling with Mantine UI
-   Updated project documentation

## [0.1.0] - 2024-01-08

### Initial Release

-   Basic project structure
-   Core functionality
-   Initial API handlers
