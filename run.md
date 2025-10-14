# NekoBox Setup & Deployment Steps

This document outlines the steps taken to successfully run NekoBox in Docker containers.

## Issues Encountered & Solutions

### 1. Docker Build Issues

**Problem**: Initial Docker build failed because configuration files were not being copied to the container.

**Solution**: Added missing configuration copy to Dockerfile:
```dockerfile
COPY --from=builder /app/NekoBox .
COPY conf ./conf  # <- Added this line
```

### 2. Configuration Parsing Problems

**Problem**: The INI configuration parser was including inline comments as part of configuration values, causing errors like:
```
lookup "nekobox-mysql          # MySQL pn�;:": no such host
```

**Solution**:
- Removed all inline comments from `conf/app.ini`
- Created clean configuration without comment text interfering with values
- The `gopkg.in/ini.v1` parser despite `IgnoreInlineComment: true` was still including comment text

### 3. reCAPTCHA Configuration

**Problem**: Application crashed on startup due to empty reCAPTCHA server key:
```
panic: recaptcha: empty secret
```

**Solution**:
1. Added configuration switches to disable both text censorship and reCAPTCHA
2. Modified code to conditionally initialize reCAPTCHA middleware only when enabled
3. Updated templates to conditionally show reCAPTCHA elements
4. Updated form handlers to conditionally verify reCAPTCHA

**Configuration options added**:
```ini
[security]
enable_text_censor = false  # Disables Aliyun content security checks
enable_recaptcha = false    # Disables reCAPTCHA verification
```

### 4. Container Networking

**Problem**: NekoBox container couldn't connect to MySQL and Redis services.

**Solution**:
- Created Docker network for inter-container communication
- All containers run within the same network for proper connectivity

## Final Deployment Steps

1. **Create Docker network**:
   ```bash
   docker network create nekobox-network
   ```

2. **Start MySQL container**:
   ```bash
   docker run --name nekobox-mysql \
     --network nekobox-network \
     -e MYSQL_ROOT_PASSWORD=password \
     -e MYSQL_DATABASE=nekobox \
     -d mysql:8.0
   ```

3. **Start Redis container**:
   ```bash
   docker run --name nekobox-redis \
     --network nekobox-network \
     -d redis:7-alpine
   ```

4. **Build and run NekoBox**:
   ```bash
   docker build -t nekobox-final .
   docker run --network nekobox-network -p 8080:8080 nekobox-final
   ```

## Configuration Used

The final working configuration (`conf/app.ini`):

```ini
[app]
production = false
title = "NekoBox"
icp = ""
external_url = "http://localhost:8080"
uptrace_dsn = ""
aliyun_access_key = ""
aliyun_access_key_secret = ""

[security]
enable_text_censor = false
enable_recaptcha = false

[server]
port = 8080
salt = "nekobox-salt-dev-123456"
xsrf_key = "nekobox-xsrf-dev-789"
xsrf_expire = 3600

[database]
user = "root"
password = "password"
host = "nekobox-mysql"
port = 3306
name = "nekobox"

[pixel]
host = "localhost:8088"

[redis]
addr = "nekobox-redis:6379"
password = ""

[recaptcha]
domain = "https://www.recaptcha.net"
site_key = "dummy"
server_key = "dummy"
turnstile_style = false

[upload]
default_avatar = ""
default_background = ""
aliyun_endpoint = ""
aliyun_access_id = ""
aliyun_access_secret = ""
aliyun_bucket = ""
aliyun_bucket_cdn_host = ""
image_endpoint = ""
image_access_id = ""
image_access_secret = ""
image_bucket = ""
image_bucket_cdn_host = ""

[mail]
account = ""
password = ""
port = 465
smtp = ""
```

## Verification

Application successfully runs and is accessible at `http://localhost:8080`. The server responds with HTTP 302 redirects indicating proper functionality.

## Key Learnings

1. **INI Parser Behavior**: The `gopkg.in/ini.v1` library can be sensitive to inline comments despite configuration options
2. **Docker Networking**: Proper container networking is crucial for multi-service applications
3. **Configuration Validation**: Applications may have strict validation requirements (like reCAPTCHA keys) even in development mode
4. **Build Context**: Ensure all necessary files are included in Docker build context and properly copied to containers