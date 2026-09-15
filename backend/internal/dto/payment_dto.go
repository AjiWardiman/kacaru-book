package dto

import "time"

// PaymentResponse adalah bentuk data pembayaran yang dikirim ke client
// setelah invoice berhasil dibuat.
type PaymentResponse struct {
	ID         int64     `json:"id"`
	BookingID  int64     `json:"booking_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	PaymentURL string    `json:"payment_url"` // link yang harus dibuka user untuk membayar
	CreatedAt  time.Time `json:"created_at"`
}

// XenditWebhookPayload adalah bentuk data yang DIKIRIM Xendit ke
// endpoint webhook kita, setiap kali ada perubahan status invoice
// (misal user berhasil membayar).
//
// Catatan: Xendit sebenarnya mengirim lebih banyak field dari ini,
// tapi kita cukup ambil yang benar-benar kita butuhkan untuk
// memproses perubahan status pembayaran.
type XenditWebhookPayload struct {
	ID            string `json:"id"`             // ID invoice dari Xendit
	ExternalID    string `json:"external_id"`    // ID transaksi versi kita (yang kita kirim saat CreateInvoice)
	Status        string `json:"status"`         // PAID / EXPIRED / dst (Xendit pakai huruf besar)
	PaymentMethod string `json:"payment_method"` // metode yang dipakai user untuk membayar
}
