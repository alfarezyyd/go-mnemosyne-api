CREATE TABLE notes
(
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id     BIGINT UNSIGNED NOT NULL,
    title       VARCHAR(255)    NOT NULL,
    content     TEXT,
    category_id BIGINT UNSIGNED                DEFAULT NULL,
    priority    ENUM ('Low', 'Medium', 'High') DEFAULT 'Low',
    due_date    DATE                           DEFAULT NULL,
    is_pinned   BOOLEAN                        DEFAULT FALSE,
    is_archived BOOLEAN                        DEFAULT FALSE,
    created_at  TIMESTAMP                      DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP                      DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (category_id) REFERENCES categories (id)
);