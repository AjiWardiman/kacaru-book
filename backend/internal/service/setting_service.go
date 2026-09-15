package service

import (
	"context"
	"errors"
	"strconv"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
)

// ErrSettingNotFound dikembalikan kalau key setting yang diminta tidak
// ada di database (misal admin salah ketik nama key-nya).
var ErrSettingNotFound = errors.New("setting not found")

// SettingService adalah kontrak logic bisnis untuk fitur pengaturan
// aplikasi (settings). Selain dipakai lewat endpoint admin, service ini
// juga dipakai secara INTERNAL oleh service lain (contoh: OTPService)
// yang butuh membaca nilai pengaturan seperti "otp_expiry_minutes",
// supaya nilai itu bisa diubah admin kapan saja tanpa perlu ubah kode.
type SettingService interface {
	List(ctx context.Context) ([]dto.SettingResponse, error)
	Get(ctx context.Context, key string) (*dto.SettingResponse, error)
	Update(ctx context.Context, userID int64, key, value string) (*dto.SettingResponse, error)
	// GetIntOrDefault mengambil nilai sebuah setting dan mengubahnya
	// jadi angka (int). Kalau key-nya tidak ditemukan ATAU nilainya
	// ternyata bukan angka yang valid, fungsi ini TIDAK mengembalikan
	// error — cukup kembalikan nilai "fallback" yang disediakan
	// pemanggil. Ini sengaja dibuat "aman gagal" (fail-safe) supaya
	// fitur lain yang bergantung pada setting (misal OTP, auto-expire
	// booking) tetap bisa jalan pakai nilai default, walau misalnya
	// data setting di database belum sempat di-migrasi atau isinya
	// rusak karena kesalahan input admin.
	GetIntOrDefault(ctx context.Context, key string, fallback int) int
}

type settingService struct {
	repo         repository.SettingRepository
	auditService AuditService
}

func NewSettingService(repo repository.SettingRepository, auditService AuditService) SettingService {
	return &settingService{repo: repo, auditService: auditService}
}

func toSettingResponse(s *model.Setting) dto.SettingResponse {
	return dto.SettingResponse{
		ID:          s.ID,
		Key:         s.Key,
		Value:       s.Value,
		Description: s.Description,
		UpdatedAt:   s.UpdatedAt,
	}
}

// List mengambil semua pengaturan aplikasi yang tersedia, supaya admin
// bisa melihat apa saja yang bisa diubah.
func (s *settingService) List(ctx context.Context) ([]dto.SettingResponse, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.SettingResponse, 0, len(list))
	for _, item := range list {
		result = append(result, toSettingResponse(&item))
	}
	return result, nil
}

// Get mengambil detail satu pengaturan berdasarkan key-nya.
func (s *settingService) Get(ctx context.Context, key string) (*dto.SettingResponse, error) {
	setting, err := s.repo.FindByKey(ctx, key)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrSettingNotFound
	}
	if err != nil {
		return nil, err
	}
	resp := toSettingResponse(setting)
	return &resp, nil
}

// Update mengubah nilai sebuah pengaturan. userID dicatat ke audit
// trail supaya ada jejak SIAPA admin yang mengubah pengaturan penting
// seperti ini (perubahan setting bisa berdampak ke seluruh aplikasi,
// jadi penting untuk dilacak).
func (s *settingService) Update(ctx context.Context, userID int64, key, value string) (*dto.SettingResponse, error) {
	// Ambil dulu nilai lama, supaya audit trail punya catatan
	// "sebelum berubah jadi apa, nilainya apa".
	existing, err := s.repo.FindByKey(ctx, key)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrSettingNotFound
	}
	if err != nil {
		return nil, err
	}
	oldSnapshot := toSettingResponse(existing)

	updated, err := s.repo.UpdateValue(ctx, key, value)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrSettingNotFound
	}
	if err != nil {
		return nil, err
	}

	resp := toSettingResponse(updated)
	s.auditService.Log(ctx, &userID, ActionUpdate, "settings", updated.ID, oldSnapshot, resp)

	return &resp, nil
}

// GetIntOrDefault lihat penjelasan lengkap di interface SettingService di atas.
func (s *settingService) GetIntOrDefault(ctx context.Context, key string, fallback int) int {
	setting, err := s.repo.FindByKey(ctx, key)
	if err != nil {
		return fallback
	}

	value, err := strconv.Atoi(setting.Value)
	if err != nil {
		return fallback
	}

	return value
}
