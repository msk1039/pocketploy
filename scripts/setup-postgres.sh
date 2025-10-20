#!/bin/bash

# PocketPloy PostgreSQL One-Time Setup Script
# This script sets up PostgreSQL database, user, and runs all migrations

set -e  # Exit on any error

echo "=========================================="
echo "PocketPloy PostgreSQL Setup"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
DB_NAME="pocketploy"
DB_USER="pocketploy_user"
DB_PASSWORD="pocketploy_secure_password_2024"  # Change this!
MIGRATIONS_DIR="$(cd "$(dirname "$0")/../backend/internal/database/migrations" && pwd)"

echo -e "${YELLOW}Configuration:${NC}"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo "  Migrations: $MIGRATIONS_DIR"
echo ""

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Error: This script must be run with sudo${NC}"
    echo "Usage: sudo ./setup-postgres.sh"
    exit 1
fi

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo -e "${RED}Error: PostgreSQL is not installed${NC}"
    echo "Install it with: sudo apt install postgresql postgresql-contrib"
    exit 1
fi

# Check if PostgreSQL is running
if ! systemctl is-active --quiet postgresql; then
    echo -e "${YELLOW}PostgreSQL is not running. Starting...${NC}"
    systemctl start postgresql
    sleep 2
fi

echo -e "${GREEN}✓ PostgreSQL is running${NC}"
echo ""

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo -e "${RED}Error: Migrations directory not found: $MIGRATIONS_DIR${NC}"
    exit 1
fi

echo "=========================================="
echo "Step 1: Creating Database and User"
echo "=========================================="
echo ""

# Create user and database
sudo -u postgres psql << EOF
-- Drop database if exists (for clean setup)
DROP DATABASE IF EXISTS $DB_NAME;

-- Drop user if exists
DROP USER IF EXISTS $DB_USER;

-- Create user
CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';

-- Create database
CREATE DATABASE $DB_NAME OWNER $DB_USER;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;

\c $DB_NAME

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO $DB_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO $DB_USER;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO $DB_USER;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\q
EOF

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Database and user created successfully${NC}"
else
    echo -e "${RED}✗ Failed to create database and user${NC}"
    exit 1
fi

echo ""
echo "=========================================="
echo "Step 2: Running Migrations"
echo "=========================================="
echo ""

# Run migrations in order
MIGRATION_FILES=(
    "001_create_users_table.sql"
    "002_create_refresh_tokens_table.sql"
    "003_create_instances_table.sql"
    "004_create_instances_archive_table.sql"
    "005_update_instances_status_constraint.sql"
)

for migration in "${MIGRATION_FILES[@]}"; do
    MIGRATION_PATH="$MIGRATIONS_DIR/$migration"
    
    if [ ! -f "$MIGRATION_PATH" ]; then
        echo -e "${RED}✗ Migration file not found: $migration${NC}"
        exit 1
    fi
    
    echo "Running: $migration"
    
    sudo -u postgres psql -d $DB_NAME -f "$MIGRATION_PATH"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $migration completed${NC}"
    else
        echo -e "${RED}✗ $migration failed${NC}"
        exit 1
    fi
    echo ""
done

echo "=========================================="
echo "Step 3: Verification"
echo "=========================================="
echo ""

# Verify tables
echo "Checking created tables..."
TABLES=$(sudo -u postgres psql -d $DB_NAME -t -c "SELECT tablename FROM pg_tables WHERE schemaname='public' ORDER BY tablename;")

if [ -z "$TABLES" ]; then
    echo -e "${RED}✗ No tables found in database${NC}"
    exit 1
fi

echo -e "${GREEN}Created tables:${NC}"
echo "$TABLES"
echo ""

# Check specific tables
REQUIRED_TABLES=("users" "refresh_tokens" "instances" "instances_archive")
MISSING_TABLES=()

for table in "${REQUIRED_TABLES[@]}"; do
    if echo "$TABLES" | grep -q "$table"; then
        echo -e "${GREEN}✓ Table '$table' exists${NC}"
    else
        MISSING_TABLES+=("$table")
        echo -e "${RED}✗ Table '$table' is missing${NC}"
    fi
done

echo ""

if [ ${#MISSING_TABLES[@]} -gt 0 ]; then
    echo -e "${RED}Error: Some required tables are missing${NC}"
    exit 1
fi

echo "=========================================="
echo "Step 4: Creating .env Template"
echo "=========================================="
echo ""

BACKEND_DIR="$(cd "$(dirname "$0")/../backend" && pwd)"
ENV_TEMPLATE="$BACKEND_DIR/.env.production.example"

cat > "$ENV_TEMPLATE" << EOF
# Server Configuration
PORT=8080
ENV=production

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
DB_SSLMODE=disable

# JWT Configuration (CHANGE THESE IN PRODUCTION!)
JWT_SECRET=your_jwt_secret_key_change_this_in_production_$(openssl rand -hex 32)
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Docker Configuration
DOCKER_HOST=unix:///var/run/docker.sock
POCKETBASE_IMAGE=ghcr.io/muchobien/pocketbase:latest

# Instance Configuration
BASE_DOMAIN=maykad.tech
TRAEFIK_NETWORK=pocketploy-network
MAX_INSTANCES_PER_USER=5

# Storage
INSTANCES_DIR=./instances
EOF

echo -e "${GREEN}✓ Created .env template at: $ENV_TEMPLATE${NC}"
echo ""

echo "=========================================="
echo "Setup Complete! ✅"
echo "=========================================="
echo ""
echo -e "${GREEN}Database setup completed successfully!${NC}"
echo ""
echo "Next steps:"
echo "1. Copy the .env template to .env:"
echo "   cp $ENV_TEMPLATE $BACKEND_DIR/.env"
echo ""
echo "2. Edit .env and update these values:"
echo "   - JWT_SECRET (use a secure random string)"
echo "   - BASE_DOMAIN (set to your domain)"
echo ""
echo "3. Test database connection:"
echo "   psql -h localhost -U $DB_USER -d $DB_NAME"
echo "   Password: $DB_PASSWORD"
echo ""
echo -e "${YELLOW}IMPORTANT:${NC}"
echo "  - Change the DB_PASSWORD in production!"
echo "  - The password is currently: $DB_PASSWORD"
echo "  - Store credentials securely"
echo ""
echo "Database connection string:"
echo "postgresql://$DB_USER:$DB_PASSWORD@localhost:5432/$DB_NAME"
echo ""
