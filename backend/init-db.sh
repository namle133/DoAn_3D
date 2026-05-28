#!/bin/bash

set -e

cd "$(dirname "$0")"

echo "🔴 Stopping all containers..."
docker-compose down -v --remove-orphans 2>/dev/null || true
docker kill vinhomes_postgres 2>/dev/null || true
docker rm vinhomes_postgres 2>/dev/null || true
docker volume rm backend_postgres_data 2>/dev/null || true

echo "✅ All containers and volumes removed"
echo ""
echo "📦 Starting fresh PostgreSQL container..."
docker-compose up -d postgres

echo "⏳ Waiting 10 seconds for database to initialize..."
sleep 10

# Check if database is ready
if docker-compose exec -T postgres pg_isready -U vinhomes_user -d vinhomes_db 2>/dev/null; then
    echo "✅ Database is ready!"
else
    echo "⚠️  Database might not be ready yet. Waiting another 5 seconds..."
    sleep 5
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "🚀 Now you can run:"
echo "   go run ./cmd/main.go"
echo ""
