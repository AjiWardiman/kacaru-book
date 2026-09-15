// Package util berisi fungsi-fungsi bantuan (helper) kecil yang dipakai
// di berbagai bagian aplikasi, seperti hashing password dan pembuatan token JWT.
package util

import "golang.org/x/crypto/bcrypt"

// HashPassword mengubah password asli (plain text) menjadi hash yang aman
// untuk disimpan di database. Kita TIDAK PERNAH menyimpan password asli
// di database — kalau database bocor, penyerang tidak akan langsung tahu
// password asli user.
//
// bcrypt.DefaultCost menentukan seberapa "berat" proses hashing ini
// (semakin berat, semakin aman tapi semakin lambat). DefaultCost adalah
// nilai yang direkomendasikan untuk kebanyakan aplikasi.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword membandingkan password yang diketik user saat login
// dengan hash yang tersimpan di database. Mengembalikan true kalau cocok.
//
// Catatan: kita tidak bisa "membalikkan" hash menjadi password asli
// (itu memang tujuan dari hashing, bersifat satu arah). Makanya untuk
// mengecek kecocokan, bcrypt punya fungsi khusus CompareHashAndPassword ini.
func CheckPassword(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
