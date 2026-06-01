#!/bin/bash

# Start backend development server (DB: Docker Postgres on host port 5433)

echo "🚀 Starting Vinhomes Backend..."
echo "================================"

cd "$(dirname "$0")" || exit

echo "📦 Checking Docker PostgreSQL (vinhomes_postgres)..."
if ! docker ps --format '{{.Names}}' | grep -qx 'vinhomes_postgres'; then
    echo "🐘 Starting PostgreSQL container on host port 5433..."
    docker-compose up -d postgres
fi

echo "⏳ Waiting for database on 127.0.0.1:5433..."
for i in {1..30}; do
    if docker exec vinhomes_postgres pg_isready -U vinhomes_user -d vinhomes_db >/dev/null 2>&1; then
        echo "✅ Docker PostgreSQL is ready"
        break
    fi
    sleep 1
done

echo ""
echo "🔧 Starting application..."
echo "📝 DB config: 127.0.0.1:5433 (Docker) — see config/.env"
echo ""
echo "🌐 Server: http://localhost:8080"
echo "   Health: curl http://localhost:8080/health"
echo ""
echo "Press Ctrl+C to stop"
echo ""

go run ./cmd/main.go
