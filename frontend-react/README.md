# Horse Tracking Frontend

A comprehensive horse management system focusing on breeding and pregnancy tracking.

## Features

- Authentication & Authorization
  - Auth0 integration
  - Secure routes
  - User-specific views
- Horse Management
  - Basic information tracking
  - Gender-specific features (Mare/Stallion/Gelding)
  - Family tree visualization with expandable nodes
  - Image support
- Pregnancy Tracking
  - Progress monitoring with visual indicators
  - Due date calculations
  - Trimester-based guidelines
  - Detailed stage information
  - Care recommendations
- Health Records
  - Medical history
  - Vaccination tracking
  - Health assessments

## Tech Stack

- React 18 with TypeScript
- Vite for fast development and building
- Mantine v7 for UI components
  - Dark theme implementation
  - Responsive layouts
  - Custom component styling
- Auth0 for authentication
- Phosphor Icons for consistent iconography
- TanStack Query for data fetching
- Day.js for date handling

## Getting Started

1. Install dependencies:

```bash
npm install
```

2. Create .env file from template:

```bash
cp env.template.txt .env
```

3. Configure environment variables in .env:

```env
VITE_AUTH0_DOMAIN=your-auth0-domain
VITE_AUTH0_CLIENT_ID=your-auth0-client-id
VITE_API_URL=http://localhost:8080
```

4. Start development server:

```bash
npm run dev
```

## Development Guidelines

- Authentication: Auth0 handles user authentication
- Icons: See [ICONS.md](docs/ICONS.md) for icon usage guidelines
- Components: Follow Mantine component patterns
- Data Fetching: Use TanStack Query for API calls
- Theming: Use dark theme consistently
- Testing: Write tests for critical components

## Testing

Run tests with:

```bash
npm test
```

## Building for Production

```bash
npm run build
```

## Contributing

1. Follow the established code style
2. Write meaningful commit messages
3. Add tests for new features
4. Update documentation as needed
