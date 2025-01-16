# Current Priorities

## High Priority (MVP Target: Mid-February 2025)

-   [x] Complete core ML features
    -   [x] Pregnancy stage calculation 
        - [x] Basic stage tracking (Early/Mid/Late/Pre-foaling)
        - [x] Mobile-friendly stage visualization
        - [x] Weather impact foundation
    -   [x] Due date tracking 
        - [x] Base calculation (340 days)
        - [x] Mobile calendar view
        - [x] Adjustable windows
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
-   [ ] Finalize Auth0 integration
    -   [ ] User roles and permissions
    -   [ ] Secure data access
-   [ ] Mobile-First UI/UX 
    -   [ ] Touch-friendly interface
    -   [ ] Bottom navigation
    -   [ ] Offline capability
    -   [ ] Swipe gestures

## Medium Priority (Q2 2025)

- [ ] Complete breeding cost tracking
- [ ] Add nutrition tracking
- [ ] Implement expense tracking
- [ ] Add breeding statistics
- [ ] Enhance dashboard visualizations

## Technical Improvements

- [x] Add more unit tests
    - [x] Pregnancy calculations
    - [x] Stage transitions
    - [x] Weather impact calculations
    - [x] Heat index formulas
    - [x] Notification triggers
    - [x] Privacy preference tests (Added 2025-01-15)
    - [x] Data retention tests
    - [x] Transaction handling tests
- [x] Improve error handling
    - [x] Transaction rollback
    - [x] Privacy change logging
    - [x] Data cleanup errors
- [x] Improve service organization
    - [x] Split interfaces into domain-specific files
    - [x] Rename generic service files to be more specific
    - [x] Move constants to domain-specific files
    - [x] Consolidate duplicate service implementations

## Future Considerations

- [ ] Add data export functionality
- [ ] Add batch operations
- [ ] Advanced statistics and reporting
- [ ] Mobile app integration

## Feature Plan

### Future Enhancements

1. Custom risk thresholds for different horse breeds
2. Integration with local weather stations for more accurate data
3. Breed-specific growth curve analysis and recommendations
4. Advanced weather impact analysis
   - Regional climate adaptation
   - Seasonal comfort adjustments
   - Historical weather pattern analysis

### Privacy & Control

-   All environmental and location features are opt-in only
-   Data is stored locally by default
-   Clear data cleanup when features are disabled
-   Transparent settings UI for feature control

## Privacy & Security

-   [x] End-to-end encryption
-   [x] Privacy dashboard
-   [x] Automated audits
-   [x] Data masking
-   [ ] Two-factor authentication
-   [ ] Backup encryption

## Testing Priority

### Unit Tests

-   [ ] Privacy Controls
    -   [ ] Data masking functions
    -   [ ] Encryption/decryption
    -   [ ] Privacy settings validation
-   [ ] ML Features
    -   [ ] Growth predictions
    -   [ ] Environmental impact calculations
    -   [ ] Health monitoring algorithms
-   [ ] Data Management
    -   [ ] Export functionality
    -   [ ] Data cleanup routines
    -   [ ] Audit logging

### Integration Tests

-   [ ] Privacy Dashboard
    -   [ ] Real-time score updates
    -   [ ] Settings synchronization
    -   [ ] Data visualization accuracy
-   [ ] Environmental Monitoring
    -   [ ] Weather data integration
    -   [ ] Impact calculations
    -   [ ] Alert system

### E2E Tests

-   [ ] Complete privacy workflow
    -   [ ] Settings configuration
    -   [ ] Data masking in UI
    -   [ ] Export/purge functionality
-   [ ] ML feature workflow
    -   [ ] Growth prediction accuracy
    -   [ ] Environmental impact assessment
    -   [ ] Health monitoring alerts

### Performance Tests

-   [ ] ML model loading times
-   [ ] Dashboard rendering performance
-   [ ] Data encryption/decryption speed
-   [ ] Audit log processing

## Completed 

### January 2025
- [x] Privacy change logging system
- [x] Data retention controls
- [x] Transaction handling for privacy updates
- [x] Comprehensive test coverage for privacy features
