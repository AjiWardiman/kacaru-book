package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/middleware"
	"tourism-booking/internal/service"
)

// DestinationHandler menangani semua request HTTP seputar destinasi wisata.
type DestinationHandler struct {
	destinationService service.DestinationService
}

func NewDestinationHandler(s service.DestinationService) *DestinationHandler {
	return &DestinationHandler{destinationService: s}
}

// Create menangani POST /admin/destinations.
// Endpoint ini SUDAH DIPROTEKSI oleh middleware Auth + RequireRole di
// router (lihat internal/router/router.go), jadi kalau kode sampai di
// sini, kita sudah pasti request ini datang dari admin yang sah.
// Create godoc
// @Summary Membuat destinasi wisata baru (Admin-only)
// @Description Membuat data destinasi wisata baru
// @Tags Destinations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateDestinationRequest true "Data Destinasi"
// @Success 201 {object} dto.APIResponse{data=dto.DestinationResponse} "Destinasi berhasil dibuat"
// @Failure 400 {object} dto.APIResponse "Data tidak valid"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 403 {object} dto.APIResponse "Forbidden (Admin only)"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /admin/destinations [post]
func (h *DestinationHandler) Create(c *gin.Context) {
	var req dto.CreateDestinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}
	if req.Name == "" || req.Location == "" {
		writeError(c, http.StatusBadRequest, "nama dan lokasi wajib diisi")
		return
	}

	// Ambil user_id dari context Gin. Nilai ini ditaruh di sana oleh
	// middleware Auth setelah berhasil memvalidasi token JWT — jadi
	// kita tidak perlu (dan tidak boleh) percaya begitu saja pada input
	// dari client soal siapa yang membuat destinasi ini.
	userID, _ := c.Get(middleware.ContextUserID)
	createdBy, _ := userID.(int64)

	result, err := h.destinationService.Create(c.Request.Context(), createdBy, req)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal membuat destinasi")
		return
	}
	writeSuccess(c, http.StatusCreated, "destinasi berhasil dibuat", result)
}

// List godoc
// @Summary Mengambil daftar destinasi wisata
// @Description Mengambil daftar destinasi wisata (Publik)
// @Tags Destinations
// @Produce json
// @Param all query bool false "Tampilkan semua destinasi termasuk non-aktif (Admin)"
// @Success 200 {object} dto.APIResponse{data=[]dto.DestinationResponse} "Daftar destinasi"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /destinations [get]
func (h *DestinationHandler) List(c *gin.Context) {
	// Secara default hanya tampilkan destinasi yang aktif.
	// Kalau ada query parameter ?all=true, tampilkan semua (termasuk
	// yang tidak aktif) — ini dipakai admin untuk keperluan pengelolaan.
	onlyActive := c.Query("all") != "true"

	result, err := h.destinationService.List(c.Request.Context(), onlyActive)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal mengambil daftar destinasi")
		return
	}
	writeSuccess(c, http.StatusOK, "berhasil", result)
}

// Get godoc
// @Summary Mengambil detail destinasi wisata
// @Description Mengambil detail destinasi wisata berdasarkan ID (Publik)
// @Tags Destinations
// @Produce json
// @Param id path int true "ID Destinasi"
// @Success 200 {object} dto.APIResponse{data=dto.DestinationResponse} "Detail destinasi"
// @Failure 400 {object} dto.APIResponse "ID tidak valid"
// @Failure 404 {object} dto.APIResponse "Destinasi tidak ditemukan"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /destinations/{id} [get]
func (h *DestinationHandler) Get(c *gin.Context) {
	// ParseInt mengubah teks dari URL (misal "5") menjadi angka (int64).
	// Kalau bukan angka valid (misal "/destinations/abc"), kita tolak.
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "id destinasi tidak valid")
		return
	}

	result, err := h.destinationService.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrDestinationNotFound) {
			writeError(c, http.StatusNotFound, "destinasi tidak ditemukan") // 404: data tidak ada
			return
		}
		writeError(c, http.StatusInternalServerError, "gagal mengambil data destinasi")
		return
	}
	writeSuccess(c, http.StatusOK, "berhasil", result)
}

// Update godoc
// @Summary Mengubah data destinasi wisata (Admin-only)
// @Description Memperbarui data destinasi wisata berdasarkan ID
// @Tags Destinations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Destinasi"
// @Param request body dto.UpdateDestinationRequest true "Data Perubahan Destinasi"
// @Success 200 {object} dto.APIResponse{data=dto.DestinationResponse} "Destinasi berhasil diubah"
// @Failure 400 {object} dto.APIResponse "Data tidak valid"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 404 {object} dto.APIResponse "Destinasi tidak ditemukan"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /admin/destinations/{id} [put]
func (h *DestinationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "id destinasi tidak valid")
		return
	}

	var req dto.UpdateDestinationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.destinationService.Update(c.Request.Context(), currentUserID(c), id, req)
	if err != nil {
		if errors.Is(err, service.ErrDestinationNotFound) {
			writeError(c, http.StatusNotFound, "destinasi tidak ditemukan")
			return
		}
		writeError(c, http.StatusInternalServerError, "gagal mengubah destinasi")
		return
	}
	writeSuccess(c, http.StatusOK, "destinasi berhasil diubah", result)
}

// Delete godoc
// @Summary Menghapus destinasi wisata (Admin-only)
// @Description Menghapus destinasi wisata berdasarkan ID
// @Tags Destinations
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Destinasi"
// @Success 200 {object} dto.APIResponse "Destinasi berhasil dihapus"
// @Failure 400 {object} dto.APIResponse "ID tidak valid"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 404 {object} dto.APIResponse "Destinasi tidak ditemukan"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /admin/destinations/{id} [delete]
func (h *DestinationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "id destinasi tidak valid")
		return
	}

	if err := h.destinationService.Delete(c.Request.Context(), currentUserID(c), id); err != nil {
		if errors.Is(err, service.ErrDestinationNotFound) {
			writeError(c, http.StatusNotFound, "destinasi tidak ditemukan")
			return
		}
		writeError(c, http.StatusInternalServerError, "gagal menghapus destinasi")
		return
	}
	writeSuccess(c, http.StatusOK, "destinasi berhasil dihapus", nil)
}
