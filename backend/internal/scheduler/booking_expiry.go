// Package scheduler berisi pekerjaan-pekerjaan yang berjalan otomatis
// di LATAR BELAKANG (background), terpisah dari alur request-response
// HTTP biasa. Contoh di sini: secara berkala memeriksa dan membersihkan
// booking yang tidak dibayar tepat waktu.
package scheduler

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"

	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
)

// BookingExpiryWorker adalah "pekerja" yang berjalan terus-menerus di
// latar belakang, memeriksa secara berkala apakah ada booking yang
// sudah lewat batas waktu bayarnya (expired_at) tapi statusnya masih
// "pending" — kalau ada, otomatis diubah jadi "expired" dan kuota yang
// tadi terpakai dikembalikan ke slot supaya bisa dipesan orang lain.
type BookingExpiryWorker struct {
	db          *gorm.DB
	bookingRepo repository.BookingRepository
	slotRepo    repository.SlotRepository
	interval    time.Duration // seberapa sering worker ini memeriksa (misal tiap 1 menit)
}

func NewBookingExpiryWorker(db *gorm.DB, bookingRepo repository.BookingRepository, slotRepo repository.SlotRepository, interval time.Duration) *BookingExpiryWorker {
	return &BookingExpiryWorker{
		db:          db,
		bookingRepo: bookingRepo,
		slotRepo:    slotRepo,
		interval:    interval,
	}
}

func (w *BookingExpiryWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("Booking expiry worker dimulai, memeriksa setiap %s", w.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Booking expiry worker dihentikan")
			return
		case <-ticker.C:
			w.processExpiredBookings(ctx)
		}
	}
}

func (w *BookingExpiryWorker) processExpiredBookings(ctx context.Context) {
	expiredBookings, err := w.bookingRepo.ListExpiredPending(ctx)
	if err != nil {
		log.Printf("gagal mengambil daftar booking yang expired: %v", err)
		return
	}

	if len(expiredBookings) == 0 {
		return
	}

	log.Printf("ditemukan %d booking yang sudah melewati batas waktu bayar, memproses...", len(expiredBookings))

	for _, booking := range expiredBookings {
		if err := w.expireOneBooking(ctx, booking); err != nil {
			log.Printf("gagal memproses expired untuk booking id=%d: %v", booking.ID, err)
			continue
		}
		log.Printf("booking id=%d (kode=%s) berhasil ditandai expired, kuota dikembalikan", booking.ID, booking.BookingCode)
	}
}

func (w *BookingExpiryWorker) expireOneBooking(ctx context.Context, booking model.Booking) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if _, err := w.slotRepo.LockForUpdate(ctx, tx, booking.SlotID); err != nil {
			return err
		}
		if err := w.slotRepo.DecrementUsedQuota(ctx, tx, booking.SlotID, booking.TotalPeople); err != nil {
			return err
		}
		return w.bookingRepo.UpdateStatusTx(ctx, tx, booking.ID, model.BookingExpired)
	})
}
