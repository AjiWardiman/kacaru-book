package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/middleware"
	"tourism-booking/internal/service"
)

// BookingHandler menangani request HTTP seputar booking/pemesanan.
type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(s service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: s}
}

// Create menangani POST /bookings.
//
// Endpoint ini butuh login (dijaga oleh middleware Auth di router),
// TAPI boleh diakses SIAPA SAJA yang sudah login — baik role "user"
// maupun "admin". Ini beda dengan endpoint /admin/... yang cuma boleh
// diakses admin, karena justru USER BIASA-lah yang seharusnya melakukan booking.
// Create godoc
// @Summary Membuat pemesanan/booking destinasi wisata
// @Description Membuat booking baru untuk slot destinasi tertentu (Wajib Login)
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateBookingRequest true "Data Pemesanan & Pengunjung"
// @Success 201 {object} dto.APIResponse{data=dto.BookingResponse} "Booking berhasil dibuat"
// @Failure 400 {object} dto.APIResponse "Data tidak valid"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 404 {object} dto.APIResponse "Slot tidak ditemukan"
// @Failure 409 {object} dto.APIResponse "Slot tidak tersedia / kuota tidak cukup"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /bookings [post]
func (h *BookingHandler) Create(c *gin.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	if req.SlotID == 0 || req.TotalPeople <= 0 {
		writeError(c, http.StatusBadRequest, "slot_id dan total_people (harus lebih dari nol) wajib diisi")
		return
	}

	// Ambil user_id dari context Gin (diisi oleh middleware Auth
	// setelah memvalidasi token), supaya booking ini tercatat sebagai
	// milik user yang benar-benar sedang login — BUKAN dari field
	// user_id yang (kalau ada) dikirim client, yang bisa saja dipalsukan.
	userIDVal, _ := c.Get(middleware.ContextUserID)
	userID, _ := userIDVal.(int64)

	result, err := h.bookingService.CreateBooking(c.Request.Context(), userID, req)
	if err != nil {
		// Setiap jenis error dari service diterjemahkan menjadi status
		// HTTP dan pesan yang paling tepat untuk dipahami client.
		switch {
		case errors.Is(err, service.ErrSlotNotFound):
			writeError(c, http.StatusNotFound, "slot tidak ditemukan")
		case errors.Is(err, service.ErrSlotNotAvailable):
			writeError(c, http.StatusConflict, "slot sedang tidak tersedia")
		case errors.Is(err, service.ErrQuotaNotEnough):
			writeError(c, http.StatusConflict, "sisa kuota tidak cukup untuk booking ini") // 409: bentrok dengan kuota yang tersisa
		case errors.Is(err, service.ErrVisitorCountMismatch):
			writeError(c, http.StatusBadRequest, "jumlah data pengunjung harus sama dengan total_people")
		default:
			writeError(c, http.StatusInternalServerError, "gagal membuat booking")
		}
		return
	}

	writeSuccess(c, http.StatusCreated, "booking berhasil dibuat", result)
}

// ListMyBookings menangani GET /bookings/my-bookings
func (h *BookingHandler) ListMyBookings(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextUserID)
	userID, _ := userIDVal.(int64)

	userRoleVal, _ := c.Get(middleware.ContextRole)
	userRole, _ := userRoleVal.(string)

	var result []dto.BookingResponse
	var err error

	if userRole == "admin" {
		result, err = h.bookingService.GetAllBookings(c.Request.Context())
	} else {
		result, err = h.bookingService.GetMyBookings(c.Request.Context(), userID)
	}

	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal mengambil daftar pemesanan")
		return
	}

	writeSuccess(c, http.StatusOK, "daftar pemesanan berhasil diambil", result)
}

// ListAllBookings menangani GET /admin/bookings
func (h *BookingHandler) ListAllBookings(c *gin.Context) {
	result, err := h.bookingService.GetAllBookings(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal mengambil seluruh pemesanan platform")
		return
	}
	writeSuccess(c, http.StatusOK, "seluruh pemesanan platform berhasil diambil", result)
}
