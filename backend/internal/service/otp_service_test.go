package service_test

import (
	"context"
	"testing"
	"time"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
	"tourism-booking/internal/service"
)

type dummySettingService struct{}

func (d *dummySettingService) List(ctx context.Context) ([]dto.SettingResponse, error) {
	return nil, nil
}
func (d *dummySettingService) Get(ctx context.Context, key string) (*dto.SettingResponse, error) {
	return nil, nil
}
func (d *dummySettingService) Update(ctx context.Context, userID int64, key, value string) (*dto.SettingResponse, error) {
	return nil, nil
}
func (d *dummySettingService) GetIntOrDefault(ctx context.Context, key string, fallback int) int {
	return fallback
}

func TestOTPService_Flow(t *testing.T) {
	db := setupTestDB(t)

	err := db.AutoMigrate(&model.OTP{})
	if err != nil {
		t.Fatalf("failed to auto migrate otp model: %v", err)
	}

	otpRepo := repository.NewOTPRepository(db)
	userRepo := repository.NewUserRepository(db)
	settingSvc := &dummySettingService{}
	notifSvc := &dummyNotifService{}

	svc := service.NewOTPService(otpRepo, userRepo, settingSvc, notifSvc)
	ctx := context.Background()

	phone := "08123456789"
	user := model.User{
		Name:     "OTP User",
		Email:    "otp@example.com",
		Phone:    &phone,
		Password: "password123",
		Role:     model.RoleUser,
	}
	db.Create(&user)

	// 1. Request OTP
	resp, err := svc.RequestPhoneVerificationOTP(ctx, user.ID)
	if err != nil {
		t.Fatalf("RequestPhoneVerificationOTP failed: %v", err)
	}
	if resp.ExpiresInMinutes <= 0 {
		t.Errorf("expected positive ExpiresInMinutes, got %d", resp.ExpiresInMinutes)
	}

	// Retrieve generated OTP from DB
	var otpRecord model.OTP
	db.Where("user_id = ?", user.ID).First(&otpRecord)
	if otpRecord.Code == "" {
		t.Fatal("OTP code in database is empty")
	}

	// 2. Verify with Wrong Code
	err = svc.VerifyPhoneOTP(ctx, user.ID, "000000")
	if err != service.ErrOTPInvalidCode {
		t.Errorf("expected ErrOTPInvalidCode for wrong code, got %v", err)
	}

	// 3. Verify with Correct Code
	err = svc.VerifyPhoneOTP(ctx, user.ID, otpRecord.Code)
	if err != nil {
		t.Fatalf("VerifyPhoneOTP failed for correct code: %v", err)
	}

	// Verify user phone_verified_at is updated
	var updatedUser model.User
	db.First(&updatedUser, user.ID)
	if updatedUser.PhoneVerifiedAt == nil {
		t.Error("expected PhoneVerifiedAt to be set after successful OTP verification")
	}

	// 4. Request OTP when already verified
	_, err = svc.RequestPhoneVerificationOTP(ctx, user.ID)
	if err != service.ErrPhoneAlreadyVerify {
		t.Errorf("expected ErrPhoneAlreadyVerify, got %v", err)
	}
}

func TestOTPService_ExpiredOTP(t *testing.T) {
	db := setupTestDB(t)
	db.AutoMigrate(&model.OTP{})

	otpRepo := repository.NewOTPRepository(db)
	userRepo := repository.NewUserRepository(db)
	settingSvc := &dummySettingService{}
	notifSvc := &dummyNotifService{}

	svc := service.NewOTPService(otpRepo, userRepo, settingSvc, notifSvc)
	ctx := context.Background()

	phone := "08987654321"
	user := model.User{
		Name:     "Expired User",
		Email:    "exp@example.com",
		Phone:    &phone,
		Password: "password123",
		Role:     model.RoleUser,
	}
	db.Create(&user)

	// Seed an expired OTP directly in DB
	expiredOTP := model.OTP{
		UserID:    user.ID,
		Code:      "654321",
		Purpose:   model.OTPPurposePhoneVerification,
		IsUsed:    false,
		ExpiredAt: time.Now().Add(-10 * time.Minute), // Expired 10 mins ago
	}
	db.Create(&expiredOTP)

	err := svc.VerifyPhoneOTP(ctx, user.ID, "654321")
	if err != service.ErrOTPExpired {
		t.Errorf("expected ErrOTPExpired for expired OTP, got %v", err)
	}
}
