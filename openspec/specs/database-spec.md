# Spesifikasi Skema Database (Database Schema Specification)

## 1. Entitas Utama

| Nama Tabel | Fungsi / Deskripsi |
| :--- | :--- |
| `users` | Akun pengguna (role: `admin` / `user`), berisikan profil, password hash, dan status verifikasi HP. |
| `categories` | Kategori tempat wisata (misal: Pantai, Gunung, Budaya). |
| `destinations` | Data destinasi wisata (nama, lokasi, deskripsi, harga tiket dasar, URL gambar, status aktif). |
| `destination_categories` | Tabel penghubung (many-to-many) antara destinasi dan kategori. |
| `slots` | Jadwal kunjungan per tanggal & kuota pengunjung (`total_quota`, `used_quota`, `is_available`). |
| `bookings` | Rekam transaksi pemesanan (`booking_code`, `user_id`, `slot_id`, `total_people`, `total_price`, `status`, `expired_at`). |
| `booking_details` | Detail nama, umur, no HP, email dari masing-masing pengunjung dalam 1 booking. |
| `payments` | Rekam pembayaran Xendit (`external_id`, `invoice_id`, `payment_url`, `status`, `paid_at`). |
| `payment_events` | Log mentah payload event callback webhook yang diterima dari Xendit. |
| `audit_logs` | Catatan histori aktivitas perubahan data oleh pengguna/admin (`user_id`, `action`, `table_name`, `old_value`, `new_value`). |
| `notifications` | Riwayat pesanan/notifikasi sistem untuk pengguna. |
| `otps` | Kode OTP verifikasi nomor HP / reset password. |
| `settings` | Konfigurasi bawaan sistem (seperti durasi waktu kadaluwarsa booking). |

---

## 2. Status Enum

- **`user_role`**: `'admin'`, `'user'`
- **`booking_status`**: `'pending'`, `'paid'`, `'expired'`, `'cancelled'`
- **`payment_status`**: `'pending'`, `'paid'`, `'expired'`, `'failed'`
- **`otp_purpose`**: `'phone_verification'`, `'login'`, `'reset_password'`

---

## 3. Constraint & Proteksi Integritas Data

- `chk_quota`: Constraint database yang memastikan `used_quota >= 0 AND used_quota <= total_quota`.
- `uq_destination_date`: Unique key pada `(destination_id, date)` mencegah duplikasi slot di tanggal yang sama.
- Foreign Key Cascade pada `booking_details`, `payments`, `payment_events` untuk pembersihan data terkait.
