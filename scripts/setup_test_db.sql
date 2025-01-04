-- Create test user
CREATE USER testuser WITH PASSWORD 'testpassword';

-- Create test database
CREATE DATABASE horse_tracking_test;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE horse_tracking_test TO testuser;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO testuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO testuser;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO testuser;
