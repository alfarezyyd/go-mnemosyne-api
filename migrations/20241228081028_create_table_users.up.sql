CREATE TABLE users
(
    id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,                                  -- ID unik untuk setiap pengguna
    name                VARCHAR(255)        NOT NULL,                                                -- Nama pengguna
    email               VARCHAR(255) UNIQUE NOT NULL,                                                -- Email pengguna, harus unik
    password            VARCHAR(255)        NOT NULL,                                                -- Hash kata sandi pengguna
    role                ENUM ('Admin', 'User') DEFAULT 'User',
    email_verified_at   TIMESTAMP              DEFAULT NULL,
    phone_number        VARCHAR(15)            DEFAULT NULL,                                         -- Nomor telepon pengguna (opsional)
    profile_picture     VARCHAR(255)           DEFAULT NULL,                                         -- URL foto profil pengguna (opsional)
    is_active           BOOLEAN                DEFAULT TRUE,                                         -- Status keaktifan pengguna
    language_preference VARCHAR(5)             DEFAULT 'id',
    created_at          TIMESTAMP              DEFAULT CURRENT_TIMESTAMP,                            -- Waktu pembuatan akun
    updated_at          TIMESTAMP              DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Waktu pembaruan terakhir
);
