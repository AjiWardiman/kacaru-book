package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTPCode membuat kode OTP acak berupa angka sepanjang "length"
// digit (misal length=6 menghasilkan kode seperti "482913").
//
// KENAPA PAKAI "crypto/rand" (bukan "math/rand" biasa)? Karena
// "math/rand" hasilnya BISA DITEBAK kalau penyerang tahu seed-nya —
// tidak aman untuk sesuatu yang sifatnya rahasia seperti kode OTP.
// "crypto/rand" adalah generator angka acak yang aman secara
// kriptografis, cocok untuk kebutuhan keamanan seperti ini.
func GenerateOTPCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)

	for i := range code {
		// rand.Int mengambil satu angka acak antara 0 sampai (10-1),
		// lalu kita jadikan indeks untuk mengambil satu karakter digit.
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("gagal membuat kode OTP acak: %w", err)
		}
		code[i] = digits[n.Int64()]
	}

	return string(code), nil
}
