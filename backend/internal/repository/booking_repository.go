package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"tourism-booking/internal/model"
)

type BookingRepository interface {
	CreateTx(ctx context.Context, tx *gorm.DB, b *model.Booking) (*model.Booking, error)
	CreateDetailsTx(ctx context.Context, tx *gorm.DB, bookingID int64, details []model.BookingDetail) error
	FindByID(ctx context.Context, id int64) (*model.Booking, error)
	ListByUser(ctx context.Context, userID int64) ([]model.Booking, error)
	ListByUserWithDetails(ctx context.Context, userID int64) ([]model.BookingWithDetails, error)
	ListAllWithDetails(ctx context.Context) ([]model.BookingWithDetails, error)
	UpdateStatusTx(ctx context.Context, tx *gorm.DB, bookingID int64, status model.BookingStatus) error
	ListExpiredPending(ctx context.Context) ([]model.Booking, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) CreateTx(ctx context.Context, tx *gorm.DB, b *model.Booking) (*model.Booking, error) {
	if err := tx.WithContext(ctx).Create(b).Error; err != nil {
		return nil, err
	}
	return b, nil
}

func (r *bookingRepository) CreateDetailsTx(ctx context.Context, tx *gorm.DB, bookingID int64, details []model.BookingDetail) error {
	for i := range details {
		details[i].BookingID = bookingID
	}
	if len(details) > 0 {
		return tx.WithContext(ctx).Create(&details).Error
	}
	return nil
}

func (r *bookingRepository) FindByID(ctx context.Context, id int64) (*model.Booking, error) {
	var b model.Booking
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&b).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *bookingRepository) ListByUser(ctx context.Context, userID int64) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) UpdateStatusTx(ctx context.Context, tx *gorm.DB, bookingID int64, status model.BookingStatus) error {
	return tx.WithContext(ctx).Model(&model.Booking{}).
		Where("id = ?", bookingID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *bookingRepository) ListExpiredPending(ctx context.Context) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.WithContext(ctx).
		Where("status = ? AND expired_at < ?", model.BookingPending, time.Now()).
		Find(&bookings).Error
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) ListByUserWithDetails(ctx context.Context, userID int64) ([]model.BookingWithDetails, error) {
	var list []model.BookingWithDetails
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			b.id,
			b.booking_code,
			b.user_id,
			b.slot_id,
			b.total_people,
			b.total_price,
			b.status,
			b.booked_at,
			b.expired_at,
			COALESCE(d.name, 'Wisata Kabupaten Bima') as destination_name,
			COALESCE(d.location, 'Kabupaten Bima, NTB') as destination_location,
			COALESCE(d.image_url, '/uploads/destinations/lariti.jpg') as destination_image,
			COALESCE(u.name, 'Wisatawan Bima') as user_name,
			COALESCE(u.email, 'customer@example.com') as user_email
		FROM bookings b
		LEFT JOIN users u ON u.id = b.user_id
		LEFT JOIN slots s ON s.id = b.slot_id
		LEFT JOIN destinations d ON d.id = s.destination_id
		WHERE b.user_id = ?
		ORDER BY b.created_at DESC
	`, userID).Scan(&list).Error
	return list, err
}

func (r *bookingRepository) ListAllWithDetails(ctx context.Context) ([]model.BookingWithDetails, error) {
	var list []model.BookingWithDetails
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			b.id,
			b.booking_code,
			b.user_id,
			b.slot_id,
			b.total_people,
			b.total_price,
			b.status,
			b.booked_at,
			b.expired_at,
			COALESCE(d.name, 'Wisata Kabupaten Bima') as destination_name,
			COALESCE(d.location, 'Kabupaten Bima, NTB') as destination_location,
			COALESCE(d.image_url, '/uploads/destinations/lariti.jpg') as destination_image,
			COALESCE(u.name, 'Wisatawan Bima') as user_name,
			COALESCE(u.email, 'customer@example.com') as user_email
		FROM bookings b
		LEFT JOIN users u ON u.id = b.user_id
		LEFT JOIN slots s ON s.id = b.slot_id
		LEFT JOIN destinations d ON d.id = s.destination_id
		ORDER BY b.id DESC
	`).Scan(&list).Error
	return list, err
}
