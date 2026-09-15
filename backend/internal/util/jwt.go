package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken dikembalikan kalau token JWT yang dikirim client
// tidak valid (rusak, sudah kadaluarsa, atau ditandatangani dengan
// secret yang salah).
var ErrInvalidToken = errors.New("invalid or expired token")

// Claims adalah "isi" dari token JWT yang kita buat.
// Selain data standar (kapan token dibuat, kapan kadaluarsa — dari
// jwt.RegisteredClaims), kita tambahkan UserID dan Role supaya nanti
// saat client mengirim token ini di request berikutnya, kita bisa
// langsung tahu siapa yang login dan apa perannya TANPA perlu query
// ulang ke database.
type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken membuat token JWT baru untuk user yang berhasil
// login/register. Token ini yang nanti dikirim balik ke client, dan
// harus disertakan client di header "Authorization: Bearer <token>"
// setiap kali mengakses endpoint yang butuh login.
//
// expiryHours menentukan berapa lama token ini berlaku sebelum harus
// login ulang (misal 24 jam).
func GenerateToken(secret string, userID int64, role string, expiryHours int) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	// HS256 adalah algoritma penandatanganan yang dipakai untuk memastikan
	// token ini benar-benar dibuat oleh server kita, bukan dipalsukan orang lain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken membaca dan memvalidasi token JWT yang dikirim client.
// Kalau token valid, kita dapatkan kembali isi Claims-nya (UserID, Role, dll).
// Kalau token rusak/kadaluarsa/palsu, akan mengembalikan ErrInvalidToken.
func ParseToken(secret, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
