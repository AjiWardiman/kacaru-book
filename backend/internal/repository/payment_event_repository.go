package repository

import (
	"context"

	"gorm.io/gorm"

	"tourism-booking/internal/model"
)

type PaymentEventRepository interface {
	CreateTx(ctx context.Context, tx *gorm.DB, e *model.PaymentEvent) error
}

type paymentEventRepository struct {
	db *gorm.DB
}

func NewPaymentEventRepository(db *gorm.DB) PaymentEventRepository {
	return &paymentEventRepository{db: db}
}

func (r *paymentEventRepository) CreateTx(ctx context.Context, tx *gorm.DB, e *model.PaymentEvent) error {
	return tx.WithContext(ctx).Create(e).Error
}
