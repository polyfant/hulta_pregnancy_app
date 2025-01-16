# Tech Stack

## Backend Architecture

### Core Components
- Go 1.23.4+ for server implementation
- PostgreSQL for data persistence
- GORM for ORM and database operations
- Gin for HTTP routing and middleware
- WebSocket for real-time updates
- JWT for authentication

### Project Structure
```
horse_tracking_go/
├── cmd/                    # Application entrypoints
│   └── server/            # Main server binary
├── internal/              # Private application code
│   ├── api/              # HTTP handlers and routes
│   ├── models/           # Domain models and types
│   ├── repository/       # Data persistence layer
│   ├── service/          # Business logic layer
│   │   ├── breeding/     # Breeding management
│   │   ├── health/       # Health records & monitoring
│   │   ├── horse/        # Horse management
│   │   ├── notification/ # Notification system
│   │   ├── pregnancy/    # Pregnancy tracking
│   │   ├── privacy/      # Privacy controls
│   │   ├── vitals/       # Vital signs monitoring
│   │   └── weather/      # Weather tracking & alerts
│   └── websocket/       # Real-time communication
├── config/               # Configuration management
├── scripts/             # Build and deployment scripts
└── tests/               # Integration tests
```

## Frontend

- React 18 with TypeScript
- Vite for blazing fast builds
- Mantine v7 for UI components
- TanStack Query for data management
- Phosphor Icons for consistent iconography
- Chart.js & Nivo for visualizations
- Day.js for date handling

## ML & Analytics

- TensorFlow.js for on-device ML
- Chart.js for data visualization
- Custom ML models for growth prediction

## Privacy & Security

- Web Crypto API for encryption
- LocalStorage for offline-first data
- Privacy-focused architecture
- Granular data controls
