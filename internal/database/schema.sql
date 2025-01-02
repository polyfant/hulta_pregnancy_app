-- Create tables
CREATE TABLE IF NOT EXISTS horses (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    breed VARCHAR(100),
    birth_date DATE,
    conception_date DATE,
    mother_id INTEGER REFERENCES horses(id),
    father_id INTEGER REFERENCES horses(id),
    is_pregnant BOOLEAN DEFAULT false,
    gender VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS health_records (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER REFERENCES horses(id),
    type VARCHAR(50) NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pregnancies (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER REFERENCES horses(id),
    user_id TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pregnancy_events (
    id SERIAL PRIMARY KEY,
    pregnancy_id INTEGER REFERENCES pregnancies(id),
    description TEXT NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pre_foaling_checklist (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER REFERENCES horses(id),
    description TEXT NOT NULL,
    is_completed BOOLEAN DEFAULT false,
    due_date DATE,
    priority VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_horses_user_id ON horses(user_id);
CREATE INDEX IF NOT EXISTS idx_pregnancies_horse_id ON pregnancies(horse_id);
CREATE INDEX IF NOT EXISTS idx_health_records_horse_id ON health_records(horse_id);
CREATE INDEX IF NOT EXISTS idx_pregnancy_events_pregnancy_id ON pregnancy_events(pregnancy_id);
CREATE INDEX IF NOT EXISTS idx_checklist_horse_id ON pre_foaling_checklist(horse_id);
