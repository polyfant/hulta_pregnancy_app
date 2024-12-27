-- +goose Up
CREATE TABLE pre_foaling_checklist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    horse_id INTEGER NOT NULL,
    description TEXT NOT NULL,
    is_completed BOOLEAN DEFAULT FALSE,
    due_date DATETIME NOT NULL,
    priority TEXT CHECK(priority IN ('HIGH', 'MEDIUM', 'LOW')) DEFAULT 'MEDIUM',
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (horse_id) REFERENCES horses(id) ON DELETE CASCADE
);

CREATE TRIGGER update_pre_foaling_checklist_timestamp 
AFTER UPDATE ON pre_foaling_checklist
BEGIN
    UPDATE pre_foaling_checklist SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;

-- +goose Down
DROP TRIGGER IF EXISTS update_pre_foaling_checklist_timestamp;
DROP TABLE IF EXISTS pre_foaling_checklist;
