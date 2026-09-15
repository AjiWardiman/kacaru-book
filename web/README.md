# 💻 Tourism Booking System — Web Frontend (React + TypeScript + Vite)

Repositori ini berisi kode sumber untuk **Antarmuka Web (Frontend)** aplikasi Tourism Booking System. Dibuat menggunakan React 19, TypeScript, Vite, dan Tailwind CSS untuk menyajikan antarmuka modern, cepat, dan responsif.

---

## 🚀 Ringkasan Fitur Web (Bahasa Awam)

1. **🏠 Katalog Destinasi Wisata:**
   - Menampilkan daftar tempat wisata lengkap dengan foto, lokasi, harga tiket, dan rating.
   - Fitur pencarian instan berdasarkan kata kunci dan filter kategori (Bahari, Gunung, Budaya).

2. **🤖 Modal AI Rekomendasi Destinasi:**
   - Interaktif pemilih rekomendasi berbasis budget dan minat wisata.

3. **🌤️ Widget Forecast Cuaca Realtime:**
   - Menampilkan prediksi cuaca 5 hari dari OpenWeatherMap API dan catatan saran keselamatan.

4. **🎫 Pemesanan Tiket & E-Ticket Digital:**
   - Pemilihan slot tanggal kunjungan dan jumlah wisatawan.
   - Tampilan **E-Ticket QR Code** yang dapat di-download atau dicetak setelah transaksi berhasil.

5. **📊 Dashboard Admin & Export Laporan:**
   - Grafik analytics pendapatan dan okupansi slot tempat wisata.
   - Tombol unduh laporan transaksi dalam format **PDF** dan **Excel**.

---

## 🛠️ Langkah-Langkah Menjalankan Frontend

1. **Masuk ke folder web:**
   ```bash
   cd web
   ```

2. **Install semua dependensi:**
   ```bash
   npm install
   ```

3. **Jalankan server mode pengembangan (Development):**
   ```bash
   npm run dev
   ```
   - Aplikasi akan terbuka di browser pada alamat **`http://localhost:5173`**.

4. **Kompilasi Versi Produksi (Build):**
   ```bash
   npm run build
   ```

---

## 📁 Struktur Folder Frontend

```
web/
├── src/
│   ├── api/                  # Axios HTTP Client & Pemanggilan Backend API
│   ├── components/           # Komponen Tampilan (Navbar, Cards, Modal, Admin View)
│   │   ├── ai/               # Komponen AI Recommendation Modal
│   │   ├── auth/             # Komponen Login, Register & OTP Verification
│   │   ├── bookings/         # Komponen Tiket Saya & Admin Booking List
│   │   ├── destinations/     # Komponen Katalog & Detail Tempat Wisata
│   │   ├── layout/           # Navbar, Footer & Sidebar Navigation
│   │   ├── reports/          # Admin Report & Export Laporan Analytics
│   │   ├── settings/         # Pengaturan Sistem & Profil Pengguna
│   │   ├── users/            # Kelola Pengguna (User Management)
│   │   └── weather/          # Widget Ramalan Cuaca OpenWeatherMap
│   ├── store/                # Zustand Store (State Auth, Booking, Notification)
│   ├── main.tsx              # Titik masuk React (App mounting)
│   └── index.css             # Tailwind CSS & Global Stylesheet
└── package.json              # Daftar library & script npm
```
