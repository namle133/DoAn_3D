#!/bin/bash

# ✅ COMPLETE BACKEND STARTUP SCRIPT
# This handles everything needed to start the backend

set -e

cd "$(dirname "$0")" || exit 1

echo ""
echo "════════════════════════════════════════════════════════"
echo "  🚀 Vinhomes Backend Startup Script"
echo "════════════════════════════════════════════════════════"
echo ""

# Step 1: Clean up
echo "🧹 Step 1/4: Cleaning up old containers..."
docker-compose down -v --remove-orphans 2>/dev/null || true
sleep 2

# Step 2: Start PostgreSQL
echo "🐘 Step 2/4: Starting PostgreSQL..."
docker-compose up -d postgres

# Step 3: Wait for PostgreSQL and init
echo "⏳ Step 3/4: Waiting for database initialization (up to 40 seconds)..."
for i in {1..40}; do
    # Check if postgres user is ready
    if docker exec vinhomes_postgres pg_isready -U postgres 2>/dev/null >/dev/null; then
        echo "   ✓ PostgreSQL ready after $i seconds"
        break
    fi
    echo -n "."
    sleep 1
done

# Give init script time to run
sleep 3

# Verify vinhomes_user was created
echo "✅ Verifying vinhomes_user exists..."
if docker exec vinhomes_postgres psql -U postgres -c "SELECT 1 FROM pg_roles WHERE rolname='vinhomes_user';" 2>/dev/null | grep -q "1"; then
    echo "   ✓ vinhomes_user created successfully"
else
    echo "   ⚠️  vinhomes_user not found, creating manually..."
    docker exec vinhomes_postgres psql -U postgres << EOF
CREATE USER vinhomes_user WITH ENCRYPTED PASSWORD 'vinhomes_pass';
ALTER USER vinhomes_user WITH SUPERUSER CREATEDB CREATEROLE;
CREATE DATABASE vinhomes_db OWNER vinhomes_user;
EOF
    echo "   ✓ Created vinhomes_user and vinhomes_db"
fi

# Step 4: Start backend
echo ""
echo "🚀 Step 4/4: Starting backend..."
echo ""
go run ./cmd/main.go
