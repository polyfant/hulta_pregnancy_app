# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

-   PostgreSQL database integration with GORM
-   Horse management system core functionality
-   Health tracking system
-   Pregnancy monitoring system
-   Financial management with expenses tracking
-   Breeding records management
-   Pre-foaling checklist functionality
-   Family tree visualization support
-   Validation system for all entities
-   Structured logging implementation
-   Error handling middleware
-   Database auto-migration system

### Changed

-   Migrated from SQLite to PostgreSQL
-   Updated database schema to use proper data types
-   Improved validation logic with dedicated validators
-   Enhanced error responses with better structure
-   Standardized model definitions with GORM tags

### Fixed

-   Foreign key constraint issues in database schema
-   Type mismatches in model definitions
-   Validation logic for dates and IDs
-   File structure organization
-   Duplicate type definitions
-   Expense type handling (changed from enum to string)

### Security

-   Added input validation for all endpoints
-   Implemented proper type checking
-   Added database constraint checks
-   Structured error handling to prevent information leakage

## [0.1.0] - 2024-01-02

### Added

-   Initial project setup
-   Basic API structure
-   Core models definition
-   Database connection handling
-   Basic CRUD operations
-   Health check endpoints
-   Logging system
-   Error handling framework
-   Configuration management
-   API documentation structure

[Unreleased]: https://github.com/polyfant/hulta_pregnancy_app/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/polyfant/hulta_pregnancy_app/releases/tag/v0.1.0
