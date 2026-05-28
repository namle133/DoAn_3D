-- PostgreSQL initialization script for Vinhomes
-- Runs automatically by docker-entrypoint-initdb.d/ via psql, so \gexec works.

-- ----------------------------------------------------------------------
-- 1. Create role vinhomes_user if it doesn't exist
--    PostgreSQL doesn't have "CREATE USER IF NOT EXISTS", use DO block.
-- ----------------------------------------------------------------------
DO
$$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'vinhomes_user') THEN
        CREATE ROLE vinhomes_user
            LOGIN
            ENCRYPTED PASSWORD 'vinhomes_pass';
    END IF;
END
$$;

-- ----------------------------------------------------------------------
-- 2. Grant superuser-ish privileges to vinhomes_user
-- ----------------------------------------------------------------------
ALTER ROLE vinhomes_user WITH SUPERUSER CREATEDB CREATEROLE;

-- ----------------------------------------------------------------------
-- 3. Create database vinhomes_db if it doesn't exist
--    CREATE DATABASE cannot run inside a transaction / DO block,
--    so we use the psql \gexec trick.
-- ----------------------------------------------------------------------
SELECT 'CREATE DATABASE vinhomes_db OWNER vinhomes_user'
WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'vinhomes_db')
\gexec

-- ----------------------------------------------------------------------
-- 4. Ensure ownership + privileges
-- ----------------------------------------------------------------------
ALTER DATABASE vinhomes_db OWNER TO vinhomes_user;
GRANT ALL PRIVILEGES ON DATABASE vinhomes_db TO vinhomes_user;
