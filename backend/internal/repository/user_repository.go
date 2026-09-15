package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"tourism-booking/internal/model"
)

// ErrNotFound dikembalikan ketika data yang dicari tidak ditemukan di database.
var ErrNotFound = errors.New("record not found")

type UserRepository interface {
	Create(ctx context.Context, u *model.User) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByPhone(ctx context.Context, phone string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindAll(ctx context.Context, search, roleStr string, limit, offset int) ([]model.User, int64, error)
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id int64) error
	SetPhoneVerified(ctx context.Context, userID int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindAll(ctx context.Context, search, roleStr string, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	if search != "" {
		s := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR phone LIKE ?", s, s, s)
	}

	if roleStr != "" && roleStr != "ALL" {
		rVal := strings.ToLower(roleStr)
		if rVal == "customer" {
			rVal = "user"
		}
		query = query.Where("role = ?", rVal)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *userRepository) Update(ctx context.Context, u *model.User) error {
	u.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Hapus kode OTP terkait agar tidak melanggar foreign key constraint otps_user_id_fkey
		if err := tx.Where("user_id = ?", id).Delete(&model.OTP{}).Error; err != nil {
			return err
		}
		// 2. Hapus notifikasi terkait
		if err := tx.Where("user_id = ?", id).Delete(&model.Notification{}).Error; err != nil {
			return err
		}
		// 3. Hapus riwayat pemesanan/booking terkait
		if err := tx.Where("user_id = ?", id).Delete(&model.Booking{}).Error; err != nil {
			return err
		}
		// 4. Hapus audit log terkait
		if err := tx.Where("user_id = ?", id).Delete(&model.AuditLog{}).Error; err != nil {
			return err
		}
		// 5. Hapus record user
		return tx.Delete(&model.User{}, id).Error
	})
}

func (r *userRepository) SetPhoneVerified(ctx context.Context, userID int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"phone_verified_at": now,
			"updated_at":        now,
		}).Error
}
