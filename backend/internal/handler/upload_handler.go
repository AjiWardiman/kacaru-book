package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// MaxUploadSize membatasi ukuran file yang diunggah maksimal 5MB.
	MaxUploadSize = 5 * 1024 * 1024 // 5 MB
	// UploadDir adalah direktori tempat menyimpan file yang diunggah.
	UploadDir = "./uploads"
)

// UploadHandler menangani pengunggahan file/gambar.
type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	// Memastikan folder uploads dan subfolder utama sudah dibuat
	_ = os.MkdirAll(UploadDir, 0755)
	_ = os.MkdirAll(filepath.Join(UploadDir, "destinations"), 0755)
	return &UploadHandler{}
}

// Upload godoc
// @Summary Mengunggah file gambar umum / profil avatar (Wajib Login - Semua Role)
// @Description Mengunggah berkas gambar (jpg, jpeg, png, webp, gif) maks 5MB untuk avatar/dokumen pendukung
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File gambar yang akan diunggah"
// @Success 201 {object} dto.APIResponse "File berhasil diunggah"
// @Failure 400 {object} dto.APIResponse "File tidak valid atau melebihi batas ukuran"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 500 {object} dto.APIResponse "Gagal menyimpan file"
// @Router /upload [post]
func (h *UploadHandler) Upload(c *gin.Context) {
	h.saveFile(c, "")
}

// UploadDestination godoc
// @Summary Mengunggah foto destinasi wisata (Wajib Login - Khusus Admin)
// @Description Mengunggah berkas gambar destinasi wisata (jpg, jpeg, png, webp, gif) maks 5MB ke folder /uploads/destinations
// @Tags Admin - Upload
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File gambar destinasi wisata yang akan diunggah"
// @Success 201 {object} dto.APIResponse "File berhasil diunggah"
// @Failure 400 {object} dto.APIResponse "File tidak valid atau melebihi batas ukuran"
// @Failure 401 {object} dto.APIResponse "Unauthorized"
// @Failure 403 {object} dto.APIResponse "Forbidden (Khusus Admin)"
// @Failure 500 {object} dto.APIResponse "Gagal menyimpan file"
// @Router /admin/destinations/upload [post]
func (h *UploadHandler) UploadDestination(c *gin.Context) {
	h.saveFile(c, "destinations")
}

// saveFile adalah fungsi pembantu internal untuk memproses dan menyimpan file ke direktori tujuan
func (h *UploadHandler) saveFile(c *gin.Context, subDir string) {
	// Batasi ukuran request body
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	// Ambil file dari form input (dukung key 'file' atau 'image')
	file, err := c.FormFile("file")
	if err != nil {
		file, err = c.FormFile("image")
	}

	if err != nil {
		writeError(c, http.StatusBadRequest, "file gambar tidak ditemukan (gunakan field 'file' atau 'image')")
		return
	}

	// Validasi Ukuran File
	if file.Size > MaxUploadSize {
		writeError(c, http.StatusBadRequest, "ukuran file terlalu besar, maksimal 5MB")
		return
	}

	// Validasi Ekstensi File
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
	}

	if !allowedExts[ext] {
		writeError(c, http.StatusBadRequest, "format file tidak didukung (hanya .jpg, .jpeg, .png, .webp, .gif)")
		return
	}

	targetDir := UploadDir
	if subDir != "" {
		targetDir = filepath.Join(UploadDir, subDir)
		_ = os.MkdirAll(targetDir, 0755)
	}

	// Buat nama file unik untuk menghindari tabrakan nama file
	cleanName := filepath.Base(file.Filename[:len(file.Filename)-len(ext)])
	cleanName = strings.ReplaceAll(cleanName, " ", "_")
	uniqueFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), cleanName, ext)
	dst := filepath.Join(targetDir, uniqueFilename)

	// Simpan file ke direktori tujuan
	if err := c.SaveUploadedFile(file, dst); err != nil {
		writeError(c, http.StatusInternalServerError, "gagal menyimpan file ke server")
		return
	}

	// Path relatif URL untuk diakses oleh client/frontend
	var fileURL string
	if subDir != "" {
		fileURL = fmt.Sprintf("/uploads/%s/%s", subDir, uniqueFilename)
	} else {
		fileURL = fmt.Sprintf("/uploads/%s", uniqueFilename)
	}

	writeSuccess(c, http.StatusCreated, "file berhasil diunggah", gin.H{
		"url":      fileURL,
		"filename": uniqueFilename,
		"size":     file.Size,
	})
}

