# Horse Tracking Application

A modern web application for tracking horses, breeding, and health records. This application helps horse owners and breeders manage their horses' information, health records, breeding costs, and pregnancy monitoring.

## Current Features

### Horse Management
- Add new horses with details (name, breed, gender, birth date)
- View all horses in a clean, sortable table
- Edit existing horse information
- Delete horses from the system

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

### Breeding Management
- Track breeding costs
- Record breeding-related expenses
- Monitor breeding history
- Breeding schedule management

### Pregnancy Monitoring
- Track pregnancy milestones
- Get pregnancy guidelines
- Monitor conception and due dates

## API Endpoints

### Currently Implemented
- `GET /api/horses` - List all horses
- `POST /api/horses` - Create a new horse
- `PUT /api/horses/{id}` - Update a horse
- `DELETE /api/horses/{id}` - Delete a horse
- `GET /api/horses/{id}` - Get details of a specific horse

### Planned Endpoints
- `GET /api/horses/{id}/family-tree` - Get family tree for a horse
- `GET /api/horses/{id}/health` - Get health summary
- `POST /api/horses/{id}/health-records` - Add health record
- `GET /api/horses/{id}/breeding-costs` - View breeding costs
- `POST /api/horses/{id}/breeding-costs` - Add breeding cost
- `GET /api/horses/{id}/pregnancy-guidelines` - Get pregnancy guidelines

## Project Structure
```
horse_tracking_go/
├── backend/              # Go backend server
│   ├── cmd/             # Application entry points
│   ├── internal/        # Internal packages
│   └── pkg/             # Public packages
└── frontend-react/      # React frontend application
    ├── src/
    │   ├── components/  # React components
    │   ├── api/         # API integration
    │   └── types/       # TypeScript types
    └── public/          # Static assets
```

## Tech Stack
### Backend
- Go 1.21+
- Gin web framework
- SQLite database

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

### Installation & Development
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd horse_tracking_go
   ```

2. Start the backend:
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

3. Start the frontend:
   ```bash
   cd frontend-react
   npm install
   npm run dev
   ```

4. Access the application:
   - Frontend: `http://localhost:5173`
   - Backend API: `http://localhost:8080`

## Testing
### Frontend Tests
```bash
cd frontend-react
npm test           # Run tests in watch mode
npm run test:ui    # Run tests with UI
npm run coverage   # Generate coverage report
```

### Backend Tests
```bash
cd backend
go test ./...
```

## Development Status
- ✅ Basic CRUD operations for horses
- ✅ Modern UI with Mantine components
- ✅ Comprehensive test suite
- 🚧 Family tree visualization (In Progress)
- 📅 Health tracking (Planned)
- 📅 Breeding management (Planned)
- 📅 Pregnancy monitoring (Planned)

## Comprehensive Breeding & Monitoring Features

### 🗓️ Pregnancy Timeline Tracking
- Interactive timeline showing current stage and upcoming milestones
- Smart notifications for critical dates and actions
- Customizable vaccination schedule with reminders
- Stage-specific nutritional guidelines with:
  - Daily feed recommendations
  - Supplement requirements
  - Hydration monitoring
  - Weight gain targets

### 🏥 Health Monitoring System
- Beginner-friendly daily health checklist
- Photo-based guide for visual health indicators
- Temperature and vital signs tracking
- Exercise planning with:
  - Safe activity recommendations by trimester
  - Warning signs to watch for
  - Automated exercise logs
- Weight tracking with visual progress charts

### 📚 Educational Hub
- Step-by-step guides for new breeders
- Video tutorials for common procedures
- Interactive learning modules covering:
  - Basic mare care
  - Pregnancy stages
  - Common complications
  - Emergency situations
- Direct vet consultation booking
- Community Q&A section

### 📸 Documentation Center
- Structured photo documentation:
  - Weekly condition photos
  - Ultrasound image storage
  - Medical record attachments
- Digital veterinary record keeping
- Expense tracking categories:
  - Routine care
  - Emergency services
  - Medications
  - Supplements

### 👶 Pre and Post-Foaling Guide
- Interactive foaling preparation checklist
- Emergency contact card generator
- Foaling kit inventory manager
- Post-birth care timeline:
  - First 24 hours checklist
  - First week monitoring
  - Mare recovery tracking
  - Foal development milestones

### 💰 Financial Planning Tools
- Comprehensive cost calculator
- Insurance requirement checker
- Emergency fund planning guide
- Budget templates for:
  - Pre-breeding costs
  - Pregnancy care
  - Foaling expenses
  - Post-birth care

All features include:
- Beginner-friendly explanations
- Visual guides and references
- Emergency action plans
- Professional veterinary guidelines
- Community support integration
- Mobile-friendly interface
- Offline access to critical information