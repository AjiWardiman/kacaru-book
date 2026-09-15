package mailer

import (
	"os"
	"testing"
)

func TestBuildMessage(t *testing.T) {
	from := "no-reply@tourism-booking.local"
	to := "user@example.com"
	subject := "Test Notifikasi"
	body := "Halo, ini pesan tes."

	msg := buildMessage(from, to, subject, body)

	expected := "From: no-reply@tourism-booking.local\r\nTo: user@example.com\r\nSubject: Test Notifikasi\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nHalo, ini pesan tes."
	if msg != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, msg)
	}
}

func TestMailtrapIntegration(t *testing.T) {
	username := os.Getenv("MAILTRAP_USERNAME")
	password := os.Getenv("MAILTRAP_PASSWORD")
	if username == "" || password == "" {
		t.Skip("Skipping Mailtrap integration test; set MAILTRAP_USERNAME and MAILTRAP_PASSWORD to run")
	}

	host := "sandbox.smtp.mailtrap.io"
	port := "2525"
	from := "no-reply@tourism-booking.local"

	m := NewMailer(host, port, username, password, from)
	err := m.Send("test-user@tourism-booking.local", "Tes Integrasi Mailtrap", "Halo! Email notifikasi berhasil dikirim melalui Mailtrap Sandbox.")
	if err != nil {
		t.Fatalf("Gagal mengirim email lewat Mailtrap: %v", err)
	}
	t.Log("Email berhasil terkirim ke Mailtrap Sandbox!")
}
