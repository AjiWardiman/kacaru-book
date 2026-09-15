package dto

// APIResponse adalah "amplop" standar yang membungkus SEMUA response API
// di aplikasi ini, baik yang berhasil maupun yang gagal. Tujuannya supaya
// client (frontend/aplikasi lain) selalu tahu bentuk response yang konsisten,
// tidak perlu menebak-nebak struktur berbeda di tiap endpoint.
//
// Contoh response sukses:
//   {"success": true, "message": "berhasil login", "data": {...}}
//
// Contoh response gagal:
//   {"success": false, "error": "email atau password salah"}
type APIResponse struct {
	Success bool        `json:"success"`           // true kalau request berhasil, false kalau gagal
	Message string      `json:"message,omitempty"` // pesan sukses, kosong kalau gagal
	Data    interface{} `json:"data,omitempty"`    // isi data hasil (bisa berupa apa saja: object, array, dll)
	Error   string      `json:"error,omitempty"`   // pesan error, kosong kalau sukses
}
