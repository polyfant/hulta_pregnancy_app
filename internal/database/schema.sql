-- Horses table
CREATE TABLE IF NOT EXISTS horses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    breed TEXT,
    date_of_birth TEXT NOT NULL,
    weight REAL,
    conception_date TEXT,
    mother_id INTEGER,
    father_id INTEGER,
    FOREIGN KEY (mother_id) REFERENCES horses(id),
    FOREIGN KEY (father_id) REFERENCES horses(id)
);

-- Health records table
CREATE TABLE IF NOT EXISTS health_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    horse_id INTEGER NOT NULL,
    date TEXT NOT NULL,
    type TEXT NOT NULL,
    notes TEXT,
    FOREIGN KEY (horse_id) REFERENCES horses(id)
);

-- Pregnancy events table
CREATE TABLE IF NOT EXISTS pregnancy_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    horse_id INTEGER NOT NULL,
    date TEXT NOT NULL,
    description TEXT NOT NULL,
    FOREIGN KEY (horse_id) REFERENCES horses(id)
);

-- Breeding costs table
CREATE TABLE IF NOT EXISTS breeding_costs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    horse_id INTEGER NOT NULL,
    description TEXT NOT NULL,
    amount REAL NOT NULL,
    date TEXT NOT NULL,
    FOREIGN KEY (horse_id) REFERENCES horses(id)
);

-- User sync table (for future use)
CREATE TABLE IF NOT EXISTS user_sync (
    user_id INTEGER PRIMARY KEY,
    last_sync TEXT NOT NULL
);
