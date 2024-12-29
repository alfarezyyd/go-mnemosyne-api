CREATE TABLE one_time_password_tokens
(
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY NOT NULL,
    user_id      BIGINT UNSIGNED                            NOT NULL,
    hashed_Token VARCHAR(255)                               NOT NULL,
    expires_at   DATETIME                                   NOT NULL,
    CONSTRAINT fk_one_time_password_tokens_users FOREIGN KEY (user_id) REFERENCES users (id)
)