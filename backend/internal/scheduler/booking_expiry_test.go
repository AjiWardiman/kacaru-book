package scheduler_test

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
	"tourism-booking/internal/scheduler"
)

func setupSchedulerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect memory sqlite database: %v", err)
	}

	err = db.AutoMigrate(
		&model.Destination{},
		&model.Slot{},
		&model.Booking{},
	)
	if err != nil {
		t.Fatalf("failed to auto migrate tables: %v", err)
	}

	return db
}

func TestBookingExpiryWorker_ProcessExpired(t *testing.T) {
	db := setupSchedulerTestDB(t)

	bookingRepo := repository.NewBookingRepository(db)
	slotRepo := repository.NewSlotRepository(db)

	worker := scheduler.NewBookingExpiryWorker(db, bookingRepo, slotRepo, 100*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Seed slot with total_quota=10, used_quota=3
	slot := model.Slot{
		DestinationID: 1,
		Date:          time.Now(),
		TotalQuota:    10,
		UsedQuota:     3,
		IsAvailable:   true,
	}
	db.Create(&slot)

	// Seed expired pending booking
	expiredBooking := model.Booking{
		BookingCode: "BK-EXP-001",
		UserID:      1,
		SlotID:      slot.ID,
		TotalPeople: 3,
		TotalPrice:  150000,
		Status:      model.BookingPending,
		BookedAt:    time.Now().Add(-1 * time.Hour),
		ExpiredAt:   time.Now().Add(-30 * time.Minute), // Expired 30 mins ago
	}
	db.Create(&expiredBooking)

	// Run worker loop in background
	worker.Start(ctx)

	// Verify booking status is now "expired"
	var updatedBooking model.Booking
	db.First(&updatedBooking, expiredBooking.ID)
	if updatedBooking.Status != model.BookingExpired {
		t.Errorf("expected booking status %s, got %s", model.BookingExpired, updatedBooking.Status)
	}

	// Verify used_quota was decremented from 3 to 0
	var updatedSlot model.Slot
	db.First(&updatedSlot, slot.ID)
	if updatedSlot.UsedQuota != 0 {
		t.Errorf("expected slot UsedQuota 0, got %d", updatedSlot.UsedQuota)
	}
}
