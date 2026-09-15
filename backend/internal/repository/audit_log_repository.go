package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"tourism-booking/internal/model"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *model.AuditLog) error
	DeleteBefore(ctx context.Context, cutoff time.Time) (int64, error)
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditLogRepository) DeleteBefore(ctx context.Context, cutoff time.Time) (int64, error) {
	res := r.db.WithContext(ctx).Where("created_at < ?", cutoff).Delete(&model.AuditLog{})
	return res.RowsAffected, res.Error
}

