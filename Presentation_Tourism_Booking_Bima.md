# Presentasi Tugas Akhir: Tourism Booking System Bima

---

## 📌 SLIDE 1: JUDUL APLIKASI
- **Judul Utama**: TOURISM BOOKING SYSTEM KABUPATEN BIMA
- **Subjudul**: Tugas Akhir Bootcamp Web Application & Backend API Development
- **Nama Aplikasi**: Kacaru Tourism Booking (Bima Explorer)
- **Tahun**: 2026

---

## 👥 SLIDE 2: TIM PENGEMBANG & PERAN SISTEM
- **Project Manager & Backend Developer**:
  - Merancang arsitektur Go (Golang) REST API, Clean Architecture, Database PostgreSQL (GORM), Keamanan JWT, Integrasi AI OpenRouter, Weather API, dan Webhook Xendit Payment Gateway.
- **Frontend Developer**:
  - Mendesain User Interface (UI) modern berbasis React JS, TypeScript, TailwindCSS, Glassmorphism, Zustand State Management, dan Responsif Mobile/Desktop.
- **QA / Tester & Documentation**:
  - Menyiapkan skenario pengujian API (Insomnia/Postman), verifikasi callback Webhook Xendit, pengujian verifikasi OTP, serta validasi kuota slot harian.
- **Tools Kolaborasi**:
  - GitHub untuk version control, Postman/Insomnia untuk dokumentasi REST API, dan VS Code untuk lingkungan pengembangan.

---

## 🎯 SLIDE 3: LATAR BELAKANG
- **Permasalahan**:
  - Pemesanan tiket wisata di Kabupaten Bima sebelumnya masih bersifat manual, antrean fisik di lokasi, dan informasi cuaca yang tidak terintegrasi.
- **Kondisi Saat Ini**:
  - Data kunjungan tidak realtime, risiko kelebihan kapasitas (over-capacity) di lokasi wisata, dan ketidakpastian wisatawan mengenai cuaca saat berkunjung.
- **Dampak**:
  - Waktu tunggu antrean lama, risiko pembatalan kunjungan secara mendadak akibat cuaca buruk, dan pencatatan laporan keuangan wisata yang kurang akurat.
- **Alasan Pengembangan**:
  - Mengembangkan sistem pemesanan tiket wisata digital terpadu berbasis AI rekomendasi cerdas dan Payment Gateway otomatis untuk memajukan pariwisata Kabupaten Bima.

---

## 💡 SLIDE 4: TUJUAN APLIKASI
- **Tujuan Utama**:
  - Menyediakan solusi platform pariwisata digital terintegrasi yang mengotomasi pemesanan tiket, memvalidasi pembayaran secara realtime, dan memberikan rekomendasi destinasi cerdas berbasis cuaca.
- **Target Pengguna & Manfaat**:
  - **Wisatawan**: Kemudahan memilih destinasi, rekomendasi AI cerdas, verifikasi OTP HP/Email, pembayaran instan QRIS/VA, serta E-Ticket digital QR Code.
  - **Pengelola / Admin**: Visibilitas data kunjungan realtime, kontrol otomatis kuota slot harian, notifikasi SSE realtime, serta laporan keuangan/transaksi yang siap diekspor ke CSV & PDF.

---

## 🛠️ SLIDE 5: TOOLS & TECH STACK
- **Backend Framework & Database**:
  - Go (Golang 1.22+), Gin Web Framework, GORM ORM, PostgreSQL Database, JWT Authentication, Bcrypt Password Hashing.
- **Frontend & Mobile Responsive**:
  - React JS (TypeScript), Vite, TailwindCSS, Zustand State Management, Lucide React Icons, Sonner Toast Notifications.
- **Integrasi API Eksternal & DevOps**:
  - **OpenRouter AI API** (Rekomendasi Cerdas Claude/Gemini), **OpenWeatherMap API** (Prakiraan Cuaca Realtime), **Xendit Payment Gateway** (QRIS & Virtual Account), **Mailtrap SMTP** (OTP & Notifikasi Email).

---

## 🔄 SLIDE 6: FLOW BISNIS
- **1. User Register & Authentication**:
  - User mendaftar akun -> Verifikasi kode OTP nomor HP/Email -> Sistem menerbitkan Token JWT aman.
- **2. Recommendations & Ticket Booking**:
  - User melihat rekomendasi AI cerdas berbasis ramalan cuaca -> Pilih destinasi & slot tanggal -> Masukkan data pengunjung rombongan -> Sistem mengunci kuota slot (`SELECT FOR UPDATE`).
- **3. Payment & E-Ticket Issuance**:
  - User melakukan pembayaran via Xendit (QRIS/VA) -> Webhook callback mengonfirmasi pembayaran lunas (`PAID`) -> Sistem menerbitkan E-Ticket QR Code digital & Notifikasi SSE.
- **4. Admin Management & Reporting**:
  - Admin mengelola tempat wisata & slot kuota harian -> Monitoring transaksi platform -> Ekspor laporan pendapatan ke CSV / PDF.

---

## 📊 SLIDE 7: ENTITY RELATIONSHIP DIAGRAM (ERD)
- **Struktur Data Utama**:
  - `users` (Role Admin & User), `destinations` (Katalog Tempat Wisata), `slots` (Kuota Harian per Tanggal), `bookings` & `booking_details` (Transaksi & Rombongan), `payments` & `payment_events` (Xendit Invoice & Log Webhook), `otps`, `notifications`, `settings`, `audit_logs`.
- **Relasi Antar Entitas**:
  - `User 1:N Bookings` (Wisatawan memesan tiket).
  - `Destination 1:N Slots` (Destinasi memiliki slot kuota harian).
  - `Slot 1:N Bookings` (Pemesanan memotong `used_quota` slot).
  - `Booking 1:1 Payment` (Pemesanan memiliki 1 invoice pembayaran Xendit).
  - `Booking 1:N BookingDetails` (1 booking menyimpan data rombongan).

---

## ⭐ SLIDE 8: LIST FITUR
- **Fitur Wisatawan (User)**:
  - Registrasi & Verifikasi OTP HP/Email
  - Rekomendasi Destinasi berbasis AI & Cuaca Realtime
  - Pemesanan Tiket Individu / Rombongan
  - Pembayaran Instan Xendit (QRIS & Virtual Account)
  - E-Ticket Digital dengan Kode QR
  - Streaming Notifikasi Realtime (SSE)
- **Fitur Pengelola (Admin)**:
  - Management Control Panel Destinasi Wisata
  - Pengaturan Slot Kuota Harian
  - Monitoring Transaksi & Manual Payment Confirmation
  - Manajemen Pengguna Platform
  - Ekspor Laporan Transaksi (Spreadsheet CSV & Printable PDF)

---

## 🖼️ SLIDE 9: SCREENSHOT UI (SHOWCASE MOCKUP)
- **Katalog Wisata & Rekomendasi AI**: Tampilan kartu destinasi modern dengan skor kecocokan AI dan kondisi cuaca realtime.
- **Form Booking & Pilihan Tanggal**: Interface pemilihan slot tanggal kunjungan dan jumlah pengunjung rombongan.
- **Modal Checkout Xendit Gateway**: Tampilan invoice digital lengkap dengan QRIS, Virtual Account Bank, dan status realtime.
- **Dashboard & Laporan Admin**: Panel kontrol eksekutif menampilkan ringkasan pendapatan, grafik transaksi, dan tombol cetak laporan.

---

## 🎬 SLIDE 10: DEMO, PENGEMBANGAN & KESIMPULAN
- **Demo Aplikasi**:
  - Demonstrasi alur end-to-end: Registrasi akun -> Rekomendasi AI -> Pemesanan Tiket -> Simulasi Bayar Lunas Xendit -> Terbit E-Ticket QR Code -> Laporan di Dashboard Admin.
- **Rencana Pengembangan Selanjutnya**:
  - Pengembangan aplikasi Mobile (React Native / Flutter), Integrasi mesin pencetak tiket thermal di lokasi wisata, serta Notifikasi Direct WhatsApp Gateway.
- **Kesimpulan & Terima Kasih**:
  - Sistem Tourism Booking Bima berhasil mentransformasi alur pemesanan tiket manual menjadi platform digital otomatis, akurat, dan transparan. Terima kasih kepada Tim Penilai dan Instruktur Juara Coding!
