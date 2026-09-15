package model

import "time"

// AuditLog merepresentasikan satu baris di tabel "audit_logs".
// Setiap kali ada perubahan data penting di sistem (create/update/delete),
// kita catat di sini: SIAPA melakukannya, APA aksinya, terhadap DATA
// MANA, dan KAPAN — ini requirement wajib dari trainer untuk keperluan
// audit/pelacakan.
type AuditLog struct {
	ID        int64
	UserID    *int64  // siapa yang melakukan aksi ini; nil kalau aksi dipicu sistem otomatis (misal webhook)
	Action    string  // "CREATE", "UPDATE", atau "DELETE"
	TableName string  // nama tabel yang terpengaruh, misal "destinations"
	RecordID  *int64  // ID baris yang terpengaruh
	OldValue  *string // data SEBELUM perubahan (format JSON), kosong untuk aksi CREATE
	NewValue  *string // data SESUDAH perubahan (format JSON), kosong untuk aksi DELETE
	CreatedAt time.Time
}

