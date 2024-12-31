CREATE TABLE shared_notes
(
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    note_id    BIGINT UNSIGNED NOT NULL,
    user_id    BIGINT UNSIGNED NOT NULL,
    permission ENUM ('Read', 'Edit')      DEFAULT 'Read',
    shared_at  TIMESTAMP                  DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP                  DEFAULT NULL,
    status     ENUM ('Active', 'Revoked') DEFAULT 'Active',
    FOREIGN KEY (note_id) REFERENCES notes (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
