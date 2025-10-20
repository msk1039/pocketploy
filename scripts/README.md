# PocketPloy Deployment Scripts

This directory contains scripts to help you deploy PocketPloy on your VPS.

## Scripts Overview

### 1. `install-dependencies.sh`
**Purpose**: Install all required dependencies on a fresh Ubuntu VPS.

**What it installs**:
- PostgreSQL 14+
- Go 1.23
- Docker & Docker Compose
- Node.js 20.x
- pnpm
- Nginx
- Essential tools (git, curl, wget, etc.)

**Usage**:
```bash
sudo ./scripts/install-dependencies.sh
```

**Note**: After running, logout and login again for Docker group to take effect.

---

### 2. `setup-postgres.sh`
**Purpose**: One-time PostgreSQL setup - creates database, user, and runs migrations.

**What it does**:
- Creates PostgreSQL user `pocketploy_user`
- Creates database `pocketploy`
- Runs all migration files (001-005)
- Generates `.env.production.example` template
- Verifies tables are created

**Usage**:
```bash
sudo ./scripts/setup-postgres.sh
```

**After running**:
1. Copy template: `cp backend/.env.production.example backend/.env`
2. Edit `.env` and update passwords and secrets
3. Test connection: `psql -h localhost -U pocketploy_user -d pocketploy`

---

### 3. `deploy-vps.sh`
**Purpose**: Automated deployment of PocketPloy application.

**What it does**:
- Checks system dependencies
- Verifies PostgreSQL setup
- Builds backend (Go binary)
- Builds frontend (Next.js)
- Creates Docker network
- Starts Traefik container
- Creates systemd services
- Starts all services

**Usage**:
```bash
./scripts/deploy-vps.sh
```

**Prerequisites**:
- Dependencies installed (`install-dependencies.sh`)
- PostgreSQL setup complete (`setup-postgres.sh`)
- `.env` files configured

---

## Deployment Workflow

### First-time Setup on Fresh VPS

```bash
# 1. SSH into your VPS
ssh root@your-vps-ip

# 2. Clone repository
git clone <your-repo-url> pocketploy
cd pocketploy

# 3. Install dependencies
sudo ./scripts/install-dependencies.sh

# 4. Logout and login again (for Docker group)
exit
ssh root@your-vps-ip
cd pocketploy

# 5. Setup PostgreSQL
sudo ./scripts/setup-postgres.sh

# 6. Configure environment
cp backend/.env.production.example backend/.env
nano backend/.env  # Update passwords and secrets

# 7. Deploy application
./scripts/deploy-vps.sh

# 8. Configure Nginx (see VPS_DEPLOYMENT_GUIDE.md)
sudo nano /etc/nginx/sites-available/pocketploy
# ... add configuration
sudo ln -s /etc/nginx/sites-available/pocketploy /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

# 9. Optional: Setup SSL
sudo certbot --nginx -d maykad.tech -d "*.maykad.tech"
```

---

## Script Execution Order

```
┌─────────────────────────────┐
│ 1. install-dependencies.sh  │  (Run once, as root)
└─────────────┬───────────────┘
              │
              ▼
      ┌───────────────┐
      │ Logout/Login  │  (For Docker group)
      └───────┬───────┘
              │
              ▼
┌─────────────────────────────┐
│ 2. setup-postgres.sh        │  (Run once, as root)
└─────────────┬───────────────┘
              │
              ▼
      ┌───────────────┐
      │ Edit .env     │  (Configure secrets)
      └───────┬───────┘
              │
              ▼
┌─────────────────────────────┐
│ 3. deploy-vps.sh            │  (Run as normal user)
└─────────────┬───────────────┘
              │
              ▼
      ┌───────────────┐
      │ Configure     │  (Nginx, SSL, etc.)
      │ Web Server    │
      └───────────────┘
```

---

## Environment Variables

### Backend `.env` (Required)

```bash
# Server
PORT=8080
ENV=production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=pocketploy_user
DB_PASSWORD=your_secure_password
DB_NAME=pocketploy
DB_SSLMODE=disable

# JWT
JWT_SECRET=your_jwt_secret_32_chars_minimum
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Docker
DOCKER_HOST=unix:///var/run/docker.sock
POCKETBASE_IMAGE=ghcr.io/muchobien/pocketbase:latest

# Instance
BASE_DOMAIN=maykad.tech
TRAEFIK_NETWORK=pocketploy-network
MAX_INSTANCES_PER_USER=5
```

### Frontend `.env.local` (Required)

```bash
NEXT_PUBLIC_API_URL=https://api.maykad.tech/api/v1
```

---

## Systemd Services

After running `deploy-vps.sh`, these services are created:

### Backend Service
```bash
sudo systemctl status pocketploy-backend
sudo systemctl restart pocketploy-backend
sudo journalctl -u pocketploy-backend -f
```

### Frontend Service
```bash
sudo systemctl status pocketploy-frontend
sudo systemctl restart pocketploy-frontend
sudo journalctl -u pocketploy-frontend -f
```

---

## Troubleshooting

### Dependencies Installation Failed
```bash
# Check error messages in the output
# Common issues:
# - Insufficient permissions: Use sudo
# - Network issues: Check internet connection
# - Package conflicts: Update system first
sudo apt update && sudo apt upgrade -y
```

### PostgreSQL Setup Failed
```bash
# Check if PostgreSQL is running
sudo systemctl status postgresql

# Check logs
sudo tail -f /var/log/postgresql/postgresql-*.log

# Verify user can connect
psql -h localhost -U postgres -l
```

### Deployment Failed
```bash
# Check build errors
# Backend:
cd backend
go build -o pocketploy-backend cmd/server/main.go

# Frontend:
cd frontend
pnpm install
pnpm build

# Check service status
sudo systemctl status pocketploy-backend
sudo systemctl status pocketploy-frontend
```

---

## Manual Steps Required

These steps are **not automated** and must be done manually:

1. **Nginx Configuration**: Create and configure `/etc/nginx/sites-available/pocketploy`
2. **SSL Certificates**: Run certbot for HTTPS
3. **Firewall**: Enable UFW with `sudo ufw enable`
4. **Environment Secrets**: Update passwords and JWT secrets
5. **Domain DNS**: Ensure wildcard DNS points to your VPS IP

Refer to `docs/VPS_DEPLOYMENT_GUIDE.md` for detailed instructions.

---

## Additional Resources

- **Complete Guide**: `docs/VPS_DEPLOYMENT_GUIDE.md`
- **Quick Start**: `docs/VPS_QUICK_START.md`
- **Testing**: `docs/TESTING_PHASE2.md`
- **API Docs**: `docs/API_INSTANCES.md`

---

## Security Notes

⚠️ **Important Security Steps**:

1. Change all default passwords in `.env` files
2. Use strong JWT secrets (minimum 32 characters)
3. Setup SSL/HTTPS with Let's Encrypt
4. Enable and configure UFW firewall
5. Disable root SSH login (use SSH keys)
6. Keep system packages updated
7. Setup automated backups

---

## Support

If you encounter issues:

1. Check service logs: `sudo journalctl -u pocketploy-backend -n 50`
2. Verify all services are running: `systemctl status pocketploy-*`
3. Check Nginx: `sudo nginx -t`
4. Review deployment guide: `docs/VPS_DEPLOYMENT_GUIDE.md`
5. Test connectivity: `curl http://localhost:8080/api/v1/health`

---

**Last Updated**: October 2025
