#!/bin/bash

# This script creates the vinhomes_user and database
# Run this once when setting up PostgreSQL

set -e

echo "Creating vinhomes_user and database..."

# Create the user with password and necessary privileges
PGPASSWORD=postgres psql -h localhost -U postgres -c "CREATE USER vinhomes_user WITH ENCRYPTED PASSWORD 'vinhomes_pass';"

# Grant privileges
PGPASSWORD=postgres psql -h localhost -U postgres -c "ALTER USER vinhomes_user WITH SUPERUSER CREATEDB CREATEROLE;"

# Create database
PGPASSWORD=postgres psql -h localhost -U postgres -c "CREATE DATABASE vinhomes_db OWNER vinhomes_user;"

echo "✅ Setup complete!"
echo "User: vinhomes_user"
echo "Password: vinhomes_pass"
echo "Database: vinhomes_db"
