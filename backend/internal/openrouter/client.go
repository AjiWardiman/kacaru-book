// Package openrouter bertugas mengelola koneksi dan komunikasi ke OpenRouter AI API.
// Paket ini mengirimkan prompt analisis yang menggabungkan prakiraan cuaca OpenWeatherMap,
// data destinasi wisata PostgreSQL, dan preferensi wisatawan ke model Claude (OpenRouter).
package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client menyimpan kredensial dan konfigurasi untuk menghubungi OpenRouter API.
type Client struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

// WeatherContext menyimpan data ramalan cuaca dari OpenWeatherMap.
type WeatherContext struct {
	Location    string  `json:"location"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
	SafetyNote  string  `json:"safety_note"`
}

// DestinationContext menyimpan data rinci destinasi wisata dari PostgreSQL.
type DestinationContext struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Location    string  `json:"location"`
	Description string  `json:"description"`
}

// UserPreference menyimpan kriteria pilihan wisatawan.
type UserPreference struct {
	MaxBudget  float64 `json:"max_budget"`
	Category   string  `json:"category"`
	Guests     int     `json:"guests"`
	Preference string  `json:"preference"`
}

// AIRecommendationItem adalah item rekomendasi dari respons AI.
type AIRecommendationItem struct {
	DestinationID int64  `json:"destination_id"`
	MatchScore    int    `json:"match_score"`
	Reason        string `json:"reason"`
}

// NewClient membuat instance baru OpenRouter Client.
func NewClient(apiKey, model string) *Client {
	if model == "" {
		model = "anthropic/claude-3-haiku"
	}
	return &Client{
		apiKey:  strings.TrimSpace(apiKey),
		model:   strings.TrimSpace(model),
		baseURL: "https://openrouter.ai/api/v1/chat/completions",
		httpClient: &http.Client{
			Timeout: 4 * time.Second,
		},
	}
}

// IsEnabled mengecek apakah API key OpenRouter sudah diisi.
func (c *Client) IsEnabled() bool {
	return c.apiKey != ""
}

// openRouterMessage mewakili struktur pesan chat completion.
type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// openRouterRequest mewakili payload ke endpoint /chat/completions.
type openRouterRequest struct {
	Model       string              `json:"model"`
	Messages    []openRouterMessage `json:"messages"`
	Temperature float64             `json:"temperature"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
}

// openRouterResponse mewakili respons dari OpenRouter.
type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

// GenerateRecommendations memanggil model Claude via OpenRouter API untuk menganalisis dan
// memberikan rekomendasi destinasi berbasis kondisi cuaca dan preferensi wisatawan.
func (c *Client) GenerateRecommendations(
	ctx context.Context,
	weather WeatherContext,
	destinations []DestinationContext,
	pref UserPreference,
) ([]AIRecommendationItem, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("openrouter API key is not configured")
	}

	// Ringkas data destinasi agar payload token kecil & eksekusi AI super cepat (< 1 detik)
	type CompactDest struct {
		ID       int64   `json:"id"`
		Name     string  `json:"name"`
		Category string  `json:"category"`
		Price    float64 `json:"price"`
		Location string  `json:"location"`
	}

	compactList := make([]CompactDest, 0, len(destinations))
	for _, d := range destinations {
		compactList = append(compactList, CompactDest{
			ID:       d.ID,
			Name:     d.Name,
			Category: d.Category,
			Price:    d.Price,
			Location: d.Location,
		})
	}

	destJSON, _ := json.Marshal(compactList)
	weatherJSON, _ := json.Marshal(weather)

	prompt := fmt.Sprintf(`Tugas: Berikan 5 rekomendasi destinasi terbaik dalam JSON.

CUACA: %s
PREFERENSI: Kat=%s, MaxBudget=Rp%.0f, Rombongan=%d

DESTINASI:
%s

INSTRUKSI:
1. Pilih MAKSIMAL 5 destinasi paling cocok.
2. Alasan (reason) RINGKAS SANGAT SINGKAT (maksimal 10-12 kata).
3. HANYA OUTPUT JSON ARRAY VALID TANPA FORMATTING MARKDOWN/TEKS LAIN.

Format JSON:
[{"destination_id":1,"match_score":95,"reason":"Alasan singkat..."}]`,
		string(weatherJSON),
		pref.Category,
		pref.MaxBudget,
		pref.Guests,
		string(destJSON),
	)

	reqPayload := openRouterRequest{
		Model: c.model,
		Messages: []openRouterMessage{
			{Role: "user", Content: prompt},
		},
		Temperature: 0.2,
		MaxTokens:   300,
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("HTTP-Referer", "https://tourism-booking.local")
	httpReq.Header.Set("X-Title", "Tourism Booking System Bima")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openrouter request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openrouter api status %d: %s", resp.StatusCode, string(respBody))
	}

	var openRouterResp openRouterResponse
	if err := json.Unmarshal(respBody, &openRouterResp); err != nil {
		return nil, fmt.Errorf("failed to parse openrouter response json: %w", err)
	}

	if openRouterResp.Error != nil {
		return nil, fmt.Errorf("openrouter error: %s", openRouterResp.Error.Message)
	}

	if len(openRouterResp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned from openrouter")
	}

	rawContent := strings.TrimSpace(openRouterResp.Choices[0].Message.Content)
	// Bersihkan markdown code block jika model menyertakan ```json ... ```
	if strings.HasPrefix(rawContent, "```") {
		lines := strings.Split(rawContent, "\n")
		if len(lines) >= 2 {
			if strings.HasPrefix(lines[0], "```") {
				lines = lines[1:]
			}
			if len(lines) > 0 && strings.HasPrefix(lines[len(lines)-1], "```") {
				lines = lines[:len(lines)-1]
			}
			rawContent = strings.TrimSpace(strings.Join(lines, "\n"))
		}
	}

	var items []AIRecommendationItem
	if err := json.Unmarshal([]byte(rawContent), &items); err != nil {
		return nil, fmt.Errorf("failed to parse AI structured json output: %w (raw: %s)", err, rawContent)
	}

	return items, nil
}
