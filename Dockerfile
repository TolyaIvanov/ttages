FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git postgresql-client \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o file-service ./cmd

# Финальный контейнер
FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/file-service /app/file-service
COPY --from=builder /app/migrations /app/migrations

RUN mkdir -p /storage && chmod 777 /storage

RUN chmod +x /app/file-service

EXPOSE 50051

CMD ["/app/file-service"]
