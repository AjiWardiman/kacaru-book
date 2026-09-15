package scheduler

import (
	"context"
	"log"
	"time"

	"tourism-booking/internal/repository"
)

// AuditLogRetentionWorker adalah background worker yang bertugas membersihkan
// data audit log tua (retensi) secara berkala (rotasi 24 jam).
type AuditLogRetentionWorker struct {
	auditRepo         repository.AuditLogRepository
	rotationInterval  time.Duration // interval rotasi (misal 24 jam)
	retentionDuration time.Duration // masa simpan retensi (misal 60 hari / 2 bulan)
}

func NewAuditLogRetentionWorker(auditRepo repository.AuditLogRepository, rotationInterval, retentionDuration time.Duration) *AuditLogRetentionWorker {
	return &AuditLogRetentionWorker{
		auditRepo:         auditRepo,
		rotationInterval:  rotationInterval,
		retentionDuration: retentionDuration,
	}
}

func (w *AuditLogRetentionWorker) Start(ctx context.Context) {
	log.Printf("Audit log retention worker dimulai: rotasi setiap %s, retensi %s", w.rotationInterval, w.retentionDuration)

	// Jalankan pemeriksaan awal saat server pertama kali dinyalakan
	w.cleanOldLogs(ctx)

	ticker := time.NewTicker(w.rotationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Audit log retention worker dihentikan")
			return
		case <-ticker.C:
			w.cleanOldLogs(ctx)
		}
	}
}

func (w *AuditLogRetentionWorker) cleanOldLogs(ctx context.Context) {
	cutoff := time.Now().Add(-w.retentionDuration)
	deletedCount, err := w.auditRepo.DeleteBefore(ctx, cutoff)
	if err != nil {
		log.Printf("gagal membersihkan audit log tua: %v", err)
		return
	}

	if deletedCount > 0 {
		log.Printf("Audit log retention worker: berhasil membersihkan %d record audit log yang lebih tua dari %s (sebelum %s)",
			deletedCount, w.retentionDuration, cutoff.Format("2006-01-02 15:04:05"))
	}
}
