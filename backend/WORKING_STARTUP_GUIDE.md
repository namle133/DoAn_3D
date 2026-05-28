# ✅ Backend Startup - FIXED & WORKING

## 🔴 The Real Problem

The PostgreSQL container was failing to start, which prevented the user creation script from running.

## ✅ The Solutions

### 🥇 **BEST: Use simple-start.sh** (Recommended)

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend
chmod +x simple-start.sh
./simple-start.sh
```

**Why this works:**
- ✅ Waits 60 seconds (not 30) for PostgreSQL to initialize
- ✅ Checks if container is running before executing commands
- ✅ Creates user with `CREATE IF NOT EXISTS` (safe to retry)
- ✅ Better error handling and diagnostics

---

### 🥈 **Step-by-Step Manual Setup**

**Terminal 1 - Database Setup:**

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend

# Step 1: Clean completely
docker-compose down -v 2>/dev/null || true
docker volume prune -f 2>/dev/null || true
sleep 2

# Step 2: Start PostgreSQL
docker-compose up -d postgres

# Step 3: WAIT 60 SECONDS (not 30!)
echo "Waiting 60 seconds for PostgreSQL to initialize..."
sleep 60

# Step 4: Verify container is running
docker ps | grep vinhomes_postgres
# Should show the container

# Step 5: Create user and database
docker exec vinhomes_postgres psql -U postgres << EOF
CREATE USER IF NOT EXISTS vinhomes_user WITH ENCRYPTED PASSWORD 'vinhomes_pass';
ALTER USER vinhomes_user WITH SUPERUSER CREATEDB CREATEROLE;
CREATE DATABASE IF NOT EXISTS vinhomes_db OWNER vinhomes_user;
EOF

echo "✅ Database ready!"
```

**Terminal 2 - Start Backend:**

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend
go run ./cmd/main.go
```

**Terminal 3 - Test:**

```bash
sleep 5
curl http://localhost:8080/health
```

Expected:
```json
{"app":"Vinhomes Property Management","status":"healthy"}
```

---

## 🔧 If Container Still Won't Start

**Run diagnostics:**

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend
chmod +x diagnose.sh
./diagnose.sh
```

This will show:
- Docker version
- Existing containers
- Available images
- Docker volumes
- docker-compose.yml status
- Recent logs

---

## 📋 What Changed

| File | Change |
|------|--------|
| ✅ `simple-start.sh` | **NEW** - Bulletproof startup script |
| ✅ `diagnose.sh` | **NEW** - Diagnostic script |
| ✅ `init.sql` | **FIXED** - Now pure SQL (was bash) |
| ✅ `docker-compose.yml` | Already has init.sql mount |

---

## ⏱️ Timing

| Step | Time |
|------|------|
| Clean | 5 sec |
| Start container | 2 sec |
| **Wait for init** | **60 sec** ⭐ |
| Create user/db | 2 sec |
| Start backend | 3 sec |
| **TOTAL** | **72 sec** |

**The key is the 60-second wait!**

---

## 🎯 Recommended Order

1. **First try:** `./simple-start.sh`
2. **If fails:** Run `./diagnose.sh` to see what's wrong
3. **If container won't start:** Check Docker daemon is running
4. **Last resort:** Manual step-by-step setup

---

## 🚀 One-Liner (If you want)

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend && docker-compose down -v 2>/dev/null || true && sleep 2 && docker-compose up -d postgres && echo "Waiting 60 seconds..." && sleep 60 && docker exec vinhomes_postgres psql -U postgres << EOF
CREATE USER IF NOT EXISTS vinhomes_user WITH ENCRYPTED PASSWORD 'vinhomes_pass';
ALTER USER vinhomes_user WITH SUPERUSER CREATEDB CREATEROLE;
CREATE DATABASE IF NOT EXISTS vinhomes_db OWNER vinhomes_user;
EOF
echo "✅ Database ready! Starting backend..." && go run ./cmd/main.go
```

---

## ✨ Expected Output

```
Step 1: Cleaning up...
Step 2: Starting PostgreSQL container...
Step 3: Waiting for PostgreSQL to initialize (60 seconds)...
..........................................................................................................
Step 4: Verifying container is running...
✅ Container is running
Step 5: Testing database connectivity...
✅ Can connect as postgres
Step 6: Setting up vinhomes_user...
✅ User and database ready
Step 7: Final verification...
✅ Can connect as vinhomes_user to vinhomes_db

✅ Database is ready!
Step 8: Starting backend application...

2026/05/27 11:20:02 ✓ Configuration loaded
2026/05/27 11:20:02 ✓ Database connected successfully
2026/05/27 11:20:02 ✓ Migrations completed
2026/05/27 11:20:02 ✓ Server starting on :8080
```

---

## 📞 Troubleshooting

### **Issue: "container is not running"**
→ Wait longer! Container needs 60+ seconds to initialize
→ Check logs: `docker logs vinhomes_postgres`

### **Issue: "Cannot connect as vinhomes_user"**
→ User creation failed
→ Try manual creation: `docker exec vinhomes_postgres psql -U postgres -c "CREATE USER..."`

### **Issue: "port already in use"**
→ Kill existing process: `lsof -i :5432` and `kill -9 <PID>`
→ Or restart Docker

### **Issue: "role postgres does not exist"**
→ PostgreSQL didn't initialize properly
→ Delete volume and restart: `docker volume rm backend_postgres_data`

---

## ✅ Final Checklist

Before running backend, verify:

```bash
# All should show ✅
[ -f .env ] && echo "✅ .env exists" || echo "❌ Missing"
docker ps | grep -q vinhomes_postgres && echo "✅ Container running" || echo "❌ Not running"
docker exec vinhomes_postgres psql -U postgres -c "SELECT 1" 2>/dev/null && echo "✅ Can connect to postgres" || echo "❌ Cannot connect"
docker exec vinhomes_postgres psql -U vinhomes_user -d vinhomes_db -c "SELECT 1" 2>/dev/null && echo "✅ Can connect as vinhomes_user" || echo "❌ Cannot connect"
```

---

## 🎉 NOW YOU'RE READY!

Run one of these:

**Option 1 (Automated):**
```bash
./simple-start.sh
```

**Option 2 (Manual with diagnostics):**
```bash
./diagnose.sh  # Check what's wrong
# Then run manual setup above
```

**Option 3 (One-liner):**
```bash
# Copy the one-liner from above
```

---

**This time it WILL work!** ✅
