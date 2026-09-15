// Package dto (Data Transfer Object) berisi struct-struct yang dipakai
// khusus untuk komunikasi dengan client (request masuk dan response keluar).
//
// KENAPA DTO DIPISAH DARI MODEL?
// Struct di package "model" itu representasi tabel database (bisa punya
// field sensitif seperti password hash). Kalau kita langsung kirim struct
// model itu sebagai response API, field sensitif bisa ikut ke-expose ke
// client tanpa sengaja. Dengan DTO, kita kontrol penuh field apa saja
// yang boleh masuk (request) dan keluar (response) API.
package dto

import "time"

// RegisterRequest adalah bentuk data yang diharapkan saat ada request
// POST /auth/register. Field "json:..." menentukan nama field di JSON
// yang dikirim client, misal client kirim {"name": "Budi", ...}.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`    // opsional, boleh dikosongkan client
	Password string `json:"password"`
}

// LoginRequest adalah bentuk data untuk request POST /auth/login.
// Identifier bisa diisi EMAIL ATAU NOMOR HP (sesuai requirement bahwa
// user boleh login pakai salah satu dari keduanya).
type LoginRequest struct {
	Identifier string `json:"identifier"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

// UserResponse adalah bentuk data user yang AMAN dikirim ke client.
// Perhatikan: TIDAK ADA field Password di sini sama sekali — ini sengaja,
// supaya hash password tidak pernah ikut terkirim ke client dalam kondisi apa pun.
type UserResponse struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Phone           *string    `json:"phone,omitempty"`             // "omitempty" = field ini disembunyikan dari JSON kalau kosong
	PhoneVerifiedAt *time.Time `json:"phone_verified_at,omitempty"`
	Role            string     `json:"role"`
	CreatedAt       time.Time  `json:"created_at"`
}

// AuthResponse adalah response yang dikembalikan setelah register/login
// berhasil. Berisi token JWT (untuk dipakai di request selanjutnya sebagai
// bukti bahwa user sudah login) dan data profil user.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
