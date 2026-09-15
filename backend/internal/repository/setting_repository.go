package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"tourism-booking/internal/model"
)

type SettingRepository interface {
	List(ctx context.Context) ([]model.Setting, error)
	FindByKey(ctx context.Context, key string) (*model.Setting, error)
	UpdateValue(ctx context.Context, key, value string) (*model.Setting, error)
}

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) List(ctx context.Context) ([]model.Setting, error) {
	var list []model.Setting
	if err := r.db.WithContext(ctx).Order("key ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *settingRepository) FindByKey(ctx context.Context, key string) (*model.Setting, error) {
	var s model.Setting
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *settingRepository) UpdateValue(ctx context.Context, key, value string) (*model.Setting, error) {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&model.Setting{}).
		Where("key = ?", key).
		Updates(map[string]interface{}{
			"value":      value,
			"updated_at": now,
		})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrNotFound
	}
	return r.FindByKey(ctx, key)
}
