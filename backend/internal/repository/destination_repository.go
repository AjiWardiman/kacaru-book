package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"tourism-booking/internal/model"
)

type DestinationRepository interface {
	Create(ctx context.Context, d *model.Destination) (*model.Destination, error)
	FindByID(ctx context.Context, id int64) (*model.Destination, error)
	List(ctx context.Context, onlyActive bool) ([]model.Destination, error)
	Update(ctx context.Context, d *model.Destination) error
	Delete(ctx context.Context, id int64) error
}

type destinationRepository struct {
	db *gorm.DB
}

func NewDestinationRepository(db *gorm.DB) DestinationRepository {
	return &destinationRepository{db: db}
}

func (r *destinationRepository) Create(ctx context.Context, d *model.Destination) (*model.Destination, error) {
	if err := r.db.WithContext(ctx).Create(d).Error; err != nil {
		return nil, err
	}
	return d, nil
}

func (r *destinationRepository) FindByID(ctx context.Context, id int64) (*model.Destination, error) {
	var d model.Destination
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *destinationRepository) List(ctx context.Context, onlyActive bool) ([]model.Destination, error) {
	var destinations []model.Destination
	query := r.db.WithContext(ctx)
	if onlyActive {
		query = query.Where("is_active = ?", true)
	}
	err := query.Order("created_at DESC").Find(&destinations).Error
	if err != nil {
		return nil, err
	}
	return destinations, nil
}

func (r *destinationRepository) Update(ctx context.Context, d *model.Destination) error {
	d.UpdatedAt = time.Now()
	res := r.db.WithContext(ctx).Model(&model.Destination{}).
		Where("id = ?", d.ID).
		Updates(map[string]interface{}{
			"name":        d.Name,
			"location":    d.Location,
			"description": d.Description,
			"category":    d.Category,
			"price":       d.Price,
			"image_url":   d.ImageURL,
			"is_active":   d.IsActive,
			"updated_at":  d.UpdatedAt,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *destinationRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Destination{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
