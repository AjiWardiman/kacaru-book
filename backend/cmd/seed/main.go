package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	"tourism-booking/internal/config"
	"tourism-booking/internal/model"
	"tourism-booking/pkg/database"
)

func main() {
	log.Println("🌱 Memulai eksekusi seeder database Super Ramai (Kabupaten Bima)...")

	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke database PostgreSQL: %v", err)
	}

	log.Println("🔌 Terhubung ke database PostgreSQL. Memproses pembersihan & pembuatan seeder...")

	// 1. Hash default password ("password123")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Gagal membuat hash password: %v", err)
	}

	now := time.Now()

	// 2. Seeder USERS (1 Admin + 24 Customers = 25 Users)
	users := []model.User{
		{
			Name:            "Admin Aji",
			Email:           "aji@example.com",
			Phone:           ptrString("081234567890"),
			PhoneVerifiedAt: &now,
			Password:        string(hashedPassword),
			Role:            model.RoleAdmin,
		},
	}

	customerNames := []string{
		"Budi Santoso", "Siti Rahmawati", "Dewi Lestari", "Andi Wijaya", "Eko Prasetyo",
		"Maya Indah", "Rudi Hartono", "Ahmad Rizky", "Fajar Nugraha", "Nia Kurnia",
		"Hendra Gunawan", "Dian Sastro", "Rahmat Hidayat", "Fitriani", "Agus Setiawan",
		"Putri Wulandari", "Bambang Pamungkas", "Tari Anggraini", "Irwan Syah", "Lestari Mbojo",
		"Syarifudin Bima", "Nurhaliza", "Zulkarnain", "Kiki Amalia",
	}

	for i, name := range customerNames {
		email := fmt.Sprintf("user%d@example.com", i+1)
		switch i {
		case 0:
			email = "budi@example.com"
		case 1:
			email = "siti@example.com"
		}
		phone := fmt.Sprintf("0899%08d", i+1000)

		users = append(users, model.User{
			Name:            name,
			Email:           email,
			Phone:           ptrString(phone),
			PhoneVerifiedAt: &now,
			Password:        string(hashedPassword),
			Role:            model.RoleUser,
		})
	}

	for i := range users {
		db.Where("email = ?", users[i].Email).FirstOrCreate(&users[i])
	}
	log.Printf("✅ Seeder 25 Akun Users (1 Admin + 24 Customers) berhasil disiapkan!")

	// 3. Seeder CATEGORIES
	categories := []struct {
		Name        string
		Description string
	}{
		{"Bahari", "Wisata bahari, pantai eksotis, dan keindahan pulau di Kabupaten Bima"},
		{"Gunung & Alam", "Wisata pegunungan, gunung berapi, dan air terjun alami"},
		{"Budaya & Geowisata", "Situs sejarah kebudayaan Mbojo dan keajaiban batuan geologi"},
	}

	for _, cat := range categories {
		var existingCount int64
		db.Table("categories").Where("name = ?", cat.Name).Count(&existingCount)
		if existingCount == 0 {
			db.Exec("INSERT INTO categories (name, description, is_active) VALUES (?, ?, true)", cat.Name, cat.Description)
		}
	}
	log.Printf("✅ Seeder 3 Kategori Wisata berhasil disiapkan!")

	// 4. Seeder DESTINATIONS (19 Destinasi Kabupaten Bima, NTB)
	adminID := users[0].ID

	destinationsData := []model.Destination{
		{
			CreatedBy:   adminID,
			Name:        "Pantai Lariti",
			Location:    "Kecamatan Sape, Kabupaten Bima, NTB",
			Description: ptrString("Pantai unik eksotis dengan fenomena laut membelah dan pasir timbul yang memukau saat laut surut."),
			Category:    ptrString("Bahari"),
			Price:       15000,
			ImageURL:    ptrString("/uploads/destinations/lariti.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Toro Maria",
			Location:    "Kecamatan Sape, Kabupaten Bima, NTB",
			Description: ptrString("Tebing karang memukau dengan lanskap bukit dan pemandangan lautan lepas perbatasan Flores."),
			Category:    ptrString("Bahari"),
			Price:       20000,
			ImageURL:    ptrString("/uploads/destinations/toro_maria.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Uma Lengge Wawo",
			Location:    "Kecamatan Wawo, Kabupaten Bima, NTB",
			Description: ptrString("Situs cagar budaya rumah adat kuno suku Bima (Mbojo) dengan struktur kerucut yang kaya sejarah."),
			Category:    ptrString("Budaya & Geowisata"),
			Price:       10000,
			ImageURL:    ptrString("/uploads/destinations/uma_lengge.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Pulau Sangiang",
			Location:    "Kecamatan Wera, Kabupaten Bima, NTB",
			Description: ptrString("Pulau gunung berapi aktif di tengah lautan Bima dengan keindahan alam bahari dan savana pesisir."),
			Category:    ptrString("Gunung & Alam"),
			Price:       50000,
			ImageURL:    ptrString("/uploads/destinations/sangiang.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Air Terjun Oi Fanda",
			Location:    "Kecamatan Ambalawi, Kabupaten Bima, NTB",
			Description: ptrString("Air terjun alami berair jernih segar yang tersembunyi di rimbunnya hutan lereng pegunungan Ambalawi."),
			Category:    ptrString("Gunung & Alam"),
			Price:       15000,
			ImageURL:    ptrString("/uploads/destinations/oi_fanda.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Pantai Pink Lambu",
			Location:    "Kecamatan Lambu, Kabupaten Bima, NTB",
			Description: ptrString("Hamparan pasir pantai berwarna merah muda alami dengan perairan jernih tenang nan fotogenik."),
			Category:    ptrString("Bahari"),
			Price:       25000,
			ImageURL:    ptrString("/uploads/destinations/pantai_pink.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Pulau Kelapa",
			Location:    "Kecamatan Lambu, Kabupaten Bima, NTB",
			Description: ptrString("Gugusan pulau eksotis nan sepi dengan lanskap perbukitan hijau dan bukit pandang bahari mempesona."),
			Category:    ptrString("Bahari"),
			Price:       75000,
			ImageURL:    ptrString("/uploads/destinations/pulau_kelapa.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Tanjung Meriam",
			Location:    "Kecamatan Lambu, Kabupaten Bima, NTB",
			Description: ptrString("Keajaiban geologi batuan basal tersusun rapi segiempat/segienam alami menyerupai monumen meriam."),
			Category:    ptrString("Budaya & Geowisata"),
			Price:       30000,
			ImageURL:    ptrString("/uploads/destinations/tanjung_meriam.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Nisa Ndoko",
			Location:    "Kecamatan Langgudu, Kabupaten Bima, NTB",
			Description: ptrString("Pulau kecil (Nisa) berpasir putih bersih dengan keindahan bawah laut terumbu karang yang terjaga."),
			Category:    ptrString("Bahari"),
			Price:       40000,
			ImageURL:    ptrString("/uploads/destinations/nisa_ndoko.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Pantai Kampa",
			Location:    "Kecamatan Sape, Kabupaten Bima, NTB",
			Description: ptrString("Pantai teluk tersembunyi berpasir putih jernih dengan air laut tenang dan pemandangan bukit karang eksotis."),
			Category:    ptrString("Bahari"),
			Price:       20000,
			ImageURL:    ptrString("/uploads/destinations/pantai_kampa.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Air Terjun Pawa",
			Location:    "Kecamatan Langgudu, Kabupaten Bima, NTB",
			Description: ptrString("Air terjun bertingkat alami di tengah rimbunnya hutan tropis pegunungan Langgudu berair jernih segar."),
			Category:    ptrString("Gunung & Alam"),
			Price:       15000,
			ImageURL:    ptrString("/uploads/destinations/air_terjun_pawa.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Gunung Tambora",
			Location:    "Kabupaten Bima - Dompu, NTB",
			Description: ptrString("Kaldera gunung berapi raksasa legendaris dunia dengan lanskap padang savana dan panorama samudra luas."),
			Category:    ptrString("Gunung & Alam"),
			Price:       50000,
			ImageURL:    ptrString("/uploads/destinations/gunung_tambora.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Pantai Kalaki",
			Location:    "Kecamatan Palibelo, Kabupaten Bima, NTB",
			Description: ptrString("Pesisir pantai favorit sunset wisatawan lokal dengan lanskap Teluk Bima dan kuliner pesisir khas Mbojo."),
			Category:    ptrString("Bahari"),
			Price:       10000,
			ImageURL:    ptrString("/uploads/destinations/pantai_kalaki.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Bukit Matompo",
			Location:    "Kecamatan Wawo, Kabupaten Bima, NTB",
			Description: ptrString("Puncak perbukitan hijau tempat menyaksikan hamparan awan pagi dan lanskap lembah pegunungan Wawo."),
			Category:    ptrString("Gunung & Alam"),
			Price:       15000,
			ImageURL:    ptrString("/uploads/destinations/bukit_matompo.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Mata Air Madapangga",
			Location:    "Kecamatan Madapangga, Kabupaten Bima, NTB",
			Description: ptrString("Pemandian alam sumber air pegunungan yang sangat jernih, sejuk, dan dikelilingi pepohonan rimbun alami."),
			Category:    ptrString("Gunung & Alam"),
			Price:       10000,
			ImageURL:    ptrString("/uploads/destinations/madapangga.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Pantai Wane",
			Location:    "Kecamatan Monta, Kabupaten Bima, NTB",
			Description: ptrString("Pantai pesisir selatan berombak eksotis dengan hamparan pasir putih dan atraksi lumba-lumba liar."),
			Category:    ptrString("Bahari"),
			Price:       20000,
			ImageURL:    ptrString("/uploads/destinations/pantai_wane.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Pulau Utop",
			Location:    "Kecamatan Sape, Kabupaten Bima, NTB",
			Description: ptrString("Pulau nirwana mungil berpasir halus dengan keindahan taman laut bawah air yang ideal untuk snorkeling."),
			Category:    ptrString("Bahari"),
			Price:       60000,
			ImageURL:    ptrString("/uploads/destinations/pulau_utop.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Istana Asi Mbojo",
			Location:    "Kota / Kabupaten Bima, NTB",
			Description: ptrString("Museum cagar budaya bekas Istana Kesultanan Bima dengan koleksi mahkota emas dan sejarah Mbojo."),
			Category:    ptrString("Budaya & Geowisata"),
			Price:       15000,
			ImageURL:    ptrString("/uploads/destinations/asi_mbojo.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Pantai Kiri Wera",
			Location:    "Kecamatan Wera, Kabupaten Bima, NTB",
			Description: ptrString("Pantai eksotis berpasir hitam berkilau khas lereng gunung api Sangiang dengan panorama perahu tradisional."),
			Category:    ptrString("Bahari"),
			Price:       15000,
			ImageURL:    ptrString("/uploads/destinations/pantai_kiri.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Nisa Dora",
			Location:    "Kecamatan Sape, Kabupaten Bima, NTB",
			Description: ptrString("Pulau nirwana kecil berpasir putih jernih untuk tempat bersantai, berenang, dan aktivitas diving."),
			Category:    ptrString("Bahari"),
			Price:       35000,
			ImageURL:    ptrString("/uploads/destinations/nisa_dora.jpg"),
			IsActive:    true,
		},
		{
			CreatedBy:   adminID,
			Name:        "Toro Jambu",
			Location:    "Kecamatan Lambu, Kabupaten Bima, NTB",
			Description: ptrString("Tanjung bukit perbukitan hijau nan megah berhadapan dengan samudra lepas, tempat favorit sunrise & camping."),
			Category:    ptrString("Bahari"),
			Price:       20000,
			ImageURL:    ptrString("/uploads/destinations/toro_jambu.jpg"),
			IsActive:    true,
		},
	}

	var savedDestinations []model.Destination
	for _, dest := range destinationsData {
		var existing model.Destination
		err := db.Where("name = ?", dest.Name).First(&existing).Error
		if err == nil {
			savedDestinations = append(savedDestinations, existing)
		} else {
			db.Create(&dest)
			savedDestinations = append(savedDestinations, dest)
		}
	}
	log.Printf("✅ Seeder %d Destinasi Wisata Kabupaten Bima berhasil disiapkan!", len(savedDestinations))

	// 5. Seeder SLOTS (36 Slot Kunjungan: 4 Tanggal per Destinasi)
	dates := []string{
		"2026-09-20", "2026-09-25", "2026-10-02", "2026-10-10",
		"2026-10-18", "2026-10-25", "2026-11-01", "2026-11-15",
	}

	var savedSlots []model.Slot
	for destIdx, dest := range savedDestinations {
		for i := 0; i < 4; i++ {
			dateStr := dates[(destIdx+i)%len(dates)]
			tDate, _ := time.Parse("2006-01-02", dateStr)

			quota := 50
			if dest.Name == "Pantai Lariti" || dest.Name == "Pantai Pink Lambu" {
				quota = 100
			}

			slot := model.Slot{
				DestinationID: dest.ID,
				Date:          tDate,
				TotalQuota:    quota,
				UsedQuota:     0,
				IsAvailable:   true,
			}

			var existingSlot model.Slot
			err := db.Where("destination_id = ? AND date = ?", dest.ID, tDate).First(&existingSlot).Error
			if err == nil {
				savedSlots = append(savedSlots, existingSlot)
			} else {
				db.Create(&slot)
				savedSlots = append(savedSlots, slot)
			}
		}
	}
	log.Printf("✅ Seeder 36 Slot Kunjungan (Jadwal & Kuota) berhasil disiapkan!")

	// 6. Seeder BOOKINGS & PAYMENTS (50 Transaksi Pemesanan & Pembayaran Super Ramai 5 Bulan Terakhir: Mei - September 2026)
	db.Exec("DELETE FROM booking_details")
	db.Exec("DELETE FROM payments")
	db.Exec("DELETE FROM bookings")
	db.Exec("UPDATE slots SET used_quota = 0")

	bookingStatuses := []model.BookingStatus{
		model.BookingPaid, model.BookingPaid, model.BookingPaid, model.BookingPaid,
		model.BookingPaid, model.BookingPaid, model.BookingPaid, model.BookingPaid,
		model.BookingPending, model.BookingCancelled,
	}

	r := rand.New(rand.NewSource(42))

	targetMonths := []time.Month{
		time.May, time.June, time.July, time.August, time.September,
	}

	paidCount := 0
	pendingCount := 0
	cancelledCount := 0

	bookingCounter := 0
	for _, m := range targetMonths {
		// 10 transaksi per bulan = 50 total transaksi 5 bulan terakhir
		for k := 1; k <= 10; k++ {
			bookingCounter++
			customer := users[1+((bookingCounter-1)%(len(users)-1))]
			slot := savedSlots[bookingCounter%len(savedSlots)]
			dest := savedDestinations[bookingCounter%len(savedDestinations)]

			st := bookingStatuses[(k-1)%len(bookingStatuses)]

			totalPeople := 1 + r.Intn(4)
			totalPrice := dest.Price * float64(totalPeople)

			dayInt := 2 + (k * 2)
			if dayInt > 28 {
				dayInt = 28
			}

			bookedAt := time.Date(2026, m, dayInt, 9+(k%8), (k*7)%60, 0, 0, time.UTC)
			expiredAt := bookedAt.Add(30 * time.Minute)

			bookingCode := fmt.Sprintf("BK-2026%02d%02d-%03d", int(m), dayInt, bookingCounter)

			b := model.Booking{
				BookingCode: bookingCode,
				UserID:      customer.ID,
				SlotID:      slot.ID,
				TotalPeople: totalPeople,
				TotalPrice:  totalPrice,
				Status:      st,
				BookedAt:    bookedAt,
				ExpiredAt:   expiredAt,
				CreatedAt:   bookedAt,
				UpdatedAt:   bookedAt,
			}

			if err := db.Create(&b).Error; err == nil {
				switch st {
				case model.BookingPaid:
					db.Model(&slot).Update("used_quota", slot.UsedQuota+totalPeople)
					paidCount++
				case model.BookingPending:
					pendingCount++
				default:
					cancelledCount++
				}

				for p := 1; p <= totalPeople; p++ {
					visitorName := customer.Name
					if p > 1 {
						visitorName = fmt.Sprintf("%s (Rombongan %d)", customer.Name, p)
					}
					detail := model.BookingDetail{
						BookingID:    b.ID,
						VisitorName:  visitorName,
						VisitorPhone: customer.Phone,
						VisitorEmail: ptrString(customer.Email),
						Age:          ptrInt(20 + r.Intn(25)),
					}
					db.Create(&detail)
				}

				extID := fmt.Sprintf("INV-%s", bookingCode)
				payStatus := model.PaymentPending
				switch st {
				case model.BookingPaid:
					payStatus = model.PaymentPaid
				case model.BookingCancelled:
					payStatus = model.PaymentFailed
				}

				pay := model.Payment{
					BookingID:     b.ID,
					ExternalID:    extID,
					InvoiceID:     ptrString(fmt.Sprintf("inv_xnd_%06d", bookingCounter*1024)),
					Amount:        totalPrice,
					PaymentMethod: ptrString("BANK_TRANSFER"),
					Status:        payStatus,
					PaymentURL:    ptrString(fmt.Sprintf("https://checkout.xendit.co/web/inv_xnd_%06d", bookingCounter*1024)),
					CreatedAt:     bookedAt,
					UpdatedAt:     bookedAt,
				}
				if st == model.BookingPaid {
					paidAt := bookedAt.Add(10 * time.Minute)
					pay.PaidAt = &paidAt
				}
				db.Create(&pay)
			}
		}
	}

	log.Printf("✅ Seeder 30 Transaksi Pemesanan Super Ramai (Lunas: %d, Pending: %d, Batal: %d) berhasil dieksekusi!", paidCount, pendingCount, cancelledCount)

	// 6. Seeder NOTIFICATIONS (Notifikasi Realistis)
	for _, u := range users {
		var uNotifCount int64
		db.Model(&model.Notification{}).Where("user_id = ?", u.ID).Count(&uNotifCount)
		if uNotifCount == 0 {
			if u.Role == model.RoleAdmin {
				db.Create(&model.Notification{
					UserID:  u.ID,
					Type:    "system_alert",
					Title:   "Selamat Datang Admin di Kacaru Book Bima",
					Message: "Sistem reservasi wisata Kabupaten Bima telah aktif. Anda dapat mengelola destinasi, slot kuota, dan memantau transaksi.",
					IsRead:  false,
					SentAt:  now.Add(-2 * time.Hour),
				})
				db.Create(&model.Notification{
					UserID:  u.ID,
					Type:    "booking_paid",
					Title:   "Pembayaran Berhasil Diterima",
					Message: "Transaksi pemesanan tiket Pantai Lariti senilai Rp 30.000 telah lunas via Xendit Bank Transfer.",
					IsRead:  false,
					SentAt:  now.Add(-1 * time.Hour),
				})
			} else {
				db.Create(&model.Notification{
					UserID:  u.ID,
					Type:    "booking_paid",
					Title:   "Pembayaran Berhasil! Tiket Elektronik Siap",
					Message: fmt.Sprintf("Halo %s, pembayaran pemesanan tiket wisata Kabupaten Bima telah terkonfirmasi lunas. Selamat berwisata!", u.Name),
					IsRead:  false,
					SentAt:  now.Add(-30 * time.Minute),
				})
				db.Create(&model.Notification{
					UserID:  u.ID,
					Type:    "reminder",
					Title:   "Pengingat Jadwal Wisata Bima",
					Message: fmt.Sprintf("Halo %s, pastikan membawa kode booking saat mengunjungi destinasi wisata. Nikmati pesona alam Mbojo!", u.Name),
					IsRead:  true,
					SentAt:  now.Add(-10 * time.Minute),
				})
			}
		}
	}

	log.Println("🎉 EKSEKUSI SEEDER DATABASE KABUPATEN BIMA SELESAI DENGAN SUKSES LENGKAP!")
}

func ptrString(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}
