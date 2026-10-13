package service_test

import (
	"context"
	"testing"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
	"tourism-booking/internal/service"
	"tourism-booking/internal/util"
)

type dummyNotifService struct{}

func (d *dummyNotifService) Notify(ctx context.Context, userID int64, userEmail, notifType, title, message string) {
}

func (d *dummyNotifService) GetNotifications(ctx context.Context, userID int64, limit, offset int) ([]model.Notification, int64, error) {
	return nil, 0, nil
}

func (d *dummyNotifService) MarkAsRead(ctx context.Context, id int64, userID int64) error {
	return nil
}

func (d *dummyNotifService) MarkAllAsRead(ctx context.Context, userID int64) error {
	return nil
}

func (d *dummyNotifService) Subscribe(userID int64) (chan model.Notification, func()) {
	return nil, func() {}
}

func TestAuthService_RegisterAndLogin(t *testing.T) {
	db := setupTestDB(t)

	userRepo := repository.NewUserRepository(db)
	notifSvc := &dummyNotifService{}
	jwtSecret := "auth-test-secret"

	svc := service.NewAuthService(userRepo, jwtSecret, 24, notifSvc)
	ctx := context.Background()

	// 1. Test Register Success
	regReq := dto.RegisterRequest{
		Name:     "Bima Sakti",
		Email:    "bima@example.com",
		Password: "password123",
	}

	authResp, err := svc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if authResp.Token == "" {
		t.Error("expected non-empty JWT token on register")
	}
	if authResp.User.Email != "bima@example.com" {
		t.Errorf("expected email bima@example.com, got %s", authResp.User.Email)
	}

	// 2. Test Register Duplicate Email
	_, err = svc.Register(ctx, regReq)
	if err != service.ErrEmailTaken {
		t.Errorf("expected ErrEmailTaken on duplicate email, got %v", err)
	}

	// 3. Test Login Success (with Email)
	loginReq := dto.LoginRequest{
		Identifier: "bima@example.com",
		Password:   "password123",
	}

	loginResp, err := svc.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if loginResp.Token == "" {
		t.Error("expected non-empty token on login")
	}

	// Validate parsed claims from returned token
	claims, err := util.ParseToken(jwtSecret, loginResp.Token)
	if err != nil {
		t.Fatalf("ParseToken failed for returned token: %v", err)
	}
	if claims.Role != string(model.RoleUser) {
		t.Errorf("expected role %s, got %s", model.RoleUser, claims.Role)
	}

	// 4. Test Login Invalid Password
	wrongPassReq := dto.LoginRequest{
		Identifier: "bima@example.com",
		Password:   "wrong-password",
	}

	_, err = svc.Login(ctx, wrongPassReq)
	if err != service.ErrInvalidCreds {
		t.Errorf("expected ErrInvalidCreds on wrong password, got %v", err)
	}

	// 5. Test Login Non-existent User
	nonExistentReq := dto.LoginRequest{
		Identifier: "nobody@example.com",
		Password:   "password123",
	}

	_, err = svc.Login(ctx, nonExistentReq)
	if err != service.ErrInvalidCreds {
		t.Errorf("expected ErrInvalidCreds on non-existent user, got %v", err)
	}
}
