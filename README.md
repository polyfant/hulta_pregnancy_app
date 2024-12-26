# Horse Tracking Application

A comprehensive web application for tracking horses, breeding, and health records. This application helps horse owners and breeders manage their horses' information, health records, breeding costs, and provides detailed pregnancy monitoring with stage-specific guidelines.

## Current Features

### Horse Management
- Add new horses with details (name, breed, gender, birth date)
- View all horses in a clean, sortable table
- Edit existing horse information
- Delete horses from the system

### Pregnancy Monitoring
- Track pregnancy stages (Early, Mid, Late, Pre-Foaling, Foaling)
- Detailed guidelines for each pregnancy stage
- Record and monitor pre-foaling signs
- Track conception dates and calculate due dates
- Monitor health events throughout pregnancy
- Nutrition and exercise recommendations by stage
- Warning signs and health checkpoints

### Breeding Management
- Record breeding events and outcomes
- Track breeding costs and expenses
- Monitor breeding history
- Manage breeding schedules

## Upcoming Features

### Family Tree Management
- Visual family tree representation
- Track lineage and breeding history
- Link horses to their parents
- Support for external parent references

### Health Tracking
- Record and monitor health assessments
- Track vaccination status
- Manage health records and checkups
- View comprehensive health summaries

## API Endpoints

### Currently Implemented
- `GET /api/horses` - List all horses
- `POST /api/horses` - Create a new horse
- `PUT /api/horses/{id}` - Update a horse
- `DELETE /api/horses/{id}` - Delete a horse
- `GET /api/horses/{id}` - Get details of a specific horse
- `POST /api/horses/{id}/pregnancy/start` - Start pregnancy tracking
- `POST /api/horses/{id}/pregnancy/events` - Add pregnancy event
- `GET /api/horses/{id}/pregnancy/status` - Get pregnancy status
- `GET /api/horses/{id}/pregnancy/guidelines` - Get stage-specific guidelines
- `POST /api/horses/{id}/pregnancy/pre-foaling-signs` - Record pre-foaling signs

### Planned Endpoints
- `GET /api/horses/{id}/family-tree` - Get family tree for a horse
- `GET /api/horses/{id}/health` - Get health summary
- `POST /api/horses/{id}/health-records` - Add health record
- `GET /api/horses/{id}/breeding-costs` - View breeding costs

## Project Structure
```
horse_tracking_go/
├── cmd/                 # Application entry points
│   └── server/         # Main server application
├── internal/           # Internal packages
│   ├── api/           # API handlers and routing
│   ├── database/      # Database implementation (SQLite)
│   ├── models/        # Data models and types
│   └── service/       # Business logic services
│       ├── pregnancy/ # Pregnancy tracking service
│       └── breeding/  # Breeding management service
├── frontend-react/    # React frontend application
└── data/             # Data storage directory
```

## Tech Stack

### Backend
- Go 1.21+
- Gin web framework for HTTP routing
- SQLite3 for data storage
- Zerolog for structured logging
- Clean architecture with service-based design

### Frontend
- React 18
- TypeScript
- Mantine UI components
- TanStack Query for API integration
- Vite for building and development

## Getting Started

### Prerequisites
1. Go 1.21 or later
2. Node.js 18+ and npm
3. Git for version control
4. SQLite3

### Installation & Development
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd horse_tracking_go
   ```

2. Install backend dependencies:
   ```bash
   go mod download
   ```

3. Start the backend:
   ```bash
   go run cmd/server/main.go
   ```

4. Install frontend dependencies:
   ```bash
   cd frontend-react
   npm install
   ```

5. Start the frontend development server:
   ```bash
   npm run dev
   ```

## Development Status

The application is currently in active development with a focus on the following areas:

### Completed
- Core horse management functionality
- Basic database schema and SQLite integration
- Pregnancy tracking service with stage-specific guidelines
- Pre-foaling sign monitoring
- Event-based pregnancy tracking

### In Progress
- Breeding management service implementation
- Frontend UI components for pregnancy tracking
- API endpoint implementation for all services
- Integration testing

### Planned
- Family tree visualization
- Health record management
- Comprehensive breeding cost tracking
- Mobile-responsive UI improvements

## Contributing

This project is under active development. Feel free to submit issues and pull requests.

## License

[Add your license information here]