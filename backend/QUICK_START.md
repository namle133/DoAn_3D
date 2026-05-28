# Quick Start Guide - Vinhomes Backend

## 5-Minute Setup

### 1. Prerequisites Check
```bash
# Verify Go is installed
go version  # Should be 1.21+

# Verify Git is installed
git --version
```

### 2. Navigate to Backend
```bash
cd backend
```

### 3. Setup Environment
```bash
# Copy environment template
cp config/.env.example config/.env
```

### 4. Download Dependencies
```bash
go mod download
go mod tidy
```

### 5. Start Database (Docker)
```bash
# Make sure Docker is running
docker-compose up -d postgres

# Wait 10 seconds for database to start
sleep 10

# Verify connection
docker-compose logs postgres
```

### 6. Run Application
```bash
go run ./cmd/main.go
```

Expected output:
```
✓ Configuration loaded
✓ Database connected successfully
✓ All migrations completed
✓ Scheduler started
✓ Starting server on :8080
```

### 7. Test API
```bash
# Check health
curl http://localhost:8080/health

# Response:
# {"app":"Vinhomes Property Management","status":"healthy"}
```

---

## Common Commands

### Build Application
```bash
go build -o vinhomes ./cmd
./vinhomes
```

### Run Tests
```bash
go test -v ./...
```

### Format Code
```bash
gofmt -s -w .
```

### Run with Makefile
```bash
make dev              # Run with hot-reload
make docker-up        # Start Docker containers
make docker-down      # Stop Docker containers
```

---

## API Quick Reference

### 1. Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123456",
    "role": "admin"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123456"
  }'

# Copy the token from response
```

### 3. Get Profile
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4. Create Building
```bash
curl -X POST http://localhost:8080/api/v1/buildings \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Vinhomes Tower A",
    "address": "Tay Ho, Hanoi",
    "year_built": 2023,
    "total_floors": 40,
    "description": "Modern residential complex"
  }'
```

---

## Project Layout

```
backend/
├── cmd/main.go                    Entry point
├── config/config.go               Configuration
├── internal/
│   ├── models/models.go          Database models
│   ├── handlers/                 API endpoints
│   ├── services/                 Business logic
│   ├── middleware/               Authentication
│   └── database/database.go      DB connection
├── docker-compose.yml            Container setup
├── Makefile                       Build commands
├── README.md                      Documentation
├── API_REFERENCE.md               Full API docs
├── DEPLOYMENT.md                  Production guide
└── IMPLEMENTATION_NOTES.md        Technical details
```

---

## Troubleshooting

### Error: "connection refused"
**Solution:** Start PostgreSQL
```bash
docker-compose up -d postgres
```

### Error: "address already in use :8080"
**Solution:** Change port in config/.env
```
SERVER_PORT=8081
```

### Database migration error
**Solution:** Ensure PostGIS is installed
```bash
docker-compose exec postgres psql -U vinhomes_user -d vinhomes_db -c "CREATE EXTENSION postgis;"
```

### Port issues with Docker
**Solution:** Check container status
```bash
docker-compose ps
docker-compose logs
```

---

## Next: Frontend Integration

After backend is running:

1. **Connect from Frontend**
   - Base URL: `http://localhost:8080/api/v1`
   - Include JWT token in Authorization header

2. **Real-time Updates** (Future)
   - Implement WebSocket for live notifications
   - Update 3D visualization in real-time

3. **GIS Integration** (Future)
   - Connect ArcGIS to backend APIs
   - Sync spatial data bidirectionally

---

## Useful Resources

- **API Documentation:** See `API_REFERENCE.md`
- **Deployment Guide:** See `DEPLOYMENT.md`
- **Technical Details:** See `IMPLEMENTATION_NOTES.md`
- **Go Documentation:** https://golang.org/doc
- **GORM Documentation:** https://gorm.io
- **PostGIS Documentation:** https://postgis.net

---

## Need Help?

1. Check logs: `docker-compose logs backend`
2. Review error message in API response
3. Check configuration in `config/.env`
4. Verify database connection: `docker-compose exec postgres pg_isready`
5. Test endpoint manually with curl/Postman

---

**Status:** ✅ Backend ready for development!
