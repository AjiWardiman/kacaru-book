# Entity Relationship Diagram (ERD) - Tourism Booking System Bima

Dokumen ini berisi struktur Entity Relationship Diagram (ERD) berbasis PostgreSQL yang dirancang dan diimplementasikan untuk aplikasi **Tourism Booking System Kabupaten Bima**.

---

## 📐 Diagram ERD (Mermaid)

```mermaid
erDiagram
    USERS ||--o{ DESTINATIONS : "admin creates"
    USERS ||--o{ BOOKINGS : "places"
    USERS ||--o{ OTPS : "receives"
    USERS ||--o{ NOTIFICATIONS : "gets"
    USERS ||--o{ AUDIT_LOGS : "performs action"

    DESTINATIONS ||--o{ SLOTS : "has daily schedules"

    SLOTS ||--o{ BOOKINGS : "reserved in"

    BOOKINGS ||--o{ BOOKING_DETAILS : "contains visitors"
    BOOKINGS ||--|| PAYMENTS : "has invoice"

    PAYMENTS ||--o{ PAYMENT_EVENTS : "logs webhook events"

    USERS {
        bigint id PK
        string name
        string email UK
        string phone UK
        string password
        enum role "user | admin"
        timestamp phone_verified_at
        timestamp created_at
        timestamp updated_at
    }

    DESTINATIONS {
        bigint id PK
        bigint created_by FK
        string name
        string location
        text description
        string category
        decimal price
        string image_url
        boolean is_active
        timestamp created_at
        timestamp updated_at
    }

    SLOTS {
        bigint id PK
        bigint destination_id FK
        date date
        int total_quota
        int used_quota
        boolean is_available
        timestamp created_at
        timestamp updated_at
    }

    BOOKINGS {
        bigint id PK
        string booking_code UK
        bigint user_id FK
        bigint slot_id FK
        int total_people
        decimal total_price
        enum status "pending | paid | expired | cancelled"
        timestamp booked_at
        timestamp expired_at
        timestamp created_at
        timestamp updated_at
    }

    BOOKING_DETAILS {
        bigint id PK
        bigint booking_id FK
        string visitor_name
        string visitor_phone
        string visitor_email
        int age
        timestamp created_at
    }

    PAYMENTS {
        bigint id PK
        bigint booking_id FK,UK
        string external_id UK
        string invoice_id
        decimal amount
        enum status "pending | paid | expired | failed"
        string payment_method
        string payment_url
        timestamp created_at
        timestamp updated_at
    }

    PAYMENT_EVENTS {
        bigint id PK
        bigint payment_id FK
        string event_type
        text payload
        timestamp created_at
    }

    OTPS {
        bigint id PK
        bigint user_id FK
        string code
        string purpose
        boolean is_used
        timestamp expired_at
        timestamp created_at
    }

    NOTIFICATIONS {
        bigint id PK
        bigint user_id FK
        string type
        string title
        text message
        boolean is_read
        timestamp sent_at
    }

    SETTINGS {
        bigint id PK
        string key UK
        string value
        text description
        timestamp updated_at
    }

    AUDIT_LOGS {
        bigint id PK
        bigint user_id FK
        string action
        string entity_type
        bigint entity_id
        jsonb old_values
        jsonb new_values
        timestamp created_at
    }
```

---

## 🔗 Penjelasan Relasi Entitas (Kardinalitas)

1. **`USERS` ➔ `DESTINATIONS` (1 : N)**
   - Satu pengguna berinisial Admin (`role = 'admin'`) dapat mendaftarkan dan mengelola banyak destinasi wisata.

2. **`DESTINATIONS` ➔ `SLOTS` (1 : N)**
   - Satu destinasi wisata memiliki banyak slot jadwal kunjungan harian (`date`, `total_quota`, `used_quota`).

3. **`USERS` & `SLOTS` ➔ `BOOKINGS` (1 : N)**
   - Seorang wisatawan (`user_id`) dapat memesan banyak tiket (`BOOKINGS`).
   - Setiap pemesanan dikaitkan dengan satu slot tanggal (`slot_id`) yang mengurangi `used_quota` pada slot tersebut.

4. **`BOOKINGS` ➔ `BOOKING_DETAILS` (1 : N)**
   - Satu kode pemesanan (`BOOKINGS`) dapat berisi 1 sampai banyak pengunjung rombongan (`BOOKING_DETAILS`), menyimpan nama, nomor HP, email, dan usia pengunjung.

5. **`BOOKINGS` ➔ `PAYMENTS` (1 : 1)**
   - Setiap satu pemesanan (`BOOKINGS`) memiliki TEPAT SATU invoice transaksi pembayaran (`PAYMENTS`) via Xendit Payment Gateway.

6. **`PAYMENTS` ➔ `PAYMENT_EVENTS` (1 : N)**
   - Setiap pembayaran mencatat riwayat log perubahan event callback webhook Xendit (seperti `PAID`, `EXPIRED`, `FAILED`).

7. **`USERS` ➔ `OTPS`, `NOTIFICATIONS`, `AUDIT_LOGS` (1 : N)**
   - Satu pengguna dapat memiliki banyak log kode OTP, notifikasi penerimaan tiket, dan rekam audit aktivitas sistem.
