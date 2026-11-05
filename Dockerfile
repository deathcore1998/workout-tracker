# Этап 1: сборка
FROM golang:1.27.1-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/main .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]