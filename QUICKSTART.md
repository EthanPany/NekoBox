# NekoBox Quick Start Guide

Get NekoBox running in 10 minutes! 🚀

## Prerequisites

- Docker installed
- Cloudflare account (for R2 storage)
- SMTP email service (SMTP2GO recommended)

---

## Step 1: Configuration (5 minutes)

### 1.1 Copy config file

```bash
cp conf/app.sample.ini conf/app.ini
```

### 1.2 Generate security keys

```bash
# Generate random salt
openssl rand -base64 32

# Generate random xsrf_key  
openssl rand -base64 32
```

### 1.3 Edit `conf/app.ini` - Update these critical sections:

```ini
[app]
production = true
external_url = "http://your-server-ip:8080"  # or your domain

[server]
salt = "YOUR_RANDOM_SALT_FROM_STEP_1.2"
xsrf_key = "YOUR_RANDOM_XSRF_KEY_FROM_STEP_1.2"

[database]
password = "YOUR_STRONG_PASSWORD"  # Choose a strong password

[mail]
account = "your-email@example.com"
password = "YOUR_SMTP_PASSWORD"
smtp = "mail.smtp2go.com"
port = 2525

[upload]
image_endpoint = "https://YOUR_ACCOUNT.r2.cloudflarestorage.com"
image_access_id = "YOUR_R2_ACCESS_KEY"
image_access_secret = "YOUR_R2_SECRET_KEY"
image_bucket = "your-bucket-name"
image_bucket_cdn_host = "your-bucket.r2.dev"
```

---

## Step 2: Start Services (3 minutes)

```bash
# Create network
docker network create nekobox-network

# Start MySQL
docker run --name nekobox-mysql \
  --network nekobox-network \
  -e MYSQL_ROOT_PASSWORD=YOUR_STRONG_PASSWORD \
  -e MYSQL_DATABASE=nekobox \
  -v nekobox-mysql-data:/var/lib/mysql \
  -d mysql:8.0

# Start Redis
docker run --name nekobox-redis \
  --network nekobox-network \
  -v nekobox-redis-data:/data \
  -d redis:7-alpine

# Wait for database to be ready (30 seconds)
sleep 30
```

---

## Step 3: Build & Deploy (2 minutes)

```bash
# Build image
docker build -t nekobox:latest .

# Start application
docker run --name nekobox-app \
  --network nekobox-network \
  -p 8080:8080 \
  -v $(pwd)/conf:/home/app/conf \
  -d nekobox:latest

# Check logs
docker logs -f nekobox-app
```

---

## Step 4: Test

Open browser: `http://your-server-ip:8080`

✅ You should see the NekoBox homepage!

---

## Common Issues

### Port already in use?
```bash
# Use different port
docker run --name nekobox-app \
  --network nekobox-network \
  -p 3000:8080 \  # Changed from 8080:8080
  -d nekobox:latest
```

### Application won't start?
```bash
# Check logs
docker logs nekobox-app

# Verify config
cat conf/app.ini
```

### Database connection failed?
```bash
# Restart MySQL and wait longer
docker restart nekobox-mysql
sleep 60
docker restart nekobox-app
```

---

## Next Steps

- Read full [DEPLOYMENT.md](./DEPLOYMENT.md) for production setup
- Set up HTTPS with Nginx
- Configure backups
- Add monitoring

---

## Quick Command Reference

```bash
# View logs
docker logs -f nekobox-app

# Restart app
docker restart nekobox-app

# Stop everything
docker stop nekobox-app nekobox-redis nekobox-mysql

# Start everything
docker start nekobox-mysql nekobox-redis nekobox-app

# Rebuild after code changes
docker stop nekobox-app
docker rm nekobox-app
docker build -t nekobox:latest .
docker run --name nekobox-app --network nekobox-network -p 8080:8080 -d nekobox:latest
```

---

**Need help?** Check [DEPLOYMENT.md](./DEPLOYMENT.md) for detailed guide!

