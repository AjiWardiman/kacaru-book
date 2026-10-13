package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"tourism-booking/internal/model"
)

var ErrDuplicateSlot = errors.New("slot for this destination and date already exists")

type SlotRepository interface {
	Create(ctx context.Context, s *model.Slot) (*model.Slot, error)
	FindByID(ctx context.Context, id int64) (*model.Slot, error)
	ListByDestination(ctx context.Context, destinationID int64) ([]model.Slot, error)

	LockForUpdate(ctx context.Context, tx *gorm.DB, slotID int64) (*model.Slot, error)
	IncrementUsedQuota(ctx context.Context, tx *gorm.DB, slotID int64, amount int) error
	DecrementUsedQuota(ctx context.Context, tx *gorm.DB, slotID int64, amount int) error
}

type slotRepository struct {
	db *gorm.DB
}

func NewSlotRepository(db *gorm.DB) SlotRepository {
	return &slotRepository{db: db}
}

func (r *slotRepository) Create(ctx context.Context, s *model.Slot) (*model.Slot, error) {
	if err := r.db.WithContext(ctx).Create(s).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrDuplicateSlot
		}
		return nil, err
	}
	return s, nil
}

func (r *slotRepository) FindByID(ctx context.Context, id int64) (*model.Slot, error) {
	var s model.Slot
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *slotRepository) ListByDestination(ctx context.Context, destinationID int64) ([]model.Slot, error) {
	var slots []model.Slot
	err := r.db.WithContext(ctx).
		Where("destination_id = ?", destinationID).
		Order("date ASC").
		Find(&slots).Error
	if err != nil {
		return nil, err
	}
	return slots, nil
}

func (r *slotRepository) LockForUpdate(ctx context.Context, tx *gorm.DB, slotID int64) (*model.Slot, error) {
	var s model.Slot
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", slotID).
		First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *slotRepository) IncrementUsedQuota(ctx context.Context, tx *gorm.DB, slotID int64, amount int) error {
	return tx.WithContext(ctx).Model(&model.Slot{}).
		Where("id = ?", slotID).
		Updates(map[string]interface{}{
			"used_quota": gorm.Expr("used_quota + ?", amount),
			"updated_at": time.Now(),
		}).Error
}

func (r *slotRepository) DecrementUsedQuota(ctx context.Context, tx *gorm.DB, slotID int64, amount int) error {
	return tx.WithContext(ctx).Model(&model.Slot{}).
		Where("id = ?", slotID).
		Updates(map[string]interface{}{
			"used_quota": gorm.Expr("CASE WHEN COALESCE(used_quota, 0) - ? < 0 THEN 0 ELSE COALESCE(used_quota, 0) - ? END", amount, amount),
			"updated_at": time.Now(),
		}).Error
}
