// Package middleware berisi fungsi-fungsi yang "menyaring" request
// SEBELUM sampai ke handler tujuannya. Contoh pemakaian di sini: memeriksa
// apakah request punya token JWT yang valid, dan apakah role user yang
// login sesuai dengan yang dibutuhkan endpoint tersebut.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/util"
)

// Nama key yang dipakai untuk menyimpan data user_id dan role di dalam
// context Gin, supaya bisa dibaca lagi oleh middleware/handler berikutnya
// dalam alur request yang sama.
const (
	ContextUserID = "user_id"
	ContextRole   = "role"
)

// Auth adalah middleware yang memeriksa token JWT di header
// "Authorization: Bearer <token>". Kalau token valid, kita simpan
// user_id dan role ke context Gin supaya handler/middleware berikutnya
// bisa memakainya tanpa perlu membaca ulang tokennya.
//
// Kalau token TIDAK ADA atau TIDAK VALID, request langsung dihentikan
// di sini (c.AbortWithStatusJSON) — tidak akan diteruskan ke handler
// tujuan sama sekali.
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		var tokenString string

		if strings.HasPrefix(header, "Bearer ") {
			tokenString = strings.TrimPrefix(header, "Bearer ")
		} else if queryToken := c.Query("token"); queryToken != "" {
			tokenString = queryToken
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "token otorisasi tidak ditemukan"})
			return
		}

		claims, err := util.ParseToken(jwtSecret, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "token tidak valid atau sudah kadaluarsa"})
			return
		}

		// Simpan info user ke context, supaya handler tujuan (atau
		// middleware RequireRole di bawah) bisa memakainya.
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)

		// c.Next() artinya "lanjutkan ke middleware/handler berikutnya".
		// Kalau baris ini tidak dipanggil, request akan berhenti di sini.
		c.Next()
	}
}

// RequireRole adalah middleware yang memeriksa apakah role user yang
// sedang login termasuk salah satu dari daftar "allowed" (role yang
// diizinkan). Middleware ini HARUS dipasang SETELAH Auth di router,
// karena dia membaca ContextRole yang diisi oleh Auth.
//
// Contoh pemakaian: RequireRole("admin") berarti hanya user dengan
// role "admin" yang boleh lewat, selain itu ditolak dengan status 403.
func RequireRole(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextRole)
		if !exists {
			// Ini seharusnya tidak pernah terjadi kalau middleware Auth
			// sudah dipasang duluan, tapi tetap dijaga untuk berjaga-jaga.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "error": "belum login"})
			return
		}

		roleStr, _ := role.(string)
		for _, a := range allowed {
			if roleStr == a {
				// Role user cocok dengan salah satu yang diizinkan → lanjutkan.
				c.Next()
				return
			}
		}

		// Role user login sah, tapi tidak punya izin untuk endpoint ini.
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "error": "tidak punya izin akses"})
	}
}
