// Package service berisi LOGIC BISNIS aplikasi.
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
	"tourism-booking/internal/util"
)

var (
	ErrEmailTaken   = errors.New("email already registered")
	ErrPhoneTaken   = errors.New("phone already registered")
	ErrInvalidCreds = errors.New("invalid email/phone or password")
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	SendPreRegisterOTP(ctx context.Context, email string) error
}

type authService struct {
	userRepo     repository.UserRepository
	jwtSecret    string
	jwtExpHours  int
	notifService NotificationService
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpHours int, notifService NotificationService) AuthService {
	return &authService{
		userRepo:     userRepo,
		jwtSecret:    jwtSecret,
		jwtExpHours:  jwtExpHours,
		notifService: notifService,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailTaken
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	if req.Phone != "" {
		if _, err := s.userRepo.FindByPhone(ctx, req.Phone); err == nil {
			return nil, ErrPhoneTaken
		} else if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var phonePtr *string
	if req.Phone != "" {
		phonePtr = &req.Phone
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    phonePtr,
		Password: hashed,
		Role:     model.RoleUser,
	}

	created, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	if s.notifService != nil {
		go s.notifService.Notify(
			context.Background(),
			created.ID,
			created.Email,
			"welcome",
			"Selamat Datang di Kacaru Book!",
			fmt.Sprintf("Selamat akun Anda (%s) telah berhasil terdaftar.", created.Name),
		)
	}

	return s.buildAuthResponse(created)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	var (
		user *model.User
		err  error
	)

	identifier := strings.TrimSpace(req.Identifier)
	user, err = s.userRepo.FindByEmail(ctx, identifier)
	if errors.Is(err, repository.ErrNotFound) {
		user, err = s.userRepo.FindByPhone(ctx, identifier)
	}

	// Dynamic normalization for Phone login (08xxx vs 628xxx vs +628xxx)
	if errors.Is(err, repository.ErrNotFound) {
		cleanPhone := strings.ReplaceAll(strings.ReplaceAll(identifier, " ", ""), "-", "")
		if strings.HasPrefix(cleanPhone, "08") {
			user, err = s.userRepo.FindByPhone(ctx, "62"+strings.TrimPrefix(cleanPhone, "0"))
			if errors.Is(err, repository.ErrNotFound) {
				user, err = s.userRepo.FindByPhone(ctx, "+62"+strings.TrimPrefix(cleanPhone, "0"))
			}
		} else if strings.HasPrefix(cleanPhone, "628") {
			user, err = s.userRepo.FindByPhone(ctx, "0"+strings.TrimPrefix(cleanPhone, "62"))
			if errors.Is(err, repository.ErrNotFound) {
				user, err = s.userRepo.FindByPhone(ctx, "+"+cleanPhone)
			}
		} else if strings.HasPrefix(cleanPhone, "+628") {
			user, err = s.userRepo.FindByPhone(ctx, "0"+strings.TrimPrefix(cleanPhone, "+62"))
			if errors.Is(err, repository.ErrNotFound) {
				user, err = s.userRepo.FindByPhone(ctx, strings.TrimPrefix(cleanPhone, "+"))
			}
		}
	}

	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCreds
	}
	if err != nil {
		return nil, err
	}

	if !util.CheckPassword(user.Password, req.Password) {
		return nil, ErrInvalidCreds
	}

	// Kirim notifikasi login berhasil secara otomatis
	if s.notifService != nil {
		go s.notifService.Notify(
			context.Background(),
			user.ID,
			user.Email,
			"login_success",
			"Selamat Datang Kembali!",
			fmt.Sprintf("Halo %s, Anda berhasil masuk ke platform Kacaru Book Wisata Bima.", user.Name),
		)
	}

	return s.buildAuthResponse(user)
}

func (s *authService) buildAuthResponse(u *model.User) (*dto.AuthResponse, error) {
	token, err := util.GenerateToken(s.jwtSecret, u.ID, string(u.Role), s.jwtExpHours)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:              u.ID,
			Name:            u.Name,
			Email:           u.Email,
			Phone:           u.Phone,
			PhoneVerifiedAt: u.PhoneVerifiedAt,
			Role:            string(u.Role),
			CreatedAt:       u.CreatedAt,
		},
	}, nil
}

func (s *authService) SendPreRegisterOTP(ctx context.Context, email string) error {
	code, err := util.GenerateOTPCode(6)
	if err != nil {
		return err
	}
	if s.notifService != nil && email != "" {
		s.notifService.Notify(
			ctx,
			0,
			email,
			"phone_verification",
			"Kode Verifikasi Nomor HP - Tourism Booking",
			fmt.Sprintf("Kode OTP verifikasi nomor HP kamu adalah: %s\n\nKode ini berlaku selama 5 menit. Jangan berikan kode ini ke siapa pun, termasuk pihak yang mengaku dari Tourism Booking.", code),
		)
	}
	return nil
}
