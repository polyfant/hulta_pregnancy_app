-- +goose Up
CREATE TABLE pre_foaling_checklist (
    id SERIAL PRIMARY KEY,
    horse_id INTEGER NOT NULL,
    description TEXT NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    due_date TIMESTAMP NOT NULL,
    priority TEXT CHECK(priority IN ('HIGH', 'MEDIUM', 'LOW')) DEFAULT 'MEDIUM',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (horse_id) REFERENCES horses(id) ON DELETE CASCADE
);

CREATE INDEX idx_pre_foaling_checklist_horse_id ON pre_foaling_checklist(horse_id);

-- +goose Down
DROP TABLE IF EXISTS pre_foaling_checklist;
