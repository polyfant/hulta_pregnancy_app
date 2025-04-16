# Multi-Tenant Database Security

## Architecture Overview

-   Unique schema per user
-   Strict access controls
-   JWT-based authentication
-   Row-level security

## Implementation Strategy

### 1. Schema Creation

```sql
-- User registration triggers schema creation
CREATE SCHEMA user_${USER_ID};
CREATE USER user_${USER_ID} WITH PASSWORD 'secure_generated_password';

-- Grant minimal, specific permissions
GRANT USAGE ON SCHEMA user_${USER_ID} TO user_${USER_ID};
GRANT ALL PRIVILEGES ON SCHEMA user_${USER_ID} TO user_${USER_ID};
```

### 2. Table Creation

```sql
-- Tables created in user-specific schema
CREATE TABLE user_${USER_ID}.horses (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    user_id UUID REFERENCES users(id),
    -- Other horse-specific fields
);

CREATE TABLE user_${USER_ID}.pregnancy_records (
    id UUID PRIMARY KEY,
    horse_id UUID REFERENCES user_${USER_ID}.horses(id),
    -- Pregnancy tracking fields
);
```

### 3. Row-Level Security

```sql
-- Enforce user can only see their own data
ALTER TABLE user_${USER_ID}.horses ENABLE ROW LEVEL SECURITY;
CREATE POLICY user_horses_policy ON user_${USER_ID}.horses
    FOR ALL
    USING (user_id = current_user_id());
```

## Authentication Flow

1. User logs in via Auth0
2. Backend generates user-specific schema
3. JWT validates user permissions
4. Database connection uses user-specific credentials

## Security Considerations

-   Dynamically generate strong passwords
-   Rotate user credentials
-   Implement connection pooling
-   Use prepared statements
-   Implement strict input validation

## Monitoring & Auditing

-   Log schema creations
-   Track user access
-   Monitor unusual database activities

🔒 SECURITY BEST PRACTICES

-   Never store plain-text passwords
-   Use connection pooling
-   Implement strict input sanitization
-   Regular security audits

## Current Priorities

### UI Migration (High Priority) 🎨

-   [ ] Set up shadcn/ui and Tailwind configuration
-   [ ] Migrate existing components from Mantine
-   [ ] Implement new horse-focused design system
-   [ ] Create component test suite for new components

### New Features Implementation 🚀

-   [ ] Breeding Calendar Integration
-   [ ] Health Monitoring Dashboard
-   [ ] Pregnancy Timeline Visualization
-   [ ] Comprehensive Breeding Checklist

## Phase 1: Core UI Components 🎯

-   [ ] Button -> shadcn/ui Button
-   [ ] Card -> shadcn/ui Card
-   [ ] Input -> shadcn/ui Input
-   [ ] Select -> shadcn/ui Select
-   [ ] Modal -> shadcn/ui Dialog
-   [ ] Tabs -> shadcn/ui Tabs

## Phase 2: Layout Components 🏗️

-   [ ] AppShell -> DashboardLayout
-   [ ] Navbar -> Navigation
-   [ ] Sidebar -> Custom Sidebar
-   [ ] Grid -> Tailwind Grid

## Phase 3: Form Components 📝

-   [ ] TextInput -> Input
-   [ ] NumberInput -> Input type="number"
-   [ ] Select -> Select/Combobox
-   [ ] MultiSelect -> MultiSelect
-   [ ] DatePicker -> Calendar

## Phase 4: Specialized Components 🐎

-   [ ] StageVisualization
-   [ ] PregnancyTimeline
-   [ ] HealthMonitor
-   [ ] BreedingCalendar
