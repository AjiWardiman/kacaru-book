package model

import "time"

// Setting merepresentasikan satu baris di tabel "settings".
//
// Tabel ini dipakai untuk menyimpan pengaturan aplikasi yang bisa
// diubah admin TANPA perlu mengubah kode dan deploy ulang, misalnya
// "berapa menit batas waktu bayar booking sebelum expired" atau
// "berapa menit kode OTP berlaku". Nilainya sengaja disimpan sebagai
// TEXT di database (kolom Value string), lalu di lapisan service kita
// ubah ke tipe data yang dibutuhkan (misal int) saat dipakai.
type Setting struct {
	ID          int64     // ID unik baris setting
	Key         string    // nama pengaturan, misal "booking_expiry_minutes" (unik, dipakai sebagai "kunci" pencarian)
	Value       string    // nilai pengaturan, selalu berbentuk teks
	Description string    // penjelasan singkat, ditampilkan di UI admin supaya jelas fungsi setting ini untuk apa
	UpdatedAt   time.Time // kapan terakhir kali nilai ini diubah
}

func (Setting) TableName() string {
	return "settings"
}

