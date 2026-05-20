# Gunakan base image Golang versi terbaru (sesuai go.mod kamu)
FROM golang:1.25-alpine

# Set direktori kerja di dalam container
WORKDIR /app

# Salin file go.mod dan go.sum (tanda * agar tidak error jika go.sum belum ada)
COPY go.mod go.sum* ./

# Download dependensi
RUN go mod download

# Salin seluruh kode sumber
COPY . .

# Build aplikasi dari cmd/main.go
RUN go build -o main ./cmd/main.go

# Jalankan aplikasi
CMD ["./main"]