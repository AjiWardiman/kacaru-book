# 🚀 PANDUAN DEMO & PRESENTASI PRAKTIS (BAHASA AWAM & MUDAH DIPAHAMI)

Dokumen ini adalah panduan lengkap langkah demi langkah untuk mempresentasikan aplikasi **Tourism Booking System** di depan Trainer, Dosen Penguji, maupun Stakeholder. Panduan ini dirancang agar demo berjalan **100% lancar, meyakinkan, dan berkesan profesional**.

---

## 📋 CHECKLIST PERSIAPAN (5 MENIT SEBELUM DEMO)

Sebelum memanggil Trainer / Penguji, pastikan hal-hal berikut sudah siap di komputer Anda:

- [ ] **Backend Golang Berjalan:** Terminal 1 menjalankan `go run ./cmd/api` (Server di `http://localhost:8080`).
- [ ] **Frontend React Berjalan:** Terminal 2 menjalankan `npm run dev` (Aplikasi di `http://localhost:5173`).
- [ ] **Database PostgreSQL Active:** Database `tourism_booking` sudah terhubung tanpa error.
- [ ] **Tab Browser Siap:** Buka browser di `http://localhost:5173` dalam posisi zoom 100%.
- [ ] **Kredensial Pengujian Siap Catat:**
  - **Akun Wisatawan (User):** `user@gmail.com` / `password123` (atau buat akun baru saat demo).
  - **Akun Pengelola (Admin):** `admin@tourism.com` / `admin123`.

---

## 🔗 TAUTAN MEDIA & SLIDE PRESENTASI

- 📊 **Slide Presentasi Kacaru Book (PPT):** [Buka Slide PPT Google Slides](https://docs.google.com/presentation/d/1xMfge3mo1Qsj4tdgSVCFg1z3gZwIddN-/edit?usp=drive_link&ouid=103647142140848675436&rtpof=true&sd=true)
- 🎨 **ER Diagram Database:** [Buka ERD Google Drive](https://drive.google.com/file/d/18pCFW9baf46Ih-uB5w8gMmOBebhcUE6p/view?usp=drive_link)
- 🎥 **Video Rekaman Demo Web:** [Buka Video Demo Google Drive](https://drive.google.com/file/d/1d9S-iIJGm-hRSf_4dwSbm_cHXEUjP5kK/view?usp=drive_link)

---

## 🎬 SKENARIO DEMO (URUTAN KLIK & SCRIPT PENJELASAN)

### 🎙️ PEMBUKAAN (30 Detik Pertama)
> *"Selamat pagi/siang Bapak/Ibu Trainer dan Penguji. Hari ini saya akan mendemonstrasikan **Tourism Booking System**, yaitu platform reservasi tiket wisata digital berbasis web yang dilengkapi dengan rekomendasi AI cerdas, integrasi cuaca realtime, sistem kuota anti-bentrok, dan e-ticket otomatis."*

---

### STAGE 1: Registrasi, Login & Verifikasi Keamanan OTP
1. **Buka Halaman Login / Registrasi:**
   - Klik tombol **"Masuk / Daftar"** di kanan atas Navbar.
   - Tunjukkan form registrasi yang bersih dan responsif.
2. **Demo Fitur OTP (One-Time Password):**
   - Masukkan nomor HP/Email, lalu klik **"Minta Kode OTP"**.
   - Tunjukkan modal verifikasi 6-digit kode OTP.
3. **Script Penjelasan ke Trainer:**
   > *"Untuk memastikan keamanan akun wisatawan dan menghindari pendaftaran palsu, sistem kami dilengkapi verifikasi OTP 6-digit berbatas waktu yang dikirim langsung ke pengguna."*

---

### STAGE 2: Katalog Destinasi, Toggle Mode Terang/Gelap & Forecast Cuaca Realtime
1. **Demo Tampilan & Switcher Light/Dark Mode:**
   - Tunjukkan katalog tempat wisata (Lengkap dengan gambar, harga tiket, rating, dan badge kategori).
   - Klik ikon **Matahari / Bulan** di Navbar untuk menunjukkan peralihan **Light Mode (Terang)** dan **Dark Mode (Gelap)** yang halus dan kontras tinggi.
2. **Buka Detail Destinasi & Widget Cuaca (OpenWeatherMap):**
   - Klik salah satu destinasi (misal: *Pantai Lariti* atau *Toro Maria*).
   - Sorot widget **OpenWeatherMap Forecast** yang menampilkan suhu (misal: `29.5°C Cerah Berawan`), tingkat kelembapan, kecepatan angin, dan catatan saran keselamatan.
3. **Script Penjelasan ke Trainer:**
   > *"Di halaman katalog, pengguna dapat mengganti tema Light/Dark Mode sesuai kenyamanan mata. Pada detail destinasi, kami mengintegrasikan OpenWeatherMap API untuk menampilkan ramalan cuaca realtime 5 hari ke depan serta rekomendasi keselamatan outdoor bagi wisatawan."*

---

### STAGE 3: Fitur Rekomendasi AI Cerdas (Point of Interest)
1. **Klik Tombol "✨ Rekomendasi AI":**
   - Di bagian atas katalog, klik tombol **"✨ Rekomendasi AI"**.
   - Geser slider **Maksimal Tiket** (misal: Rp 50.000) dan pilih **Kategori Kesukaan** (misal: "Bahari").
   - Klik **"Cari Rekomendasi"**.
2. **Tunjukkan Hasil Rekomendasi:**
   - Tunjukkan kartu destinasi hasil rekomendasi lengkap dengan persentase kecocokan (misal: **98% Cocok**) dan alasan analisis dari AI.
3. **Script Penjelasan ke Trainer:**
   > *"Fitur AI Rekomendasi ini membantu wisatawan yang bingung memilih destinasi. Dengan memasukkan budget dan minat wisata, AI akan secara otomatis menganalisis dan mencocokkan tempat wisata terbaik di database secara realtime."*

---

### STAGE 4: Pemesanan Anti-Bentrok Kuota, Xendit Payment & E-Ticket Digital
1. **Pilih Tanggal & Kuota Kunjungan:**
   - Pilih slot tanggal kunjungan yang kuotanya tersedia.
   - Isi jumlah pengunjung (misal: 2 orang) dan klik **"Pesan Sekarang"**.
2. **Proses Pembayaran Digital (Xendit Gateway):**
   - Tunjukkan invoice pembayaran yang terbit otomatis.
   - Klik simulasi bayar lunas.
3. **Buka E-Ticket Digital ber-QR Code:**
   - Masuk ke menu **"Tiket Saya"**, klik **"Lihat Tiket"**.
   - Tunjukkan tampilan **QR Code E-Ticket** digital yang siap di-scan petugas di pintu masuk dan tombol **Cetak / PDF**.
4. **Script Penjelasan ke Trainer:**
   > *"Sistem kami menggunakan algoritma Row Locking (SELECT FOR UPDATE) di tingkat database, sehingga menjamin kuota tiket anti-bentrok meskipun dipesan oleh puluhan wisatawan secara bersamaan. Setelah bayar via Xendit, E-Ticket ber-QR Code langsung terbit otomatis."*

---

### STAGE 5: Dashboard Admin, Manajemen Kuota & Export Laporan Analytics
1. **Login Akun Admin:**
   - Masuk sebagai Admin (`admin@tourism.com`).
2. **Buka Dashboard Laporan / Analytics:**
   - Tunjukkan ringkasan statistik: **Total Pendapatan**, **Jumlah Transaksi**, **Persentase Okupansi**, dan **Grafik Bulanan**.
3. **Demo Export Laporan (PDF & Excel):**
   - Klik tombol **"📥 Export PDF"** dan **"📊 Export Excel"**.
   - Buka file hasil unduhan untuk menunjukkan format laporan yang rapi dan siap cetak.
4. **Script Penjelasan ke Trainer:**
   > *"Bagi pengelola/admin, kami menyediakan Dashboard Analytics untuk memantau pendapatan dan keterisian tempat wisata secara realtime. Laporan keuangan dan pengunjung juga dapat di-export langsung ke file PDF dan Excel hanya dalam satu kali klik."*

---

## 🛡️ TANGKISAN PERTANYAAN DEMO PROYEK (Q&A DEFENSE FOR TRAINER)

Berikut adalah pertanyaan yang paling sering ditanyakan Penguji dan jawaban terbaiknya:

### ❓ Q1: *"Bagaimana cara sistem mencegah 2 orang memesan sisa 1 tiket secara bersamaan?"*
> **Jawaban:**  
> *"Kami menerapkan mekanisme **Database Row Locking (`SELECT FOR UPDATE`)** pada transaksi pemesanan di Golang. Ketika Request A masuk, database langsung mengunci baris data slot tersebut secara atomik sampai transaksi A selesai, baru Request B diproses. Hasilnya, tidak akan terjadi double-booking atau kuota minus."*

### ❓ Q2: *"Bagaimana jika pengunjung sudah pesan tapi tidak melakukan pembayaran?"*
> **Jawaban:**  
> *"Backend Golang kami memiliki **Background Scheduler Worker** yang berjalan otomatis setiap 1 menit. Worker ini akan memeriksa transaksi yang melewati batas waktu bayar (expired), membatalkan pesanan secara otomatis, dan mengembalikan kuota tiket ke database agar dapat dipesan orang lain."*

### ❓ Q3: *"Bagaimana jika koneksi API Cuaca atau AI mengalami gangguan (offline/limit)?"*
> **Jawaban:**  
> *"Sistem dibangun dengan **Fallback Mechanism & Resilience**. Jika API eksternal gagal merespons, aplikasi tidak akan crash/error, melainkan secara otomatis menyajikan data perkiraan lokal yang aman sehingga pengalaman pengguna tetap berjalan lancar."*

### ❓ Q4: *"Seberapa aman sistem autentikasi dan data di aplikasi ini?"*
> **Jawaban:**  
> *"Password pengguna dienkripsi menggunakan **bcrypt** dengan salt. Akses API dilindungi oleh **JSON Web Token (JWT)** dengan expiration period, serta otorisasi Role-Based Access Control (RBAC) yang memisahkan hak akses User biasa dan Admin."*

---

## 💻 TIPS MEMBUKA "INSPECT ELEMENT" (F12) SAAT DEMO APLIKASI

Jika Trainer/Penguji meminta Anda membuka **F12 / Inspect Element**:

1. **Tab Console:**  
   - Tunjukkan bahwa Console bersih **bebas dari error merah**.
2. **Tab Network:**  
   - Lakukan satu aksi (misal: klik Rekomendasi AI atau ambil data destinasi).
   - Tunjukkan request API mengembalikan respon **`200 OK`** dengan payload JSON yang rapi.
3. **Penjelasan:**  
   - *"Semua komunikasi antarmuka frontend dan server backend terstandarisasi lewat REST API berformat JSON dengan exception handling yang ketat."*
