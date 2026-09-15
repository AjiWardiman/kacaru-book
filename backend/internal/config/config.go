// Package config bertugas membaca semua pengaturan aplikasi (alamat
// database, secret JWT, kredensial email, dll) dari file .env atau
// dari environment variable komputer/server.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menyimpan semua pengaturan yang dibutuhkan aplikasi supaya
// bisa jalan. Dikumpulkan di satu tempat biar gampang dicari dan diubah.
type Config struct {
	AppPort string // di port berapa server ini akan berjalan, misal "8080"

	// Pengaturan koneksi ke database PostgreSQL
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Pengaturan untuk membuat token JWT (dipakai saat login)
	JWTSecret     string // kunci rahasia untuk menandatangani token, JANGAN dibagikan ke siapa pun
	JWTExpiryHour int    // token akan kadaluarsa setelah berapa jam

	// Pengaturan SMTP Mailtrap, dipakai untuk mengirim email notifikasi
	// (dipakai di tahap selanjutnya, belum aktif di kode saat ini)
	MailtrapHost     string
	MailtrapPort     string
	MailtrapUsername string
	MailtrapPassword string
	MailFromAddress  string

	// Pengaturan Xendit, dipakai untuk proses pembayaran
	XenditSecretKey     string
	XenditCallbackToken string
	XenditBaseURL       string // alamat API Xendit; kosongkan untuk pakai alamat resmi

	// Pengaturan OpenRouter AI API
	OpenRouterAPIKey string
	OpenRouterModel  string
}

// Load membaca file .env (kalau ada) lalu mengumpulkan semua pengaturan
// menjadi satu struct Config yang siap dipakai di seluruh aplikasi.
//
// Cara kerjanya:
//  1. Coba baca file .env di folder project. Kalau file .env tidak
//     ditemukan, itu TIDAK dianggap error fatal — aplikasi akan tetap
//     jalan menggunakan environment variable asli komputer/server, atau
//     nilai default kalau memang tidak ada sama sekali (praktis untuk
//     dijalankan di server produksi yang memang tidak pakai file .env).
//  2. Untuk tiap pengaturan, ambil nilainya dari environment variable.
//     Kalau tidak ada, pakai nilai default yang sudah disediakan
//     (lihat fungsi getEnv di bawah).
func Load() *Config {
	if err := godotenv.Load(".env", "../.env"); err != nil {
		log.Println("file .env tidak ditemukan, memakai environment variable asli / nilai default")
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "admin"),
		DBName:     getEnv("DB_NAME", "tourism_booking"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:     getEnv("JWT_SECRET", "change-this-secret-in-production"),
		JWTExpiryHour: 24,

		MailtrapHost:     getEnv("MAILTRAP_HOST", "sandbox.smtp.mailtrap.io"),
		MailtrapPort:     getEnv("MAILTRAP_PORT", "2525"),
		MailtrapUsername: getEnv("MAILTRAP_USERNAME", ""),
		MailtrapPassword: getEnv("MAILTRAP_PASSWORD", ""),
		MailFromAddress:  getEnv("MAIL_FROM_ADDRESS", "no-reply@tourism-booking.local"),

		XenditSecretKey:     getEnv("XENDIT_SECRET_KEY", ""),
		XenditCallbackToken: getEnv("XENDIT_CALLBACK_TOKEN", ""),
		XenditBaseURL:       getEnv("XENDIT_BASE_URL", ""), // kosong = pakai alamat resmi Xendit

		OpenRouterAPIKey: getEnv("OPENROUTER_API_KEY", ""),
		OpenRouterModel:  getEnv("OPENROUTER_MODEL", "anthropic/claude-3-haiku"),
	}
}

// getEnv adalah fungsi bantuan kecil: ambil nilai environment variable
// berdasarkan key-nya. Kalau environment variable itu tidak ada atau
// kosong, kembalikan nilai "fallback" (nilai cadangan/default) yang sudah
// kita siapkan, supaya aplikasi tetap bisa jalan meski file .env belum
// lengkap diisi.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
