# 🚀 Complete Backend Setup Guide

## Problem Identified & Fixed ✅

**Issue:** Two different `.env` files with conflicting credentials
- ❌ `backend/.env` had `vinhomes_user`
- ❌ `backend/config/.env` had `postgres` user
- ❌ Docker was trying to create `vinhomes_user` but database didn't exist yet

**Solution:** ✅ Both files now use `vinhomes_user:vinhomes_pass`

---

## 🎯 Setup Steps (Copy & Paste)

### Step 1: Verify Credentials Match

```bash
echo "=== Checking .env files ==="
echo "Root .env:"
grep DB_USER /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend/.env

echo ""
echo "Config .env:"
grep DB_USER /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend/config/.env

echo ""
echo "Docker-compose:"
grep -A 5 "POSTGRES_USER" /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend/docker-compose.yml | head -4
```

All should show: `vinhomes_user` ✅

### Step 2: Clean Everything

**IMPORTANT: Run these commands EXACTLY:**

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend

# Kill and remove EVERYTHING
docker-compose down -v --remove-orphans 2>/dev/null || true
docker kill vinhomes_postgres 2>/dev/null || true
docker rm vinhomes_postgres 2>/dev/null || true
docker volume rm backend_postgres_data 2>/dev/null || true

# Verify all removed
docker ps | grep vinhomes || echo "✅ All containers removed"
docker volume ls | grep postgres_data || echo "✅ All volumes removed"
```

### Step 3: Start Fresh Database

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend

# Start fresh
docker-compose up -d postgres

# Wait for it to initialize
echo "⏳ Waiting 10 seconds..."
sleep 10

# Verify it's ready
docker-compose exec postgres pg_isready -U vinhomes_user -d vinhomes_db
```

Expected output: `accepting connections`

### Step 4: Check Database User Exists

```bash
# Connect to postgres and check users
docker-compose exec postgres psql -U vinhomes_user -d vinhomes_db -c "\du"

# Should show: vinhomes_user | Superuser, Create role, Create DB
```

### Step 5: Run Backend

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend
go run ./cmd/main.go
```

Expected output:
```
2026/05/27 11:15:16 ✓ Configuration loaded
2026/05/27 11:15:16 ✓ Database connected successfully
2026/05/27 11:15:16 ✓ Migrations completed
2026/05/27 11:15:16 ✓ Server starting on :8080
```

### Step 6: Test API

**In another terminal:**

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"app":"Vinhomes Property Management","status":"healthy"}
```

---

## ⚡ One-Command Setup (Copy & Paste ALL)

If you want to do it all at once, paste this entire script:

```bash
#!/bin/bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend

echo "🔴 Cleaning up..."
docker-compose down -v --remove-orphans 2>/dev/null || true
docker kill vinhomes_postgres 2>/dev/null || true
docker rm vinhomes_postgres 2>/dev/null || true
docker volume rm backend_postgres_data 2>/dev/null || true

echo "📦 Starting database..."
docker-compose up -d postgres

echo "⏳ Waiting for database..."
sleep 10

echo "✅ Database ready! Starting backend..."
go run ./cmd/main.go
```

---

## 🔍 Troubleshooting

### If you get "role vinhomes_user does not exist"

**Solution:**
```bash
# Stop everything
docker-compose down -v

# Delete the volume completely
docker volume rm backend_postgres_data

# Start fresh
docker-compose up -d postgres
sleep 10

# Verify user exists
docker-compose exec postgres psql -U vinhomes_user -d vinhomes_db -c "\l"
```

### If you get "connection refused"

**Solution:**
```bash
# Check if postgres is running
docker ps | grep postgres

# If not running:
docker-compose up -d postgres

# Wait longer
sleep 15

# Check logs
docker-compose logs postgres | tail -20
```

### If you get "database does not exist"

**Solution:**
```bash
# Make sure you're using the right credentials
cat .env | grep DB_

# Should show: vinhomes_user and vinhomes_pass

# If not, update it:
# Edit .env file and verify DB_USER=vinhomes_user
```

---

## ✅ Verification Checklist

Before running `go run ./cmd/main.go`, check:

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend

# 1. Check .env file exists and has correct credentials
[ -f .env ] && echo "✅ .env exists" || echo "❌ .env missing"
grep "vinhomes_user" .env && echo "✅ Correct DB_USER" || echo "❌ Wrong DB_USER"

# 2. Check docker-compose has correct credentials
grep "POSTGRES_USER: vinhomes_user" docker-compose.yml && echo "✅ Correct POSTGRES_USER" || echo "❌ Wrong POSTGRES_USER"

# 3. Check PostgreSQL is running
docker ps | grep postgres && echo "✅ PostgreSQL running" || echo "❌ PostgreSQL not running"

# 4. Check database is accessible
docker exec vinhomes_postgres psql -U vinhomes_user -d vinhomes_db -c "SELECT version();" && echo "✅ Database accessible" || echo "❌ Database not accessible"

# 5. Check backend code compiles
go build ./cmd/main.go 2>&1 | head -3 && echo "✅ Code compiles" || echo "❌ Compilation error"
```

All should show ✅

---

## 🎯 Expected Timeline

| Step | Time | Status |
|------|------|--------|
| 1. Clean up | 5 sec | ⏳ |
| 2. Start database | 2 sec | ⏳ |
| 3. Wait for init | 10 sec | ⏳ |
| 4. Run backend | 5 sec | ⏳ |
| **TOTAL** | **22 sec** | ⏱️ |

---

## 📝 Summary of Changes

✅ Updated `/backend/config/.env` - Changed `postgres` → `vinhomes_user`
✅ Verified `/backend/.env` - Already had correct credentials
✅ Created `/backend/init-db.sh` - Automated setup script
✅ All credentials now match docker-compose.yml

---

## 🚀 Ready to Go!

Follow the **Step 2-6** above or use the **One-Command Setup** 

Your backend will be running in **~22 seconds** ⚡

Good luck! 🎉
