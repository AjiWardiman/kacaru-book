package service

import (
	"context"
	"encoding/json"
	"log"

	"tourism-booking/internal/model"
	"tourism-booking/internal/repository"
)

// Konstanta jenis aksi yang dicatat di audit trail, supaya penulisannya
// konsisten di seluruh kode (tidak ada yang salah ketik "Create" vs "CREATE").
const (
	ActionCreate = "CREATE"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
)

// AuditService adalah kontrak untuk mencatat audit trail — dipakai
// oleh service lain (Destination, Slot, Booking, Payment) setiap kali
// mereka melakukan perubahan data penting.
type AuditService interface {
	// Log mencatat satu aksi ke audit trail.
	//   - userID: siapa yang melakukan aksi (nil kalau dipicu sistem
	//     otomatis, misal webhook Xendit atau scheduler auto-expire)
	//   - action: salah satu dari konstanta Action... di atas
	//   - tableName: nama tabel yang terpengaruh
	//   - recordID: ID baris yang terpengaruh
	//   - oldValue, newValue: data sebelum/sesudah perubahan (boleh nil
	//     kalau tidak relevan, misal oldValue kosong untuk aksi CREATE)
	//
	// Fungsi ini SENGAJA tidak mengembalikan error ke pemanggilnya —
	// audit trail adalah catatan pendukung, kegagalan mencatatnya tidak
	// boleh sampai membatalkan aksi bisnis utama yang sudah berhasil.
	// Kegagalan cukup dicatat di log server untuk diperiksa nanti.
	Log(ctx context.Context, userID *int64, action, tableName string, recordID int64, oldValue, newValue interface{})
}

type auditService struct {
	repo repository.AuditLogRepository
}

func NewAuditService(repo repository.AuditLogRepository) AuditService {
	return &auditService{repo: repo}
}

func (s *auditService) Log(ctx context.Context, userID *int64, action, tableName string, recordID int64, oldValue, newValue interface{}) {
	oldJSON := toJSONPtr(oldValue)
	newJSON := toJSONPtr(newValue)

	entry := &model.AuditLog{
		UserID:    userID,
		Action:    action,
		TableName: tableName,
		RecordID:  &recordID,
		OldValue:  oldJSON,
		NewValue:  newJSON,
	}

	if err := s.repo.Create(ctx, entry); err != nil {
		log.Printf("gagal mencatat audit trail (aksi=%s, tabel=%s, record_id=%d): %v", action, tableName, recordID, err)
	}
}

// toJSONPtr mengubah data apa pun (struct, map, dll) menjadi string
// JSON, dibungkus pointer. Mengembalikan nil kalau input-nya nil,
// supaya kolom di database tersimpan sebagai NULL (bukan literal
// teks "null") saat memang tidak ada data yang relevan.
func toJSONPtr(v interface{}) *string {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("gagal mengubah data audit trail menjadi JSON: %v", err)
		return nil
	}
	s := string(b)
	return &s
}
