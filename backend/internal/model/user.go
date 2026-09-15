// Package model berisi struct-struct yang merepresentasikan baris data
// di setiap tabel database. Struct di file ini SENGAJA dibuat mirip
// dengan tabel di database, TIDAK boleh langsung dipakai sebagai
// response API (untuk itu gunakan DTO di package dto).
package model

import "time"

// UserRole adalah tipe data untuk peran/role user.
// Kita batasi cuma boleh "admin" atau "user" (lihat konstanta di bawah),
// supaya tidak ada typo seperti "Admin" atau "admn" yang bisa lolos ke database.
type UserRole string

const (
	RoleAdmin UserRole = "admin" // role untuk pengelola aplikasi (bisa CRUD destinasi, slot, dll)
	RoleUser  UserRole = "user"  // role untuk pengguna biasa (bisa booking, bayar, dll)
)

// User merepresentasikan satu baris di tabel "users".
//
// Catatan penting:
//   - Phone dan PhoneVerifiedAt dibuat pointer (*string, *time.Time) karena
//     kolom ini boleh kosong (NULL) di database — tidak semua user mengisi
//     nomor HP saat daftar.
//   - Password di sini adalah HASH bcrypt, BUKAN password asli/mentah.
//     Password asli tidak pernah disimpan di mana pun demi keamanan.
type User struct {
	ID              int64      // ID unik user, otomatis bertambah dari database
	Name            string     // nama lengkap user
	Email           string     // email, harus unik (tidak boleh sama dengan user lain)
	Phone           *string    // nomor HP, boleh kosong (nil) kalau user tidak mengisi
	PhoneVerifiedAt *time.Time // waktu nomor HP diverifikasi, kosong (nil) kalau belum diverifikasi
	Password        string     // hash bcrypt dari password, BUKAN password asli
	Role            UserRole   // peran user: "admin" atau "user"
	CreatedAt       time.Time  // waktu akun dibuat
	UpdatedAt       time.Time  // waktu data user terakhir diubah
}

func (User) TableName() string {
	return "users"
}

