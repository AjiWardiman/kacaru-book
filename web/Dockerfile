# ==============================================================================
# Multi-Stage Dockerfile untuk Frontend Web (React + Vite + Nginx)
# ==============================================================================

# ------------------------------------------------------------------------------
# Stage 1: Build Aplikasi Frontend Vite
# ------------------------------------------------------------------------------
FROM node:20-alpine AS builder

WORKDIR /app

# Cache dependensi npm agar proses build lebih cepat
COPY package.json package-lock.json ./
RUN npm ci

# Salin seluruh berkas kode frontend dan lakukan proses build produksi
COPY . .
RUN npm run build

# ------------------------------------------------------------------------------
# Stage 2: Runtime Image Ringan Menggunakan Nginx
# ------------------------------------------------------------------------------
FROM nginx:alpine

# Salin hasil kompilasi statis (folder dist) ke direktori root Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Salin konfigurasi Nginx untuk menangani routing SPA (React Router) & proxy
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Buka port standar web HTTP
EXPOSE 80

# Jalankan Nginx di foreground
CMD ["nginx", "-g", "daemon off;"]
