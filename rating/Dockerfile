# ==========================================
# STAGE 1: BUILDER (Merakit Kode)
# ==========================================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Pastikan lokasi main.go sudah benar (apakah ./cmd/main.go atau ./main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o satset-rating ./cmd/main.go

# ==========================================
# STAGE 2: RUNNER (Menjalankan Aplikasi)
# ==========================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Mengambil hasil rakitan dari Stage 1 yang ukurannya sangat kecil
COPY --from=builder /app/satset-rating .

# Port khusus Rating Service
EXPOSE 8087

CMD ["./satset-rating"]