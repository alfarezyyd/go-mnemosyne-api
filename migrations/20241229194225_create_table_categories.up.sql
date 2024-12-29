CREATE TABLE categories
(
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,                      -- ID unik untuk setiap kategori
    user_id     BIGINT UNSIGNED NOT NULL,                                        -- Relasi ke tabel `users`
    name        VARCHAR(100)    NOT NULL,                                        -- Nama kategori
    description TEXT      DEFAULT NULL,                                          -- Deskripsi kategori (opsional)
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                             -- Waktu kategori dibuat
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Waktu kategori diperbarui
    FOREIGN KEY (user_id) REFERENCES users (id)                                  -- Relasi ke tabel `users`
);
