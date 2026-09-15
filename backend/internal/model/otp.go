package model

import "time"

// OTPPurpose adalah tipe data untuk tujuan/keperluan sebuah kode OTP
// dibuat. Dibatasi hanya boleh 3 pilihan ini (sesuai ENUM "otp_purpose"
// di database), supaya tidak ada typo yang bisa lolos ke database.
type OTPPurpose string

const (
	// OTPPurposePhoneVerification dipakai saat user memverifikasi
	// nomor HP-nya. Di aplikasi ini, pengiriman kodenya DISIMULASIKAN
	// lewat email (bukan SMS asli) sesuai kesepakatan tugas akhir.
	OTPPurposePhoneVerification OTPPurpose = "phone_verification"
	OTPPurposeLogin             OTPPurpose = "login"
	OTPPurposeResetPassword     OTPPurpose = "reset_password"
)

// OTP merepresentasikan satu baris di tabel "otps".
//
// Catatan penting:
//   - Code disimpan APA ADANYA (bukan hash), karena kode OTP ini
//     memang cuma berlaku sangat singkat (lihat ExpiredAt) dan sekali
//     pakai (lihat IsUsed) — berbeda dengan password yang harus di-hash
//     karena disimpan permanen.
//   - Satu user bisa punya banyak baris OTP seiring waktu (setiap kali
//     minta kode baru), makanya IsUsed dan ExpiredAt penting untuk
//     memastikan hanya kode yang PALING BARU dan BELUM KADALUARSA yang
//     bisa dipakai berhasil verifikasi.
type OTP struct {
	ID        int64      // ID unik baris OTP
	UserID    int64      // milik user yang mana
	Code      string     // kode OTP, 6 digit angka
	Purpose   OTPPurpose // untuk keperluan apa kode ini dibuat
	IsUsed    bool       // true kalau kode ini sudah pernah berhasil dipakai verifikasi
	ExpiredAt time.Time  // batas waktu kode ini masih boleh dipakai
	CreatedAt time.Time  // kapan kode ini dibuat
}

func (OTP) TableName() string {
	return "otps"
}

