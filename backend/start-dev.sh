#!/bin/bash

# Start backend development server

echo "🚀 Starting Vinhomes Backend..."
echo "================================"

# Navigate to backend directory
cd "$(dirname "$0")" || exit

# Check if docker containers are running
echo "📦 Checking Docker containers..."
if ! docker-compose ps | grep -q "vinhomes_postgres"; then
    echo "🐘 Starting PostgreSQL container..."
    docker-compose up -d postgres
    echo "⏳ Waiting for database to be ready..."
    sleep 5
fi

# Run the application
echo ""
echo "🔧 Starting application..."
echo "📝 Configuration loaded from .env"
echo ""
echo "🌐 Server will be available at: http://localhost:8080"
echo "📊 API documentation at: http://localhost:8080/docs (if swagger enabled)"
echo ""
echo "To test the API:"
echo "  curl http://localhost:8080/health"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

go run ./cmd/main.go
