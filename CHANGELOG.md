# Changelog

## [Unreleased]

### Fixed

-   Fixed type error in StageVisualization component by correcting property access from `status.stage` to `status.currentStage`
-   Updated Tooltip and Progress components to use correct stage property

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
-   Comprehensive role-based authentication middleware
-   Auth0 JWT token validation
-   Simplified role structure (user, admin, owner, farm_manager)
-   Yarn package manager integration
-   Test environment setup with Vitest
-   MSW for API mocking
-   Basic component tests structure
-   Mantine UI provider setup
-   React Query integration

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
-   Refactored authentication middleware in `internal/middleware/auth_middleware.go`
-   Updated route handling to support role-based access control
-   Improved test suite for rate limiting and authentication

### Fixed

-   Improved data encryption
-   Added comprehensive audit logging
-   Enhanced privacy controls
-   Added environmental monitoring
-   Fixed service pointer vs value handling
-   Fixed database configuration structure
-   Fixed repository constructor implementations
-   Resolved unused parameter warnings in test files
-   Standardized error handling in authentication middleware
-   Resolved potential data synchronization issue in frontend offline charts
-   Identified and prevented schema mismatch between IndexedDB and SQL data structures
-   Implemented robust data transformation for offline chart caching
-   Package dependencies
-   Build configuration
-   TypeScript configurations

### Dependencies

-   Added `github.com/golang-jwt/jwt/v4`
-   Added `github.com/auth0/go-jwt-middleware/v2`
-   Added `golang.org/x/time`

### Security

-   Enhanced token validation
-   Implemented role-based route protection

## [0.3.1] - 2025-01-19

### Fixed

-   Fixed database migration handling using goose 🔧
-   Standardized API paths to use /api/v1 prefix consistently
-   Updated frontend API client to use consistent versioned endpoints
-   Fixed AddHorse component to use proper API client with auth headers

## [0.3.0] - 2025-01-15

### Added

-   Privacy change logging system 🔒
-   Data retention controls for weather and health data
-   Automated data cleanup based on retention settings
-   Transaction support for privacy preference updates

### Changed

-   Updated privacy preferences to use upsert functionality
-   Improved test database setup and cleanup
-   Enhanced error handling in privacy repository
-   Optimized weather data retention period to 30 days

### Fixed

-   Fixed privacy change log table creation
-   Improved transaction rollback handling
-   Fixed weather data deletion test
-   Updated test assertions for better reliability

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
