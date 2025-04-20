# Database Verification Guide

## 1. Start Docker Containers
```bash
docker-compose up -d
```

## 2. Check Container Status
```bash
docker-compose ps
```
Verify `postgres` service is "Up"

## 3. Connect to PostgreSQL
```bash
# Connect to running container
docker-compose exec postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}

# Inside psql, run:
\l  # List databases
\c horse_tracking_db  # Connect to database
\dt  # List tables
```

## 4. Troubleshooting Checklist
- Confirm `.env` variables match:
  * `POSTGRES_DB`
  * `POSTGRES_USER`
  * `POSTGRES_PASSWORD`

- Check Docker logs
```bash
docker-compose logs postgres
```

## 5. Potential Issues to Watch
- Permissions
- Volume mounting
- Network configuration

## 6. Recommended Test Queries
```sql
-- Create a test table
CREATE TABLE test_horses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100)
);

-- Insert test data
INSERT INTO test_horses (name) VALUES ('Test Horse');

-- Verify insertion
SELECT * FROM test_horses;
```

## 7. Backup Strategy
- Regularly backup `postgres_data` volume
- Consider point-in-time recovery setup

🚨 SECURITY NOTE:
Never commit `.env` files with real credentials to version control!
