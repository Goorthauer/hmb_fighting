# Базовый образ для сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o hmb_fighting ./server/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/hmb_fighting .
COPY server/migrations /app/migrations
RUN apk add --no-cache curl && \
    curl -sfL https://github.com/pressly/goose/releases/download/v3.15.0/goose_linux_x86_64 -o /usr/local/bin/goose && \
    chmod +x /usr/local/bin/goose
EXPOSE 8080
CMD ["sh", "-c", "goose -dir /app/migrations postgres \"$DATABASE_URL\" up && ./hmb_fighting"]