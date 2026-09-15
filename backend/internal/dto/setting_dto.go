package dto

import "time"

// SettingResponse adalah bentuk data pengaturan aplikasi yang dikirim
// ke client (admin).
type SettingResponse struct {
	ID          int64     `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateSettingRequest adalah bentuk data yang diharapkan saat admin
// mengubah nilai satu pengaturan lewat PUT /admin/settings/:key.
// Cuma butuh "value" karena key-nya sendiri sudah ada di URL, dan
// description tidak boleh diubah lewat API (sifatnya cuma dokumentasi
// bawaan dari migrasi database, bukan sesuatu yang perlu admin ubah).
type UpdateSettingRequest struct {
	Value string `json:"value"`
}
