package service

import (
	"context"
	"errors"
	"time"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
)

// ErrDestinationNotFound dikembalikan kalau destinasi dengan ID yang
// diminta tidak ada di database.
var ErrDestinationNotFound = errors.New("destination not found")

// DestinationService adalah kontrak untuk semua logic bisnis seputar
// destinasi wisata: membuat, membaca, mengubah, dan menghapus.
//
// PERHATIAN: Update dan Delete sekarang menerima parameter userID —
// ini dipakai untuk mencatat SIAPA yang melakukan perubahan ke audit
// trail (requirement wajib dari trainer).
type DestinationService interface {
	Create(ctx context.Context, createdBy int64, req dto.CreateDestinationRequest) (*dto.DestinationResponse, error)
	Get(ctx context.Context, id int64) (*dto.DestinationResponse, error)
	List(ctx context.Context, onlyActive bool) ([]dto.DestinationResponse, error)
	Update(ctx context.Context, userID, id int64, req dto.UpdateDestinationRequest) (*dto.DestinationResponse, error)
	Delete(ctx context.Context, userID, id int64) error
}

type destinationService struct {
	repo         repository.DestinationRepository
	slotRepo     repository.SlotRepository
	auditService AuditService
}

func NewDestinationService(repo repository.DestinationRepository, slotRepo repository.SlotRepository, auditService AuditService) DestinationService {
	return &destinationService{repo: repo, slotRepo: slotRepo, auditService: auditService}
}

// toDestinationResponse mengubah struct model.Destination (representasi
// tabel database) menjadi dto.DestinationResponse (bentuk aman untuk
// dikirim ke client). Field pointer (*string) di model diubah menjadi
// string biasa di response, dengan nilai kosong "" kalau memang nil.
func toDestinationResponse(d *model.Destination) dto.DestinationResponse {
	resp := dto.DestinationResponse{
		ID:        d.ID,
		Name:      d.Name,
		Location:  d.Location,
		Price:     d.Price,
		IsActive:  d.IsActive,
		CreatedAt: d.CreatedAt,
	}
	if d.Description != nil {
		resp.Description = *d.Description
	}
	if d.Category != nil {
		resp.Category = *d.Category
	}
	if d.ImageURL != nil {
		resp.ImageURL = *d.ImageURL
	}
	return resp
}

// Create membuat destinasi wisata baru. createdBy adalah ID admin yang
// sedang login (diambil dari token JWT di lapisan handler/middleware),
// dicatat supaya kita tahu siapa yang bertanggung jawab membuat data ini.
func (s *destinationService) Create(ctx context.Context, createdBy int64, req dto.CreateDestinationRequest) (*dto.DestinationResponse, error) {
	d := &model.Destination{
		CreatedBy:   createdBy,
		Name:        req.Name,
		Location:    req.Location,
		Description: strPtr(req.Description),
		Category:    strPtr(req.Category),
		Price:       req.Price,
		ImageURL:    strPtr(req.ImageURL),
		IsActive:    true,
	}
	created, err := s.repo.Create(ctx, d)
	if err != nil {
		return nil, err
	}

	// Auto-create default slots for the new destination for the next 14 days so visitors can immediately book
	if s.slotRepo != nil {
		now := time.Now()
		for i := 0; i < 14; i++ {
			slotDate := now.AddDate(0, 0, i)
			_, _ = s.slotRepo.Create(ctx, &model.Slot{
				DestinationID: created.ID,
				Date:          time.Date(slotDate.Year(), slotDate.Month(), slotDate.Day(), 0, 0, 0, 0, time.UTC),
				TotalQuota:    100,
				UsedQuota:     0,
				IsAvailable:   true,
			})
		}
	}

	resp := toDestinationResponse(created)

	// Catat ke audit trail: aksi CREATE, tidak ada "data lama" (nil)
	// karena ini data baru, "data baru" adalah destinasi yang baru dibuat.
	s.auditService.Log(ctx, &createdBy, ActionCreate, "destinations", created.ID, nil, resp)

	return &resp, nil
}

// Get mengambil detail satu destinasi berdasarkan ID.
// Dipakai baik oleh endpoint publik (lihat detail destinasi) maupun admin.
func (s *destinationService) Get(ctx context.Context, id int64) (*dto.DestinationResponse, error) {
	d, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrDestinationNotFound
	}
	if err != nil {
		return nil, err
	}
	resp := toDestinationResponse(d)
	return &resp, nil
}

// List mengambil daftar destinasi. onlyActive=true dipakai untuk halaman
// publik (hanya tampilkan destinasi yang aktif), onlyActive=false dipakai
// admin untuk melihat semua termasuk yang sudah dinonaktifkan.
func (s *destinationService) List(ctx context.Context, onlyActive bool) ([]dto.DestinationResponse, error) {
	list, err := s.repo.List(ctx, onlyActive)
	if err != nil {
		return nil, err
	}
	result := make([]dto.DestinationResponse, 0, len(list))
	for _, d := range list {
		result = append(result, toDestinationResponse(&d))
	}
	return result, nil
}

// Update mengubah data destinasi yang sudah ada.
//
// PENTING: ini "partial update" — setiap field di request yang dikirim
// KOSONG (string kosong atau 0) TIDAK akan mengubah nilai lama di
// database. Jadi admin bisa mengirim cuma field yang mau diubah saja,
// tanpa perlu mengirim ulang semua data destinasi.
func (s *destinationService) Update(ctx context.Context, userID, id int64, req dto.UpdateDestinationRequest) (*dto.DestinationResponse, error) {
	// Ambil dulu data yang sudah ada, supaya field yang tidak diubah
	// tetap mempertahankan nilai lamanya, DAN supaya kita punya "data
	// sebelum perubahan" untuk dicatat di audit trail.
	existing, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrDestinationNotFound
	}
	if err != nil {
		return nil, err
	}
	oldSnapshot := toDestinationResponse(existing)

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Location != "" {
		existing.Location = req.Location
	}
	if req.Description != "" {
		existing.Description = &req.Description
	}
	if req.Category != "" {
		existing.Category = &req.Category
	}
	if req.Price > 0 {
		existing.Price = req.Price
	}
	if req.ImageURL != "" {
		existing.ImageURL = &req.ImageURL
	}
	// IsActive dicek beda caranya (pakai pointer) karena "false" adalah
	// nilai yang valid dan berbeda maknanya dari "tidak dikirim sama sekali".
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrDestinationNotFound
		}
		return nil, err
	}

	resp := toDestinationResponse(existing)

	// Catat ke audit trail: simpan snapshot data SEBELUM dan SESUDAH
	// perubahan, supaya kalau ada pertanyaan "kenapa harga destinasi
	// ini berubah?", ada jejak lengkapnya.
	s.auditService.Log(ctx, &userID, ActionUpdate, "destinations", id, oldSnapshot, resp)

	return &resp, nil
}

// Delete menghapus destinasi secara permanen dari database.
func (s *destinationService) Delete(ctx context.Context, userID, id int64) error {
	// Ambil dulu datanya SEBELUM dihapus, supaya kita masih punya
	// catatan "data apa yang tadinya ada di sini" untuk audit trail
	// (setelah dihapus, datanya sudah tidak bisa diambil lagi dari database).
	existing, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrDestinationNotFound
	}
	if err != nil {
		return err
	}
	oldSnapshot := toDestinationResponse(existing)

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrDestinationNotFound
		}
		return err
	}

	// Catat ke audit trail: aksi DELETE, tidak ada "data baru" (nil)
	// karena datanya sudah dihapus.
	s.auditService.Log(ctx, &userID, ActionDelete, "destinations", id, oldSnapshot, nil)

	return nil
}

// strPtr adalah fungsi bantuan kecil: ubah string menjadi pointer,
// tapi kalau stringnya kosong, kembalikan nil (bukan pointer ke string
// kosong). Ini supaya kolom di database yang boleh NULL benar-benar
// diisi NULL, bukan string kosong "".
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
