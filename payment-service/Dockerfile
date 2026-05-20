
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY internal/ ./internal/
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/payment-service .

EXPOSE 8080

CMD ["./payment-service"]