package repository

import (
	"context"

	"gorm.io/gorm"

	"tourism-booking/internal/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *model.Notification) error
	FindByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.Notification, error)
	CountUnread(ctx context.Context, userID int64) (int64, error)
	MarkAsRead(ctx context.Context, id int64, userID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n *model.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID int64, limit, offset int) ([]model.Notification, error) {
	var notifications []model.Notification
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("sent_at DESC")
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	err := query.Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) CountUnread(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id int64, userID int64) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}
