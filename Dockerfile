# Базовый образ для сборки
FROM golang:1.23-alpine AS builder

# Установка необходимых инструментов, если требуется CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -ldflags="-s -w" -o hmb_fighting ./server/cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates curl bash && \
    curl -sfL https://github.com/pressly/goose/releases/download/v3.15.0/goose_linux_x86_64 -o /usr/local/bin/goose && \
    chmod +x /usr/local/bin/goose && \
    curl -sfL https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh -o /usr/local/bin/wait-for-it.sh && \
    chmod +x /usr/local/bin/wait-for-it.sh
WORKDIR /app
COPY --from=builder /app/hmb_fighting .
COPY server/migrations /app/migrations
EXPOSE 8080
CMD ["/bin/bash", "-c", "wait-for-it.sh db:5432 --timeout=30 && wait-for-it.sh redis:6379 --timeout=30 && goose -dir /app/migrations postgres \"$DATABASE_URL\" up && exec ./hmb_fighting"]