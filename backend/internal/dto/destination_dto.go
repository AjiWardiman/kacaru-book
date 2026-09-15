package dto

import "time"

// CreateDestinationRequest adalah data yang dikirim admin saat membuat
// destinasi wisata baru (POST /admin/destinations).
type CreateDestinationRequest struct {
	Name        string  `json:"name"`
	Location    string  `json:"location"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
}

// UpdateDestinationRequest adalah data untuk mengubah destinasi yang
// sudah ada (PUT /admin/destinations/:id).
//
// Catatan: semua field ini OPSIONAL. Kalau field dikirim kosong ("" atau 0),
// nilai lama di database TIDAK akan diubah (lihat logic di service).
// Khusus IsActive dibuat pointer (*bool) supaya kita bisa membedakan antara
// "user tidak mengirim field ini" (nil) dengan "user mengirim false" (isi false).
type UpdateDestinationRequest struct {
	Name        string  `json:"name"`
	Location    string  `json:"location"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}

// DestinationResponse adalah bentuk data destinasi yang dikirim ke client
// (baik untuk publik maupun admin).
type DestinationResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description string    `json:"description,omitempty"`
	Category    string    `json:"category,omitempty"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateSlotRequest adalah data untuk membuat slot/jadwal baru pada
// sebuah destinasi (POST /admin/destinations/:id/slots).
type CreateSlotRequest struct {
	Date       string `json:"date"`        // format tanggal: "YYYY-MM-DD", contoh "2026-09-20"
	TotalQuota int    `json:"total_quota"` // jumlah maksimal pengunjung yang boleh booking di tanggal ini
}

// SlotResponse adalah bentuk data slot yang dikirim ke client.
// RemainingQuota dihitung otomatis (TotalQuota - UsedQuota) supaya
// client tidak perlu menghitung sendiri sisa kuota.
type SlotResponse struct {
	ID             int64  `json:"id"`
	DestinationID  int64  `json:"destination_id"`
	Date           string `json:"date"`
	TotalQuota     int    `json:"total_quota"`
	UsedQuota      int    `json:"used_quota"`
	RemainingQuota int    `json:"remaining_quota"` // sisa kuota yang masih bisa dibooking
	IsAvailable    bool   `json:"is_available"`
}
