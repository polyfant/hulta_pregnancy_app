# Horse Pregnancy Tracking Application

## Overview

A comprehensive horse pregnancy tracking system with advanced stage monitoring and due date prediction.

## Key Features

### Pregnancy Tracking

-   Accurate stage calculation (Early, Mid, Late, Overdue)
-   Due date prediction with viable windows (320-370 days)
-   Overdue monitoring and high-risk detection
-   Progress tracking with detailed statistics
-   Comprehensive pregnancy timeline

## Core Features

-   Horse pregnancy tracking and monitoring
-   Pre-foaling checklist management
-   Health record keeping
-   Secure user authentication with Auth0
-   Real-time pregnancy status dashboard
-   Privacy-focused data management
-   API versioning (v1)

## Getting Started

### Prerequisites

-   Go 1.23.4+
-   Node.js 18+
-   PostgreSQL 13+
-   Auth0 account

### Environment Setup

1. Copy `.env.example` to `.env` and configure:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=horsetracker
   DB_NAME=horse_tracking_db
   DB_PASSWORD=your_password
   
   AUTH0_DOMAIN=your-domain.auth0.com
   AUTH0_AUDIENCE=your-api-identifier
   ```

2. Copy `frontend-react/env.template.txt` to `frontend-react/.env` and configure:
   ```env
   VITE_AUTH0_DOMAIN=your-domain.auth0.com
   VITE_AUTH0_CLIENT_ID=your-client-id
   VITE_AUTH0_AUDIENCE=your-api-identifier
   VITE_API_URL=http://localhost:8080
   ```

### Quick Start

1. Start PostgreSQL:
   ```bash
   # Using Docker
   docker run -d --name postgres \
     -e POSTGRES_USER=horsetracker \
     -e POSTGRES_PASSWORD=your_password \
     -e POSTGRES_DB=horse_tracking_db \
     -p 5432:5432 \
     postgres:13
   ```

2. Start the backend:
   ```bash
   go run cmd/server/main.go
   ```

3. Start the frontend:
   ```bash
   cd frontend-react
   npm install
   npm run dev
   ```

4. Access the application:
   - Frontend: http://localhost:5173
   - API: http://localhost:8080/api/v1

## Latest Updates (2025-01-19)

- ✨ Standardized API versioning with /api/v1 prefix
- 🔧 Improved database migrations using goose
- 🔒 Enhanced Auth0 integration
- 🚀 Updated frontend API client

For a complete list of changes, see [CHANGELOG.md](CHANGELOG.md)
For development roadmap, see [TODO.md](TODO.md)

## Documentation

- [API Documentation](docs/api.md)
- [Development Guide](docs/development.md)
- [Deployment Guide](docs/deployment.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
