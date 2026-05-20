FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy file module agar bisa download dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy SELURUH kode sumber kamu ke dalam Docker
COPY . .

# Build aplikasinya (Perhatikan: kita tunjuk cmd/main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/payment-service .

# Sesuai dengan absen port Payment Service milikmu
EXPOSE 8089

CMD ["./payment-service"]