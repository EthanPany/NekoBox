# NekoBox Deployment Guide

Complete step-by-step guide to deploy NekoBox in production.

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Initial Setup](#initial-setup)
3. [Configuration](#configuration)
4. [Database Setup](#database-setup)
5. [Storage Setup (Cloudflare R2)](#storage-setup-cloudflare-r2)
6. [Email Setup](#email-setup)
7. [Docker Deployment](#docker-deployment)
8. [Testing](#testing)
9. [Production Checklist](#production-checklist)
10. [Troubleshooting](#troubleshooting)

---

## Prerequisites

Before you begin, ensure you have:

- ✅ Docker and Docker Compose installed
- ✅ A server with at least 1GB RAM
- ✅ Domain name (optional, but recommended)
- ✅ Cloudflare account (for R2 storage)
- ✅ SMTP email service (we recommend SMTP2GO)

---

## Initial Setup

### 1. Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/NekoBox.git
cd NekoBox
```

### 2. Copy Configuration Template

```bash
cp conf/app.sample.ini conf/app.ini
```

---

## Configuration

### 1. Edit `conf/app.ini`

Open the file and update the following sections:

#### App Settings

```ini
[app]
production = true  # IMPORTANT: Set to true for production
title = "NekoBox"
icp = ""  # Your ICP filing number (if in China)
external_url = "https://your-domain.com"  # Your actual domain
```

#### Security Settings

```ini
[security]
enable_text_censor = false  # Enable if you have Aliyun/Qiniu API
enable_recaptcha = false    # Enable if you want reCAPTCHA
```

#### Server Settings

```ini
[server]
port = 8080
salt = "CHANGE_THIS_TO_RANDOM_STRING_32_CHARS"
xsrf_key = "CHANGE_THIS_TO_ANOTHER_RANDOM_STRING"
xsrf_expire = 3600
```

**⚠️ IMPORTANT:** Generate strong random strings for `salt` and `xsrf_key`:

```bash
# Generate random salt
openssl rand -base64 32

# Generate random xsrf_key
openssl rand -base64 32
```

---

## Database Setup

### 1. Create Docker Network

```bash
docker network create nekobox-network
```

### 2. Start MySQL Database

```bash
docker run --name nekobox-mysql \
  --network nekobox-network \
  -e MYSQL_ROOT_PASSWORD=YOUR_STRONG_PASSWORD \
  -e MYSQL_DATABASE=nekobox \
  -v nekobox-mysql-data:/var/lib/mysql \
  -d mysql:8.0
```

### 3. Update Database Configuration

```ini
[database]
user = "root"
password = "YOUR_STRONG_PASSWORD"  # Match the password above
host = "nekobox-mysql"
port = 3306
name = "nekobox"
```

### 4. Start Redis Cache

```bash
docker run --name nekobox-redis \
  --network nekobox-network \
  -v nekobox-redis-data:/data \
  -d redis:7-alpine
```

---

## Storage Setup (Cloudflare R2)

### 1. Create R2 Bucket

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com)
2. Navigate to **R2 Object Storage**
3. Click **Create bucket**
4. Name it (e.g., `my-nekobox`)
5. Choose a location

### 2. Generate R2 API Tokens

1. Go to **R2** > **Manage R2 API Tokens**
2. Click **Create API token**
3. Give it **Edit** permissions
4. Copy the **Access Key ID** and **Secret Access Key**

### 3. Configure R2 in `app.ini`

```ini
[upload]
# Default images for new users
default_avatar = "https://api.dicebear.com/7.x/pixel-art/svg"
default_background = "https://images.unsplash.com/photo-1557683316-973673baf926?w=1200&h=400&fit=crop"

# Cloudflare R2 Configuration
image_endpoint = "https://YOUR_ACCOUNT_ID.r2.cloudflarestorage.com"
image_access_id = "YOUR_ACCESS_KEY_ID"
image_access_secret = "YOUR_SECRET_ACCESS_KEY"
image_bucket = "my-nekobox"
image_bucket_cdn_host = "YOUR_PUBLIC_R2_URL.r2.dev"
```

**Finding Your R2 Values:**
- `image_endpoint`: Found in R2 bucket settings as "S3 API"
- `image_bucket_cdn_host`: Found as "Public bucket URL" (enable public access first)

---

## Email Setup

### Option 1: SMTP2GO (Recommended)

1. Sign up at [SMTP2GO](https://www.smtp2go.com/)
2. Verify your email
3. Get your SMTP credentials

```ini
[mail]
account = "your-email@example.com"
password = "YOUR_SMTP2GO_API_KEY"
port = 2525
smtp = "mail.smtp2go.com"
```

### Option 2: Gmail

```ini
[mail]
account = "your-gmail@gmail.com"
password = "YOUR_APP_PASSWORD"  # Generate at myaccount.google.com
port = 587
smtp = "smtp.gmail.com"
```

### Option 3: Custom SMTP

```ini
[mail]
account = "noreply@your-domain.com"
password = "YOUR_SMTP_PASSWORD"
port = 587
smtp = "smtp.your-domain.com"
```

---

## Docker Deployment

### 1. Build the Docker Image

```bash
docker build -t nekobox:latest .
```

This takes 5-10 minutes on first build.

### 2. Start the Application

```bash
docker run --name nekobox-app \
  --network nekobox-network \
  -p 8080:8080 \
  -v $(pwd)/conf:/home/app/conf \
  -d nekobox:latest
```

**With custom port (e.g., 80 for HTTP):**

```bash
docker run --name nekobox-app \
  --network nekobox-network \
  -p 80:8080 \
  -v $(pwd)/conf:/home/app/conf \
  -d nekobox:latest
```

### 3. Check Application Logs

```bash
docker logs -f nekobox-app
```

You should see:
```
[Flamego] Listening on 0.0.0.0:8080 (production)
time="2025-10-18T15:00:00+08:00" level=info msg="Starting web server"
```

---

## Testing

### 1. Check Container Status

```bash
docker ps --filter "name=nekobox"
```

All three containers should be running:
- `nekobox-app`
- `nekobox-mysql`
- `nekobox-redis`

### 2. Test the Application

```bash
curl http://localhost:8080
```

Should return HTML content.

### 3. Register a Test User

1. Open your browser to `http://your-server-ip:8080`
2. Click "注册" (Register)
3. Fill in the form
4. Check if you receive a confirmation email
5. Try uploading an avatar

### 4. Test Image Upload

1. Go to user profile
2. Upload an avatar
3. Verify it appears in your R2 bucket
4. Check the image loads correctly

---

## Production Checklist

### Security

- [ ] Changed `salt` and `xsrf_key` to random values
- [ ] Set `production = true` in config
- [ ] Changed MySQL root password
- [ ] Enabled HTTPS (see Nginx setup below)
- [ ] Set up firewall (only ports 80/443 open)
- [ ] Regular backups configured

### Performance

- [ ] Database persistent volume configured
- [ ] Redis persistent volume configured
- [ ] CDN configured for R2 (optional)
- [ ] Monitoring set up (optional)

### Optional: Nginx Reverse Proxy with HTTPS

Create `nginx.conf`:

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Get SSL certificate:

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d your-domain.com
```

---

## Maintenance

### Update Application

```bash
# Stop and remove old container
docker stop nekobox-app
docker rm nekobox-app

# Pull latest code
git pull

# Rebuild image
docker build -t nekobox:latest .

# Start new container
docker run --name nekobox-app \
  --network nekobox-network \
  -p 8080:8080 \
  -v $(pwd)/conf:/home/app/conf \
  -d nekobox:latest
```

### Backup Database

```bash
# Backup
docker exec nekobox-mysql mysqldump -u root -p nekobox > backup-$(date +%Y%m%d).sql

# Restore
docker exec -i nekobox-mysql mysql -u root -p nekobox < backup-20251018.sql
```

### View Logs

```bash
# Application logs
docker logs -f nekobox-app

# Last 100 lines
docker logs --tail 100 nekobox-app

# MySQL logs
docker logs nekobox-mysql

# Redis logs
docker logs nekobox-redis
```

### Restart Services

```bash
# Restart application
docker restart nekobox-app

# Restart database
docker restart nekobox-mysql

# Restart cache
docker restart nekobox-redis
```

---

## Troubleshooting

### Application won't start

```bash
# Check logs
docker logs nekobox-app

# Common issues:
# 1. Database not ready - wait 30 seconds and restart
# 2. Config file errors - check conf/app.ini syntax
# 3. Port already in use - change port mapping
```

### Database connection errors

```bash
# Check if MySQL is running
docker ps | grep mysql

# Check MySQL logs
docker logs nekobox-mysql

# Test connection
docker exec -it nekobox-mysql mysql -u root -p
```

### Images not uploading

```bash
# Check R2 credentials
# Verify image_access_id and image_access_secret

# Check application logs for upload errors
docker logs nekobox-app | grep -i "upload\|r2\|s3"
```

### Email not sending

```bash
# Check SMTP configuration
# Test with: curl -v telnet://mail.smtp2go.com:2525

# Check app logs
docker logs nekobox-app | grep -i "mail\|smtp"
```

### Container keeps restarting

```bash
# Check what's wrong
docker logs nekobox-app

# Common causes:
# - Missing conf/app.ini file
# - Invalid configuration
# - Port conflict
```

---

## Quick Commands Reference

```bash
# Start all services
docker network create nekobox-network
docker start nekobox-mysql nekobox-redis nekobox-app

# Stop all services
docker stop nekobox-app nekobox-redis nekobox-mysql

# View status
docker ps --filter "name=nekobox"

# Clean up (WARNING: Removes data)
docker stop nekobox-app nekobox-redis nekobox-mysql
docker rm nekobox-app nekobox-redis nekobox-mysql
docker network rm nekobox-network
```

---

## Support

If you encounter issues:

1. Check the logs: `docker logs nekobox-app`
2. Verify configuration in `conf/app.ini`
3. Ensure all containers are running: `docker ps`
4. Check this troubleshooting guide
5. Open an issue on GitHub

---

## Production Recommendations

### For Small Scale (< 1000 users)
- **Server**: 2GB RAM, 2 CPU cores
- **Database**: Default MySQL settings
- **Redis**: Default settings
- **Estimated cost**: $5-10/month

### For Medium Scale (1000-10000 users)
- **Server**: 4GB RAM, 4 CPU cores
- **Database**: Increase MySQL buffer pool
- **Redis**: Add persistence
- **CDN**: Enable for static assets
- **Estimated cost**: $20-40/month

### For Large Scale (> 10000 users)
- **Server**: 8GB+ RAM, 8+ CPU cores
- **Database**: Separate MySQL server, read replicas
- **Redis**: Separate Redis server, clustering
- **CDN**: Required for R2
- **Load Balancer**: Multiple app instances
- **Estimated cost**: $100+/month

---

**🎉 Congratulations! Your NekoBox is now deployed!**

Access it at: `http://your-server-ip:8080` or `https://your-domain.com`

