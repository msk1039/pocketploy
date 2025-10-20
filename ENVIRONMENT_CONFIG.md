# Environment Configuration Guide

This project supports both **local development** and **production VPS** deployments with different configurations.

## Local Development (HTTP with nip.io)

### Backend `.env`
```bash
ENV=development
BASE_DOMAIN=192.168.197.165.nip.io  # Your local IP
ALLOWED_ORIGINS=http://localhost:3000
DOCKER_HOST=unix:///Users/meow/.colima/default/docker.sock  # Or your Docker socket
```

### Frontend `.env.local`
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Traefik & Docker Compose
Use the default files:
- `traefik.yml` (HTTP only on port 80)
- `docker-compose.yml` (HTTP on port 80, dashboard on 8081)

### Instance URLs
Instances will be accessible at:
```
http://username-instancename.192.168.197.165.nip.io
```

### Start Development
```bash
# Start Traefik
docker-compose up -d

# Start backend
cd backend
go run cmd/server/main.go

# Start frontend
cd frontend
pnpm dev
```

---

## Production VPS (HTTPS with domain)

### Backend `.env`
```bash
ENV=production
PORT=8080
HOST=0.0.0.0
BASE_DOMAIN=pocketploy.maykad.tech
ALLOWED_ORIGINS=https://pocketploy.maykad.tech
DOCKER_HOST=unix:///var/run/docker.sock
DB_USER=pocketploy_user
DB_PASSWORD=your_secure_password
```

### Frontend `.env.local`
```bash
NEXT_PUBLIC_API_URL=https://pocketploy.maykad.tech/api/v1
```

### Traefik & Docker Compose
Use production files:
```bash
# On VPS, use production config
docker-compose -f docker-compose.production.yml up -d
```

This will:
- Run HTTP on port 8090 (with redirect to HTTPS)
- Run HTTPS on port 8091
- Run dashboard on port 8081

### Nginx Configuration
Copy the production Nginx config:
```bash
sudo cp nginx.production.conf /etc/nginx/sites-available/pocketploy
sudo ln -sf /etc/nginx/sites-available/pocketploy /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Instance URLs
Instances will be accessible at:
```
https://username-instancename.pocketploy.maykad.tech
```

### Deploy to Production
```bash
# Pull latest changes
git pull origin main

# Rebuild backend
cd backend
go build -o pocketploy-backend cmd/server/main.go
sudo systemctl restart pocketploy-backend

# Rebuild frontend
cd ../frontend
pnpm install
pnpm build
sudo systemctl restart pocketploy-frontend

# Restart Traefik
cd ..
docker-compose -f docker-compose.production.yml down
docker-compose -f docker-compose.production.yml up -d
```

---

## Key Differences

| Feature | Development | Production |
|---------|------------|------------|
| Protocol | HTTP | HTTPS |
| Domain | nip.io | pocketploy.maykad.tech |
| Traefik Ports | 80 (HTTP) | 8090 (HTTP), 8091 (HTTPS) |
| SSL/TLS | None | Cloudflare Origin Certificates |
| Frontend URL | http://localhost:3000 | https://pocketploy.maykad.tech |
| API URL | http://localhost:8080 | https://pocketploy.maykad.tech/api/ |
| Instance URLs | http://user-name.IP.nip.io | https://user-name.pocketploy.maykad.tech |

---

## Switching Between Environments

The code automatically detects the environment based on the `ENV` variable in `.env`:
- `ENV=development` → HTTP, nip.io URLs
- `ENV=production` → HTTPS, domain URLs

No code changes needed - just update your `.env` files!
