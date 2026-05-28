#!/bin/bash

# 🔍 Diagnostic script to debug PostgreSQL issues

echo "🔍 Vinhomes Backend - Diagnostic Report"
echo "========================================"
echo ""

echo "1️⃣  Docker Running?"
docker --version
echo ""

echo "2️⃣  Existing Containers:"
docker ps -a | grep vinhomes || echo "None found"
echo ""

echo "3️⃣  Docker Images:"
docker images | grep postgres || echo "No postgres image found"
echo ""

echo "4️⃣  Docker Volumes:"
docker volume ls | grep postgres || echo "No volumes found"
echo ""

echo "5️⃣  docker-compose.yml Check:"
if [ -f docker-compose.yml ]; then
    echo "✅ File exists"
    echo "   Services:"
    grep "^  " docker-compose.yml | head -10
else
    echo "❌ File not found!"
fi
echo ""

echo "6️⃣  .env File Check:"
if [ -f .env ]; then
    echo "✅ File exists"
    echo "   DB_USER=$(grep DB_USER .env)"
    echo "   DB_NAME=$(grep DB_NAME .env)"
else
    echo "❌ File not found!"
fi
echo ""

echo "7️⃣  Available Disk Space:"
df -h | head -3
echo ""

echo "8️⃣  Docker Logs (if container exists):"
if docker ps -a | grep -q vinhomes_postgres; then
    docker logs --tail 20 vinhomes_postgres 2>/dev/null || echo "Could not retrieve logs"
else
    echo "Container doesn't exist yet"
fi
echo ""

echo "✅ Diagnostic complete"
