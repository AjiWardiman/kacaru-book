package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/service"
)

// AuthHandler menangani request HTTP yang berhubungan dengan autentikasi
// (register dan login).
type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register menangani POST /auth/register.
// Alurnya: baca body JSON → bersihkan input (trim spasi, lowercase
// email) → validasi dasar → serahkan ke service → kirim response.
// Register godoc
// @Summary Register user baru
// @Description Pendaftaran akun pengguna baru (role default 'user')
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Data Registrasi"
// @Success 201 {object} dto.APIResponse{data=dto.AuthResponse} "Pendaftaran berhasil"
// @Failure 400 {object} dto.APIResponse "Data tidak valid"
// @Failure 409 {object} dto.APIResponse "Email atau nomor HP sudah terdaftar"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	// ShouldBindJSON otomatis membaca body request dan mengubahnya
	// menjadi struct RegisterRequest. Kalau JSON-nya rusak/tidak valid,
	// akan mengembalikan error di sini.
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	// Bersihkan input dari spasi berlebih, dan samakan email jadi huruf
	// kecil semua (supaya "Budi@Email.com" dan "budi@email.com"
	// dianggap sama, tidak bisa dipakai untuk daftar dua kali).
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)

	// Validasi dasar di level handler (validasi yang lebih ke bisnis,
	// seperti "email sudah dipakai", dilakukan di service).
	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeError(c, http.StatusBadRequest, "nama, email, dan password wajib diisi")
		return
	}
	if len(req.Password) < 8 {
		writeError(c, http.StatusBadRequest, "password minimal 8 karakter")
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		// errors.Is dipakai untuk mengecek JENIS error secara pasti,
		// supaya kita bisa memberi status HTTP dan pesan yang sesuai
		// untuk masing-masing kasus.
		switch {
		case errors.Is(err, service.ErrEmailTaken):
			writeError(c, http.StatusConflict, "email sudah terdaftar") // 409: data sudah ada
		case errors.Is(err, service.ErrPhoneTaken):
			writeError(c, http.StatusConflict, "nomor HP sudah terdaftar")
		default:
			writeError(c, http.StatusInternalServerError, "gagal mendaftarkan user") // 500: error tak terduga
		}
		return
	}

	writeSuccess(c, http.StatusCreated, "pendaftaran berhasil", result) // 201: data berhasil dibuat
}

// Login godoc
// @Summary Login user / admin
// @Description Authentikasi menggunakan email/nomor HP dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Data Login"
// @Success 200 {object} dto.APIResponse{data=dto.AuthResponse} "Login berhasil"
// @Failure 400 {object} dto.APIResponse "Data tidak valid"
// @Failure 401 {object} dto.APIResponse "Kredensial tidak valid"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	if req.Identifier == "" && req.Email != "" {
		req.Identifier = req.Email
	}
	req.Identifier = strings.TrimSpace(req.Identifier)
	if req.Identifier == "" || req.Password == "" {
		writeError(c, http.StatusBadRequest, "email/nomor HP dan password wajib diisi")
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCreds) {
			writeError(c, http.StatusUnauthorized, "email/nomor HP atau password salah") // 401: identitas/password salah
			return
		}
		writeError(c, http.StatusInternalServerError, "gagal login")
		return
	}

	writeSuccess(c, http.StatusOK, "login berhasil", result) // 200: berhasil
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Email) == "" {
		writeError(c, http.StatusBadRequest, "email wajib diisi")
		return
	}
	if err := h.authService.SendPreRegisterOTP(c.Request.Context(), strings.TrimSpace(req.Email)); err != nil {
		writeError(c, http.StatusInternalServerError, "gagal mengirim OTP")
		return
	}
	writeSuccess(c, http.StatusOK, "kode OTP berhasil dikirim ke Mailtrap", nil)
}
