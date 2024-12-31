CREATE TABLE settings
(
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id     BIGINT UNSIGNED DEFAULT NULL,
    `key`       VARCHAR(255)                                  NOT NULL,
    value       TEXT                                          NOT NULL,
    type        ENUM ('String', 'Boolean', 'Integer', 'JSON') NOT NULL,
    is_global   BOOLEAN         DEFAULT FALSE,
    `group`     VARCHAR(100)    DEFAULT NULL,
    description TEXT            DEFAULT NULL,
    created_at  TIMESTAMP       DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP       DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
