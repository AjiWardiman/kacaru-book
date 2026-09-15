package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tourism-booking/internal/dto"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
	"tourism-booking/internal/util"
)

// Daftar error yang bisa terjadi di proses OTP, biar handler bisa
// mengecek jenisnya secara pasti (pakai errors.Is) dan memberi status
// HTTP + pesan yang tepat ke client.
var (
	ErrPhoneNotSet        = errors.New("user has no phone number to verify")
	ErrPhoneAlreadyVerify = errors.New("phone already verified")
	ErrOTPNotRequested    = errors.New("no active otp requested")
	ErrOTPExpired         = errors.New("otp expired")
	ErrOTPInvalidCode     = errors.New("otp code invalid")
)

// otpCodeLength adalah panjang kode OTP yang dibuat, 6 digit seperti
// kebanyakan aplikasi pada umumnya (cukup aman, tapi masih gampang
// diketik ulang manual oleh user).
const otpCodeLength = 6

// defaultOTPExpiryMinutes dipakai HANYA kalau setting
// "otp_expiry_minutes" entah kenapa tidak ada/rusak di database — lihat
// SettingService.GetIntOrDefault.
const defaultOTPExpiryMinutes = 5

// OTPService adalah kontrak logic bisnis untuk verifikasi nomor HP
// lewat kode OTP (One Time Password).
//
// CATATAN SIMULASI: sesuai kesepakatan tugas akhir, kode OTP di sini
// dikirim lewat EMAIL (bukan SMS beneran), memakai infrastruktur
// pengiriman email (Mailtrap) yang sudah ada untuk fitur notifikasi.
// Alur bisnisnya (generate kode, simpan, cek kadaluarsa, sekali pakai)
// tetap sama persis seperti OTP asli lewat SMS.
type OTPService interface {
	// RequestPhoneVerificationOTP membuat kode OTP baru untuk
	// verifikasi HP milik user, menyimpannya ke database, lalu
	// mengirimkannya lewat email.
	RequestPhoneVerificationOTP(ctx context.Context, userID int64) (*dto.RequestOTPResponse, error)
	// VerifyPhoneOTP mencocokkan kode yang dimasukkan user dengan kode
	// yang tersimpan. Kalau cocok dan belum kadaluarsa, nomor HP user
	// akan ditandai sudah terverifikasi.
	VerifyPhoneOTP(ctx context.Context, userID int64, code string) error
}

type otpService struct {
	otpRepo             repository.OTPRepository
	userRepo            repository.UserRepository
	settingService      SettingService
	notificationService NotificationService
}

func NewOTPService(
	otpRepo repository.OTPRepository,
	userRepo repository.UserRepository,
	settingService SettingService,
	notificationService NotificationService,
) OTPService {
	return &otpService{
		otpRepo:             otpRepo,
		userRepo:            userRepo,
		settingService:      settingService,
		notificationService: notificationService,
	}
}

// RequestPhoneVerificationOTP menjalankan seluruh proses pembuatan dan
// pengiriman kode OTP untuk verifikasi HP:
//  1. Pastikan user memang punya nomor HP yang mau diverifikasi, dan
//     belum pernah diverifikasi sebelumnya.
//  2. Baca "otp_expiry_minutes" dari tabel settings (bisa diubah admin
//     kapan saja lewat endpoint /admin/settings, tanpa perlu ubah kode).
//  3. Buat kode acak 6 digit, dan matikan (invalidate) kode-kode lama
//     milik user ini yang masih aktif — supaya cuma SATU kode yang
//     berlaku dalam satu waktu.
//  4. Simpan kode baru ke database.
//  5. Kirim kodenya lewat email (simulasi SMS) memakai NotificationService
//     yang sudah ada, sekaligus mencatat riwayat pengirimannya.
func (s *otpService) RequestPhoneVerificationOTP(ctx context.Context, userID int64) (*dto.RequestOTPResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Phone == nil || *user.Phone == "" {
		return nil, ErrPhoneNotSet
	}
	if user.PhoneVerifiedAt != nil {
		return nil, ErrPhoneAlreadyVerify
	}

	// Ambil batas waktu berlaku OTP dari settings, dengan nilai
	// cadangan 5 menit kalau setting-nya entah kenapa tidak ada.
	expiryMinutes := s.settingService.GetIntOrDefault(ctx, "otp_expiry_minutes", defaultOTPExpiryMinutes)

	code, err := util.GenerateOTPCode(otpCodeLength)
	if err != nil {
		return nil, err
	}

	// Matikan dulu semua kode lama yang masih aktif, supaya user tidak
	// bisa "salah pakai" kode yang lebih tua kalau dia minta kode baru
	// berkali-kali.
	if err := s.otpRepo.InvalidateActive(ctx, userID, model.OTPPurposePhoneVerification); err != nil {
		return nil, err
	}

	otp := &model.OTP{
		UserID:    userID,
		Code:      code,
		Purpose:   model.OTPPurposePhoneVerification,
		ExpiredAt: time.Now().Add(time.Duration(expiryMinutes) * time.Minute),
	}
	if _, err := s.otpRepo.Create(ctx, otp); err != nil {
		return nil, err
	}

	// Kirim kodenya lewat email. Dipanggil sinkron (bukan goroutine
	// terpisah) karena user perlu tahu SEGERA kalau ternyata
	// pengiriman gagal (beda kasusnya dengan notifikasi booking yang
	// sifatnya cuma info tambahan) — walau NotificationService sendiri
	// tetap tidak mengembalikan error fatal (lihat NotificationService.Notify).
	subject := "Kode Verifikasi Nomor HP - Tourism Booking"
	body := fmt.Sprintf(
		"Kode OTP verifikasi nomor HP kamu adalah: %s\n\nKode ini berlaku selama %d menit. Jangan berikan kode ini ke siapa pun, termasuk pihak yang mengaku dari Tourism Booking.",
		code, expiryMinutes,
	)
	s.notificationService.Notify(ctx, user.ID, user.Email, "phone_verification_otp", subject, body)

	return &dto.RequestOTPResponse{
		Message:          "kode OTP telah dikirim ke email kamu (simulasi verifikasi HP)",
		ExpiresInMinutes: expiryMinutes,
	}, nil
}

// VerifyPhoneOTP mencocokkan kode yang dikirim user dengan kode OTP
// paling baru yang masih aktif di database, lalu:
//   - kalau TIDAK ADA kode aktif sama sekali → ErrOTPNotRequested
//     (user belum pernah minta kode, atau kode lamanya sudah dipakai)
//   - kalau ADA tapi sudah lewat batas waktu → ErrOTPExpired
//   - kalau ADA dan belum kadaluarsa, tapi kodenya beda → ErrOTPInvalidCode
//   - kalau cocok semua → tandai kode "sudah dipakai" dan tandai nomor
//     HP user sebagai terverifikasi.
func (s *otpService) VerifyPhoneOTP(ctx context.Context, userID int64, code string) error {
	otp, err := s.otpRepo.FindLatestActive(ctx, userID, model.OTPPurposePhoneVerification)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrOTPNotRequested
	}
	if err != nil {
		return err
	}

	if time.Now().After(otp.ExpiredAt) {
		return ErrOTPExpired
	}

	if otp.Code != code {
		return ErrOTPInvalidCode
	}

	// Tandai kode ini sudah dipakai SEBELUM menandai HP terverifikasi,
	// supaya kalau langkah menandai HP gagal (misal database error),
	// kode ini tetap dianggap "terpakai" dan tidak bisa dicoba lagi
	// berkali-kali oleh orang lain yang mungkin menebak-nebak kode.
	if err := s.otpRepo.MarkUsed(ctx, otp.ID); err != nil {
		return err
	}

	return s.userRepo.SetPhoneVerified(ctx, userID)
}
