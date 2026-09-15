# Tourism Booking System — Backend (Golang)

Tugas Akhir Bootcamp Golang Backend — Sistem Booking Destinasi Wisata dengan
Reservasi Anti-Bentrok dan Pencarian Paralel.

## Status Progress

Sudah selesai dan **teruji berjalan** (build + test end-to-end dengan
PostgreSQL sungguhan, termasuk test konkruensi nyata dan simulasi Xendit):

- [x] Setup project (clean architecture: handler → service → repository)
- [x] Migrasi database PostgreSQL (13 tabel sesuai ERD)
- [x] Auth: register & login (email atau nomor HP), JWT, bcrypt
- [x] CRUD Destinations (admin) + browsing publik
- [x] CRUD Slots per destinasi (admin), dengan constraint anti-duplikat tanggal
- [x] Middleware role-based (admin vs user)
- [x] **Booking dengan anti-bentrok** (row locking `SELECT FOR UPDATE`) —
      sudah diuji dengan 10 request booking BERSAMAAN ke slot berkuota 5,
      hasilnya tepat 5 berhasil dan 5 ditolak, data database 100% konsisten
- [x] **Integrasi Xendit** (create invoice + terima webhook) — sudah diuji
      pakai server tiruan yang berperan sebagai Xendit, alur lengkap
      booking → bayar → webhook → status berubah, semuanya teruji
- [x] **Auto-expire booking** — worker latar belakang yang otomatis
      membatalkan booking yang tidak dibayar tepat waktu dan
      mengembalikan kuota; sudah diuji berjalan otomatis
- [x] **Notifikasi email (Mailtrap)** — email otomatis dikirim ke user
      saat pembayaran berhasil atau booking dibatalkan/expired; sudah
      diuji dengan server SMTP tiruan, format email dan pemicunya
      tervalidasi bekerja
- [x] **Audit trail** — mencatat siapa melakukan aksi apa (create/update/
      delete) terhadap data apa, lengkap dengan snapshot data sebelum
      dan sesudah perubahan (format JSON); sudah diuji di semua titik
      penting (destination, slot, booking, perubahan status payment)
- [x] Semua pesan response API dalam Bahasa Indonesia

- [x] **Settings** — endpoint admin untuk melihat dan mengubah
      pengaturan aplikasi (misal `otp_expiry_minutes`) tanpa perlu
      ubah kode/deploy ulang, lengkap dengan audit trail tiap ada perubahan
- [x] **OTP** — verifikasi nomor HP pakai kode 6 digit yang dikirim
      lewat email (simulasi SMS sesuai kesepakatan tugas akhir); kode
      sekali pakai, punya batas waktu (dibaca dari tabel settings), dan
      kode lama otomatis tidak berlaku lagi begitu user minta kode baru

Belum dikerjakan (tahap selanjutnya):

- [ ] Swagger (dokumentasi API)
- [ ] Playwright (API testing)

## Menjalankan di Komputer Sendiri

### 1. Install dependency

```bash
go mod tidy
```

> **Catatan**: file `go.mod` di sini sudah versi bersih (tanpa `replace`
> directive). Perintah di atas akan otomatis mengambil semua dependency
> langsung dari sumber resminya — tidak perlu langkah tambahan apa pun,
> selama komputer kamu punya akses internet normal ke proxy.golang.org.

### 2. Setup database PostgreSQL

Buat database baru:

```bash
createdb tourism_booking
```

Jalankan migrasi:

```bash
psql -d tourism_booking -f migrations/001_init.sql
```

### 3. Setup environment variables

Copy `.env.example` menjadi `.env` dan sesuaikan (terutama `DB_PASSWORD`,
`DB_PORT` kalau PostgreSQL tidak pakai port default 5432, dan kredensial
Xendit — lihat bagian "Setup Xendit" di bawah):

```bash
cp .env.example .env
```

File `.env` ini **otomatis terbaca** oleh aplikasi (pakai library
`godotenv`) — tidak perlu di-export manual sebelum menjalankan `go run`.

### 4. Jalankan server

```bash
go run ./cmd/api
```

Server berjalan di `http://localhost:8080`. Worker auto-expire booking
otomatis mulai berjalan di latar belakang, memeriksa booking yang
kadaluarsa setiap 1 menit.

### 5. Test cepat

Gunakan koleksi Insomnia yang sudah disiapkan (lihat bagian di bawah),
atau langsung lewat curl:

```bash
curl http://localhost:8080/health
```

## Setup Xendit (untuk fitur pembayaran)

1. Daftar akun gratis di [dashboard.xendit.co](https://dashboard.xendit.co),
   pastikan berada di **mode Test/Sandbox**.
2. Ambil **Secret API Key** (test mode) dari menu Settings → API Keys.
3. Buat **Callback Token** di menu Settings → Webhooks — ini token
   rahasia yang dipakai untuk memvalidasi bahwa webhook yang masuk
   benar-benar dari Xendit.
4. Isi `.env`:
   ```
   XENDIT_SECRET_KEY=xnd_development_xxxxx
   XENDIT_CALLBACK_TOKEN=isi_dengan_callback_token_dari_dashboard
   XENDIT_BASE_URL=
   ```
   (`XENDIT_BASE_URL` dikosongkan supaya otomatis pakai alamat resmi Xendit.)
5. Untuk testing webhook di komputer sendiri (localhost), Xendit tidak
   bisa "mengetuk" balik ke `localhost` secara langsung dari internet.
   Gunakan tool seperti **ngrok** (`ngrok http 8080`) untuk membuat
   `localhost` bisa diakses sementara dari internet, lalu daftarkan
   URL ngrok itu (contoh: `https://xxxx.ngrok.io/webhooks/xendit`) di
   dashboard Xendit sebagai URL webhook.

## Setup Mailtrap (untuk fitur notifikasi email)

1. Daftar akun gratis di [mailtrap.io](https://mailtrap.io), buat inbox baru.
2. Buka inbox itu, pilih tab **SMTP Settings**, catat kredensial yang muncul.
3. Isi `.env`:
   ```
   MAILTRAP_HOST=sandbox.smtp.mailtrap.io
   MAILTRAP_PORT=2525
   MAILTRAP_USERNAME=isi_dari_dashboard_mailtrap
   MAILTRAP_PASSWORD=isi_dari_dashboard_mailtrap
   MAIL_FROM_ADDRESS=no-reply@tourism-booking.local
   ```
4. Semua email yang "terkirim" akan muncul di dashboard Mailtrap
   (bukan ke inbox asli), jadi aman untuk testing.

## Membuat User Admin

Admin dibuat manual lewat database (belum ada endpoint publik untuk
membuat admin, sengaja untuk keamanan):

```sql
UPDATE users SET role = 'admin' WHERE email = 'email_kamu@example.com';
```

Login ulang setelah ini untuk mendapatkan token JWT baru dengan role `admin`.

## Struktur Folder

```
cmd/api/          entry point aplikasi
internal/config/  loader environment variables
internal/model/   struct yang merepresentasikan tabel database
internal/dto/     request/response struct (DTO pattern)
internal/repository/  lapisan akses database (SQL query)
internal/service/     business logic
internal/handler/     HTTP handler (Gin)
internal/middleware/  JWT auth & role-based authorization
internal/router/      routing dan wiring endpoint
internal/xendit/      klien HTTP untuk komunikasi dengan API Xendit
internal/mailer/      klien SMTP untuk mengirim email lewat Mailtrap
internal/scheduler/   worker latar belakang (auto-expire booking)
pkg/database/     koneksi PostgreSQL
migrations/       file SQL migrasi schema
insomnia/         koleksi request Insomnia siap import
```

## Endpoint yang Sudah Ada

| Method | Path | Akses | Keterangan |
|---|---|---|---|
| GET | /health | Publik | Health check |
| POST | /auth/register | Publik | Registrasi user baru |
| POST | /auth/login | Publik | Login (email atau nomor HP) |
| GET | /destinations | Publik | List destinasi aktif |
| GET | /destinations/:id | Publik | Detail destinasi |
| GET | /destinations/:id/slots | Publik | List slot/jadwal destinasi |
| POST | /admin/destinations | Admin | Buat destinasi baru |
| PUT | /admin/destinations/:id | Admin | Update destinasi |
| DELETE | /admin/destinations/:id | Admin | Hapus destinasi |
| POST | /admin/destinations/:id/slots | Admin | Buat slot/jadwal baru |
| POST | /bookings | Login (user/admin) | Buat booking baru (anti-bentrok) |
| POST | /bookings/:id/pay | Login (user/admin) | Buat invoice pembayaran (Xendit) |
| POST | /otp/request | Login (user/admin) | Minta kode OTP verifikasi HP (dikirim lewat email) |
| POST | /otp/verify | Login (user/admin) | Verifikasi kode OTP yang diterima |
| POST | /webhooks/xendit | Xendit (bukan user) | Terima notifikasi status pembayaran |
| GET | /admin/settings | Admin | Lihat semua pengaturan aplikasi |
| GET | /admin/settings/:key | Admin | Lihat detail satu pengaturan |
| PUT | /admin/settings/:key | Admin | Ubah nilai satu pengaturan |

## Cara Import Koleksi Insomnia

Folder `insomnia/tourism-booking-insomnia.json` berisi koleksi request
siap pakai. Cara import: buka Insomnia → **Import** → **From File** →
pilih file itu. Semua request admin sudah pakai variabel `{{ _.token }}`
— tinggal isi token hasil login di environment supaya tidak perlu
copy-paste manual ke tiap request.

## Menguji Fitur Anti-Bentrok Booking

Buat slot dengan kuota kecil (misal 5), lalu kirim BANYAK request
`POST /bookings` ke slot yang sama SECARA BERSAMAAN (pakai script
goroutine Go, atau tool seperti `hey`/`k6`). Hasil yang benar: jumlah
booking yang berhasil TIDAK PERNAH melebihi kuota, walau request
dikirim bersamaan.

## Menguji Auto-Expire Booking

Karena defaultnya worker memeriksa tiap 1 menit, untuk testing cepat:
1. Buat booking seperti biasa.
2. Paksa kolom `expired_at` booking itu ke masa lalu lewat database:
   ```sql
   UPDATE bookings SET expired_at = now() - interval '1 minute' WHERE id = <id_booking>;
   ```
3. Tunggu maksimal 1 menit, cek status booking — harus otomatis
   berubah jadi `expired` dan kuota di slot terkait kembali bertambah.
