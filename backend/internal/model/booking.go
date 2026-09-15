package model

import "time"

// BookingStatus adalah status sebuah booking sepanjang "hidupnya",
// dari awal dibuat sampai selesai (dibayar) atau gagal (expired/dibatalkan).
type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"   // baru dibuat, menunggu user membayar
	BookingPaid      BookingStatus = "paid"      // sudah dibayar dan dikonfirmasi (dipakai di tahap integrasi Xendit)
	BookingExpired   BookingStatus = "expired"   // tidak dibayar sampai batas waktu, otomatis dianggap gagal
	BookingCancelled BookingStatus = "cancelled" // dibatalkan (misal oleh user sendiri)
)

// Booking merepresentasikan satu baris di tabel "bookings".
// Satu Booking bisa untuk BEBERAPA ORANG sekaligus (rombongan) — data
// tiap orangnya disimpan terpisah di BookingDetail.
type Booking struct {
	ID          int64
	BookingCode string        // kode yang ditampilkan ke user, misal "BK-20260920-A1B2C3"
	UserID      int64         // siapa (user mana) yang melakukan booking ini
	SlotID      int64         // slot/jadwal destinasi mana yang dipesan
	TotalPeople int           // jumlah orang dalam booking ini
	TotalPrice  float64       // total harga = harga destinasi x jumlah orang
	Status      BookingStatus // lihat konstanta di atas
	BookedAt    time.Time     // waktu booking ini dibuat
	ExpiredAt   time.Time     // batas waktu bayar — lewat dari ini, booking otomatis dianggap gagal
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// BookingDetail merepresentasikan data SATU pengunjung dalam sebuah
// booking. Kalau booking untuk 3 orang, akan ada 3 baris BookingDetail
// dengan BookingID yang sama.
type BookingDetail struct {
	ID           int64
	BookingID    int64   // booking mana yang punya pengunjung ini
	VisitorName  string
	VisitorPhone *string // boleh kosong
	VisitorEmail *string // boleh kosong
	Age          *int    // boleh kosong
	CreatedAt    time.Time
}

func (Booking) TableName() string {
	return "bookings"
}

func (BookingDetail) TableName() string {
	return "booking_details"
}

// BookingWithDetails merepresentasikan data gabungan booking dengan user dan destinasi
type BookingWithDetails struct {
	ID              int64         `gorm:"column:id"`
	BookingCode     string        `gorm:"column:booking_code"`
	UserID          int64         `gorm:"column:user_id"`
	SlotID          int64         `gorm:"column:slot_id"`
	TotalPeople     int           `gorm:"column:total_people"`
	TotalPrice      float64       `gorm:"column:total_price"`
	Status          BookingStatus `gorm:"column:status"`
	BookedAt        time.Time     `gorm:"column:booked_at"`
	ExpiredAt       time.Time     `gorm:"column:expired_at"`
	DestinationName     string        `gorm:"column:destination_name"`
	DestinationLocation string        `gorm:"column:destination_location"`
	DestinationImage    string        `gorm:"column:destination_image"`
	UserName            string        `gorm:"column:user_name"`
	UserEmail           string        `gorm:"column:user_email"`
}

