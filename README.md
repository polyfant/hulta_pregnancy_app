# Hulta Pregnancy App 🐎

A specialized horse breeding and pregnancy tracking application for Hulta Farm, designed to help horse breeders manage breeding programs, monitor pregnancies, and track health metrics.

## Features 🚀

-   **Horse Management**

    -   Comprehensive horse profiles with detailed information
    -   Pedigree tracking
    -   Health status monitoring
    -   Premium and champion designation

-   **Pregnancy Tracking**

    -   Timeline visualization with key milestones
    -   Stage-based progress tracking
    -   Due date calculation (340 days)
    -   Pre-foaling checklist

-   **Breeding Management**

    -   Breeding records and history
    -   Success rate statistics
    -   Mare and stallion performance metrics
    -   Optimal breeding time recommendations

-   **Health Monitoring**

    -   Veterinary appointment scheduling
    -   Health check reminders
    -   Vital sign tracking
    -   ML-powered growth predictions

-   **Environmental Monitoring**

    -   Weather impact analysis
    -   Temperature and humidity tracking
    -   Environmental risk assessment
    -   Care recommendations based on conditions

-   **Notifications**

    -   Multi-channel alerts (WebSocket, Email, SMS)
    -   Customizable notification preferences
    -   Critical alerts for health issues
    -   Stage change notifications

-   **Privacy Controls**
    -   Granular data sharing settings
    -   Data retention management
    -   Privacy auditing
    -   Data export and purging

## Tech Stack 💻

### Frontend

-   React 18 with TypeScript
-   shadcn/ui components
-   React Query for state management
-   Framer Motion for animations
-   Zod for form validation
-   Auth0 for authentication

### Backend

-   Go 1.23.4+
-   Gin Web Framework
-   GORM for database interactions
-   PostgreSQL database
-   Auth0 JWT validation
-   Rate limiting middleware

### DevOps

-   Docker for containerization
-   GitHub Actions for CI/CD (planned)
-   Automated backups
-   Structured logging

## Getting Started 🏁

### Prerequisites

-   Node.js 20+
-   Go 1.23.4+
-   Docker and Docker Compose
-   PostgreSQL

### Setup Instructions

1. Clone the repository:



2. Set up environment variables:

    ```bash
    cp .env.example .env
    # Edit .env with your database and Auth0 credentials
    ```

3. Start the backend:

    ```bash
    cd backend
    go mod download
    go run cmd/main.go
    ```

4. Start the frontend:

    ```bash
    cd frontend-react
    npm install
    npm run dev
    ```

5. Alternatively, use Docker Compose:
    ```bash
    docker-compose up
    ```

### Auth0 Configuration

To use authentication features, you'll need to set up an Auth0 application and API:

1. Create an Auth0 application (Single Page Application)
2. Set up an API in Auth0
3. Configure the following environment variables:
    - `VITE_AUTH0_DOMAIN`
    - `VITE_AUTH0_CLIENT_ID`
    - `VITE_AUTH0_AUDIENCE`

## Project Structure 📂

```
├── backend/
│   ├── cmd/                # Application entrypoints
│   ├── internal/           # Internal packages
│   │   ├── api/            # API handlers
│   │   ├── middleware/     # HTTP middleware
│   │   ├── models/         # Database models
│   │   ├── repositories/   # Data access layer
│   │   └── services/       # Business logic
│   └── migrations/         # Database migrations
├── frontend-react/
│   ├── public/             # Static assets
│   ├── src/
│   │   ├── api/            # API client and hooks
│   │   ├── auth/           # Authentication
│   │   ├── components/     # UI components
│   │   ├── contexts/       # React contexts
│   │   ├── hooks/          # Custom hooks
│   │   ├── pages/          # Page components
│   │   ├── types/          # TypeScript type definitions
│   │   └── utils/          # Utility functions
│   └── index.html          # HTML entry point
└── docker-compose.yml      # Docker Compose configuration
```


## License 📝

This project is proprietary and not licensed for public use or redistribution.

## Acknowledgements 👏

-   The Hulta Farm team for their domain expertise
-   All contributors who have helped build this application
