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
