// Package mailer bertugas mengirim email lewat protokol SMTP.
// Dipakai untuk mengirim notifikasi ke user (misal "booking berhasil
// dibayar") lewat layanan Mailtrap.
//
// KENAPA PAKAI net/smtp BAWAAN GO (bukan library pihak ketiga seperti
// gomail)? Supaya tidak menambah dependency eksternal untuk kebutuhan
// yang sebenarnya sederhana — kirim email teks biasa lewat SMTP. Go
// sudah punya semua yang dibutuhkan untuk ini di package standarnya.
package mailer

import (
	"fmt"
	"net/smtp"
)

// Mailer menyimpan pengaturan koneksi SMTP (host, port, kredensial)
// yang dipakai untuk mengirim email.
type Mailer struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewMailer(host, port, username, password, from string) *Mailer {
	return &Mailer{Host: host, Port: port, Username: username, Password: password, From: from}
}

// Send mengirim satu email teks biasa (plain text) ke satu alamat tujuan.
//
// Cara kerjanya:
//  1. Susun pesan email lengkap dengan header standar (From, To,
//     Subject, Content-Type) — ini format RFC 5322 yang dipahami
//     semua server email.
//  2. Siapkan kredensial autentikasi SMTP (username + password dari Mailtrap).
//  3. Panggil smtp.SendMail bawaan Go, yang otomatis menangani proses
//     koneksi, autentikasi, dan pengiriman ke server SMTP.
func (m *Mailer) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	message := buildMessage(m.From, to, subject, body)

	return smtp.SendMail(addr, auth, m.From, []string{to}, []byte(message))
}

// buildMessage menyusun isi email lengkap dengan header, mengikuti
// format standar email (RFC 5322). "\r\n" (bukan cuma "\n") WAJIB
// dipakai sebagai pemisah baris karena itu aturan protokol SMTP.
func buildMessage(from, to, subject, body string) string {
	return fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, to, subject, body,
	)
}
