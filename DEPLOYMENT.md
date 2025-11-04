# NekoBox Docker Deployment Guide

## Overview
This document details the complete setup and deployment process for NekoBox using Docker, including all debugging steps and solutions encountered during deployment on a cloud server with BT Panel.

## System Requirements
- Docker installed
- MySQL 5.7+ or PostgreSQL 6.0+
- Redis
- Linux server (tested on Ubuntu with BT Panel)

## Initial Setup

### 1. Database Configuration

#### MySQL Setup
```sql
-- Create database
CREATE DATABASE IF NOT EXISTS nekobox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create user with local access
CREATE USER 'nekobox'@'localhost' IDENTIFIED BY 'LLWAJRM7WL2xiSEG';
GRANT ALL PRIVILEGES ON nekobox.* TO 'nekobox'@'localhost';

-- IMPORTANT: For Docker bridge network access, also grant from Docker network
CREATE USER IF NOT EXISTS 'nekobox'@'172.17.%' IDENTIFIED BY 'LLWAJRM7WL2xiSEG';
GRANT ALL PRIVILEGES ON nekobox.* TO 'nekobox'@'172.17.%';

-- OR grant from all hosts (simpler but less secure)
CREATE USER IF NOT EXISTS 'nekobox'@'%' IDENTIFIED BY 'LLWAJRM7WL2xiSEG';
GRANT ALL PRIVILEGES ON nekobox.* TO 'nekobox'@'%';

FLUSH PRIVILEGES;
```

**Note:** The database tables will be created automatically by GORM's AutoMigrate feature on first startup.

#### Redis Configuration

Edit `/www/server/redis/redis.conf` (or your Redis config path):

```conf
# Allow connections from localhost and Docker bridge network
bind 127.0.0.1 172.17.0.1

# Disable protected mode to allow Docker network access
protected-mode no

# Port configuration
port 6379
```

Restart Redis after configuration:
```bash
systemctl restart redis
# or
service redis-server restart
```

### 2. Application Configuration

Edit `conf/app.ini` with the following settings for Docker bridge networking:

```ini
[app]
production = true
title = "NekoBox"
external_url = "box.ethanpan.me"

[server]
port = 8080

[database]
user = "nekobox"
password = "LLWAJRM7WL2xiSEG"
host = "172.17.0.1"  # Docker bridge gateway
port = 3306
name = "nekobox"

[redis]
addr = "172.17.0.1:6379"  # Docker bridge gateway
password = ""
```

**Key Points:**
- `172.17.0.1` is the default Docker bridge gateway IP
- This allows containers to access services running on the host
- For `--network=host` mode, use `127.0.0.1` instead

### 3. Dockerfile Configuration

The Dockerfile must include the configuration files:

```dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0

ARG GITHUB_SHA=dev

COPY . .

RUN go mod tidy
RUN go build -v -ldflags "-w -s -extldflags '-static' -X 'github.com/wuhan005/NekoBox/internal/conf.BuildCommit=$GITHUB_SHA'" -o NekoBox ./cmd/

FROM alpine:latest

RUN apk update && apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone

WORKDIR /home/app

COPY --from=builder /app/NekoBox .
COPY --from=builder /app/conf ./conf  # IMPORTANT: Include config directory

RUN chmod 777 /home/app/NekoBox

ENTRYPOINT ["./NekoBox", "web"]
EXPOSE 8080
```

## Building and Running

### Build the Docker Image

```bash
cd /www/dk_project/dk_app/NekoBox
docker build -t nekobox:latest .
```

### Run the Container

#### Option 1: Bridge Networking (Recommended for BT Panel)
```bash
docker run -d \
  -p 8001:8080 \
  --restart=always \
  --name nekobox \
  nekobox:latest
```

**Port Mapping:**
- Host port: `8001`
- Container port: `8080`
- Access via: `http://localhost:8001`

#### Option 2: Host Networking (Alternative)
```bash
docker run -d \
  --network=host \
  --restart=always \
  --name nekobox \
  nekobox:latest
```

**Note:** For host networking, update `conf/app.ini` to use `127.0.0.1` for database and Redis hosts.

## Verification

### Check Container Status
```bash
docker ps | grep nekobox
```

Expected output:
```
CONTAINER ID   IMAGE             COMMAND           CREATED        STATUS        PORTS                    NAMES
xxxxxx         nekobox:latest    "./NekoBox web"   X minutes ago  Up X minutes  0.0.0.0:8001->8080/tcp   nekobox
```

### Check Logs
```bash
docker logs nekobox
```

Expected output:
```
time="2025-11-05T00:58:19+08:00" level=info msg="Starting web server" external_url=box.ethanpan.me
[Flamego] Listening on 0.0.0.0:8080 (production)
```

### Test HTTP Access
```bash
curl -I http://localhost:8001
```

Expected response:
```
HTTP/1.1 200 OK
```

## Troubleshooting

### Issue 1: Missing Configuration Files
**Error:** `load configuration: parse "conf/app.ini": no such file or directory`

**Solution:** Ensure the Dockerfile includes the config directory copy:
```dockerfile
COPY --from=builder /app/conf ./conf
```

### Issue 2: MySQL Connection Refused
**Error:** `dial tcp: lookup nekobox on X.X.X.X:53: no such host`

**Cause:** Hostname resolution fails in Docker container.

**Solution:** Use IP addresses instead of hostnames:
- For bridge networking: `172.17.0.1` (Docker gateway)
- For host networking: `127.0.0.1`

### Issue 3: MySQL Access Denied
**Error:** `Error 1130: Host '172.17.0.2' is not allowed to connect to this MySQL server`

**Cause:** MySQL user only has `localhost` access.

**Solution:** Grant access from Docker network:
```sql
GRANT ALL PRIVILEGES ON nekobox.* TO 'nekobox'@'172.17.%' IDENTIFIED BY 'LLWAJRM7WL2xiSEG';
FLUSH PRIVILEGES;
```

### Issue 4: Redis Connection Refused
**Error:** `dial tcp 172.17.0.1:6379: connect: connection refused`

**Cause:** Redis only listening on `127.0.0.1`.

**Solution:** Update Redis bind address in `/www/server/redis/redis.conf`:
```conf
bind 127.0.0.1 172.17.0.1
```

### Issue 5: Redis Protected Mode Error
**Error:** `DENIED Redis is running in protected mode...`

**Cause:** Redis protected mode blocks external connections without authentication.

**Solution:** Disable protected mode in Redis config:
```conf
protected-mode no
```

Then restart Redis:
```bash
systemctl restart redis
```

### Issue 6: HTTP 500 Internal Server Error
**Symptoms:** Container runs but returns HTTP 500.

**Debugging steps:**
1. Check container logs: `docker logs nekobox`
2. Look for panic messages or errors
3. Common causes:
   - Redis connection issues
   - Database connection issues
   - Missing session storage

## Network Modes Comparison

### Bridge Mode (Default)
**Pros:**
- Standard Docker networking
- Port mapping visible in BT Panel
- Isolated from host network

**Cons:**
- Requires IP-based host access (`172.17.0.1`)
- Requires Redis/MySQL to accept external connections

**Configuration:**
```ini
[database]
host = "172.17.0.1"

[redis]
addr = "172.17.0.1:6379"
```

### Host Mode
**Pros:**
- Direct localhost access
- Simpler configuration
- Better performance

**Cons:**
- Port mapping not visible in Docker
- Shares host network namespace
- Less isolation

**Configuration:**
```ini
[database]
host = "127.0.0.1"

[redis]
addr = "localhost:6379"
```

## BT Panel Integration

### Reverse Proxy Setup
1. Go to BT Panel → Website
2. Add site or configure existing site
3. Set up reverse proxy:
   - Proxy name: `NekoBox`
   - Target URL: `http://127.0.0.1:8001`
   - Enable WebSocket support if needed

### Firewall Configuration
If accessing directly via IP:
1. Go to BT Panel → Security
2. Add port: `8001`
3. Protocol: `TCP`
4. Allow from: `All` or specific IPs

## Maintenance

### View Logs
```bash
# Real-time logs
docker logs -f nekobox

# Last 50 lines
docker logs --tail 50 nekobox
```

### Restart Container
```bash
docker restart nekobox
```

### Stop Container
```bash
docker stop nekobox
```

### Remove Container
```bash
docker stop nekobox
docker rm nekobox
```

### Update Application
```bash
# Pull latest code
cd /www/dk_project/dk_app/NekoBox
git pull

# Rebuild image
docker build -t nekobox:latest .

# Recreate container
docker stop nekobox
docker rm nekobox
docker run -d -p 8001:8080 --restart=always --name nekobox nekobox:latest
```

## Database Schema

NekoBox uses GORM with AutoMigrate, which automatically creates these tables:

### Tables
1. **users** - User accounts and profiles
2. **questions** - Anonymous questions
3. **censor_logs** - Content moderation logs
4. **upload_images** - Uploaded image metadata
5. **upload_image_questions** - Image-question relationships

**Note:** No manual migration is required. Tables are created/updated automatically on application startup.

## Security Considerations

1. **Redis Security:**
   - Binding to Docker network (`172.17.0.1`) only exposes Redis to local Docker containers
   - For production, consider enabling Redis authentication
   - Ensure Redis is not exposed to the public internet

2. **MySQL Security:**
   - Use strong passwords
   - Limit access to specific networks (`172.17.%`)
   - Regularly update and patch MySQL

3. **Application Security:**
   - Keep secrets in environment variables or secure config management
   - Regularly update dependencies
   - Monitor logs for suspicious activity

## Environment Details

- **Server OS:** Ubuntu Linux
- **Docker Version:** Latest
- **MySQL Version:** 5.7+
- **Redis Version:** Latest
- **Panel:** BT Panel (宝塔面板)
- **Go Version:** 1.20 (build)
- **Runtime:** Alpine Linux (container)

## Debugging Commands Reference

```bash
# Check all containers
docker ps -a

# Check container logs
docker logs nekobox

# Check port bindings
netstat -tlnp | grep 8001
# or
ss -tlnp | grep 8001

# Check MySQL connectivity from container
docker exec nekobox ping -c 3 172.17.0.1

# Check Redis connectivity
redis-cli -h 172.17.0.1 ping

# Enter container shell
docker exec -it nekobox sh

# Check MySQL grants
mysql -u root -p -e "SHOW GRANTS FOR 'nekobox'@'172.17.%';"

# Check Redis config
grep -E "^bind|^protected-mode" /www/server/redis/redis.conf

# Test application directly
curl -I http://localhost:8001
```

## Quick Reference

### Complete Deployment Checklist
- [ ] MySQL database created
- [ ] MySQL user created with Docker network access
- [ ] Redis configured to accept Docker network connections
- [ ] Redis protected mode disabled
- [ ] `conf/app.ini` configured with correct hosts
- [ ] Docker image built successfully
- [ ] Container running without errors
- [ ] HTTP 200 response received
- [ ] Port accessible from external network (if needed)
- [ ] Reverse proxy configured (if using domain)

### File Locations
- Application code: `/www/dk_project/dk_app/NekoBox`
- Config file: `/www/dk_project/dk_app/NekoBox/conf/app.ini`
- Redis config: `/www/server/redis/redis.conf`
- MySQL config: `/etc/my.cnf` or `/www/server/mysql/my.cnf`

### Container Settings (for BT Panel GUI)
- **Container Name:** `nekobox`
- **Image:** `nekobox:latest`
- **Port Mapping:** `8001:8080`
- **Restart Policy:** `always`
- **Network:** `bridge` (default)
- **Command:** Leave empty (uses image default)

## Support

For issues and questions:
- GitHub: https://github.com/wuhan005/NekoBox
- Documentation: This file

---

**Last Updated:** 2025-11-05
**Deployment Status:** ✅ Working
**Access URL:** http://box.ethanpan.me
