package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"tourism-booking/internal/model"
)

type PaymentRepository interface {
	Create(ctx context.Context, p *model.Payment) (*model.Payment, error)
	FindByExternalID(ctx context.Context, externalID string) (*model.Payment, error)
	FindByBookingID(ctx context.Context, bookingID int64) (*model.Payment, error)
	UpdateStatusTx(ctx context.Context, tx *gorm.DB, paymentID int64, status model.PaymentStatus, paymentMethod *string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, p *model.Payment) (*model.Payment, error) {
	if err := r.db.WithContext(ctx).Create(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *paymentRepository) FindByExternalID(ctx context.Context, externalID string) (*model.Payment, error) {
	var p model.Payment
	err := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) FindByBookingID(ctx context.Context, bookingID int64) (*model.Payment, error) {
	var p model.Payment
	err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) UpdateStatusTx(ctx context.Context, tx *gorm.DB, paymentID int64, status model.PaymentStatus, paymentMethod *string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if paymentMethod != nil && *paymentMethod != "" {
		updates["payment_method"] = *paymentMethod
	}
	if status == model.PaymentPaid {
		updates["paid_at"] = time.Now()
	}
	return tx.WithContext(ctx).Model(&model.Payment{}).
		Where("id = ?", paymentID).
		Updates(updates).Error
}
