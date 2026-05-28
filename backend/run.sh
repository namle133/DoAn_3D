#!/bin/bash

set -e

cd "$(dirname "$0")" || exit 1

echo "🧹 Cleaning up old containers and volumes..."
docker-compose down -v --remove-orphans 2>/dev/null || true
sleep 2

echo "🐘 Starting PostgreSQL..."
docker-compose up -d postgres

echo "⏳ Waiting for PostgreSQL to initialize (up to 30 seconds)..."
for i in {1..30}; do
    if docker exec vinhomes_postgres pg_isready -U vinhomes_user -d vinhomes_db 2>/dev/null >/dev/null; then
        echo "✅ PostgreSQL is ready!"
        sleep 2  # Extra buffer
        break
    fi
    echo -n "."
    sleep 1
done

echo ""
echo "🚀 Starting backend..."
go run ./cmd/main.go
