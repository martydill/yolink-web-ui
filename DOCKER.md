# Docker Setup for YoLink Web UI

This document explains how to run the YoLink Web UI using Docker and Docker Compose.

## Quick Start

### Production Mode
```bash
# Start both services in production mode
make docker-up

# Access the application
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
```

### Development Mode
```bash
# Start both services in development mode (with hot reload)
make docker-dev

# Access the application
# Frontend: http://localhost:5173 (Vite dev server)
# Backend API: http://localhost:8080
```

### Stop Services
```bash
# Stop and remove all containers
make docker-down
```

## Environment Setup

1. Copy the example environment file:
```bash
cp backend/.env.example backend/.env
```

2. Edit `backend/.env` with your YoLink API credentials:
```bash
YOLINK_ACCESS_TOKEN=your_actual_token
YOLINK_SECRET_KEY=your_actual_secret
MQTT_USERNAME=your_mqtt_username
MQTT_PASSWORD=your_mqtt_password
```

## Available Commands

| Command | Description |
|---------|-------------|
| `make docker-up` | Start in production mode (frontend on port 3000) |
| `make docker-dev` | Start in development mode (frontend on port 5173) |
| `make docker-down` | Stop and remove all containers |
| `make docker-build` | Build Docker images |
| `make docker-logs` | View container logs |

## Architecture

### Production Mode
- **Backend**: Go application in Alpine Linux container
- **Frontend**: Built Svelte app served by Nginx
- **Networking**: Services communicate via Docker network
- **Volumes**: Environment file mounted read-only

### Development Mode
- **Backend**: Go application with live reload
- **Frontend**: Vite dev server with hot module replacement
- **Volumes**: Source code mounted for live development

## Port Mapping

| Service | Production Port | Development Port |
|---------|----------------|------------------|
| Frontend | 3000 | 5173 |
| Backend | 8080 | 8080 |

## Troubleshooting

### Check Service Status
```bash
# View logs
make docker-logs

# Check if containers are running
docker-compose ps
```

### Rebuild Images
```bash
# Force rebuild without cache
docker-compose build --no-cache
```

### Clean Up
```bash
# Remove all containers, networks, and images
docker-compose down --rmi all --volumes --remove-orphans
```

## Health Checks

Both services include health checks:
- **Backend**: Checks `/api/devices` endpoint
- **Frontend**: Checks root path availability

Containers will be marked as unhealthy if they fail to respond within the configured timeouts.