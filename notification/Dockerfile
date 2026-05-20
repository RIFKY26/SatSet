FROM golang:1.25-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# UPDATE BARIS INI: Arahkan build ke folder cmd/main.go yang benar
RUN go build -o satset-notification ./cmd/main.go

CMD ["./satset-notification"]