# Horse Tracking Application

A modern web application for tracking horses, breeding, and health records.

## Features
- Horse Management (Add, Edit, Delete horses)
- Breeding Tracking
- Pregnancy Monitoring
- Health Records
- Milestone Tracking
- Mobile-friendly Interface

## Tech Stack
- Backend: Go with Gin framework
- Database: SQLite
- Frontend: (To be implemented with a modern framework)

## Project Structure
```
horse_tracking_go/
├── api/            # API handlers and routes
├── internal/       # Internal packages
│   ├── models/     # Data models
│   ├── database/   # Database operations
│   └── service/    # Business logic
├── web/           # Frontend assets (to be added)
└── cmd/           # Application entry points
```

## Getting Started
1. Install Go 1.21 or later
2. Clone the repository
3. Run `go mod download`
4. Start the server with `go run cmd/server/main.go`
