# Horse Tracking Application

A modern web application for tracking horses, breeding, and health records. This application helps horse owners and breeders manage their horses' information, health records, breeding costs, and pregnancy monitoring.

## Features

### Horse Management
- Add new horses with details like name, breed, and birth date
- View all horses in the system
- Get detailed information about specific horses

### Health Tracking
- Record and monitor health assessments
- Track vaccination status
- Manage health records and checkups
- View comprehensive health summaries

### Breeding Management
- Track breeding costs
- Record breeding-related expenses
- Monitor breeding history

### Pregnancy Monitoring
- Track pregnancy milestones
- Get pregnancy guidelines
- Monitor conception and due dates

## API Endpoints

### Horse Endpoints
- `GET /horses` - List all horses
- `POST /horses` - Create a new horse
- `GET /horses/{id}` - Get details of a specific horse

### Health Endpoints
- `GET /horses/{id}/health` - Get health summary for a horse
- `POST /horses/{id}/health-records` - Add a new health record

### Breeding Endpoints
- `GET /horses/{id}/breeding-costs` - View breeding costs
- `POST /horses/{id}/breeding-costs` - Add breeding cost record

### Pregnancy Endpoints
- `GET /horses/{id}/pregnancy-guidelines` - Get pregnancy guidelines

## Project Structure
```
horse_tracking_go/
├── cmd/                    # Application entry points
│   ├── migrate/           # Database migration tool
│   └── server/            # Main server application
├── internal/              # Internal packages
│   ├── api/              # API handlers and routes
│   ├── database/         # Database operations
│   ├── models/           # Data models
│   └── service/          # Business logic services
│       ├── breeding/     # Breeding-related services
│       ├── health/       # Health monitoring services
│       └── pregnancy/    # Pregnancy tracking services
└── web/                  # Frontend application
```

## Tech Stack
- **Backend**: 
  - Go 1.21+
  - Gin web framework
  - SQLite database
- **API Documentation**: 
  - Swagger/OpenAPI
  - Available at `/swagger/*`

## Getting Started

### Prerequisites
1. Install Go 1.21 or later
2. Git for version control

### Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd horse_tracking_go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

4. Access the API:
   - The server runs on `http://localhost:8080`
   - API documentation available at `http://localhost:8080/swagger/index.html`

## Development

### Running Tests
```bash
go test ./...
```

### API Documentation
The API is documented using Swagger. You can view the interactive API documentation by:
1. Starting the server
2. Visiting `http://localhost:8080/swagger/index.html`

## Database
The application uses SQLite for data storage. The database file (`horse_tracker.db`) is automatically created in the root directory when you first run the application.
