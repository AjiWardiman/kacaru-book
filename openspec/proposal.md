# Proposal Fitur Baru: AI Rekomendasi, Integrasi Cuaca OpenWeatherMap & Sistem Laporan

## 1. Latar Belakang & Tujuan
Sesuai arahan Trainer, aplikasi Tourism Booking perlu ditingkatkan dengan fitur cerdas (AI) dan fitur analitik bisnis yang memberikan nilai tambah nyata bagi wisatawan dan pengelola:
1. **AI Rekomendasi Destinasi:** Membantu wisatawan mendapatkan saran destinasi wisata sesuai preferensi (budget, kategori kesukaan, lokasi).
2. **Prediksi Cuaca (OpenWeatherMap API):** Menampilkan prakiraan cuaca di halaman detail destinasi & pilih slot tanggal untuk membantu wisatawan merencanakan kunjungan yang aman dari hujan/badai.
3. **Fitur Laporan & Analytics (Admin):** Menyajikan grafik statistik pendapatan & grafik pengunjung, serta kemampuan ekspor laporan transaksi ke PDF dan Excel/CSV.

---

## 2. Rencana Arsitektur & Perubahan Komponen

### A. Fitur 1: AI Rekomendasi Destinasi Wisata
- **Pendekatan:** Content-Based Filtering / Hybrid AI (Integrasi Gemini API / Rule-based Scoring System di Go).
- **Backend (Go):**
  - Endpoint `POST /api/v1/recommendations`: Menerima input preferensi (budget maks, kategori favorit, lokasi/jarak).
  - Algoritma pencocokan data destinasi dari database berdasarkan kriteria preferensi pengguna.
- **Frontend (React):**
  - Widget / Form "Cari Rekomendasi AI" di beranda/halaman destinasi.

### B. Fitur 2: Integrasi Cuaca (OpenWeatherMap API)
- **Backend (Go):**
  - Service `weather_service.go` yang melakukan fetch data ke OpenWeatherMap Forecast API (`api.openweathermap.org/data/2.5/forecast`).
  - Caching data cuaca (1-3 jam) untuk menghemat kuota panggil API.
  - Endpoint `GET /api/v1/destinations/:id/weather`: Mengembalikan prakiraan cuaca 5 hari ke depan di lokasi destinasi.
- **Frontend (React):**
  - Komponen **Weather Badge / Widget** di halaman detail destinasi dan pada saat memilih slot tanggal pemesanan (misal: "☀️ Cerah 28°C" atau "🌧️ Hujan Deras").

### C. Fitur 3: Sistem Laporan & Analytics (Admin Dashboard)
- **Backend (Go):**
  - Endpoint `GET /api/v1/admin/reports/summary`: Ringkasan Total Pendapatan, Total Booking, Destinasi Terpopuler, Okupansi Slot.
  - Endpoint `GET /api/v1/admin/reports/export/pdf`: Generate file PDF laporan keuangan & pengunjung.
  - Endpoint `GET /api/v1/admin/reports/export/excel`: Export file Excel (.xlsx/csv) data transaksi booking.
- **Frontend (React):**
  - Halaman Dashboard Admin dilengkapi grafik (misal menggunakan Recharts / Chart.js).
  - Tombol **Export PDF** dan **Export Excel**.

---

## 3. Hasil yang Diharapkan
- Nilai akademis Tugas Akhir meningkat pesat karena memiliki komponen **AI**, **Integrasi Third-Party Realtime API (OpenWeatherMap)**, dan **Business Intelligence / Reporting**.
- Pengalaman pengguna (UX) wisatawan menjadi jauh lebih interaktif dan informatif.
