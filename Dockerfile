# ── Этап сборки ──────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем зависимости отдельным слоем — кеш не инвалидируется при изменении кода
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO выключен для статического бинаря; -w -s убирает debug-символы (меньше размер)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o fencing-club ./cmd/main.go

# ── Минимальный рантайм-образ ─────────────────────────────────────────────────
FROM alpine:3.19

# ca-certificates нужны для HTTPS-запросов (если сервис их делает)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/fencing-club .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

# Запускаем без root — минимальные привилегии
USER nobody

CMD ["./fencing-club"]
