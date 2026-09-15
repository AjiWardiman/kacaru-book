package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/service"
)

// SlotHandler menangani semua request HTTP seputar slot (jadwal + kuota
// kunjungan) pada destinasi.
type SlotHandler struct {
	slotService service.SlotService
}

func NewSlotHandler(s service.SlotService) *SlotHandler {
	return &SlotHandler{slotService: s}
}

// Create menangani POST /admin/destinations/:id/slots (khusus admin).
// ":id" di sini adalah ID destinasi, diambil dari path URL.
// Create godoc
// @Summary Membuat slot jadwal kunjungan baru (Admin-only)
// @Description Membuat slot kuota kunjungan untuk destinasi tertentu di tanggal spesifik
// @Tags Slots
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Destinasi"
// @Param request body dto.CreateSlotRequest true "Data Slot (Tanggal & Kuota)"
// @Success 201 {object} dto.APIResponse{data=dto.SlotResponse} "Slot berhasil dibuat"
// @Failure 400 {object} dto.APIResponse "Data tidak valid"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 409 {object} dto.APIResponse "Slot untuk tanggal tersebut sudah ada"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /admin/destinations/{id}/slots [post]
func (h *SlotHandler) Create(c *gin.Context) {
	destinationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "id destinasi tidak valid")
		return
	}

	var req dto.CreateSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.slotService.Create(c.Request.Context(), currentUserID(c), destinationID, req)
	if err != nil {
		// Setiap jenis error dari service diterjemahkan menjadi status
		// HTTP dan pesan yang paling sesuai untuk client.
		switch {
		case errors.Is(err, service.ErrInvalidDate):
			writeError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInvalidQuota):
			writeError(c, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrSlotDuplicate):
			writeError(c, http.StatusConflict, "slot untuk tanggal ini sudah ada") // 409: bentrok dengan data yang sudah ada
		default:
			writeError(c, http.StatusInternalServerError, "gagal membuat slot")
		}
		return
	}
	writeSuccess(c, http.StatusCreated, "slot berhasil dibuat", result)
}

// ListByDestination godoc
// @Summary Mengambil daftar slot kunjungan destinasi
// @Description Mengambil daftar slot jadwal dan sisa kuota untuk destinasi tertentu (Publik)
// @Tags Slots
// @Produce json
// @Param id path int true "ID Destinasi"
// @Success 200 {object} dto.APIResponse{data=[]dto.SlotResponse} "Daftar slot kunjungan"
// @Failure 400 {object} dto.APIResponse "ID destinasi tidak valid"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /destinations/{id}/slots [get]
func (h *SlotHandler) ListByDestination(c *gin.Context) {
	destinationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "id destinasi tidak valid")
		return
	}

	result, err := h.slotService.ListByDestination(c.Request.Context(), destinationID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal mengambil daftar slot")
		return
	}
	writeSuccess(c, http.StatusOK, "berhasil", result)
}
