# ==============================================================================
# Multi-Stage Dockerfile untuk Backend Golang (Tourism Booking System)
# ==============================================================================

# ------------------------------------------------------------------------------
# Stage 1: Kompilasi Binary Go (Versi Go sesuai go.mod: 1.25)
# ------------------------------------------------------------------------------
ARG GO_VERSION=1.25
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

# Pengaturan toolchain otomatis jika dibutuhkan
ENV GOTOOLCHAIN=auto

# Install paket penting untuk compile dan download dependensi
RUN apk add --no-cache git ca-certificates tzdata

# Cache modul dependensi Go agar build berikutnya lebih cepat
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh source code backend
COPY . .

# Kompilasi binary API server dan database seeder secara statis tanpa CGO
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/seed ./cmd/seed

# ------------------------------------------------------------------------------
# Stage 2: Runtime Image Minimal (Alpine Linux)
# ------------------------------------------------------------------------------
FROM alpine:3.21

WORKDIR /app

# Pasang sertifikat SSL publik dan data zona waktu
RUN apk add --no-cache ca-certificates tzdata

# Salin binary hasil kompilasi dari stage builder
COPY --from=builder /app/api /app/api
COPY --from=builder /app/seed /app/seed

# Buat folder penyimpanan uploads file
RUN mkdir -p /app/uploads

# Port server REST API
EXPOSE 8080

# Jalankan server utama saat container menyala
CMD ["/app/api"]
