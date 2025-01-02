# Horse Tracking Frontend

A comprehensive horse management system with focus on breeding and pregnancy tracking.

## Features

- Horse Management
  - Basic information tracking
  - Gender-specific features
  - Family tree visualization
- Pregnancy Tracking
  - Progress monitoring with visual indicators
  - Due date calculations
  - Pregnancy stage guidelines
- Health Records
  - Medical history
  - Vaccination tracking
  - Health assessments

## Tech Stack

- React 18 with TypeScript
- Vite for fast development and building
- Mantine v7 for UI components
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

- Icons: See [ICONS.md](docs/ICONS.md) for icon usage guidelines
- Components: Follow Mantine component patterns
- Data Fetching: Use TanStack Query for API calls
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
