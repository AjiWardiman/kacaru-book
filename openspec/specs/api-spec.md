# Spesifikasi API (API Specification Overview)

## 1. Endpoint Autentikasi & OTP (`/auth` & `/otp`)
- `POST /auth/register` - Pendaftaran user baru (memicu notifikasi welcome)
- `POST /auth/login` - Login & dapatkan token JWT
- `POST /auth/send-otp` - [Publik] Kirim kode OTP 6-digit pra-registrasi ke Mailtrap
- `POST /otp/request` - [Auth] Minta kode OTP verifikasi nomor HP baru via Mailtrap
- `POST /otp/verify` - [Auth] Verifikasi kode OTP 6-digit (update `phone_verified_at`)

## 2. Endpoint Destinasi (`/destinations`)
- `GET /destinations` - Daftar destinasi publik (dapat difilter kategori/lokasi)
- `GET /destinations/:id` - Detail destinasi wisata
- `POST /admin/destinations` - [Admin] Tambah destinasi wisata baru
- `PUT /admin/destinations/:id` - [Admin] Update data destinasi
- `DELETE /admin/destinations/:id` - [Admin] Hapus/Nonaktifkan destinasi

## 3. Endpoint Slot & Kuota (`/destinations/:id/slots`)
- `GET /destinations/:id/slots` - Cek ketersediaan slot per tanggal
- `POST /admin/slots` - [Admin] Buat / atur kuota slot kunjungan baru
- `PUT /admin/slots/:id` - [Admin] Ubah kuota atau ketersediaan slot

## 4. Endpoint Pemesanan (`/bookings`)
- `POST /bookings` - Buat pemesanan tiket baru (kunci kuota slot & status pending)
- `GET /bookings/my-bookings` - Riwayat pemesanan milik user aktif
- `POST /bookings/:id/pay` - Buat invoice & URL pembayaran Xendit

## 5. Endpoint Pembayaran & Webhook (`/webhooks`)
- `POST /webhooks/xendit` - Callback webhook dari Xendit (proses status `paid`/`expired`)

## 6. Endpoint Fitur Tambahan (AI, Weather, Notifikasi, Export)
- `POST /recommendations` - Dapatkan rekomendasi destinasi berbasis AI
- `GET /weather` - Informasi prakiraan cuaca lokasi wisata (OpenWeatherMap + Caching)
- `GET /notifications` - Daftar notifikasi user
- `GET /notifications/stream` - Server-Sent Events (SSE) stream notifikasi realtime
- `GET /admin/reports/summary` - [Admin] Laporan ringkasan statistik
- `GET /admin/reports/export/excel` - [Admin] Ekspor laporan format Excel (`.xlsx`)
- `GET /admin/reports/export/pdf` - [Admin] Ekspor laporan format PDF
