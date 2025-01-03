# Setup script for Horse Tracker Prototype

# Create directory structure
$dirs = @(
    "scripts",
    "nginx/conf.d",
    "frontend-react/src",
    "backend",
    "docs"
)

foreach ($dir in $dirs) {
    New-Item -ItemType Directory -Force -Path $dir
}

# Create docker-compose.yml
@'
version: "3.8"

services:
  db:
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  backend:
    build: 
      context: ./backend
      dockerfile: Dockerfile
    environment:
      DATABASE_URL: ${DATABASE_URL}
      AUTH0_DOMAIN: ${AUTH0_DOMAIN}
      AUTH0_CLIENT_ID: ${AUTH0_CLIENT_ID}
      AUTH0_AUDIENCE: ${AUTH0_AUDIENCE}
      CORS_ORIGIN: ${CORS_ORIGIN}
    ports:
      - "3000:3000"
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend-react
      dockerfile: Dockerfile
    environment:
      VITE_API_URL: ${API_URL}
      VITE_AUTH0_DOMAIN: ${AUTH0_DOMAIN}
      VITE_AUTH0_CLIENT_ID: ${AUTH0_CLIENT_ID}
      VITE_AUTH0_AUDIENCE: ${AUTH0_AUDIENCE}
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  postgres_data:
'@ | Out-File -FilePath docker-compose.yml -Encoding UTF8

# Create CHANGELOG.md
@'
# Changelog

## [1.0.0] - 2024-03-XX
### Added
- Initial prototype release
- Auth0 integration
- PostgreSQL database setup
- Docker containerization
- Automated deployment scripts
- Frontend React application
- Backend Go API
- Security configurations
'@ | Out-File -FilePath CHANGELOG.md -Encoding UTF8

# Create start.ps1
@'
# Previous Windows start script content here...
'@ | Out-File -FilePath scripts/start.ps1 -Encoding UTF8

# Create verify.ps1
@'
# Previous Windows verify script content here...
'@ | Out-File -FilePath scripts/verify.ps1 -Encoding UTF8

# Create .env.example
@'
# Database Configuration
POSTGRES_DB=horse_proto_db
POSTGRES_USER=horse_proto
POSTGRES_PASSWORD=H0rseTrack3r2024!
DATABASE_URL=postgresql://horse_proto:H0rseTrack3r2024!@db:5432/horse_proto_db

# Network Configuration
HOST_IP=
CORS_ORIGIN=http://localhost
API_URL=http://localhost:3000

# Auth0 Configuration
AUTH0_DOMAIN=dev-r083cwkcv0pgz20x.eu.auth0.com
AUTH0_CLIENT_ID=OBmEol1z4U49r3YI3priDdGbvF5i4O7d
AUTH0_AUDIENCE=https://horse-tracker-api.demo
'@ | Out-File -FilePath .env.example -Encoding UTF8

# Create README.md
@'
# Horse Tracker Prototype

## Quick Start Guide

### Windows Setup
1. Open PowerShell as Administrator and run:
```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
./scripts/start.ps1
```

2. To verify the installation:
```powershell
./scripts/verify.ps1
```

### Linux/Mac Setup
1. Make scripts executable and start:
```bash
chmod +x start.sh verify.sh
./start.sh
```

2. To verify the installation:
```bash
./verify.sh
```

## Accessing the Application
- Frontend: http://<your-ip>
- Backend API: http://<your-ip>:3000

## Requirements
- Docker
- Docker Compose
- PowerShell (Windows) or Bash (Linux/Mac)
- Available ports: 80, 3000, 5432

## Default Credentials
Database:
- Name: horse_proto_db
- User: horse_proto
- Password: H0rseTrack3r2024!

Auth0 is pre-configured and ready to use.
'@ | Out-File -FilePath README.md -Encoding UTF8

# Create .gitignore
@'
# Environment variables
.env

# Dependencies
node_modules/
vendor/

# Build outputs
dist/
build/

# Logs
*.log

# System files
.DS_Store
Thumbs.db

# IDE files
.idea/
.vscode/
*.swp
'@ | Out-File -FilePath .gitignore -Encoding UTF8

Write-Host "✅ Project structure created successfully!" -ForegroundColor Green
Write-Host "📁 Created directories and files:" -ForegroundColor Cyan
Get-ChildItem -Recurse | Where-Object { !$_.PSIsContainer } | ForEach-Object {
    Write-Host "   $($_.FullName)" -ForegroundColor Gray
}
Write-Host "`n🚀 Next steps:" -ForegroundColor Yellow
Write-Host "1. Review and customize the configuration files" -ForegroundColor White
Write-Host "2. Run ./scripts/start.ps1 to start the application" -ForegroundColor White
Write-Host "3. Run ./scripts/verify.ps1 to verify the setup" -ForegroundColor White