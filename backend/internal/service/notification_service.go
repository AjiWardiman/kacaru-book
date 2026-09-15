package service

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"tourism-booking/internal/mailer"
	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
)

// NotificationService adalah kontrak untuk mengirim notifikasi
// (lewat email + real-time SSE streamer) ke user, sekaligus mencatat riwayatnya.
type NotificationService interface {
	Notify(ctx context.Context, userID int64, userEmail, notifType, title, message string)
	GetNotifications(ctx context.Context, userID int64, limit, offset int) ([]model.Notification, int64, error)
	MarkAsRead(ctx context.Context, id int64, userID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
	Subscribe(userID int64) (chan model.Notification, func())
}

type notificationService struct {
	mailer           *mailer.Mailer
	notificationRepo repository.NotificationRepository

	mu        sync.RWMutex
	listeners map[int64][]chan model.Notification
}

func NewNotificationService(m *mailer.Mailer, notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{
		mailer:           m,
		notificationRepo: notificationRepo,
		listeners:        make(map[int64][]chan model.Notification),
	}
}

// Subscribe membuat channel listener SSE baru untuk user tertentu.
func (s *notificationService) Subscribe(userID int64) (chan model.Notification, func()) {
	s.mu.Lock()
	ch := make(chan model.Notification, 10)
	s.listeners[userID] = append(s.listeners[userID], ch)
	s.mu.Unlock()

	unsubscribe := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		chs := s.listeners[userID]
		for i, listener := range chs {
			if listener == ch {
				s.listeners[userID] = append(chs[:i], chs[i+1:]...)
				close(ch)
				break
			}
		}
		if len(s.listeners[userID]) == 0 {
			delete(s.listeners, userID)
		}
	}

	return ch, unsubscribe
}

// Notify mencoba mengirim email ke user, menyimpan riwayatnya ke database,
// DAN menyiarkan (broadcast) via SSE channel secara real-time.
func (s *notificationService) Notify(ctx context.Context, userID int64, userEmail, notifType, title, message string) {
	// Hanya kirim email lewat Mailtrap jika tipe notifikasi adalah OTP / Verifikasi
	// Notifikasi login_success, welcome, dll. tetap disimpan di DB & disiarkan via SSE tanpa membuang kuota Mailtrap.
	isOTPEmail := strings.Contains(strings.ToLower(notifType), "otp") || strings.Contains(strings.ToLower(notifType), "verification")
	if userEmail != "" && isOTPEmail {
		// Jalankan di goroutine terpisah (asinkron) agar API respon INSTAN (< 50ms) tanpa tertahan koneksi SMTP
		go func(toEmail, msgTitle, msgBody string) {
			time.Sleep(300 * time.Millisecond) // Hindari rate limit Mailtrap Sandbox (Too many emails per second)
			err := s.mailer.Send(toEmail, msgTitle, msgBody)
			if err != nil {
				log.Printf("gagal mengirim email OTP ke %s: %v", toEmail, err)
			}
		}(userEmail, title, message)
	}

	notification := &model.Notification{
		UserID:  userID,
		Type:    notifType,
		Title:   title,
		Message: message,
		SentAt:  time.Now(),
	}
	if dbErr := s.notificationRepo.Create(ctx, notification); dbErr != nil {
		log.Printf("gagal menyimpan riwayat notifikasi untuk user_id=%d: %v", userID, dbErr)
	}

	// Broadcast ke semua koneksi SSE aktif milik user ini
	s.mu.RLock()
	defer s.mu.RUnlock()
	if chs, exists := s.listeners[userID]; exists {
		for _, ch := range chs {
			select {
			case ch <- *notification:
			default:
				// Jika channel penuh, tidak menahan goroutine
			}
		}
	}
}

func (s *notificationService) GetNotifications(ctx context.Context, userID int64, limit, offset int) ([]model.Notification, int64, error) {
	notifications, err := s.notificationRepo.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	unread, err := s.notificationRepo.CountUnread(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return notifications, unread, nil
}

func (s *notificationService) MarkAsRead(ctx context.Context, id int64, userID int64) error {
	return s.notificationRepo.MarkAsRead(ctx, id, userID)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}
