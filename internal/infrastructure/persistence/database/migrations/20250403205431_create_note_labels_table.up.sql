CREATE TABLE note_labels (
    note_id VARCHAR(255) NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    label_id VARCHAR(255) NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    PRIMARY KEY (note_id, label_id)
);