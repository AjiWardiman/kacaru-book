package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
	"tourism-booking/internal/service"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect memory sqlite database: %v", err)
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Destination{},
		&model.Slot{},
		&model.Booking{},
		&model.BookingDetail{},
		&model.Payment{},
		&model.AuditLog{},
	)
	if err != nil {
		t.Fatalf("failed to auto migrate tables: %v", err)
	}

	return db
}

type dummyAuditService struct{}

func (d *dummyAuditService) Log(ctx context.Context, userID *int64, action, tableName string, recordID int64, oldValue, newValue interface{}) {
}

func TestBookingService_CreateBooking_ValidationErrors(t *testing.T) {
	db := setupTestDB(t)

	slotRepo := repository.NewSlotRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	destRepo := repository.NewDestinationRepository(db)
	auditSvc := &dummyAuditService{}

	svc := service.NewBookingService(db, slotRepo, bookingRepo, destRepo, auditSvc)
	ctx := context.Background()

	t.Run("visitor count mismatch", func(t *testing.T) {
		req := dto.CreateBookingRequest{
			SlotID:      1,
			TotalPeople: 2,
			Visitors: []dto.VisitorInput{
				{Name: "John Doe", Phone: "0812345", Email: "john@example.com", Age: 25},
			},
		}

		_, err := svc.CreateBooking(ctx, 10, req)
		if err != service.ErrVisitorCountMismatch {
			t.Errorf("expected ErrVisitorCountMismatch, got %v", err)
		}
	})

	t.Run("slot not found", func(t *testing.T) {
		req := dto.CreateBookingRequest{
			SlotID:      999,
			TotalPeople: 1,
			Visitors: []dto.VisitorInput{
				{Name: "John Doe", Phone: "0812345", Email: "john@example.com", Age: 25},
			},
		}

		_, err := svc.CreateBooking(ctx, 10, req)
		if err != service.ErrSlotNotFound {
			t.Errorf("expected ErrSlotNotFound, got %v", err)
		}
	})

	t.Run("slot not available", func(t *testing.T) {
		dest := model.Destination{Name: "Pantai Indah", Location: "Bali", Price: 50000, IsActive: true}
		db.Create(&dest)

		slot := model.Slot{
			DestinationID: dest.ID,
			Date:          time.Now().AddDate(0, 0, 1),
			TotalQuota:    10,
			UsedQuota:     0,
			IsAvailable:   false,
		}
		db.Create(&slot)

		req := dto.CreateBookingRequest{
			SlotID:      slot.ID,
			TotalPeople: 1,
			Visitors: []dto.VisitorInput{
				{Name: "John Doe", Phone: "0812345", Email: "john@example.com", Age: 25},
			},
		}

		_, err := svc.CreateBooking(ctx, 10, req)
		if err != service.ErrSlotNotAvailable {
			t.Errorf("expected ErrSlotNotAvailable, got %v", err)
		}
	})

	t.Run("quota not enough", func(t *testing.T) {
		dest := model.Destination{Name: "Gunung Bromo", Location: "Malang", Price: 100000, IsActive: true}
		db.Create(&dest)

		slot := model.Slot{
			DestinationID: dest.ID,
			Date:          time.Now().AddDate(0, 0, 2),
			TotalQuota:    5,
			UsedQuota:     4, // Only 1 quota left
			IsAvailable:   true,
		}
		db.Create(&slot)

		req := dto.CreateBookingRequest{
			SlotID:      slot.ID,
			TotalPeople: 2, // Requesting 2
			Visitors: []dto.VisitorInput{
				{Name: "User A", Phone: "111", Email: "a@ex.com", Age: 20},
				{Name: "User B", Phone: "222", Email: "b@ex.com", Age: 22},
			},
		}

		_, err := svc.CreateBooking(ctx, 10, req)
		if err != service.ErrQuotaNotEnough {
			t.Errorf("expected ErrQuotaNotEnough, got %v", err)
		}
	})
}

func TestBookingService_CreateBooking_Success(t *testing.T) {
	db := setupTestDB(t)

	slotRepo := repository.NewSlotRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	destRepo := repository.NewDestinationRepository(db)
	auditSvc := &dummyAuditService{}

	svc := service.NewBookingService(db, slotRepo, bookingRepo, destRepo, auditSvc)
	ctx := context.Background()

	// Seed user & destination
	user := model.User{Name: "Test User", Email: "test@example.com", Role: model.RoleUser}
	db.Create(&user)

	dest := model.Destination{Name: "Candi Borobudur", Location: "Magelang", Price: 75000, IsActive: true}
	db.Create(&dest)

	slot := model.Slot{
		DestinationID: dest.ID,
		Date:          time.Now().AddDate(0, 0, 3),
		TotalQuota:    10,
		UsedQuota:     2,
		IsAvailable:   true,
	}
	db.Create(&slot)

	req := dto.CreateBookingRequest{
		SlotID:      slot.ID,
		TotalPeople: 2,
		Visitors: []dto.VisitorInput{
			{Name: "Alice", Phone: "08111", Email: "alice@ex.com", Age: 25},
			{Name: "Bob", Phone: "08222", Email: "bob@ex.com", Age: 30},
		},
	}

	resp, err := svc.CreateBooking(ctx, user.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.TotalPeople != 2 {
		t.Errorf("expected TotalPeople 2, got %d", resp.TotalPeople)
	}

	expectedTotalPrice := float64(2 * 75000)
	if resp.TotalPrice != expectedTotalPrice {
		t.Errorf("expected TotalPrice %f, got %f", expectedTotalPrice, resp.TotalPrice)
	}

	if resp.Status != string(model.BookingPending) {
		t.Errorf("expected Status %s, got %s", model.BookingPending, resp.Status)
	}

	// Verify used_quota updated in database
	var updatedSlot model.Slot
	db.First(&updatedSlot, slot.ID)
	if updatedSlot.UsedQuota != 4 { // 2 + 2 = 4
		t.Errorf("expected slot UsedQuota 4, got %d", updatedSlot.UsedQuota)
	}
}
