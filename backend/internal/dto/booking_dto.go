package dto

import "time"

// VisitorInput adalah data SATU pengunjung yang dikirim client saat
// membuat booking. Satu booking bisa punya banyak visitor (rombongan
// keluarga/teman yang berangkat bersama-sama).
type VisitorInput struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// CreateBookingRequest adalah data yang dikirim user saat membuat
// booking baru (POST /bookings).
//
// PENTING: jumlah data di array Visitors HARUS SAMA PERSIS dengan
// TotalPeople. Ini divalidasi di service (lihat ErrVisitorCountMismatch)
// — tujuannya supaya data pengunjung yang tersimpan selalu konsisten
// dengan jumlah kuota yang benar-benar dipakai. Misal TotalPeople=3
// tapi cuma kirim 2 data Visitors, request ini akan ditolak.
type CreateBookingRequest struct {
	SlotID      int64          `json:"slot_id"`
	TotalPeople int            `json:"total_people"`
	Visitors    []VisitorInput `json:"visitors"`
}

// BookingResponse adalah bentuk data booking yang dikirim ke client
// setelah booking berhasil dibuat.
type BookingResponse struct {
	ID              int64     `json:"id"`
	BookingCode     string    `json:"booking_code"`
	SlotID          int64     `json:"slot_id"`
	TotalPeople     int       `json:"total_people"`
	TotalPrice      float64   `json:"total_price"`
	Status          string    `json:"status"`
	BookedAt        time.Time `json:"booked_at"`
	ExpiredAt       time.Time `json:"expired_at"` // batas waktu bayar sebelum booking otomatis gagal
	DestinationName     string    `json:"destination_name,omitempty"`
	DestinationLocation string    `json:"destination_location,omitempty"`
	DestinationImage    string    `json:"destination_image,omitempty"`
	UserName            string    `json:"user_name,omitempty"`
	UserEmail           string    `json:"user_email,omitempty"`
	PaymentURL          string    `json:"payment_url,omitempty"`
}
