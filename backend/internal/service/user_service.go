package service

import (
	"context"
	"fmt"
	"strings"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
	"tourism-booking/internal/util"
)

type UserService interface {
	GetUsers(ctx context.Context, search, roleStr string, limit, offset int) ([]dto.UserResponse, int64, error)
	GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	CreateUser(ctx context.Context, adminID int64, req dto.RegisterRequest, role string) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, adminID int64, id int64, req dto.RegisterRequest, role string) (*dto.UserResponse, error)
	DeleteUser(ctx context.Context, adminID int64, id int64) error
}

type userService struct {
	userRepo     repository.UserRepository
	auditService AuditService
	notifService NotificationService
}

func NewUserService(userRepo repository.UserRepository, auditService AuditService, notifService NotificationService) UserService {
	return &userService{
		userRepo:     userRepo,
		auditService: auditService,
		notifService: notifService,
	}
}

func (s *userService) GetUsers(ctx context.Context, search, roleStr string, limit, offset int) ([]dto.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAll(ctx, search, roleStr, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var res []dto.UserResponse
	for _, u := range users {
		res = append(res, dto.UserResponse{
			ID:              u.ID,
			Name:            u.Name,
			Email:           u.Email,
			Phone:           u.Phone,
			PhoneVerifiedAt: u.PhoneVerifiedAt,
			Role:            string(u.Role),
			CreatedAt:       u.CreatedAt,
		})
	}
	return res, total, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Phone:           u.Phone,
		PhoneVerifiedAt: u.PhoneVerifiedAt,
		Role:            string(u.Role),
		CreatedAt:       u.CreatedAt,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, adminID int64, req dto.RegisterRequest, role string) (*dto.UserResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)

	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailTaken
	}

	if req.Phone != "" {
		if _, err := s.userRepo.FindByPhone(ctx, req.Phone); err == nil {
			return nil, ErrPhoneTaken
		}
	}

	if req.Password == "" {
		req.Password = "password123"
	}

	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var phonePtr *string
	if req.Phone != "" {
		phonePtr = &req.Phone
	}

	userRole := model.RoleUser
	if strings.ToLower(role) == "admin" {
		userRole = model.RoleAdmin
	}

	u := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    phonePtr,
		Password: hashed,
		Role:     userRole,
	}

	created, err := s.userRepo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	// Record Audit Log
	s.auditService.Log(ctx, &adminID, "CREATE", "users", created.ID, nil, created)

	// Send Real-Time SSE Notification to new user & active admin
	if s.notifService != nil {
		go s.notifService.Notify(
			context.Background(),
			created.ID,
			created.Email,
			"account_created",
			"Akun Pengguna Baru Dibuat",
			fmt.Sprintf("Halo %s, akun Anda dengan role %s telah dibuat oleh Admin.", created.Name, created.Role),
		)
		if adminID > 0 {
			go s.notifService.Notify(
				context.Background(),
				adminID,
				"",
				"user_management",
				"Pengguna Baru Berhasil Ditambahkan",
				fmt.Sprintf("Pengguna %s (%s) dengan role %s telah berhasil ditambahkan ke sistem.", created.Name, created.Email, created.Role),
			)
		}
	}

	res := &dto.UserResponse{
		ID:              created.ID,
		Name:            created.Name,
		Email:           created.Email,
		Phone:           created.Phone,
		PhoneVerifiedAt: created.PhoneVerifiedAt,
		Role:            string(created.Role),
		CreatedAt:       created.CreatedAt,
	}
	return res, nil
}

func (s *userService) UpdateUser(ctx context.Context, adminID int64, id int64, req dto.RegisterRequest, role string) (*dto.UserResponse, error) {
	existing, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	oldUser := *existing

	if req.Name != "" {
		existing.Name = strings.TrimSpace(req.Name)
	}
	if req.Email != "" && req.Email != existing.Email {
		cleanEmail := strings.TrimSpace(strings.ToLower(req.Email))
		if foundUser, err := s.userRepo.FindByEmail(ctx, cleanEmail); err == nil && foundUser.ID != existing.ID {
			return nil, ErrEmailTaken
		}
		existing.Email = cleanEmail
	}
	if req.Phone != "" {
		phoneStr := strings.TrimSpace(req.Phone)
		if existing.Phone == nil || *existing.Phone != phoneStr {
			if foundUser, err := s.userRepo.FindByPhone(ctx, phoneStr); err == nil && foundUser.ID != existing.ID {
				return nil, fmt.Errorf("Nomor WhatsApp/HP %s sudah terdaftar di akun lain", phoneStr)
			}
			existing.Phone = &phoneStr
		}
	}
	if role != "" {
		if strings.ToLower(role) == "admin" {
			existing.Role = model.RoleAdmin
		} else {
			existing.Role = model.RoleUser
		}
	}
	if req.Password != "" {
		hashed, err := util.HashPassword(req.Password)
		if err == nil {
			existing.Password = hashed
		}
	}

	if err := s.userRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	// Record Audit Log
	s.auditService.Log(ctx, &adminID, "UPDATE", "users", existing.ID, oldUser, existing)

	// Send Real-Time SSE Notification to target user & active admin
	if s.notifService != nil {
		go s.notifService.Notify(
			context.Background(),
			existing.ID,
			existing.Email,
			"account_updated",
			"Profil Akun Diperbarui",
			fmt.Sprintf("Halo %s, data profil Anda telah berhasil diperbarui.", existing.Name),
		)
		if adminID > 0 && adminID != existing.ID {
			go s.notifService.Notify(
				context.Background(),
				adminID,
				"",
				"user_management",
				"Data Pengguna Diperbarui",
				fmt.Sprintf("Data pengguna %s (#%d) telah berhasil diperbarui di sistem.", existing.Name, existing.ID),
			)
		}
	}

	return &dto.UserResponse{
		ID:              existing.ID,
		Name:            existing.Name,
		Email:           existing.Email,
		Phone:           existing.Phone,
		PhoneVerifiedAt: existing.PhoneVerifiedAt,
		Role:            string(existing.Role),
		CreatedAt:       existing.CreatedAt,
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, adminID int64, id int64) error {
	existing, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Record Audit Log
	s.auditService.Log(ctx, &adminID, "DELETE", "users", id, existing, nil)

	// Send Real-Time SSE Notification to active admin
	if s.notifService != nil && adminID > 0 {
		go s.notifService.Notify(
			context.Background(),
			adminID,
			"",
			"user_management",
			"Pengguna Telah Dihapus",
			fmt.Sprintf("Pengguna %s (#%d) telah dihapus dari sistem.", existing.Name, id),
		)
	}
	return nil
}

