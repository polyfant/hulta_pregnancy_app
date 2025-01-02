# Technology Stack & Standards 🦍

## Frontend Stack

### Core

-   **Framework**: React 18
-   **Language**: TypeScript 5.2+
-   **Build Tool**: Vite
-   **Package Manager**: npm (standardized across team)

### UI & Styling

-   **Component Library**: Mantine v7
    -   No Material UI or other UI libraries
    -   Consistent use of Mantine hooks and utilities
-   **Icons**

    -   **Current Library**: @phosphor-icons/react

        -   Modern, consistent design
        -   Tree-shakeable imports
        -   Direct imports from package (no aliasing)
        -   Example:

            ```typescript
            import { Plus, Heart, Horse } from '@phosphor-icons/react';

            // In JSX
            <Button leftSection={<Plus size={16} />}>Add New</Button>;
            ```

-   **Date Handling**: Day.js
-   **CSS Strategy**:
    -   Mantine's built-in styling system
    -   CSS Modules for custom components

### State Management & Data Fetching

-   **API Client**: TanStack Query (React Query) v5
-   **Forms**: Mantine form hooks
-   **Global State**: React Context (small scale)
    -   Consider Zustand if needed for larger state

### Testing

-   **Unit Testing**: Vitest
-   **Component Testing**: React Testing Library
-   **E2E Testing**: Playwright
-   **Coverage**: Istanbul

### Development Tools

-   **Linting**: ESLint with TypeScript rules
-   **Formatting**: Prettier
-   **Git Hooks**: husky + lint-staged
-   **VSCode Extensions**:
    -   ESLint
    -   Prettier
    -   TypeScript + JavaScript
    -   Mantine snippets

## Backend Stack

### Core

-   **Language**: Go 1.21+
-   **Framework**: Gin
-   **Build Tool**: Go build system

### Database

-   **Primary DB**: SQLite (with LiteStream for replication)
-   **Query Builder**: Native SQL (no ORM)
-   **Migrations**: golang-migrate

### API

-   **Style**: RESTful
-   **Documentation**: OpenAPI/Swagger
-   **Validation**: go-playground/validator
-   **Middleware**: Gin built-in

### Testing

-   **Testing Framework**: Go testing package
-   **Mocking**: go-mock
-   **Coverage**: Go test coverage

### Development Tools

-   **Hot Reload**: Air
-   **Linting**: golangci-lint
-   **Formatting**: gofmt
-   **VSCode Extensions**:
    -   Go
    -   Go Test Explorer
    -   SQLite Viewer

## Infrastructure

### Development Environment

-   **OS Support**: Cross-platform (Windows, macOS, Linux)
-   **Container**: Docker (development only)
-   **Local DB**: SQLite
-   **API Testing**: Postman/Thunder Client

### Production Environment (Planned)

-   **Hosting**: DigitalOcean/Linode
-   **Database**: SQLite + LiteStream
-   **Backup**: S3-compatible storage
-   **Monitoring**: Prometheus + Grafana
-   **Logging**: Zerolog

## Standards & Conventions

### Code Style

-   **Frontend**:

    ```typescript
    // File naming
    components / MyComponent.tsx;
    hooks / useMyHook.ts;
    utils / myUtil.ts;

    // Component structure
    export function ComponentName({ prop1, prop2 }: Props) {
    	// hooks first
    	// business logic
    	// render
    }
    ```

-   **Backend**:

    ```go
    // File naming
    handler_name.go
    model_name.go

    // Function naming
    func (s *Service) HandleSomething() error
    ```

### Git Workflow

-   **Branch Naming**:
    -   feature/description
    -   fix/description
    -   refactor/description
-   **Commit Messages**: Conventional Commits
    -   feat: description
    -   fix: description
    -   refactor: description

### API Standards

-   **Endpoints**:
    -   RESTful naming
    -   Versioned (/api/v1/...)
-   **Response Format**:
    ```json
    {
    	"data": {},
    	"error": null,
    	"metadata": {}
    }
    ```

### Error Handling

-   **Frontend**:
    -   Global error boundary
    -   Toast notifications
    -   Form-level validation
-   **Backend**:
    -   Structured error responses
    -   Logging with correlation IDs
    -   Panic recovery middleware

### Performance Standards

-   **Frontend**:
    -   Lighthouse score > 90
    -   Bundle size < 250KB initial load
    -   First paint < 1.5s
-   **Backend**:
    -   Response time < 100ms
    -   Memory usage < 100MB
    -   CPU usage < 50%

## Development Workflow

1. **Starting New Feature**

    ```bash
    # Frontend
    cd frontend-react
    npm install
    npm run dev

    # Backend
    cd ../
    go mod download
    go run cmd/server/main.go
    ```

2. **Testing**

    ```bash
    # Frontend
    npm run test
    npm run test:e2e

    # Backend
    go test ./...
    ```

3. **Building**

    ```bash
    # Frontend
    npm run build

    # Backend
    go build -o bin/server cmd/server/main.go
    ```

## Notes

-   Always check this document before introducing new dependencies
-   Update this document when making architectural decisions
-   Discuss major deviations with the team
-   Keep dependencies up to date, but stable

## Version Control

-   Last Updated: 2024-12-29
-   Next Review: 2025-01-29
