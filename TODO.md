# Current Priorities

## High Priority (MVP Target: Mid-February 2025)
Last Updated: 2025-01-19

-   [ ] Complete core ML features
    -   [ ] Pregnancy stage calculation 
        - [ ] Basic stage tracking (Early/Mid/Late/Pre-foaling)
        - [ ] Mobile-friendly stage visualization
        - [ ] Weather impact foundation
    -   [ ] Due date tracking 
        - [ ] Base calculation (340 days)
        - [ ] Mobile calendar view
        - [ ] Adjustable windows
    -   [ ] Pre-foaling checklist
        - [ ] Season-specific items
        - [ ] Mobile-friendly checkboxes
    -   [ ] Health monitoring during pregnancy
        - [ ] Vital signs tracking
        - [ ] Mobile-optimized data entry
-   [x] Implement privacy controls
    -   [x] Stage change alerts
    -   [x] Pre-foaling notifications
    -   [x] Health check reminders
    -   [x] Privacy change logging (Added 2025-01-15)
    -   [x] Data retention management
    -   [x] Automated cleanup
    -   [x] Transaction handling
-   [ ] Add notification system
    -   [ ] Pregnancy timeline view
    -   [ ] Horse status dashboard
    -   [ ] Easy data entry forms
-   [x] Finalize Auth0 integration (2025-01-19)
    -   [x] User authentication flow
    -   [x] API endpoint protection
    -   [x] Frontend token handling
    -   [ ] User roles and permissions (Pending)
    -   [ ] Role-based access control
-   [x] Mobile-First UI/UX 
    -   [x] Touch-friendly interface (Using Mantine UI)
    -   [x] Bottom navigation
    -   [ ] Offline capability
    -   [x] Swipe gestures (Mantine components)

## Immediate Next Steps (Week of Jan 20-26, 2025)

1. Frontend Improvements
   - [x] Add loading states (Using React Query)
   - [ ] Implement error boundaries
   - [x] Add form validation (horseValidation.ts implemented)
   - [x] Improve mobile responsiveness (Using Mantine UI)

2. Backend Enhancements
   - [x] Add request logging middleware (auth_middleware.go)
   - [ ] Implement rate limiting
   - [ ] Add API documentation using Swagger
   - [x] Set up automated backups (backup.go)

3. Testing & Quality
   - [ ] Add integration tests for API endpoints
   - [ ] Set up E2E testing with Cypress
   - [ ] Implement API response time monitoring
   - [ ] Add error tracking (e.g., Sentry)

## Medium Priority (Q2 2025)

- [ ] Complete breeding cost tracking
- [ ] Add nutrition tracking
- [ ] Implement expense tracking
- [ ] Add breeding statistics
- [ ] Enhance dashboard visualizations

## Specialized Horse Breeding Features (Q3-Q4 2025)

### Horse-Specific Timeline 🐎
- [ ] Mare's cycle visualization
  - [ ] Heat cycle tracking
  - [ ] Optimal breeding time calculator
  - [ ] Cycle history and patterns
- [ ] Accurate gestation tracking (340 days)
  - [ ] Trimester-based milestones
  - [ ] Expected foaling date range
  - [ ] Development stages visualization
- [ ] Season-aware breeding recommendations
  - [ ] Seasonal fertility patterns
  - [ ] Climate impact analysis
  - [ ] Breeding season planning

### Multi-Horse Management 🏡
- [ ] Herd Overview Dashboard
  - [ ] Status at a glance
  - [ ] Group management
  - [ ] Health alerts
- [ ] Advanced Breeding Management
  - [ ] Breeding pairs tracking
  - [ ] Genetic compatibility checks
  - [ ] Success rate analytics
- [ ] Family Tree Visualization
  - [ ] Interactive pedigree charts
  - [ ] Genetic trait tracking
  - [ ] Inbreeding coefficient calculator
- [ ] Stallion Management
  - [ ] Breeding schedule
  - [ ] Performance history
  - [ ] Offspring tracking

### Specialized Tools 🛠️
- [ ] Weather Impact Analysis
  - [ ] Local weather integration
  - [ ] Temperature stress monitoring
  - [ ] Exercise recommendations
- [ ] Veterinary Care
  - [ ] Appointment scheduling
  - [ ] Vaccination tracking
  - [ ] Health record timeline
- [ ] Foaling Preparation
  - [ ] Dynamic checklists
  - [ ] Supply inventory
  - [ ] Emergency contact info
- [ ] Competition Integration
  - [ ] Show schedule management
  - [ ] Training milestone tracking
  - [ ] Rest period calculations

### Breed-Specific Features 📊
- [ ] Breed Standards Database
  - [ ] Growth charts by breed
  - [ ] Expected gestation variations
  - [ ] Common health considerations
- [ ] Custom Health Monitoring
  - [ ] Breed-specific vital ranges
  - [ ] Known genetic issues
  - [ ] Preventive care recommendations

## Technical Improvements

- [x] Add more unit tests
- [x] Standardize API versioning (2025-01-19)
- [x] Implement proper database migrations (2025-01-19)
- [ ] Set up CI/CD pipeline
- [ ] Set up monitoring and alerting

## Documentation Needs

- [ ] API Documentation
- [ ] User Guide
- [ ] Development Setup Guide
- [ ] Deployment Guide
- [ ] Contributing Guidelines

## Completed 

- [x] Privacy dashboard implementation
- [x] Data retention controls
- [x] Transaction handling for privacy updates
- [x] Comprehensive test coverage for privacy features
- [x] Finalize Auth0 integration (2025-01-19)
- [x] Standardize API versioning (2025-01-19)
- [x] Implement proper database migrations (2025-01-19)
- [x] Mobile-first UI with Mantine
- [x] Form validation system
- [x] Automated database backups
