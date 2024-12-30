# Horse Tracking Application

A comprehensive web application for tracking horses, breeding, and health records. This application helps horse owners and breeders manage their horses' information, health records, breeding costs, and provides detailed pregnancy monitoring with stage-specific guidelines.

## Project Status: Active Development

This is a private, closed-source project under active development. All code, designs, and documentation are proprietary and confidential.

## Current Features

### Horse Management
- Add new horses with details (name, breed, gender, birth date)
- View all horses in a clean, sortable list
- Track gender-specific information (stallions, mares, geldings)
- Support for external parent references
- Dark theme with excellent readability

### Pregnancy Monitoring
- Track pregnancy stages (Early, Mid, Late, Due Soon)
- Calculate and display pregnancy progress
- Show stage-specific recommendations
- Record conception date and father information
- Support for both internal and external fathers
- Timeline view of pregnancy events
- Track vet checks and important milestones
- Warning system for health concerns

### Breeding Management
- Record breeding events
- Track conception dates
- Manage stallion information
- Support for external stallion records
- Validate breeding eligibility

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
- `GET /api/horses/{id}` - Get details of a specific horse
- `GET /api/horses/{id}/pregnancy/status` - Get pregnancy status and stage
- `POST /api/horses/{id}/pregnancy/events` - Add pregnancy event

### Planned Endpoints
- `GET /api/horses/{id}/family-tree` - Get family tree data
- `GET /api/horses/{id}/health` - Get health records
- `POST /api/horses/{id}/health-records` - Add health record
- `GET /api/horses/stallions` - List available stallions

## Tech Stack

### Backend
- Go 1.21+
- Gin web framework
- SQLite3 database
- Clean architecture

### Frontend
- React 18 with TypeScript
- Mantine v7 UI framework
- TanStack Query for data fetching
- React Router v6
- Day.js for date handling
- Dark theme optimized for readability

## Development Setup

### Prerequisites
- Node.js 18+
- Go 1.21+
- Git

### Getting Started
1. Clone the repository
2. Install frontend dependencies:
   ```bash
   cd frontend-react
   npm install
   ```
3. Start the frontend development server:
   ```bash
   npm run dev
   ```
4. In a new terminal, start the backend server:
   ```bash
   cd ../backend
   go run main.go
   ```

## Contributing

This is a private project. Please contact the project maintainers for contribution guidelines.
