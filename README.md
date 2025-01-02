# Horse Management System

A comprehensive system for managing horse breeding, health records, and expenses.

## Features

-   🐎 Horse Management

    -   Basic information tracking
    -   Breeding records
    -   Family tree visualization
    -   Age and breeding eligibility calculations

-   🏥 Health Tracking

    -   Health records
    -   Vaccination schedules
    -   Dental care tracking
    -   Pre-foaling monitoring

-   👶 Pregnancy Management

    -   Conception tracking
    -   Pregnancy event logging
    -   Due date calculations
    -   Pre-foaling checklist

-   💰 Financial Management
    -   Expense tracking
    -   Recurring expenses
    -   Breeding costs
    -   Financial summaries

## Tech Stack

-   Backend:
    -   Go 1.21+
    -   Gin Web Framework
    -   GORM
    -   PostgreSQL

## Getting Started

1. Prerequisites:

    ```bash
    - Go 1.21+
    - PostgreSQL 14+
    ```

2. Database Setup:

    ```sql
    CREATE DATABASE HE_horse_db;
    ```

3. Run the server:
    ```bash
    go run cmd/server/main.go
    ```

## Project Structure

```
.
├── cmd/
│   ├── server/        # Main application entry
│   └── migrate/       # Database migration tools
├── internal/
│   ├── api/          # HTTP handlers and routes
│   ├── database/     # Database access layer
│   ├── models/       # Data models
│   ├── service/      # Business logic
│   ├── validation/   # Input validation
│   └── middleware/   # HTTP middleware
```

## API Documentation

[API documentation to be added]
