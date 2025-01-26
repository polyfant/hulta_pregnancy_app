-- +goose Up
-- Create pregnancies table (depends on users and horses)
CREATE TABLE IF NOT EXISTS pregnancies (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    horse_id INTEGER NOT NULL,
    start_date DATE NOT NULL,
    expected_foaling_date DATE,
    actual_foaling_date DATE,
    status VARCHAR(50) DEFAULT 'active',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pregnancies_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_pregnancies_horse FOREIGN KEY (horse_id) REFERENCES horses(id) ON DELETE CASCADE
);

-- Create pregnancy events table (depends on pregnancies)
CREATE TABLE IF NOT EXISTS pregnancy_events (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    pregnancy_id INTEGER NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    event_date TIMESTAMP NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pregnancy_events_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_pregnancy_events_pregnancy FOREIGN KEY (pregnancy_id) REFERENCES pregnancies(id) ON DELETE CASCADE
);

-- Create breeding records table (depends on users and horses)
CREATE TABLE IF NOT EXISTS breeding_records (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    mare_id INTEGER NOT NULL,
    stallion_id INTEGER,
    external_stallion_name VARCHAR(255),
    breeding_date TIMESTAMP NOT NULL,
    success_status VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_breeding_records_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_breeding_records_mare FOREIGN KEY (mare_id) REFERENCES horses(id) ON DELETE CASCADE,
    CONSTRAINT fk_breeding_records_stallion FOREIGN KEY (stallion_id) REFERENCES horses(id) ON DELETE SET NULL
);

-- Create pre-foaling checklist table (depends on horses)
CREATE TABLE IF NOT EXISTS pre_foaling_checklist (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER NOT NULL,
    description TEXT NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    due_date TIMESTAMP NOT NULL,
    priority TEXT CHECK(priority IN ('HIGH', 'MEDIUM', 'LOW')) DEFAULT 'MEDIUM',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pre_foaling_checklist_horse FOREIGN KEY (horse_id) REFERENCES horses(id) ON DELETE CASCADE
);

-- Add indexes for pregnancy-related tables
CREATE INDEX idx_pregnancies_user_id ON pregnancies(user_id);
CREATE INDEX idx_pregnancies_horse_id ON pregnancies(horse_id);
CREATE INDEX idx_pregnancies_status ON pregnancies(status);
CREATE INDEX idx_pregnancy_events_user_id ON pregnancy_events(user_id);
CREATE INDEX idx_pregnancy_events_pregnancy_id ON pregnancy_events(pregnancy_id);
CREATE INDEX idx_pregnancy_events_type ON pregnancy_events(event_type);
CREATE INDEX idx_breeding_records_user_id ON breeding_records(user_id);
CREATE INDEX idx_breeding_records_mare_id ON breeding_records(mare_id);
CREATE INDEX idx_breeding_records_stallion_id ON breeding_records(stallion_id);
CREATE INDEX idx_breeding_records_date ON breeding_records(breeding_date);
CREATE INDEX idx_pre_foaling_checklist_horse_id ON pre_foaling_checklist(horse_id);
CREATE INDEX idx_pre_foaling_checklist_due_date ON pre_foaling_checklist(due_date);

-- +goose Down
DROP TABLE IF EXISTS pre_foaling_checklist;
DROP TABLE IF EXISTS breeding_records;
DROP TABLE IF EXISTS pregnancy_events;
DROP TABLE IF EXISTS pregnancies;
