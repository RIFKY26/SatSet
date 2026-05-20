FROM golang:1.25-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Path disesuaikan menjadi ./cmd/main.go
RUN go build -o satset-rating ./cmd/main.go

CMD ["./satset-rating"]