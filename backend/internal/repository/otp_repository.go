package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"tourism-booking/internal/model"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *model.OTP) (*model.OTP, error)
	FindLatestActive(ctx context.Context, userID int64, purpose model.OTPPurpose) (*model.OTP, error)
	MarkUsed(ctx context.Context, id int64) error
	InvalidateActive(ctx context.Context, userID int64, purpose model.OTPPurpose) error
}

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(ctx context.Context, otp *model.OTP) (*model.OTP, error) {
	if err := r.db.WithContext(ctx).Create(otp).Error; err != nil {
		return nil, err
	}
	return otp, nil
}

func (r *otpRepository) FindLatestActive(ctx context.Context, userID int64, purpose model.OTPPurpose) (*model.OTP, error) {
	var o model.OTP
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND purpose = ? AND is_used = false", userID, purpose).
		Order("created_at DESC").
		First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *otpRepository) MarkUsed(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.OTP{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}

func (r *otpRepository) InvalidateActive(ctx context.Context, userID int64, purpose model.OTPPurpose) error {
	return r.db.WithContext(ctx).Model(&model.OTP{}).
		Where("user_id = ? AND purpose = ? AND is_used = false", userID, purpose).
		Update("is_used", true).Error
}
