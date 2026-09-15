// Package handler adalah lapisan yang langsung "berhadapan" dengan
// request HTTP dari client. Tugasnya: membaca request (body JSON,
// parameter URL, header), memanggil service yang sesuai, lalu menyusun
// response yang dikirim balik ke client. Handler TIDAK BOLEH berisi
// logic bisnis (itu tugas service) dan TIDAK BOLEH langsung memanggil database
// (itu tugas repository).
package handler

import (
	"github.com/gin-gonic/gin"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/middleware"
)

// writeSuccess mengirim response sukses ke client dengan format yang
// konsisten (lihat dto.APIResponse). Dipakai oleh semua handler supaya
// bentuk response selalu seragam di seluruh aplikasi.
func writeSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, dto.APIResponse{Success: true, Message: message, Data: data})
}

// writeError mengirim response gagal ke client dengan format yang konsisten.
func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, dto.APIResponse{Success: false, Error: message})
}

// currentUserID adalah fungsi bantuan untuk mengambil user_id dari
// context Gin (diisi oleh middleware Auth setelah memvalidasi token).
// Dipakai di berbagai handler untuk mencatat SIAPA yang melakukan
// sebuah aksi, misalnya untuk keperluan audit trail.
func currentUserID(c *gin.Context) int64 {
	val, _ := c.Get(middleware.ContextUserID)
	id, _ := val.(int64)
	return id
}
