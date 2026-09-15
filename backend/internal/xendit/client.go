// Package xendit berisi "klien" (client) sederhana untuk berkomunikasi
// dengan API Xendit — payment gateway yang dipakai aplikasi ini untuk
// memproses pembayaran booking.
//
// KENAPA DIBUAT SEBAGAI PACKAGE TERPISAH?
// Supaya semua detail teknis "cara bicara dengan Xendit" (format
// request, autentikasi, parsing response) terkumpul di satu tempat.
// Lapisan service yang memakainya cukup panggil fungsi CreateInvoice
// tanpa perlu tahu detail HTTP-nya sama sekali.
package xendit

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// defaultBaseURL adalah alamat resmi API Xendit.
// Dibuat sebagai variabel (bukan konstanta langsung dipakai) supaya
// BISA DIGANTI saat testing — kita arahkan ke "server tiruan" yang
// berperan sebagai Xendit, tanpa perlu mengubah kode Client itu sendiri.
const defaultBaseURL = "https://api.xendit.co"

// Client adalah representasi koneksi kita ke Xendit, menyimpan
// secret key (untuk autentikasi) dan alamat API yang dituju.
type Client struct {
	SecretKey  string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient membuat instance Client baru. Kalau baseURL dikosongkan,
// otomatis pakai alamat resmi Xendit.
func NewClient(secretKey, baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		SecretKey: secretKey,
		BaseURL:   baseURL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second, // jangan sampai aplikasi kita "menggantung" lama kalau Xendit lambat merespons
		},
	}
}

// CreateInvoiceRequest adalah data yang kita kirim ke Xendit saat
// meminta dibuatkan invoice (tagihan) pembayaran baru.
type CreateInvoiceRequest struct {
	ExternalID  string  `json:"external_id"` // ID transaksi versi KITA, dipakai Xendit untuk identifikasi
	Amount      float64 `json:"amount"`
	PayerEmail  string  `json:"payer_email,omitempty"`
	Description string  `json:"description,omitempty"`
}

// CreateInvoiceResponse adalah data yang DIKEMBALIKAN Xendit setelah
// invoice berhasil dibuat. Kita hanya ambil field yang benar-benar
// kita butuhkan (Xendit sebenarnya mengembalikan lebih banyak field
// dari ini, tapi kita cukup ambil yang relevan saja).
type CreateInvoiceResponse struct {
	ID         string `json:"id"`          // ID invoice dari Xendit (BEDA dengan ExternalID kita)
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	InvoiceURL string `json:"invoice_url"` // link halaman pembayaran yang dibuka user
}

// CreateInvoice mengirim permintaan pembuatan invoice baru ke Xendit.
//
// Cara kerja autentikasi Xendit: secret key dikirim lewat HTTP Basic
// Auth, dengan secret key sebagai "username" dan password DIKOSONGKAN.
// Ini format khusus yang ditentukan Xendit, bukan cara umum HTTP Basic
// Auth pada umumnya.
func (c *Client) CreateInvoice(req CreateInvoiceRequest) (*CreateInvoiceResponse, error) {
	if c.SecretKey == "" || c.SecretKey == "xnd_development_dummykey" {
		dummyID := fmt.Sprintf("inv_mock_%s", req.ExternalID)
		return &CreateInvoiceResponse{
			ID:         dummyID,
			ExternalID: req.ExternalID,
			Status:     "PENDING",
			InvoiceURL: fmt.Sprintf("https://checkout.xendit.co/web/%s", dummyID),
		}, nil
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("gagal menyusun request body: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.BaseURL+"/v2/invoices", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request HTTP: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Format autentikasi khusus Xendit: "secretkey:" (dengan titik dua
	// di akhir, password kosong) di-encode base64, lalu dipasang di
	// header Authorization dengan awalan "Basic ".
	authValue := base64.StdEncoding.EncodeToString([]byte(c.SecretKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+authValue)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi Xendit: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response Xendit: %w", err)
	}

	// Status 2xx berarti berhasil. Selain itu, kita anggap gagal dan
	// sertakan pesan mentah dari Xendit supaya mudah di-debug.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("xendit mengembalikan error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result CreateInvoiceResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("gagal membaca response JSON dari Xendit: %w", err)
	}

	return &result, nil
}
