# Spesifikasi Alur Sistem (System Flow Specification)

## 1. Peran Pengguna (User Roles)

### A. Wisatawan (User)
- Membuka katalog destinasi wisata dan kategori.
- **Menggunakan AI Rekomendasi Destinasi** untuk mendapatkan saran wisata sesuai preferensi budget & kategori.
- **Melihat Prakiraan Cuaca (OpenWeatherMap)** di destinasi sebelum memilih tanggal kunjungan.
- Melakukan pendaftaran (Register), login JWT, dan verifikasi OTP.
- Memilih tanggal kunjungan dan melihat ketersediaan kuota slot.
- Melakukan pemesanan tiket dengan mengisi data detail rombongan/pengunjung.
- Melakukan pembayaran via Xendit Payment Gateway (QRIS, VA, E-Wallet, Kartu Kredit).
- Menerima notifikasi dan mengakses E-Ticket (kode booking & QR Code).

### B. Pengelola (Admin)
- Kelola Destinasi Wisata (CRUD data & upload gambar).
- Kelola Kategori Wisata.
- Kelola Slot Kunjungan (Menentukan tanggal, total kuota, dan mengunci slot manual).
- **Dashboard Analytics & Laporan:** Melihat grafik pendapatan, grafik pengunjung, serta melakukan ekspor laporan ke format PDF & Excel.
- Monitoring Transaksi & Audit Trail Log.

---

## 2. Alur Transaksi Pemesanan (Booking Lifecycle) + Fitur AI & Cuaca

```
[User Browses / Asks AI Recommendation]
             │
             ▼
[View Destination & Weather Forecast (OpenWeatherMap)]
             │
             ▼
[User Selects Date & Checks Slot Quota] 
             │
             ▼
[User Fills Visitors Detail & Submits]
             │
             ▼
[Create Booking (status: pending) & Lock Slot Quota]
             │
             ▼
[Create Xendit Invoice Payment URL]
             │
             ▼
      ┌──────┴──────┐
      │             │
  (Payment Paid) (Expired/Failed)
      │             │
      ▼             ▼
[Update Booking]  [Update Booking]
[Status: 'paid']  [Status: 'expired']
      │             │
      ▼             ▼
[Notify User]     [Release Slot Quota]
[Generate QR Code Ticket]
```

---

## 3. Aturan Disiplin Sistem & Integrasi Eksternal

   - **Async Email Dispatch & Speed Optimization:** Pengiriman email dilakukan secara *non-blocking (asynchronous goroutine)* sehingga respon API backend dan pembukaan modal OTP di frontend berlangsung secara instan (< 50ms) tanpa tertahan oleh latensi handshake SMTP.
   - **Rate Limit Protection:** Pengiriman email di background menggunakan jeda aman `300ms` antar email untuk mencegah penolakan HTTP/SMTP `550 Too many emails per second` dari Mailtrap Free Tier.
   - **Keunikan Identitas:** Email dan Nomor HP harus unik per pengguna. Backend menolak pendaftaran jika email atau nomor HP sudah terdaftar.

2. **Alur Verifikasi OTP (One-Time Password):**
   - Kode OTP berupa **6-digit angka** acak yang dibuat oleh backend (`GenerateOTPCode`).
   - Berlaku selama durasi `otp_expiry_minutes` (default 5-10 menit, dikelola via tabel `settings`).
   - Kode OTP dikirimkan secara otomatis via Mailtrap ke email pengguna.
   - Verifikasi OTP wajib dilakukan sebelum akun ditandai terverifikasi penuh (`phone_verified_at`).

3. **AI Recommendation Engine:**
   - Algoritma pencocokan preferensi wisatawan dengan destinasi berbasis scoring kriteria (budget, lokasi, tags kategori) atau LLM Prompt.

4. **OpenWeatherMap Integration:**
   - Integrasi API cuaca realtime dengan sistem caching backend (1-3 jam) untuk optimasi kuota panggil API.

5. **Admin Reporting & Export System:**
   - Visualisasi chart statistik di frontend React.
   - Generator dokumen server-side di Go untuk ekspor file PDF & Excel (`.xlsx`).
