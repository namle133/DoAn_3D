#!/bin/bash

# Khởi tạo DB trên Docker Postgres (host port 5433)
# Chạy từ thư mục backend sau khi: docker-compose up -d postgres

set -e

cd "$(dirname "$0")" || exit 1

DB_HOST=127.0.0.1
DB_PORT=5433
export PGPASSWORD=postgres

echo "Connecting to Docker Postgres at ${DB_HOST}:${DB_PORT}..."

psql -h "$DB_HOST" -p "$DB_PORT" -U postgres -d postgres -v ON_ERROR_STOP=1 <<'EOF'
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

echo "✅ vinhomes_db ready on Docker (127.0.0.1:5433)"
