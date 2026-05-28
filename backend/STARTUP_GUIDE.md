# 🚀 Backend Startup Guide - Step by Step

## ✅ What's been completed

- ✅ All 34 files created
- ✅ All compilation errors fixed
- ✅ Docker support configured for ARM64 Mac
- ✅ Database ready (PostgreSQL running in Docker)
- ✅ Backend application ready to start

---

## 🛠️ Step 1: Verify PostgreSQL is Running

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend

# Check if PostgreSQL container is running
docker ps | grep postgres

# Start if not running:
docker-compose up -d postgres
```

Expected output:
```
CONTAINER ID   IMAGE              STATUS
abc123def456   postgres:16-alpine   Up 2 minutes
```

---

## 🏃 Step 2: Start the Backend Server

### Option A: Run directly (Recommended for development)

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend
go run ./cmd/main.go
```

### Option B: Build and run binary

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend
go build -o vinhomes ./cmd/main.go
./vinhomes
```

### Option C: Use Makefile

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend
make run
```

### Option D: With hot-reload (development)

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend
make dev  # Requires 'air' tool
```

---

## ✅ Step 3: Verify Backend is Running

Open a new terminal and test:

```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"app":"Vinhomes Property Management","status":"healthy"}
```

---

## 🧪 Step 4: Test API Endpoints

### Register a user:
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin@test",
    "email": "admin@test.com",
    "password": "password123",
    "role": "admin"
  }'
```

### Login:
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin@test",
    "password": "password123"
  }'
```

Expected response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {...}
}
```

### Get profile (replace TOKEN with actual token):
```bash
curl http://localhost:8080/profile \
  -H "Authorization: Bearer TOKEN"
```

---

## 📋 Step 5: Check Database

Connect to PostgreSQL:

```bash
docker exec -it vinhomes_postgres psql -U vinhomes_user -d vinhomes_db

# In psql, check tables:
\dt

# Exit:
\q
```

Expected tables:
```
- users
- buildings
- floors
- apartments
- contracts
- residents
- invoices
- maintenance_requests
- notifications
- audit_logs
- (and more...)
```

---

## 📊 Quick Status Report

| Component | Status | Port |
|-----------|--------|------|
| PostgreSQL | ✅ Running | 5432 |
| Backend | ⏳ Ready to start | 8080 |
| API Docs | 📖 See API_REFERENCE.md | - |

---

## 🔧 Troubleshooting

### Issue: "port 5432 already in use"
```bash
docker-compose down
docker-compose up -d postgres
```

### Issue: "connection refused on port 8080"
```bash
# Check if backend is running
lsof -i :8080

# Kill existing process if needed
kill -9 <PID>

# Restart backend
go run ./cmd/main.go
```

### Issue: Database errors
```bash
# Check database logs
docker-compose logs postgres | tail -20

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

### Issue: Go module errors
```bash
go mod tidy
go mod download
```

---

## 📚 Next Steps

1. ✅ Run the backend locally
2. ✅ Test API endpoints
3. ⏳ Connect frontend to backend
4. ⏳ Deploy to production

---

## 📞 API Endpoints Available

**Authentication:**
- `POST /auth/register` - Register user
- `POST /auth/login` - Login & get JWT
- `GET /profile` - Get user profile

**Buildings:**
- `POST /buildings` - Create building
- `GET /buildings` - List buildings
- `GET /buildings/:id` - Get building details

**Contracts:**
- `POST /contracts` - Create contract
- `GET /contracts` - List contracts

**Invoices:**
- `POST /invoices` - Create invoice
- `GET /invoices` - List invoices

**Maintenance:**
- `POST /maintenance-requests` - Create maintenance request
- `GET /maintenance-requests` - List requests

And 40+ more endpoints! See API_REFERENCE.md for complete list.

---

## 🎯 You're all set!

Your backend is **ready to run**. Choose one command from Step 2 above to start!

```bash
cd /Users/namle/Desktop/uit/DOAN_3d/DoAn_3D/backend
go run ./cmd/main.go
```

Then in another terminal:
```bash
curl http://localhost:8080/health
```

Should return:
```json
{"app":"Vinhomes Property Management","status":"healthy"}
```

**Enjoy! 🚀**
