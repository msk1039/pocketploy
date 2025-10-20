#!/bin/bash

# PocketPloy VPS Dependencies Installation Script
# Installs all required dependencies on Ubuntu 20.04/22.04

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
cat << "EOF"
╔═══════════════════════════════════════════╗
║   PocketPloy Dependencies Installer      ║
╚═══════════════════════════════════════════╝
EOF
echo -e "${NC}"

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Error: This script must be run with sudo${NC}"
    echo "Usage: sudo ./install-dependencies.sh"
    exit 1
fi

# Check if running on Ubuntu
if ! grep -q "Ubuntu" /etc/os-release; then
    echo -e "${YELLOW}Warning: This script is designed for Ubuntu${NC}"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo "=========================================="
echo "Step 1: System Update"
echo "=========================================="
echo ""

apt update
apt upgrade -y

echo -e "${GREEN}✓ System updated${NC}"
echo ""

echo "=========================================="
echo "Step 2: Install Essential Tools"
echo "=========================================="
echo ""

apt install -y \
    git \
    curl \
    wget \
    vim \
    nano \
    htop \
    net-tools \
    build-essential \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    gnupg \
    lsb-release

echo -e "${GREEN}✓ Essential tools installed${NC}"
echo ""

echo "=========================================="
echo "Step 3: Install PostgreSQL"
echo "=========================================="
echo ""

apt install -y postgresql postgresql-contrib

# Start and enable PostgreSQL
systemctl start postgresql
systemctl enable postgresql

echo -e "${GREEN}✓ PostgreSQL installed and started${NC}"
echo ""

echo "=========================================="
echo "Step 4: Install Go 1.23"
echo "=========================================="
echo ""

# Remove old Go installation
rm -rf /usr/local/go

# Download and install Go
cd /tmp
GO_VERSION="1.23.0"
wget "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
rm "go${GO_VERSION}.linux-amd64.tar.gz"

# Add Go to PATH for all users
cat > /etc/profile.d/go.sh << 'EOF'
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
EOF

chmod +x /etc/profile.d/go.sh
source /etc/profile.d/go.sh

# Verify installation
if /usr/local/go/bin/go version; then
    echo -e "${GREEN}✓ Go installed: $(/usr/local/go/bin/go version)${NC}"
else
    echo -e "${RED}✗ Go installation failed${NC}"
    exit 1
fi

echo ""

echo "=========================================="
echo "Step 5: Install Docker"
echo "=========================================="
echo ""

# Add Docker's official GPG key
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

# Add the repository to Apt sources
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

# Update package index
apt update

# Install Docker
apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Start and enable Docker
systemctl start docker
systemctl enable docker

# Add current user to docker group
if [ ! -z "$SUDO_USER" ]; then
    usermod -aG docker $SUDO_USER
    echo -e "${YELLOW}Note: User $SUDO_USER added to docker group${NC}"
    echo -e "${YELLOW}You need to logout and login again for this to take effect${NC}"
fi

echo -e "${GREEN}✓ Docker installed${NC}"
echo ""

echo "=========================================="
echo "Step 6: Install Docker Compose"
echo "=========================================="
echo ""

# Docker Compose v2 is already installed as a plugin
# Create symlink for docker-compose command
ln -sf /usr/libexec/docker/cli-plugins/docker-compose /usr/local/bin/docker-compose

# Verify installation
if docker-compose version; then
    echo -e "${GREEN}✓ Docker Compose installed${NC}"
else
    echo -e "${RED}✗ Docker Compose installation failed${NC}"
    exit 1
fi

echo ""

echo "=========================================="
echo "Step 7: Install Node.js 20.x"
echo "=========================================="
echo ""

# Add NodeSource repository
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -

# Install Node.js
apt install -y nodejs

# Verify installation
if node --version; then
    echo -e "${GREEN}✓ Node.js installed: $(node --version)${NC}"
else
    echo -e "${RED}✗ Node.js installation failed${NC}"
    exit 1
fi

echo ""

echo "=========================================="
echo "Step 8: Install pnpm"
echo "=========================================="
echo ""

# Install pnpm
npm install -g pnpm

# Verify installation
if pnpm --version; then
    echo -e "${GREEN}✓ pnpm installed: $(pnpm --version)${NC}"
else
    echo -e "${RED}✗ pnpm installation failed${NC}"
    exit 1
fi

echo ""

echo "=========================================="
echo "Step 9: Install Nginx"
echo "=========================================="
echo ""

apt install -y nginx

# Start and enable Nginx
systemctl start nginx
systemctl enable nginx

echo -e "${GREEN}✓ Nginx installed and started${NC}"
echo ""

echo "=========================================="
echo "Step 10: Pull PocketBase Docker Image"
echo "=========================================="
echo ""

docker pull ghcr.io/muchobien/pocketbase:latest

echo -e "${GREEN}✓ PocketBase image downloaded${NC}"
echo ""

echo "=========================================="
echo "Step 11: Configure Firewall (UFW)"
echo "=========================================="
echo ""

apt install -y ufw

# Allow SSH
ufw allow OpenSSH

# Allow HTTP and HTTPS
ufw allow 'Nginx Full'

# Don't enable yet, let user do it manually
echo -e "${YELLOW}UFW installed but not enabled${NC}"
echo "To enable: sudo ufw enable"
echo ""

echo "=========================================="
echo "Installation Summary"
echo "=========================================="
echo ""

echo "Installed versions:"
echo "  PostgreSQL: $(psql --version | head -n1)"
echo "  Go: $(/usr/local/go/bin/go version)"
echo "  Docker: $(docker --version)"
echo "  Docker Compose: $(docker-compose version --short)"
echo "  Node.js: $(node --version)"
echo "  pnpm: $(pnpm --version)"
echo "  Nginx: $(nginx -v 2>&1)"
echo ""

echo -e "${GREEN}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}Installation Complete! ✅${NC}"
echo -e "${GREEN}═══════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Important:${NC}"
echo "1. Logout and login again for Docker group to take effect"
echo "2. Run 'source /etc/profile.d/go.sh' to load Go paths"
echo "3. Verify Docker: 'docker ps'"
echo "4. Verify Go: 'go version'"
echo ""

echo "Next steps:"
echo "1. Logout and login: exit"
echo "2. Run PostgreSQL setup: sudo ./scripts/setup-postgres.sh"
echo "3. Deploy application: ./scripts/deploy-vps.sh"
echo ""
