package service

import (
	"context"
	"errors"
	"time"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
)

// Daftar error yang bisa terjadi di proses pembuatan/pengambilan slot.
var (
	ErrSlotNotFound  = errors.New("slot tidak ditemukan")
	ErrInvalidDate   = errors.New("format tanggal tidak valid, gunakan YYYY-MM-DD")
	ErrSlotDuplicate = repository.ErrDuplicateSlot // "alias" ke error yang sama dari repository
	ErrInvalidQuota  = errors.New("total_quota harus lebih besar dari nol")
)

// SlotService adalah kontrak untuk logic bisnis seputar slot (jadwal +
// kuota kunjungan) pada sebuah destinasi.
type SlotService interface {
	Create(ctx context.Context, userID, destinationID int64, req dto.CreateSlotRequest) (*dto.SlotResponse, error)
	ListByDestination(ctx context.Context, destinationID int64) ([]dto.SlotResponse, error)
}

type slotService struct {
	repo         repository.SlotRepository
	auditService AuditService
}

func NewSlotService(repo repository.SlotRepository, auditService AuditService) SlotService {
	return &slotService{repo: repo, auditService: auditService}
}

// toSlotResponse mengubah model.Slot menjadi dto.SlotResponse.
// RemainingQuota (sisa kuota) dihitung di sini supaya client tidak perlu
// menghitung sendiri (TotalQuota - UsedQuota).
func toSlotResponse(s *model.Slot) dto.SlotResponse {
	return dto.SlotResponse{
		ID:             s.ID,
		DestinationID:  s.DestinationID,
		Date:           s.Date.Format("2006-01-02"), // format tanggal Go yang aneh ini artinya "YYYY-MM-DD"
		TotalQuota:     s.TotalQuota,
		UsedQuota:      s.UsedQuota,
		RemainingQuota: s.TotalQuota - s.UsedQuota,
		IsAvailable:    s.IsAvailable,
	}
}

// Create membuat slot/jadwal baru untuk sebuah destinasi.
// Ada 2 validasi sebelum data disimpan:
//  1. total_quota harus lebih besar dari 0 (tidak masuk akal ada slot
//     dengan kuota 0 atau minus).
//  2. Format tanggal harus benar ("YYYY-MM-DD"), kalau tidak, kita
//     tolak lebih awal daripada membiarkan error aneh dari database.
func (s *slotService) Create(ctx context.Context, userID, destinationID int64, req dto.CreateSlotRequest) (*dto.SlotResponse, error) {
	if req.TotalQuota <= 0 {
		return nil, ErrInvalidQuota
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, ErrInvalidDate
	}

	slot := &model.Slot{
		DestinationID: destinationID,
		Date:          date,
		TotalQuota:    req.TotalQuota,
		IsAvailable:   true,
	}

	created, err := s.repo.Create(ctx, slot)
	if err != nil {
		// err di sini bisa saja ErrDuplicateSlot (kalau tanggal sudah
		// dipakai slot lain di destinasi yang sama) — diteruskan apa
		// adanya ke handler untuk diterjemahkan jadi response HTTP yang tepat.
		return nil, err
	}

	resp := toSlotResponse(created)

	s.auditService.Log(ctx, &userID, ActionCreate, "slots", created.ID, nil, resp)

	return &resp, nil
}

// ListByDestination mengambil semua slot/jadwal milik satu destinasi.
// Jika belum ada slot untuk 14 hari ke depan, otomatis disiapkan slot default.
func (s *slotService) ListByDestination(ctx context.Context, destinationID int64) ([]dto.SlotResponse, error) {
	slots, err := s.repo.ListByDestination(ctx, destinationID)
	if err != nil {
		return nil, err
	}

	dateMap := make(map[string]bool)
	for _, sl := range slots {
		dateMap[sl.Date.Format("2006-01-02")] = true
	}

	now := time.Now()
	hasNew := false
	for i := 0; i < 14; i++ {
		targetDateStr := now.AddDate(0, 0, i).Format("2006-01-02")
		if !dateMap[targetDateStr] {
			tDate, _ := time.Parse("2006-01-02", targetDateStr)
			newSlot := &model.Slot{
				DestinationID: destinationID,
				Date:          tDate,
				TotalQuota:    50,
				UsedQuota:     0,
				IsAvailable:   true,
			}
			if _, createErr := s.repo.Create(ctx, newSlot); createErr == nil {
				hasNew = true
			}
		}
	}

	if hasNew {
		slots, err = s.repo.ListByDestination(ctx, destinationID)
		if err != nil {
			return nil, err
		}
	}

	result := make([]dto.SlotResponse, 0, len(slots))
	for _, sl := range slots {
		result = append(result, toSlotResponse(&sl))
	}
	return result, nil
}
