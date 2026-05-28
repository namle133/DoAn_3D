#!/bin/bash

echo "🔄 Stopping old containers..."
docker kill vinhomes_postgres 2>/dev/null
docker rm vinhomes_postgres 2>/dev/null
docker volume rm backend_postgres_data 2>/dev/null

echo "📦 Starting fresh PostgreSQL..."
docker-compose up -d postgres

echo "⏳ Waiting for database to be ready..."
sleep 3

echo "✅ Database ready!"
echo ""
echo "🚀 Now run the backend in another terminal:"
echo "   go run ./cmd/main.go"
echo ""
