FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Sesuaikan lokasi main.go (apakah ./main.go atau ./cmd/main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o promo-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/promo-service .

# UBAH PORT KE JATAH PROMO SERVICE
EXPOSE 8086

CMD ["./promo-service"]