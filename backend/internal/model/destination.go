package model

import "time"

// Category merepresentasikan satu baris di tabel "categories".
// Ini adalah kategori destinasi wisata, misalnya "Pantai", "Gunung", "Budaya", dll.
type Category struct {
	ID          int64
	Name        string
	Description *string // boleh kosong (nil), tidak semua kategori butuh deskripsi
	IsActive    bool    // kalau false, kategori ini tidak ditampilkan lagi tapi datanya tetap ada
}

// Destination merepresentasikan satu baris di tabel "destinations".
// Ini adalah data destinasi wisata (misal: Pantai Lariti, Kabupaten Bima).
type Destination struct {
	ID          int64
	CreatedBy   int64     // ID user (admin) yang membuat destinasi ini, mengacu ke tabel users
	Name        string    // nama destinasi, misal "Pantai Lariti"
	Location    string    // lokasi/alamat destinasi
	Description *string   // deskripsi destinasi, boleh kosong
	Category    *string   // kategori destinasi (versi simpel: teks bebas, bukan relasi tabel)
	Price       float64   // harga tiket masuk per orang
	ImageURL    *string   // link gambar destinasi, boleh kosong kalau belum ada foto
	IsActive    bool      // kalau false, destinasi disembunyikan dari daftar publik (soft delete)
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Slot merepresentasikan satu baris di tabel "slots".
// Slot adalah "jadwal kunjungan" untuk destinasi tertentu di tanggal tertentu,
// lengkap dengan kuota (batas jumlah pengunjung) di tanggal itu.
//
// Slot inilah yang jadi kunci dari fitur "anti-bentrok": setiap kali ada
// booking baru, sistem akan mengecek dan mengunci baris slot ini supaya
// tidak ada dua booking yang "berebut" kuota yang sama secara bersamaan.
type Slot struct {
	ID            int64
	DestinationID int64     // destinasi mana yang punya slot ini
	Date          time.Time // tanggal kunjungan
	TotalQuota    int       // jumlah maksimal pengunjung yang boleh booking di tanggal ini
	UsedQuota     int       // jumlah pengunjung yang SUDAH booking (bertambah tiap ada booking baru)
	IsAvailable   bool      // kalau false, slot ini ditutup manual oleh admin (misal karena cuaca buruk)
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (Category) TableName() string {
	return "categories"
}

func (Destination) TableName() string {
	return "destinations"
}

func (Slot) TableName() string {
	return "slots"
}

