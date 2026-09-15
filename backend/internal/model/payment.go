package model

import "time"

// PaymentStatus adalah status pembayaran, mengikuti status yang
// dikirim balik oleh Xendit lewat webhook.
type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending" // invoice dibuat, menunggu user membayar
	PaymentPaid    PaymentStatus = "paid"    // user sudah berhasil membayar
	PaymentExpired PaymentStatus = "expired" // invoice kadaluarsa karena tidak dibayar
	PaymentFailed  PaymentStatus = "failed"  // pembayaran gagal karena sebab lain
)

// Payment merepresentasikan satu baris di tabel "payments".
// Satu Booking punya TEPAT SATU Payment (relasi one-to-one).
type Payment struct {
	ID            int64
	BookingID     int64
	ExternalID    string  // ID yang KITA kirim ke Xendit untuk mengenali transaksi ini (biasanya = booking_code)
	InvoiceID     *string // ID invoice yang DIBERIKAN Xendit setelah invoice berhasil dibuat
	Amount        float64
	PaymentMethod *string // metode pembayaran yang dipakai user (diisi setelah bayar, dari webhook)
	Status        PaymentStatus
	PaymentURL    *string // link halaman pembayaran Xendit yang dibuka user
	PaidAt        *time.Time
	ExpiredAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// PaymentEvent merepresentasikan satu baris di tabel "payment_events".
// Setiap kali Xendit mengirim webhook (notifikasi perubahan status),
// kita simpan datanya di sini sebagai catatan/log — berguna untuk
// audit dan investigasi kalau ada masalah di kemudian hari.
type PaymentEvent struct {
	ID         int64
	PaymentID  int64
	EventType  string // contoh: "invoice.paid", "invoice.expired"
	Payload    string // data mentah (JSON) yang dikirim Xendit, disimpan apa adanya
	ReceivedAt time.Time
}

func (Payment) TableName() string {
	return "payments"
}

func (PaymentEvent) TableName() string {
	return "payment_events"
}

