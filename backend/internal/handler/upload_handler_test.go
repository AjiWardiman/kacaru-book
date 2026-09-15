package handler

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUploadHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUploadHandler()
	defer os.RemoveAll(UploadDir)

	r := gin.New()
	r.POST("/upload", h.Upload)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test_image.png")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	_, _ = io.WriteString(part, "fake-image-binary-data")
	_ = writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestUploadHandler_InvalidExtension(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUploadHandler()
	defer os.RemoveAll(UploadDir)

	r := gin.New()
	r.POST("/upload", h.Upload)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test_script.exe")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	_, _ = io.WriteString(part, "malicious-executable-content")
	_ = writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestUploadHandler_UploadDestination_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewUploadHandler()
	defer os.RemoveAll(UploadDir)

	r := gin.New()
	r.POST("/admin/destinations/upload", h.UploadDestination)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "pantai.jpg")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	_, _ = io.WriteString(part, "fake-image-binary-data")
	_ = writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/admin/destinations/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

