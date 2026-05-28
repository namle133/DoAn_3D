# Backend Setup & Deployment Guide

## Local Development Setup

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 13+ with PostGIS extension
- Make (optional, for using Makefile)
- Git

### Step 1: Clone Repository & Setup

```bash
# Navigate to backend directory
cd backend

# Copy environment template
cp config/.env.example config/.env

# Download Go dependencies
go mod download
go mod tidy
```

### Step 2: Configure Environment

Edit `config/.env`:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=vinhomes_user
DB_PASSWORD=your_secure_password
DB_NAME=vinhomes_db

# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# JWT Configuration
JWT_SECRET=your-very-secure-secret-key-change-in-production
JWT_EXPIRY_HOURS=24

# Email Configuration (for notifications)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

APP_ENV=development
```

### Step 3: Setup PostgreSQL Database

#### Option A: Using Docker (Recommended)

```bash
# Start PostgreSQL with PostGIS
docker-compose up -d postgres

# Verify connection
docker-compose ps

# View logs
docker-compose logs postgres
```

#### Option B: Manual Installation (macOS)

```bash
# Using Homebrew
brew install postgresql@15
brew install postgis

# Start PostgreSQL service
brew services start postgresql@15

# Create database user
psql -U postgres -c "CREATE USER vinhomes_user WITH PASSWORD 'vinhomes_pass';"

# Create database
psql -U postgres -c "CREATE DATABASE vinhomes_db OWNER vinhomes_user;"

# Connect to database and enable PostGIS
psql -U postgres -d vinhomes_db -c "CREATE EXTENSION IF NOT EXISTS postgis;"
psql -U postgres -d vinhomes_db -c "CREATE EXTENSION IF NOT EXISTS hstore;"
```

### Step 4: Run Application

#### Option A: Direct Execution

```bash
# Run migrations and start server
go run ./cmd/main.go
```

#### Option B: Using Make

```bash
# Install development tools (first time)
make install-tools

# Run with hot-reload
make dev

# Or build and run
make build
./vinhomes
```

#### Option C: Docker

```bash
# Build and start all services
docker-compose up -d

# Check services
docker-compose ps

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down
```

### Step 5: Verify Installation

```bash
# Check API health
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy","app":"Vinhomes Property Management"}
```

---

## Testing

### Run Unit Tests

```bash
make test

# Or directly
go test -v ./...
```

### Run Tests with Coverage

```bash
make test-coverage

# View coverage report
open coverage.html
```

---

## Database Management

### Run Migrations

Migrations run automatically on application startup.

### Seed Sample Data

```bash
make seed-db

# Or manually create test data:
# 1. Create admin user (see API examples)
# 2. Create buildings, floors, apartments
# 3. Create residents and contracts
```

### Check Database

```bash
# Connect to PostgreSQL
psql -U vinhomes_user -d vinhomes_db

# List tables
\dt

# View table schema
\d buildings

# Exit
\q
```

---

## Production Deployment

### Prerequisites
- Docker & Docker Compose
- Linux server (Ubuntu 20.04+ recommended)
- Domain name with SSL certificate
- Environment variables configured

### Step 1: Prepare Server

```bash
# Update system
sudo apt-get update
sudo apt-get upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### Step 2: Clone & Configure

```bash
# Clone repository
git clone <repository-url>
cd DoAn_3D/backend

# Copy and configure .env for production
cp config/.env.example config/.env

# Edit .env with production values
nano config/.env
```

**Production .env Settings:**
```
DB_HOST=postgres
DB_PORT=5432
DB_USER=vinhomes_user
DB_PASSWORD=<STRONG_RANDOM_PASSWORD>
DB_NAME=vinhomes_db

SERVER_PORT=8080
GIN_MODE=release

JWT_SECRET=<STRONG_RANDOM_SECRET_KEY>
JWT_EXPIRY_HOURS=24

SMTP_HOST=<YOUR_SMTP_HOST>
SMTP_PORT=587
SMTP_USER=<YOUR_EMAIL>
SMTP_PASSWORD=<YOUR_PASSWORD>

APP_ENV=production
```

### Step 3: Deploy with Docker Compose

```bash
# Build and start containers
docker-compose -f docker-compose.yml up -d

# Verify services
docker-compose ps

# View logs
docker-compose logs -f
```

### Step 4: Setup Reverse Proxy (Nginx)

```bash
# Install Nginx
sudo apt-get install -y nginx

# Create Nginx config
sudo nano /etc/nginx/sites-available/vinhomes

# Add configuration:
```

```nginx
upstream vinhomes_backend {
    server localhost:8080;
}

server {
    listen 80;
    server_name api.vinhomes.com;
    client_max_body_size 50M;

    location / {
        proxy_pass http://vinhomes_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/vinhomes /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Restart Nginx
sudo systemctl restart nginx
```

### Step 5: Setup SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt-get install -y certbot python3-certbot-nginx

# Generate certificate
sudo certbot --nginx -d api.vinhomes.com

# Auto-renewal verification
sudo certbot renew --dry-run
```

### Step 6: Setup Monitoring

```bash
# View application logs
docker-compose logs -f backend

# Check database status
docker-compose exec postgres pg_isready -U vinhomes_user

# Monitor resources
docker stats
```

---

## Troubleshooting

### Database Connection Error

```
error: "Failed to initialize database: connection refused"
```

**Solution:**
- Verify PostgreSQL is running: `docker-compose ps postgres`
- Check DB credentials in `.env` file
- Verify database exists: `psql -U vinhomes_user -d vinhomes_db -c "SELECT 1"`

### Port Already in Use

```
error: "listen tcp :8080: bind: address already in use"
```

**Solution:**
```bash
# Kill process using port 8080
sudo lsof -ti:8080 | xargs kill -9

# Or change port in .env
SERVER_PORT=8081
```

### JWT Token Expired

```
error: "Invalid token"
```

**Solution:**
- Generate new token by logging in again
- Increase JWT_EXPIRY_HOURS in .env if needed

### Database Migration Failed

**Solution:**
- Check if PostgreSQL has PostGIS extension: `psql -d vinhomes_db -c "SELECT version()"`
- Enable extension if missing: `psql -d vinhomes_db -c "CREATE EXTENSION postgis"`

---

## Performance Optimization

### Database Indexing

```sql
-- Create indexes for frequently queried columns
CREATE INDEX idx_contracts_status ON contracts(status);
CREATE INDEX idx_invoices_contract_id ON invoices(contract_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_maintenance_status ON maintenance_requests(status);
CREATE INDEX idx_residents_user_id ON residents(user_id);
```

### Connection Pooling

Edit `docker-compose.yml` to add PgBouncer:

```yaml
pgbouncer:
  image: pgbouncer/pgbouncer:latest
  environment:
    PGBOUNCER_DATABASES_HOST: postgres
    PGBOUNCER_DATABASES_PORT: 5432
```

### Caching

Consider adding Redis for caching:

```yaml
redis:
  image: redis:7-alpine
  ports:
    - "6379:6379"
```

---

## Backup & Recovery

### Database Backup

```bash
# Backup database
docker-compose exec postgres pg_dump -U vinhomes_user vinhomes_db > backup_$(date +%Y%m%d).sql

# Backup with compression
docker-compose exec postgres pg_dump -U vinhomes_user vinhomes_db | gzip > backup_$(date +%Y%m%d).sql.gz
```

### Database Restore

```bash
# Restore from backup
docker-compose exec -T postgres psql -U vinhomes_user vinhomes_db < backup_YYYYMMDD.sql

# From compressed backup
gunzip < backup_YYYYMMDD.sql.gz | docker-compose exec -T postgres psql -U vinhomes_user vinhomes_db
```

---

## Security Checklist

- [ ] Change JWT_SECRET to strong random value
- [ ] Change DB_PASSWORD to strong random password
- [ ] Use HTTPS/SSL certificates
- [ ] Enable database authentication
- [ ] Setup firewall rules
- [ ] Restrict API access to authorized IPs
- [ ] Regular security updates
- [ ] Backup database regularly
- [ ] Monitor application logs
- [ ] Setup alerting for errors
