package handler

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/middleware"
	"tourism-booking/internal/service"
)

// PaymentHandler menangani request HTTP seputar pembayaran, termasuk
// menerima webhook dari Xendit.
type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(s service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: s}
}

// CreatePayment menangani POST /bookings/:id/pay.
// Endpoint ini butuh login — dipanggil user setelah booking dibuat,
// untuk mendapatkan link pembayaran dari Xendit.
// CreatePayment godoc
// @Summary Membuat link pembayaran Xendit untuk booking
// @Description Membuat invoice pembayaran Xendit baru untuk booking yang berstatus pending (Wajib Login)
// @Tags Payments
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Booking"
// @Success 201 {object} dto.APIResponse{data=dto.PaymentResponse} "Link pembayaran berhasil dibuat"
// @Failure 400 {object} dto.APIResponse "ID booking tidak valid"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 403 {object} dto.APIResponse "Booking bukan milik Anda"
// @Failure 404 {object} dto.APIResponse "Booking tidak ditemukan"
// @Failure 409 {object} dto.APIResponse "Status booking tidak valid / pembayaran sudah ada"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /bookings/{id}/pay [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	bookingID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "id booking tidak valid")
		return
	}

	userIDVal, _ := c.Get(middleware.ContextUserID)
	userID, _ := userIDVal.(int64)

	result, err := h.paymentService.CreatePayment(c.Request.Context(), userID, bookingID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookingNotFound):
			writeError(c, http.StatusNotFound, "booking tidak ditemukan")
		case errors.Is(err, service.ErrBookingNotOwned):
			writeError(c, http.StatusForbidden, "booking ini bukan milik Anda")
		case errors.Is(err, service.ErrBookingNotPending):
			writeError(c, http.StatusConflict, "booking sudah tidak berstatus menunggu pembayaran")
		case errors.Is(err, service.ErrPaymentAlreadyExists):
			writeError(c, http.StatusConflict, "pembayaran untuk booking ini sudah pernah dibuat")
		default:
			writeError(c, http.StatusInternalServerError, "gagal membuat pembayaran")
		}
		return
	}

	writeSuccess(c, http.StatusCreated, "link pembayaran berhasil dibuat", result)
}

// XenditWebhook godoc
// @Summary Receive Xendit payment webhook callback
// @Description Callback endpoint yang dipanggil oleh server Xendit saat ada perubahan status invoice
// @Tags Payments
// @Accept json
// @Produce json
// @Param x-callback-token header string true "Xendit Callback Token"
// @Param payload body object true "Xendit Callback Payload JSON"
// @Success 200 {object} dto.APIResponse "Webhook berhasil diproses"
// @Failure 400 {object} dto.APIResponse "Body webhook tidak valid"
// @Failure 401 {object} dto.APIResponse "Callback token tidak valid"
// @Failure 404 {object} dto.APIResponse "Data pembayaran tidak ditemukan"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /webhooks/xendit [post]
func (h *PaymentHandler) XenditWebhook(c *gin.Context) {
	callbackToken := c.GetHeader("x-callback-token")

	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		writeError(c, http.StatusBadRequest, "gagal membaca body webhook")
		return
	}

	err = h.paymentService.HandleXenditWebhook(c.Request.Context(), rawBody, callbackToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCallbackToken):
			// 401: token tidak cocok, kemungkinan ini bukan request asli dari Xendit
			writeError(c, http.StatusUnauthorized, "callback token tidak valid")
		case errors.Is(err, service.ErrPaymentNotFound):
			writeError(c, http.StatusNotFound, "data pembayaran tidak ditemukan")
		default:
			writeError(c, http.StatusInternalServerError, "gagal memproses webhook")
		}
		return
	}

	// Xendit mengharapkan response 200 OK sebagai tanda webhook
	// berhasil diterima dan diproses — kalau tidak, Xendit akan
	// mencoba mengirim ulang webhook yang sama berkali-kali.
	writeSuccess(c, http.StatusOK, "webhook berhasil diproses", nil)
}
