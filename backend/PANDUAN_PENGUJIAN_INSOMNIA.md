# Panduan Pengujian API Tourism Booking System menggunakan Insomnia

Panduan ini berisi langkah-langkah praktis dan urutan pengujian seluruh endpoint API **Tourism Booking System Bima** menggunakan aplikasi **Insomnia**.

---

## 1. Persiapan Server & Database

1. Pastikan database PostgreSQL Anda sudah berjalan (default port `5433` atau `5432` sesuai `.env`).
2. Jalankan aplikasi Go:
   ```powershell
   go run ./cmd/api
   ```
   *(Pastikan muncul log: `Server listening on :8080`)*

---

## 2. Cara Impor Koleksi Request ke Insomnia

Aplikasi ini sudah menyediakan file ekspor koleksi request yang siap pakai di folder proyek:
📁 **`insomnia/tourism-booking-insomnia.json`**

### Langkah Import:
1. Buka aplikasi **Insomnia**.
2. Klik tombol **Create / Import** (atau tombol `+` / menu **Preferences/Settings** -> **Data** -> **Import Data**).
3. Pilih **From File**.
4. Arahkan ke file [`insomnia/tourism-booking-insomnia.json`](file:///e:/Golang/tugas-akhir/tourism-booking/insomnia/tourism-booking-insomnia.json).
5. Klik **Import**. Koleksi bernama **"Tourism Booking API"** akan otomatis muncul di Insomnia dengan 11 folder request lengkap.

---

## 3. Konfigurasi Environment Variable di Insomnia

Insomnia menggunakan Environment Variable agar Anda tidak perlu mengetik alamat URL dan token secara manual berulang-ulang.

1. Di sebelah kiri atas Insomnia, klik dropdown **Base Environment** -> pilih **Manage Environments**.
2. Isikan data berikut:
   ```json
   {
     "base_url": "http://localhost:8080",
     "token": "PASTE_TOKEN_JWT_DI_SINI",
     "destination_id": "1",
     "booking_id": "1",
     "user_id": "1"
   }
   ```
3. Klik **Done**.

---

## 4. Urutan Langkah-Langkah Pengujian (Step-by-Step)

Jalankan request di Insomnia sesuai urutan bernomor di bawah ini:

### 🔹 Langkah 1: Health Check (Memastikan Server Hidup)
- **Folder**: `1. Health Check`
- **Request**: `Health Check Status Server`
- **URL**: `{{ _.base_url }}/health`
- **Ekspektasi Status**: `200 OK`
- **Respon**: `{"success": true, "message": "layanan berjalan normal"}`

---

### 🔹 Langkah 2: Autentikasi (Register & Login)

#### 2a. Register User Baru (`POST /auth/register`)
- **Folder**: `2. Auth` -> `Register User Baru`
- **Body JSON**:
  ```json
  {
    "name": "Budi Santoso",
    "email": "budi.user@example.com",
    "phone": "089988776655",
    "password": "password123"
  }
  ```
- **Ekspektasi Status**: `201 Created`

#### 2b. Login (`POST /auth/login`)
- **Folder**: `2. Auth` -> `Login - Admin (Aji)` / `Login - Customer (Budi)`
- **Body JSON**:
  ```json
  {
    "identifier": "aji@example.com",
    "password": "password123"
  }
  ```
- **Ekspektasi Status**: `200 OK`
- **Salin Token**: Dari respon JSON di bagian `data.token`, **salin string token tersebut** lalu tempel ke Environment Variable Insomnia pada field `"token"`.

---

### 🔹 Langkah 3: Pengujian Katalog Wisata & Slot (Publik)
- **Folder**: `3. Destinations & Slots (Publik)`
- **Request**:
  - `GET Get List Destinations` (`/destinations`)
  - `GET Get Destination Detail by ID` (`/destinations/{{ _.destination_id }}`)
  - `GET Get Slots by Destination ID` (`/destinations/{{ _.destination_id }}/slots`)
- **Ekspektasi Status**: `200 OK`

---

### 🔹 Langkah 4: Pengujian Fitur AI & Weather
- **Folder**: `4. Fitur AI & Weather (Publik)`
- **Request**:
  - `GET Get Weather Forecast` (`/weather?location=Bima, NTB`)
  - `POST AI Destination Recommendation` (`/recommendations`)
    ```json
    {
      "max_budget": 50000,
      "category": "Bahari",
      "guests": 2,
      "preference": "pantai pasir putih bersunset indah"
    }
    ```
- **Ekspektasi Status**: `200 OK`

---

### 🔹 Langkah 5: Pengujian Booking & Payment Gateway Xendit
- **Folder**: `6. Booking & Payment Xendit (Wajib Login)`
- **Request**:
  1. `POST Create Booking Baru` (`/bookings`)
  2. `POST Create Payment Link Xendit` (`/bookings/{{ _.booking_id }}/pay`)
  3. `POST Xendit Callback Webhook` (`/webhooks/xendit`)
- **Ekspektasi Status**: `200 OK` / `201 Created`

---

### 🔹 Langkah 6: Pengujian Laporan Analytics, Ekspor Excel & PDF (Khusus Admin)
- **Folder**: `8. Laporan & Analytics (Admin Only)`
- **Header**: `Authorization: Bearer {{ _.token }}` (Wajib Admin Token)
- **Request**:
  - `GET Get Report Summary` (`/admin/reports/summary`)
  - `GET Export Report Excel` (`/admin/reports/export/excel?token={{ _.token }}`)
  - `GET Export Report PDF` (`/admin/reports/export/pdf?token={{ _.token }}`)
- **Ekspektasi Status**: `200 OK`

---

### 🔹 Langkah 7: Pengujian Manajemen User (Khusus Admin)
- **Folder**: `9. User Management (Admin Only)`
- **Request**:
  - `GET Get List Users` (`/admin/users`)
  - `POST Create User Baru` (`/admin/users`)
  - `PUT Update User` (`/admin/users/{{ _.user_id }}`)
  - `DELETE Delete User` (`/admin/users/{{ _.user_id }}`)
- **Ekspektasi Status**: `200 OK` / `201 Created`

---

### 🔹 Langkah 8: Pengujian Settings & Config (Khusus Admin)
- **Folder**: `10. Settings / Pengaturan (Admin Only)`
- **Request**:
  - `GET Get All Settings` (`/admin/settings`)
  - `PUT Update Setting Value` (`/admin/settings/otp_expiry_minutes`)
- **Ekspektasi Status**: `200 OK`
