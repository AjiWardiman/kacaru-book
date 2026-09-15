-- ============================================================================
-- Tourism Booking System - Skema Database Awal
-- PostgreSQL
--
-- File ini berisi perintah untuk membuat SEMUA tabel yang dibutuhkan
-- aplikasi, sesuai dengan ERD yang sudah dirancang. Jalankan file ini
-- SEKALI SAJA saat pertama kali setup database (lihat README.md untuk caranya).
-- ============================================================================

-- Tipe data ENUM adalah cara PostgreSQL membatasi sebuah kolom hanya
-- boleh diisi dengan pilihan-pilihan tertentu saja (mirip "dropdown"),
-- supaya tidak ada typo atau nilai aneh yang bisa masuk ke database.

CREATE TYPE user_role AS ENUM ('admin', 'user');
CREATE TYPE booking_status AS ENUM ('pending', 'paid', 'expired', 'cancelled');
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'expired', 'failed');
CREATE TYPE otp_purpose AS ENUM ('phone_verification', 'login', 'reset_password');

-- Tabel akun pengguna (baik user biasa maupun admin).
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,       -- UNIQUE: tidak boleh ada 2 user dengan email sama
    phone VARCHAR(20) UNIQUE,                  -- boleh kosong (NULL), tapi kalau diisi harus unik
    phone_verified_at TIMESTAMPTZ,             -- diisi otomatis saat verifikasi HP berhasil
    password VARCHAR(255) NOT NULL,            -- ini HASH bcrypt, bukan password asli
    role user_role NOT NULL DEFAULT 'user',    -- default semua user baru adalah role "user" biasa
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel kategori destinasi wisata (misal "Pantai", "Gunung", "Budaya").
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Tabel destinasi wisata itu sendiri.
CREATE TABLE destinations (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NOT NULL REFERENCES users(id), -- admin yang membuat data ini
    name VARCHAR(200) NOT NULL,
    location VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),      -- kategori versi teks bebas (versi simpel, bukan relasi tabel categories)
    price NUMERIC(12,2) NOT NULL DEFAULT 0,
    image_url VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT true, -- kalau false, disembunyikan dari daftar publik
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel penghubung many-to-many antara destinations dan categories,
-- untuk versi yang lebih terstruktur (satu destinasi bisa punya banyak
-- kategori, satu kategori bisa dipakai banyak destinasi).
CREATE TABLE destination_categories (
    destination_id BIGINT NOT NULL REFERENCES destinations(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (destination_id, category_id)
);

-- Tabel slot: jadwal kunjungan + kuota untuk setiap destinasi di
-- tanggal tertentu. Inilah tabel kunci untuk fitur "anti-bentrok booking".
CREATE TABLE slots (
    id BIGSERIAL PRIMARY KEY,
    destination_id BIGINT NOT NULL REFERENCES destinations(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_quota INT NOT NULL,           -- kuota maksimal pengunjung di tanggal ini
    used_quota INT NOT NULL DEFAULT 0,  -- kuota yang sudah terpakai booking
    is_available BOOLEAN NOT NULL DEFAULT true, -- admin bisa menutup slot manual (misal cuaca buruk)
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    -- UNIQUE ini mencegah admin membuat 2 slot untuk destinasi + tanggal yang sama.
    CONSTRAINT uq_destination_date UNIQUE (destination_id, date),
    -- CHECK ini adalah "jaring pengaman" di level database: kuota terpakai
    -- tidak boleh negatif dan tidak boleh melebihi kuota total, APAPUN
    -- yang terjadi di kode aplikasi (jaga-jaga kalau ada bug di masa depan).
    CONSTRAINT chk_quota CHECK (used_quota >= 0 AND used_quota <= total_quota)
);

-- Tabel booking/pemesanan yang dibuat user untuk slot tertentu.
CREATE TABLE bookings (
    id BIGSERIAL PRIMARY KEY,
    booking_code VARCHAR(30) NOT NULL UNIQUE, -- kode booking yang ditampilkan ke user, misal "BK-20260920-001"
    user_id BIGINT NOT NULL REFERENCES users(id),
    slot_id BIGINT NOT NULL REFERENCES slots(id),
    total_people INT NOT NULL,       -- jumlah orang yang ikut dalam booking ini
    total_price NUMERIC(12,2) NOT NULL,
    status booking_status NOT NULL DEFAULT 'pending', -- pending / paid / expired / cancelled
    booked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expired_at TIMESTAMPTZ NOT NULL, -- batas waktu booking harus dibayar sebelum otomatis dibatalkan
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel detail data setiap pengunjung dalam satu booking (satu booking
-- bisa untuk banyak orang/rombongan).
CREATE TABLE booking_details (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    visitor_name VARCHAR(150) NOT NULL,
    visitor_phone VARCHAR(20),
    visitor_email VARCHAR(150),
    age INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel pembayaran, akan terhubung dengan Xendit (payment gateway).
CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    external_id VARCHAR(100) NOT NULL UNIQUE, -- ID yang kita kirim ke Xendit untuk mengenali transaksi ini
    invoice_id VARCHAR(150),                   -- ID invoice yang dikembalikan oleh Xendit
    amount NUMERIC(12,2) NOT NULL,
    payment_method VARCHAR(50),
    status payment_status NOT NULL DEFAULT 'pending', -- pending / paid / expired / failed
    payment_url VARCHAR(500),   -- link halaman pembayaran dari Xendit yang dibuka user
    paid_at TIMESTAMPTZ,
    expired_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel log setiap event/callback (webhook) yang dikirim Xendit ke
-- sistem kita, misal notifikasi "pembayaran berhasil".
CREATE TABLE payment_events (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- misal "invoice.paid", "invoice.expired"
    payload TEXT,                     -- data mentah (JSON) yang dikirim Xendit, disimpan untuk keperluan audit
    received_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel audit trail: mencatat SIAPA melakukan aksi APA terhadap data
-- MANA dan KAPAN. Ini requirement wajib dari trainer untuk melacak
-- histori perubahan data penting di sistem.
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT REFERENCES users(id),
    action VARCHAR(20) NOT NULL,
    table_name VARCHAR(100) NOT NULL,
    record_id BIGINT,
    old_value JSONB,
    new_value JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE settings (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL UNIQUE,   -- nama pengaturan, misal "booking_expiry_minutes"
    value TEXT,                          -- nilai pengaturan (disimpan sebagai teks, diubah tipe di kode Go)
    description TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel notifikasi (riwayat notifikasi yang dikirim ke user, misal
-- "booking kamu berhasil dibayar").
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    message TEXT,
    is_read BOOLEAN NOT NULL DEFAULT false,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel kode OTP (One Time Password), dipakai untuk verifikasi HP
-- (dikirim lewat email via Mailtrap sebagai simulasi, sesuai kesepakatan).
CREATE TABLE otps (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    code VARCHAR(6) NOT NULL,
    purpose otp_purpose NOT NULL, -- phone_verification / login / reset_password
    is_used BOOLEAN NOT NULL DEFAULT false, -- supaya satu kode OTP tidak bisa dipakai berulang kali
    expired_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ============================================================================
-- INDEX: mempercepat pencarian data yang sering dipakai. Tanpa index,
-- database harus memeriksa SEMUA baris satu per satu (lambat kalau data
-- sudah banyak); dengan index, database bisa "loncat" langsung ke baris
-- yang relevan.
-- ============================================================================
CREATE INDEX idx_slots_destination_date ON slots(destination_id, date);
CREATE INDEX idx_bookings_user ON bookings(user_id);
CREATE INDEX idx_bookings_slot ON bookings(slot_id);
CREATE INDEX idx_bookings_status_expired ON bookings(status, expired_at); -- dipakai untuk mencari booking yang perlu di-expired otomatis
CREATE INDEX idx_payments_booking ON payments(booking_id);
CREATE INDEX idx_audit_logs_table_record ON audit_logs(table_name, record_id);
CREATE INDEX idx_otps_user_purpose ON otps(user_id, purpose);

-- ============================================================================
-- DATA AWAL: nilai default untuk tabel settings, supaya aplikasi
-- langsung punya pengaturan dasar begitu database selesai di-migrasi,
-- tanpa perlu diisi manual dulu.
-- ============================================================================
INSERT INTO settings (key, value, description) VALUES
    ('booking_expiry_minutes', '30', 'Batas waktu pembayaran booking sebelum otomatis expired (menit)'),
    ('otp_expiry_minutes', '5', 'Batas waktu berlaku kode OTP (menit)'),
    ('admin_contact_email', 'admin@tourism-booking.local', 'Email kontak admin untuk notifikasi sistem');
