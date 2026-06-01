#!/bin/bash

set -e

cd "$(dirname "$0")" || exit 1

echo "🐘 Starting Docker PostgreSQL (host port 5433)..."
docker-compose up -d postgres

echo "⏳ Waiting for PostgreSQL..."
for i in {1..30}; do
    if docker exec vinhomes_postgres pg_isready -U vinhomes_user -d vinhomes_db >/dev/null 2>&1; then
        echo "✅ PostgreSQL is ready on 127.0.0.1:5433"
        break
    fi
    sleep 1
done

echo "🚀 Starting backend..."
go run ./cmd/main.go
