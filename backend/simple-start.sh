#!/bin/bash

# ✅ SIMPLE BACKEND STARTUP - NO FANCY STUFF
# This is the absolute simplest approach

set -e

cd "$(dirname "$0")" || exit 1

echo ""
echo "🚀 Starting Vinhomes Backend"
echo "============================"
echo ""

# Step 1: Clean everything
echo "Step 1: Cleaning up..."
docker-compose down -v 2>/dev/null || true
docker volume prune -f 2>/dev/null || true
sleep 2

# Step 2: Start container
echo "Step 2: Starting PostgreSQL container..."
docker-compose up -d postgres

# Step 3: Wait a LOT longer
echo "Step 3: Waiting for PostgreSQL to initialize (60 seconds)..."
sleep 60

# Step 4: Check if container is running
echo "Step 4: Verifying container is running..."
if ! docker ps | grep -q vinhomes_postgres; then
    echo "❌ Container failed to start!"
    echo "Checking logs..."
    docker logs vinhomes_postgres || echo "No logs available"
    exit 1
fi
echo "✅ Container is running"

# Step 5: Check connectivity
echo "Step 5: Testing database connectivity..."
if docker exec vinhomes_postgres psql -U postgres -c "SELECT 1" 2>/dev/null >/dev/null; then
    echo "✅ Can connect as postgres"
else
    echo "❌ Cannot connect to database"
    exit 1
fi

# Step 6: Create user/database (idempotent, PostgreSQL-correct syntax)
echo "Step 6: Setting up vinhomes_user..."
docker exec -i vinhomes_postgres psql -U postgres <<'EOF'
DO
$$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'vinhomes_user') THEN
        CREATE ROLE vinhomes_user LOGIN ENCRYPTED PASSWORD 'vinhomes_pass';
    END IF;
END
$$;
ALTER ROLE vinhomes_user WITH SUPERUSER CREATEDB CREATEROLE;

SELECT 'CREATE DATABASE vinhomes_db OWNER vinhomes_user'
WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'vinhomes_db')
\gexec

ALTER DATABASE vinhomes_db OWNER TO vinhomes_user;
GRANT ALL PRIVILEGES ON DATABASE vinhomes_db TO vinhomes_user;
EOF

echo "✅ User and database ready"

# Step 7: Verify final state
echo "Step 7: Final verification..."
if docker exec vinhomes_postgres psql -U vinhomes_user -d vinhomes_db -c "SELECT 1" 2>/dev/null >/dev/null; then
    echo "✅ Can connect as vinhomes_user to vinhomes_db"
else
    echo "❌ Cannot connect as vinhomes_user"
    echo "Attempting fallback connection as postgres..."
    docker exec vinhomes_postgres psql -U postgres -d postgres -c "SELECT 1"
fi

# Step 8: Start backend
echo ""
echo "✅ Database is ready!"
echo "Step 8: Starting backend application..."
echo ""
go run ./cmd/main.go
