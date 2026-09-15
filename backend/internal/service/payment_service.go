package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
	"tourism-booking/internal/xendit"
)

var (
	ErrBookingNotOwned      = errors.New("booking does not belong to this user")
	ErrBookingNotPending    = errors.New("booking is not in pending status")
	ErrPaymentAlreadyExists = errors.New("payment already exists for this booking")
	ErrInvalidCallbackToken = errors.New("invalid webhook callback token")
	ErrPaymentNotFound      = errors.New("payment not found")
)

type PaymentService interface {
	CreatePayment(ctx context.Context, userID, bookingID int64) (*dto.PaymentResponse, error)
	HandleXenditWebhook(ctx context.Context, rawPayload []byte, callbackToken string) error
}

type paymentService struct {
	db                  *gorm.DB
	xenditClient        *xendit.Client
	callbackToken       string
	paymentRepo         repository.PaymentRepository
	eventRepo           repository.PaymentEventRepository
	bookingRepo         repository.BookingRepository
	slotRepo            repository.SlotRepository
	userRepo            repository.UserRepository
	notificationService NotificationService
	auditService        AuditService
}

func NewPaymentService(
	db *gorm.DB,
	xenditClient *xendit.Client,
	callbackToken string,
	paymentRepo repository.PaymentRepository,
	eventRepo repository.PaymentEventRepository,
	bookingRepo repository.BookingRepository,
	slotRepo repository.SlotRepository,
	userRepo repository.UserRepository,
	notificationService NotificationService,
	auditService AuditService,
) PaymentService {
	return &paymentService{
		db:                  db,
		xenditClient:        xenditClient,
		callbackToken:       callbackToken,
		paymentRepo:         paymentRepo,
		eventRepo:           eventRepo,
		bookingRepo:         bookingRepo,
		slotRepo:            slotRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
		auditService:        auditService,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, userID, bookingID int64) (*dto.PaymentResponse, error) {
	booking, err := s.bookingRepo.FindByID(ctx, bookingID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrBookingNotFound
	}
	if err != nil {
		return nil, err
	}

	if booking.UserID != userID {
		return nil, ErrBookingNotOwned
	}
	if booking.Status != model.BookingPending {
		return nil, ErrBookingNotPending
	}

	if _, err := s.paymentRepo.FindByBookingID(ctx, bookingID); err == nil {
		return nil, ErrPaymentAlreadyExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	invoiceResp, err := s.xenditClient.CreateInvoice(xendit.CreateInvoiceRequest{
		ExternalID:  booking.BookingCode,
		Amount:      booking.TotalPrice,
		Description: fmt.Sprintf("Pembayaran booking %s", booking.BookingCode),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal membuat invoice di Xendit: %w", err)
	}

	payment := &model.Payment{
		BookingID:  bookingID,
		ExternalID: booking.BookingCode,
		InvoiceID:  &invoiceResp.ID,
		Amount:     booking.TotalPrice,
		Status:     model.PaymentPending,
		PaymentURL: &invoiceResp.InvoiceURL,
	}
	created, err := s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	return &dto.PaymentResponse{
		ID:         created.ID,
		BookingID:  created.BookingID,
		Amount:     created.Amount,
		Status:     string(created.Status),
		PaymentURL: *created.PaymentURL,
		CreatedAt:  created.CreatedAt,
	}, nil
}

func (s *paymentService) HandleXenditWebhook(ctx context.Context, rawPayload []byte, callbackToken string) error {
	if s.callbackToken != "" && callbackToken != s.callbackToken {
		return ErrInvalidCallbackToken
	}

	var payload dto.XenditWebhookPayload
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return fmt.Errorf("payload webhook tidak valid: %w", err)
	}

	payment, err := s.paymentRepo.FindByExternalID(ctx, payload.ExternalID)
	if errors.Is(err, repository.ErrNotFound) {
		var booking *model.Booking
		var bookingID int64

		if strings.HasPrefix(payload.ExternalID, "INV-BIMA-") {
			idStr := strings.TrimPrefix(payload.ExternalID, "INV-BIMA-")
			bookingID, _ = strconv.ParseInt(idStr, 10, 64)
		}

		if bookingID > 0 {
			booking, _ = s.bookingRepo.FindByID(ctx, bookingID)
		} else {
			var b model.Booking
			if err := s.db.WithContext(ctx).Where("booking_code = ?", payload.ExternalID).First(&b).Error; err == nil {
				booking = &b
			} else if idNum, parseErr := strconv.ParseInt(payload.ExternalID, 10, 64); parseErr == nil && idNum > 0 {
				booking, _ = s.bookingRepo.FindByID(ctx, idNum)
			}
		}

		if booking != nil {
			mockUrl := "https://checkout.xendit.co/web/inv_" + payload.ExternalID
			mockInv := "inv_" + payload.ExternalID
			p := &model.Payment{
				BookingID:  booking.ID,
				ExternalID: payload.ExternalID,
				InvoiceID:  &mockInv,
				Amount:     booking.TotalPrice,
				Status:     model.PaymentPending,
				PaymentURL: &mockUrl,
			}
			payment, err = s.paymentRepo.Create(ctx, p)
			if err != nil {
				return err
			}
		} else {
			return ErrPaymentNotFound
		}
	} else if err != nil {
		return err
	}

	booking, err := s.bookingRepo.FindByID(ctx, payment.BookingID)
	if err != nil {
		return err
	}

	newStatus := model.PaymentStatus(strings.ToLower(payload.Status))

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var paymentMethodPtr *string
		if payload.PaymentMethod != "" {
			paymentMethodPtr = &payload.PaymentMethod
		}
		if err := s.paymentRepo.UpdateStatusTx(ctx, tx, payment.ID, newStatus, paymentMethodPtr); err != nil {
			return err
		}

		switch newStatus {
		case model.PaymentPaid:
			if err := s.bookingRepo.UpdateStatusTx(ctx, tx, booking.ID, model.BookingPaid); err != nil {
				return err
			}
		case model.PaymentExpired, model.PaymentFailed:
			if _, err := s.slotRepo.LockForUpdate(ctx, tx, booking.SlotID); err != nil {
				return err
			}
			if err := s.slotRepo.DecrementUsedQuota(ctx, tx, booking.SlotID, booking.TotalPeople); err != nil {
				return err
			}
			if err := s.bookingRepo.UpdateStatusTx(ctx, tx, booking.ID, model.BookingExpired); err != nil {
				return err
			}
		}

		event := &model.PaymentEvent{
			PaymentID: payment.ID,
			EventType: payload.Status,
			Payload:   string(rawPayload),
		}
		return s.eventRepo.CreateTx(ctx, tx, event)
	})

	if err != nil {
		return err
	}

	s.auditService.Log(ctx, nil, ActionUpdate, "payments", payment.ID,
		map[string]string{"status": string(payment.Status)},
		map[string]string{"status": string(newStatus)})

	s.sendPaymentNotification(booking, newStatus)

	return nil
}

func (s *paymentService) sendPaymentNotification(booking *model.Booking, newStatus model.PaymentStatus) {
	go func() {
		ctx := context.Background()

		user, err := s.userRepo.FindByID(ctx, booking.UserID)
		if err != nil {
			log.Printf("gagal mengambil data user untuk notifikasi booking id=%d: %v", booking.ID, err)
			return
		}

		var title, message string
		var notifType string
		switch newStatus {
		case model.PaymentPaid:
			notifType = "booking_paid"
			title = "Pembayaran Berhasil - " + booking.BookingCode
			message = fmt.Sprintf(
				"Halo %s,\n\nPembayaran untuk booking %s sebesar Rp%.0f telah berhasil dikonfirmasi. Terima kasih telah memesan!\n\nSalam,\nTourism Booking System",
				user.Name, booking.BookingCode, booking.TotalPrice,
			)
		case model.PaymentExpired, model.PaymentFailed:
			notifType = "booking_expired"
			title = "Booking Dibatalkan - " + booking.BookingCode
			message = fmt.Sprintf(
				"Halo %s,\n\nMohon maaf, booking %s telah dibatalkan karena pembayaran tidak diselesaikan tepat waktu. Kuota yang tadi dipesan sudah dikembalikan.\n\nSalam,\nTourism Booking System",
				user.Name, booking.BookingCode,
			)
		default:
			return
		}

		s.notificationService.Notify(ctx, user.ID, user.Email, notifType, title, message)
	}()
}
