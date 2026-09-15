package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/repository"
	"tourism-booking/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUserRequest merepresentasikan request body untuk admin membuat/mengedit user.
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// List menangani GET /admin/users.
func (h *UserHandler) List(c *gin.Context) {
	search := c.Query("search")
	roleStr := c.Query("role")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "500"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, total, err := h.userService.GetUsers(c.Request.Context(), search, roleStr, limit, offset)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal mengambil daftar pengguna")
		return
	}

	writeSuccess(c, http.StatusOK, "daftar pengguna berhasil diambil", gin.H{
		"users": users,
		"total": total,
	})
}

// Get menangani GET /admin/users/:id.
func (h *UserHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID pengguna tidak valid")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(c, http.StatusNotFound, "pengguna tidak ditemukan")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal mengambil data pengguna")
		return
	}

	writeSuccess(c, http.StatusOK, "detail pengguna berhasil diambil", user)
}

// Create menangani POST /admin/users.
func (h *UserHandler) Create(c *gin.Context) {
	adminID := currentUserID(c)

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	if req.Name == "" || req.Email == "" {
		writeError(c, http.StatusBadRequest, "nama dan email wajib diisi")
		return
	}

	dtoReq := dto.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	res, err := h.userService.CreateUser(c.Request.Context(), adminID, dtoReq, req.Role)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			writeError(c, http.StatusConflict, "email sudah terdaftar")
			return
		}
		if errors.Is(err, service.ErrPhoneTaken) {
			writeError(c, http.StatusConflict, "nomor HP sudah terdaftar")
			return
		}
		writeError(c, http.StatusInternalServerError, "gagal menambahkan pengguna")
		return
	}

	writeSuccess(c, http.StatusCreated, "pengguna baru berhasil ditambahkan", res)
}

// Update menangani PUT /admin/users/:id.
func (h *UserHandler) Update(c *gin.Context) {
	adminID := currentUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID pengguna tidak valid")
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	dtoReq := dto.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	res, err := h.userService.UpdateUser(c.Request.Context(), adminID, id, dtoReq, req.Role)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(c, http.StatusNotFound, "pengguna tidak ditemukan")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal memperbarui data pengguna")
		return
	}

	writeSuccess(c, http.StatusOK, "data pengguna berhasil diperbarui", res)
}

// Delete menangani DELETE /admin/users/:id.
func (h *UserHandler) Delete(c *gin.Context) {
	adminID := currentUserID(c)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID pengguna tidak valid")
		return
	}

	if id == adminID {
		writeError(c, http.StatusBadRequest, "tidak dapat menghapus akun admin yang sedang login")
		return
	}

	err = h.userService.DeleteUser(c.Request.Context(), adminID, id)
	if errors.Is(err, repository.ErrNotFound) {
		writeError(c, http.StatusNotFound, "pengguna tidak ditemukan")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal menghapus pengguna")
		return
	}

	writeSuccess(c, http.StatusOK, "pengguna berhasil dihapus", nil)
}

// UpdateSelfProfile menangani PUT /users/profile untuk pengguna yang sedang login.
func (h *UserHandler) UpdateSelfProfile(c *gin.Context) {
	userID := currentUserID(c)
	if userID == 0 {
		writeError(c, http.StatusUnauthorized, "belum login")
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	dtoReq := dto.RegisterRequest{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	res, err := h.userService.UpdateUser(c.Request.Context(), userID, userID, dtoReq, "")
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, "profil berhasil diperbarui", res)
}

// UpdateSelfPassword menangani PUT /users/password untuk mengubah kata sandi pengguna yang sedang login.
func (h *UserHandler) UpdateSelfPassword(c *gin.Context) {
	userID := currentUserID(c)
	if userID == 0 {
		writeError(c, http.StatusUnauthorized, "belum login")
		return
	}

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.NewPassword == "" {
		writeError(c, http.StatusBadRequest, "password baru wajib diisi (minimal 6 karakter)")
		return
	}

	dtoReq := dto.RegisterRequest{
		Password: req.NewPassword,
	}

	res, err := h.userService.UpdateUser(c.Request.Context(), userID, userID, dtoReq, "")
	if err != nil {
		writeError(c, http.StatusInternalServerError, "gagal memperbarui password")
		return
	}

	writeSuccess(c, http.StatusOK, "kata sandi berhasil diperbarui", res)
}
