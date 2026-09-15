package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/service"
)

// SettingHandler menangani request HTTP seputar pengaturan aplikasi
// (settings). Semua endpoint di sini KHUSUS ADMIN — lihat pemasangan
// middleware-nya di router/router.go.
type SettingHandler struct {
	settingService service.SettingService
}

func NewSettingHandler(settingService service.SettingService) *SettingHandler {
	return &SettingHandler{settingService: settingService}
}

// List menangani GET /admin/settings — menampilkan semua pengaturan
// aplikasi yang ada, supaya admin tahu apa saja yang bisa diubah.
// List godoc
// @Summary Mengambil daftar pengaturan aplikasi (Admin-only)
// @Description Menampilkan semua daftar pengaturan aplikasi
// @Tags Settings
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=[]dto.SettingResponse} "Daftar pengaturan aplikasi"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 403 {object} dto.APIResponse "Forbidden (Admin only)"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /admin/settings [get]
func (h *SettingHandler) List(c *gin.Context) {
	result, err := h.settingService.List(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal mengambil daftar pengaturan")
		return
	}
	writeSuccess(c, http.StatusOK, "berhasil mengambil daftar pengaturan", result)
}

// Get godoc
// @Summary Mengambil detail pengaturan aplikasi berdasarkan Key (Admin-only)
// @Description Menampilkan detail satu pengaturan berdasarkan nama Key
// @Tags Settings
// @Produce json
// @Security BearerAuth
// @Param key path string true "Nama Key Setting (misal: otp_expiry_minutes)"
// @Success 200 {object} dto.APIResponse{data=dto.SettingResponse} "Detail pengaturan aplikasi"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 404 {object} dto.APIResponse "Setting tidak ditemukan"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /admin/settings/{key} [get]
func (h *SettingHandler) Get(c *gin.Context) {
	key := c.Param("key")

	result, err := h.settingService.Get(c.Request.Context(), key)
	if err != nil {
		if errors.Is(err, service.ErrSettingNotFound) {
			writeError(c, http.StatusNotFound, "pengaturan dengan key tersebut tidak ditemukan")
			return
		}
		writeError(c, http.StatusInternalServerError, "gagal mengambil pengaturan")
		return
	}
	writeSuccess(c, http.StatusOK, "berhasil mengambil pengaturan", result)
}

// Update godoc
// @Summary Memperbarui nilai pengaturan aplikasi (Admin-only)
// @Description Mengubah nilai dari satu pengaturan aplikasi berdasarkan Key
// @Tags Settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "Nama Key Setting (misal: otp_expiry_minutes)"
// @Param request body dto.UpdateSettingRequest true "Nilai baru setting"
// @Success 200 {object} dto.APIResponse{data=dto.SettingResponse} "Setting berhasil diubah"
// @Failure 400 {object} dto.APIResponse "Data tidak valid"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 404 {object} dto.APIResponse "Setting tidak ditemukan"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /admin/settings/{key} [put]
func (h *SettingHandler) Update(c *gin.Context) {
	key := c.Param("key")

	var req dto.UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}
	if req.Value == "" {
		writeError(c, http.StatusBadRequest, "value wajib diisi")
		return
	}

	result, err := h.settingService.Update(c.Request.Context(), currentUserID(c), key, req.Value)
	if err != nil {
		if errors.Is(err, service.ErrSettingNotFound) {
			writeError(c, http.StatusNotFound, "pengaturan dengan key tersebut tidak ditemukan")
			return
		}
		writeError(c, http.StatusInternalServerError, "gagal mengubah pengaturan")
		return
	}
	writeSuccess(c, http.StatusOK, "berhasil mengubah pengaturan", result)
}
