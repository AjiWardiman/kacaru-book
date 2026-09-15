package dto

// RequestOTPResponse dikembalikan setelah sistem berhasil MEMBUAT dan
// MENGIRIM kode OTP. Perhatikan: kode OTP-nya sendiri TIDAK ikut
// dikirim di response ini (sengaja, demi keamanan) — kode yang
// sebenarnya cuma dikirim lewat email. ExpiresInMinutes cuma
// informasi supaya user tahu berapa lama dia harus segera memasukkan kodenya.
type RequestOTPResponse struct {
	Message          string `json:"message"`
	ExpiresInMinutes int    `json:"expires_in_minutes"`
}

// VerifyOTPRequest adalah bentuk data saat user memasukkan kode OTP
// yang diterimanya lewat email, untuk verifikasi POST /otp/verify.
type VerifyOTPRequest struct {
	Code string `json:"code"`
}
