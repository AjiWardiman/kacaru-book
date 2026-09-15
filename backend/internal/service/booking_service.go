package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
)

var (
	ErrSlotNotAvailable     = errors.New("slot is not available")
	ErrQuotaNotEnough       = errors.New("remaining quota is not enough for this booking")
	ErrVisitorCountMismatch = errors.New("number of visitors must match total_people")
	ErrBookingNotFound      = errors.New("booking not found")
)

const bookingExpiryMinutes = 30

type BookingService interface {
	CreateBooking(ctx context.Context, userID int64, req dto.CreateBookingRequest) (*dto.BookingResponse, error)
	GetMyBookings(ctx context.Context, userID int64) ([]dto.BookingResponse, error)
	GetAllBookings(ctx context.Context) ([]dto.BookingResponse, error)
}

type bookingService struct {
	db              *gorm.DB
	slotRepo        repository.SlotRepository
	bookingRepo     repository.BookingRepository
	destinationRepo repository.DestinationRepository
	auditService    AuditService
}

func NewBookingService(
	db *gorm.DB,
	slotRepo repository.SlotRepository,
	bookingRepo repository.BookingRepository,
	destinationRepo repository.DestinationRepository,
	auditService AuditService,
) BookingService {
	return &bookingService{
		db:              db,
		slotRepo:        slotRepo,
		bookingRepo:     bookingRepo,
		destinationRepo: destinationRepo,
		auditService:    auditService,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, userID int64, req dto.CreateBookingRequest) (*dto.BookingResponse, error) {
	if len(req.Visitors) != req.TotalPeople {
		return nil, ErrVisitorCountMismatch
	}

	var created *model.Booking

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		slot, err := s.slotRepo.LockForUpdate(ctx, tx, req.SlotID)
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSlotNotFound
		}
		if err != nil {
			return err
		}

		if !slot.IsAvailable {
			return ErrSlotNotAvailable
		}
		remainingQuota := slot.TotalQuota - slot.UsedQuota
		if req.TotalPeople > remainingQuota {
			return ErrQuotaNotEnough
		}

		destination, err := s.destinationRepo.FindByID(ctx, slot.DestinationID)
		if err != nil {
			return err
		}

		booking := &model.Booking{
			BookingCode: generateBookingCode(),
			UserID:      userID,
			SlotID:      req.SlotID,
			TotalPeople: req.TotalPeople,
			TotalPrice:  destination.Price * float64(req.TotalPeople),
			Status:      model.BookingPending,
			BookedAt:    time.Now(),
			ExpiredAt:   time.Now().Add(bookingExpiryMinutes * time.Minute),
		}
		var createErr error
		created, createErr = s.bookingRepo.CreateTx(ctx, tx, booking)
		if createErr != nil {
			return createErr
		}

		details := make([]model.BookingDetail, 0, len(req.Visitors))
		for _, v := range req.Visitors {
			details = append(details, model.BookingDetail{
				VisitorName:  v.Name,
				VisitorPhone: strPtr(v.Phone),
				VisitorEmail: strPtr(v.Email),
				Age:          intPtr(v.Age),
			})
		}
		if err := s.bookingRepo.CreateDetailsTx(ctx, tx, created.ID, details); err != nil {
			return err
		}

		if err := s.slotRepo.IncrementUsedQuota(ctx, tx, req.SlotID, req.TotalPeople); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	resp := &dto.BookingResponse{
		ID:          created.ID,
		BookingCode: created.BookingCode,
		SlotID:      created.SlotID,
		TotalPeople: created.TotalPeople,
		TotalPrice:  created.TotalPrice,
		Status:      string(created.Status),
		BookedAt:    created.BookedAt,
		ExpiredAt:   created.ExpiredAt,
	}

	s.auditService.Log(ctx, &userID, ActionCreate, "bookings", created.ID, nil, resp)

	return resp, nil
}

func (s *bookingService) GetMyBookings(ctx context.Context, userID int64) ([]dto.BookingResponse, error) {
	bookings, err := s.bookingRepo.ListByUserWithDetails(ctx, userID)
	if err != nil {
		return nil, err
	}

	resps := make([]dto.BookingResponse, 0, len(bookings))
	for _, b := range bookings {
		resps = append(resps, dto.BookingResponse{
			ID:                  b.ID,
			BookingCode:         b.BookingCode,
			SlotID:              b.SlotID,
			TotalPeople:         b.TotalPeople,
			TotalPrice:          b.TotalPrice,
			Status:              string(b.Status),
			BookedAt:            b.BookedAt,
			ExpiredAt:           b.ExpiredAt,
			DestinationName:     b.DestinationName,
			DestinationLocation: b.DestinationLocation,
			DestinationImage:    b.DestinationImage,
			UserName:            b.UserName,
			UserEmail:           b.UserEmail,
		})
	}
	return resps, nil
}

func (s *bookingService) GetAllBookings(ctx context.Context) ([]dto.BookingResponse, error) {
	bookings, err := s.bookingRepo.ListAllWithDetails(ctx)
	if err != nil {
		return nil, err
	}

	resps := make([]dto.BookingResponse, 0, len(bookings))
	for _, b := range bookings {
		resps = append(resps, dto.BookingResponse{
			ID:                  b.ID,
			BookingCode:         b.BookingCode,
			SlotID:              b.SlotID,
			TotalPeople:         b.TotalPeople,
			TotalPrice:          b.TotalPrice,
			Status:              string(b.Status),
			BookedAt:            b.BookedAt,
			ExpiredAt:           b.ExpiredAt,
			DestinationName:     b.DestinationName,
			DestinationLocation: b.DestinationLocation,
			DestinationImage:    b.DestinationImage,
			UserName:            b.UserName,
			UserEmail:           b.UserEmail,
		})
	}
	return resps, nil
}

func generateBookingCode() string {
	randomBytes := make([]byte, 3)
	_, _ = rand.Read(randomBytes)
	randomPart := strings.ToUpper(fmt.Sprintf("%x", randomBytes))
	return fmt.Sprintf("BK-%s-%s", time.Now().Format("20060102"), randomPart)
}

func intPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}
