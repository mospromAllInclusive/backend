# syntax=docker/dockerfile:1.6

ARG GO_VERSION=1.24.3
ARG MIGRATE_VERSION=4.17.1
ARG APP_DIR=./

# --- Build stage ---
FROM golang:${GO_VERSION}-alpine AS builder
RUN apk add --no-cache build-base git ca-certificates
WORKDIR /src

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Сборка бинаря
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ${APP_DIR}

# --- Download golang-migrate (arch-aware) ---
FROM alpine:3.20 AS migrate-bin
ARG MIGRATE_VERSION
RUN set -eux; \
    apk add --no-cache ca-certificates curl tar; \
    arch="$(apk --print-arch)"; \
    case "$arch" in \
      x86_64)  pkg="migrate.linux-amd64.tar.gz" ;; \
      aarch64) pkg="migrate.linux-arm64.tar.gz" ;; \
      armv7)   pkg="migrate.linux-armv6.tar.gz" ;; \
      *) echo "unsupported arch: $arch" >&2; exit 1 ;; \
    esac; \
    curl -fsSL -o /tmp/migrate.tgz \
      "https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/${pkg}"; \
    tar -xzf /tmp/migrate.tgz -C /usr/local/bin; \
    chmod +x /usr/local/bin/migrate

# --- Runtime ---
FROM alpine:3.20

# Устанавливаем зависимости и создаём пользователя
RUN apk add --no-cache ca-certificates bash postgresql-client \
 && adduser -D -g '' app

WORKDIR /app

# Бинарь сервера и migrate
COPY --from=builder     /out/server            /usr/local/bin/server
COPY --from=migrate-bin /usr/local/bin/migrate /usr/local/bin/migrate

# Миграции внутрь образа
COPY migrations/postgres ./migrations/postgres

# Entrypoint: ждём БД, применяем миграции, запускаем сервер
COPY --chmod=0755 <<'EOF' /usr/local/bin/entrypoint.sh
#!/usr/bin/env sh
set -euo pipefail

# Берём DATABASE_URL, если нет — POSTGRES_URL
DB_URL="${DATABASE_URL:-}"
if [ -z "$DB_URL" ]; then
  DB_URL="${POSTGRES_URL:-}"
fi
if [ -z "$DB_URL" ]; then
  echo "ERROR: neither DATABASE_URL nor POSTGRES_URL is set"; exit 1
fi

# Красивый лог без пароля
mask_url () {
  echo "$1" | sed -E 's#(://[^:/@]+:)[^@]+@#\1******@#'
}
echo "Using database: $(mask_url "$DB_URL")"

# Ждём Postgres
export PGCONNECT_TIMEOUT="${PGCONNECT_TIMEOUT:-3}"
echo "Waiting for Postgres to be ready..."
until pg_isready -d "$DB_URL" >/dev/null 2>&1; do
  sleep 2
done
echo "Postgres is ready."

# Накатываем миграции (файлы вида 20251017000000_init.up.sql поддерживаются)
MIGRATIONS_PATH="${MIGRATIONS_PATH:-/app/migrations/postgres}"
echo "Applying migrations from $MIGRATIONS_PATH ..."
migrate -verbose -path "$MIGRATIONS_PATH" -database "$DB_URL" up

echo "Starting server..."
exec server
EOF

# Права на рабочую директорию и непривилегированный запуск
RUN chown -R app:app /app
USER app

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD []
