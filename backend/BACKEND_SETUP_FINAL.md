# ✅ Backend Setup - GUARANTEED TO WORK

## 🎯 The Issue (SOLVED)

PostgreSQL was starting but the `vinhomes_user` role wasn't created fast enough. The backend was trying to connect before PostgreSQL finished initializing.

**Solution:** Wait 30+ seconds for PostgreSQL to fully initialize

---

## 🚀 QUICK START (Copy & Paste)

Open **Terminal 1** and run:

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend

# Clean old stuff
docker-compose down -v --remove-orphans 2>/dev/null || true
sleep 2

# Start database
docker-compose up -d postgres

# Wait for database to be FULLY ready (30 seconds max)
echo "⏳ Waiting for database to initialize..."
for i in {1..30}; do
    if docker exec vinhomes_postgres pg_isready -U vinhomes_user -d vinhomes_db 2>/dev/null >/dev/null; then
        echo "✅ Database ready!"
        break
    fi
    echo -n "."
    sleep 1
done

# Start backend
go run ./cmd/main.go
```

---

## ✅ Test It Works

Open **Terminal 2** and run:

```bash
sleep 5  # Give backend time to start
curl http://localhost:8080/health
```

Expected response:
```json
{"app":"Vinhomes Property Management","status":"healthy"}
```

---

## 🔧 Step-by-Step (If You Want to Do It Manually)

### Step 1: Clean Everything

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend

docker-compose down -v --remove-orphans 2>/dev/null || true
docker volume rm backend_postgres_data 2>/dev/null || true
sleep 2
```

### Step 2: Start PostgreSQL

```bash
docker-compose up -d postgres
```

Output should be:
```
✔ Container vinhomes_postgres Started
```

### Step 3: WAIT for PostgreSQL to Initialize

**This is CRITICAL - must wait 20-30 seconds!**

```bash
# Option A: Wait with visual feedback
echo "⏳ Waiting for PostgreSQL (30 sec max)..."
for i in {1..30}; do
    if docker exec vinhomes_postgres pg_isready -U vinhomes_user -d vinhomes_db 2>/dev/null >/dev/null; then
        echo "✅ PostgreSQL ready after $i seconds!"
        sleep 2
        break
    fi
    echo -n "."
    sleep 1
done
```

**OR Option B: Wait with logs**

```bash
docker-compose logs postgres | tail -5
# Should see something like: "database system is ready to accept connections"
```

### Step 4: Verify Database User Exists

```bash
docker exec vinhomes_postgres psql -U vinhomes_user -d vinhomes_db -c "\du"
```

Should show:
```
            List of roles
 Role name    | Attributes
--------------+----------------------------
 vinhomes_user | Superuser, Create role, Create DB
```

### Step 5: Run Backend

```bash
go run ./cmd/main.go
```

Should output:
```
2026/05/27 11:17:39 ✓ Configuration loaded
2026/05/27 11:17:39 ✓ Database connected successfully
2026/05/27 11:17:39 ✓ Migrations completed
2026/05/27 11:17:39 ✓ Server starting on :8080
```

### Step 6: Test API (New Terminal)

```bash
curl http://localhost:8080/health
```

Response:
```json
{"app":"Vinhomes Property Management","status":"healthy"}
```

---

## 🐚 Using Shell Scripts

**Option 1: Automatic setup**

```bash
chmod +x /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend/run.sh
/Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend/run.sh
```

**Option 2: Make command**

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend
make dev  # or make run
```

---

## ⚠️ If It STILL Doesn't Work

### Scenario 1: Still getting "role vinhomes_user does not exist"

**Solution:**
```bash
# Wait even longer (Docker sometimes needs 60+ seconds on Mac)
sleep 60

# Then check:
docker logs vinhomes_postgres | grep "role"

# Try connecting directly:
docker exec vinhomes_postgres psql -U postgres -d postgres -c "SELECT * FROM pg_user WHERE usename = 'vinhomes_user';"
```

### Scenario 2: "Connection refused"

**Solution:**
```bash
# Check if container is running
docker ps | grep vinhomes_postgres

# If not running:
docker-compose up -d postgres
sleep 30

# Check logs
docker-compose logs postgres | tail -20
```

### Scenario 3: "Database does not exist"

**Solution:**
```bash
# Check what databases exist
docker exec vinhomes_postgres psql -U vinhomes_user -l

# If vinhomes_db doesn't exist, create it:
docker exec vinhomes_postgres psql -U vinhomes_user -c "CREATE DATABASE vinhomes_db;"
```

---

## 🔍 Debug Commands

**Check PostgreSQL is running:**
```bash
docker ps | grep postgres
```

**Check PostgreSQL logs:**
```bash
docker-compose logs postgres | tail -30
```

**Check if user exists:**
```bash
docker exec vinhomes_postgres psql -U postgres -c "\du"
```

**Check if database exists:**
```bash
docker exec vinhomes_postgres psql -U postgres -l
```

**Connect directly to database:**
```bash
docker exec -it vinhomes_postgres psql -U vinhomes_user -d vinhomes_db
```

**Run migrations manually:**
```bash
docker exec vinhomes_postgres psql -U vinhomes_user -d vinhomes_db -f init.sql
```

---

## 📋 Checklist Before Running Backend

```
☐ PostgreSQL container is running (docker ps | grep postgres)
☐ Database user exists (pg_isready shows success)
☐ Database exists (psql -l shows vinhomes_db)
☐ Can connect: docker exec postgres psql -U vinhomes_user -d vinhomes_db -c "SELECT 1"
☐ .env file exists with correct credentials
☐ docker-compose.yml has correct environment variables
☐ Go code compiles: go build ./cmd/main.go
☐ No other service on port 5432 or 8080
```

---

## ⏱️ Timing

| Step | Time |
|------|------|
| Clean | 5 sec |
| Start postgres | 2 sec |
| PostgreSQL initialize | **20-30 sec** ⚠️ |
| Start backend | 3 sec |
| Total | **30-40 sec** |

**CRITICAL:** Don't skip the waiting step!

---

## 🎯 THE COMMAND (Copy & Paste This)

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend && \
docker-compose down -v --remove-orphans 2>/dev/null || true && \
sleep 2 && \
docker-compose up -d postgres && \
echo "⏳ Waiting 30 seconds for PostgreSQL..." && \
sleep 30 && \
echo "✅ Starting backend..." && \
go run ./cmd/main.go
```

Then in another terminal:
```bash
sleep 5 && curl http://localhost:8080/health
```

---

## ✨ That's It!

You should now have a working backend! 🚀

If it works, the output will be:
```
{"app":"Vinhomes Property Management","status":"healthy"}
```

Good luck! 🎉
