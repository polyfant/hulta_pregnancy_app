# Tech Stack

## Frontend

### Core Technologies

-   React 18
-   TypeScript
-   shadcn/ui (UI components)
-   React Query
-   Vite
-   Yarn
-   MSW (Mock Service Worker)
-   Vitest

### UI & Styling

-   **Component Library**: shadcn/ui
    -   Based on Radix UI primitives
    -   TailwindCSS for styling
    -   Consistent theming with CSS variables
    -   Dark theme implementation
-   **Icons**
    -   **Current Library**: lucide-react
        -   Lightweight and consistent icon set
        -   Minimal bundle size
        -   Consistent styling across components
        -   Standardized naming convention
-   **Date Handling**: Day.js
-   **CSS Strategy**:
    -   TailwindCSS
    -   CSS Variables for theming
    -   Component-specific styles using cn utility

### Authentication

-   **Provider**: Auth0
    -   User authentication
    -   Protected routes
    -   User profile management
    -   Role-based access control

### State Management & Data Fetching

-   **API Client**: TanStack Query (React Query) v5
-   **Forms**: react-hook-form with zod validation
-   **Global State**: React Context (small scale)
    -   Consider Zustand if needed for larger state

## Backend

-   Go 1.23.4+ for high-performance API
-   Gin Web Framework
-   GORM for ORM
-   PostgreSQL for database
-   Redis for caching
-   Zap for structured logging
-   Docker for containerization
-   GitHub Actions for CI/CD
-   Auth0 for authentication & authorization
-   Swagger for API documentation
-   Rate limiting middleware
-   CORS protection
-   SSL/TLS encryption
-   Secure headers middleware

## ML & Analytics

-   TensorFlow.js for on-device ML
-   Chart.js for data visualization
-   Custom ML models for growth prediction

## Privacy & Security

-   Web Crypto API for encryption
-   LocalStorage for offline-first data
-   Privacy-focused architecture
-   Granular data controls

## Testing

-   Vitest
-   React Testing Library
-   MSW for API mocking
