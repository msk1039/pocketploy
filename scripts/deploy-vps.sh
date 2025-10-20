#!/bin/bash

# PocketPloy VPS Deployment Helper Script
# This script helps deploy PocketPloy on a fresh VPS

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
cat << "EOF"
╔═══════════════════════════════════════════╗
║   PocketPloy VPS Deployment Helper       ║
╚═══════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Get current directory
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "Project root: $PROJECT_ROOT"
echo ""

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print step
print_step() {
    echo ""
    echo -e "${BLUE}═══════════════════════════════════════════${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════${NC}"
    echo ""
}

# Check if running on VPS (not macOS)
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo -e "${RED}Error: This script is for Linux VPS deployment${NC}"
    echo "You are running on macOS. This script should be run on your VPS."
    exit 1
fi

print_step "Step 1: System Check"

# Check for required commands
REQUIRED_COMMANDS=("docker" "docker-compose" "psql" "go" "node" "pnpm" "nginx")
MISSING_COMMANDS=()

for cmd in "${REQUIRED_COMMANDS[@]}"; do
    if command_exists "$cmd"; then
        echo -e "${GREEN}✓ $cmd installed${NC}"
    else
        echo -e "${RED}✗ $cmd not found${NC}"
        MISSING_COMMANDS+=("$cmd")
    fi
done

if [ ${#MISSING_COMMANDS[@]} -gt 0 ]; then
    echo ""
    echo -e "${RED}Missing dependencies: ${MISSING_COMMANDS[*]}${NC}"
    echo "Please follow the VPS_DEPLOYMENT_GUIDE.md to install dependencies"
    exit 1
fi

print_step "Step 2: PostgreSQL Setup Check"

# Check if PostgreSQL is running
if ! systemctl is-active --quiet postgresql; then
    echo -e "${YELLOW}PostgreSQL is not running${NC}"
    echo "Starting PostgreSQL..."
    sudo systemctl start postgresql
fi

# Check if database exists
if sudo -u postgres psql -lqt | cut -d \| -f 1 | grep -qw pocketploy; then
    echo -e "${GREEN}✓ Database 'pocketploy' exists${NC}"
else
    echo -e "${YELLOW}Database 'pocketploy' not found${NC}"
    echo "Running PostgreSQL setup script..."
    sudo "$SCRIPT_DIR/setup-postgres.sh"
fi

print_step "Step 3: Backend Configuration"

cd "$PROJECT_ROOT/backend"

# Check if .env exists
if [ ! -f ".env" ]; then
    if [ -f ".env.production.example" ]; then
        echo "Creating .env from template..."
        cp .env.production.example .env
        echo -e "${YELLOW}⚠ Please edit backend/.env and update:${NC}"
        echo "  - JWT_SECRET"
        echo "  - DB_PASSWORD"
        echo "  - BASE_DOMAIN"
        read -p "Press Enter when done..."
    else
        echo -e "${RED}No .env or template found${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✓ .env file exists${NC}"
fi

# Build backend
echo "Building backend..."
go build -o pocketploy-backend cmd/server/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Backend built successfully${NC}"
else
    echo -e "${RED}✗ Backend build failed${NC}"
    exit 1
fi

print_step "Step 4: Frontend Configuration"

cd "$PROJECT_ROOT/frontend"

# Check if .env.local exists
if [ ! -f ".env.local" ]; then
    echo "Creating .env.local..."
    cat > .env.local << EOF
NEXT_PUBLIC_API_URL=https://api.maykad.tech/api/v1
EOF
    echo -e "${YELLOW}⚠ Update NEXT_PUBLIC_API_URL in frontend/.env.local if needed${NC}"
else
    echo -e "${GREEN}✓ .env.local exists${NC}"
fi

# Install dependencies and build
echo "Installing frontend dependencies..."
pnpm install

echo "Building frontend..."
pnpm build

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Frontend built successfully${NC}"
else
    echo -e "${RED}✗ Frontend build failed${NC}"
    exit 1
fi

print_step "Step 5: Docker Setup"

cd "$PROJECT_ROOT"

# Create Docker network
if docker network ls | grep -q pocketploy-network; then
    echo -e "${GREEN}✓ Docker network 'pocketploy-network' exists${NC}"
else
    echo "Creating Docker network..."
    docker network create pocketploy-network
fi

# Start Traefik
echo "Starting Traefik..."
docker-compose up -d traefik

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Traefik started${NC}"
else
    echo -e "${RED}✗ Failed to start Traefik${NC}"
    exit 1
fi

print_step "Step 6: Systemd Services"

# Backend service
echo "Creating backend systemd service..."
sudo tee /etc/systemd/system/pocketploy-backend.service > /dev/null << EOF
[Unit]
Description=PocketPloy Backend Service
After=network.target postgresql.service docker.service

[Service]
Type=simple
User=$USER
WorkingDirectory=$PROJECT_ROOT/backend
ExecStart=$PROJECT_ROOT/backend/pocketploy-backend
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=pocketploy-backend

Environment="PATH=/usr/local/go/bin:/usr/bin:/bin"
Environment="GOPATH=$HOME/go"

[Install]
WantedBy=multi-user.target
EOF

# Frontend service
echo "Creating frontend systemd service..."
sudo tee /etc/systemd/system/pocketploy-frontend.service > /dev/null << EOF
[Unit]
Description=PocketPloy Frontend Service
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$PROJECT_ROOT/frontend
ExecStart=/usr/bin/pnpm start
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=pocketploy-frontend

Environment="NODE_ENV=production"
Environment="PORT=3000"

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
sudo systemctl daemon-reload

# Enable services
echo "Enabling services..."
sudo systemctl enable pocketploy-backend
sudo systemctl enable pocketploy-frontend

# Start services
echo "Starting services..."
sudo systemctl start pocketploy-backend
sudo systemctl start pocketploy-frontend

# Wait a bit
sleep 3

# Check status
if systemctl is-active --quiet pocketploy-backend; then
    echo -e "${GREEN}✓ Backend service is running${NC}"
else
    echo -e "${RED}✗ Backend service failed to start${NC}"
    echo "Check logs: sudo journalctl -u pocketploy-backend -n 50"
fi

if systemctl is-active --quiet pocketploy-frontend; then
    echo -e "${GREEN}✓ Frontend service is running${NC}"
else
    echo -e "${RED}✗ Frontend service failed to start${NC}"
    echo "Check logs: sudo journalctl -u pocketploy-frontend -n 50"
fi

print_step "Step 7: Nginx Configuration"

echo "Nginx configuration needs to be set up manually."
echo "Please refer to VPS_DEPLOYMENT_GUIDE.md Step 6 for Nginx setup."
echo ""
echo "Quick summary:"
echo "1. Create /etc/nginx/sites-available/pocketploy"
echo "2. Link to sites-enabled"
echo "3. Test: sudo nginx -t"
echo "4. Reload: sudo systemctl reload nginx"

print_step "Deployment Complete!"

echo -e "${GREEN}PocketPloy has been deployed!${NC}"
echo ""
echo "Services status:"
echo "  Backend:  $(systemctl is-active pocketploy-backend)"
echo "  Frontend: $(systemctl is-active pocketploy-frontend)"
echo "  Traefik:  $(docker inspect -f '{{.State.Status}}' pocketploy-traefik 2>/dev/null || echo 'not found')"
echo ""
echo "Next steps:"
echo "1. Configure Nginx (see VPS_DEPLOYMENT_GUIDE.md)"
echo "2. Setup SSL with certbot"
echo "3. Configure firewall"
echo "4. Test the application"
echo ""
echo "Useful commands:"
echo "  Backend logs:   sudo journalctl -u pocketploy-backend -f"
echo "  Frontend logs:  sudo journalctl -u pocketploy-frontend -f"
echo "  Traefik logs:   docker logs -f pocketploy-traefik"
echo "  Restart backend: sudo systemctl restart pocketploy-backend"
echo ""
