# Vinhomes Property Management Backend

## Project Setup

### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- Docker & Docker Compose (optional)
- Git

### Installation

1. **Clone and Setup**
```bash
cd backend
cp config/.env.example config/.env
```

2. **Edit `.env` file** with your configuration

3. **Install dependencies**
```bash
go mod download
go mod tidy
```

4. **Run Database (Docker)**
```bash
docker-compose up -d postgres
```

Or install PostgreSQL locally with PostGIS extension.

5. **Run Application**
```bash
go run ./cmd/main.go
```

Or build and run:
```bash
go build -o vinhomes ./cmd
./vinhomes
```

6. **Docker Compose (Full Stack)**
```bash
docker-compose up -d
```
