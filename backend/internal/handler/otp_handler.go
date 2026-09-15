package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/service"
)

// OTPHandler menangani request HTTP seputar verifikasi HP lewat kode
// OTP. Semua endpoint di sini WAJIB LOGIN (butuh token JWT) — dipasang
// lewat middleware.Auth di router/router.go, karena OTP di sini
// khusus untuk memverifikasi HP milik AKUN YANG SEDANG LOGIN, bukan
// akun sembarangan.
type OTPHandler struct {
	otpService service.OTPService
}

func NewOTPHandler(otpService service.OTPService) *OTPHandler {
	return &OTPHandler{otpService: otpService}
}

// RequestOTP menangani POST /otp/request — meminta kode OTP baru untuk
// verifikasi nomor HP milik user yang sedang login. Kode dikirim lewat
// email (simulasi SMS, sesuai kesepakatan tugas akhir).
// RequestOTP godoc
// @Summary Meminta kode OTP untuk verifikasi nomor HP
// @Description Mengirim kode OTP baru ke email user untuk verifikasi HP (Wajib Login)
// @Tags OTP
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.RequestOTPResponse} "Kode OTP berhasil dikirim"
// @Failure 400 {object} dto.APIResponse "Nomor HP belum diisi"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 409 {object} dto.APIResponse "Nomor HP sudah terverifikasi"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /otp/request [post]
func (h *OTPHandler) RequestOTP(c *gin.Context) {
	userID := currentUserID(c)

	result, err := h.otpService.RequestPhoneVerificationOTP(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPhoneNotSet):
			writeError(c, http.StatusBadRequest, "akun kamu belum mengisi nomor HP, lengkapi dulu sebelum verifikasi")
		case errors.Is(err, service.ErrPhoneAlreadyVerify):
			writeError(c, http.StatusConflict, "nomor HP kamu sudah terverifikasi")
		default:
			writeError(c, http.StatusInternalServerError, "gagal mengirim kode OTP")
		}
		return
	}

	writeSuccess(c, http.StatusOK, result.Message, result)
}

// VerifyOTP godoc
// @Summary Verifikasi kode OTP nomor HP
// @Description Mencocokkan kode OTP yang dikirim user dengan kode aktif di database (Wajib Login)
// @Tags OTP
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.VerifyOTPRequest true "Kode OTP"
// @Success 200 {object} dto.APIResponse "Nomor HP berhasil diverifikasi"
// @Failure 400 {object} dto.APIResponse "Kode OTP salah / kadaluarsa / tidak aktif"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /otp/verify [post]
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	userID := currentUserID(c)

	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}
	req.Code = strings.TrimSpace(req.Code)
	if req.Code == "" {
		writeError(c, http.StatusBadRequest, "kode OTP wajib diisi")
		return
	}

	err := h.otpService.VerifyPhoneOTP(c.Request.Context(), userID, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOTPNotRequested):
			writeError(c, http.StatusBadRequest, "belum ada kode OTP aktif, silakan minta kode baru dulu")
		case errors.Is(err, service.ErrOTPExpired):
			writeError(c, http.StatusBadRequest, "kode OTP sudah kadaluarsa, silakan minta kode baru")
		case errors.Is(err, service.ErrOTPInvalidCode):
			writeError(c, http.StatusBadRequest, "kode OTP salah")
		default:
			writeError(c, http.StatusInternalServerError, "gagal memverifikasi kode OTP")
		}
		return
	}

	writeSuccess(c, http.StatusOK, "nomor HP berhasil diverifikasi", nil)
}
