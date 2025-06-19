# Installation Guide

This guide covers comprehensive installation scenarios for Finala, from development to production deployments.

## Prerequisites

### System Requirements

- **Minimum**:
  - 2GB RAM
  - 2GB free disk space
  - Docker 20.10+
  - Docker Compose 1.29+

- **Recommended**:
  - 4GB RAM (for large AWS environments)
  - 10GB free disk space
  - Docker 24.0+
  - Docker Compose 2.0+

### Network Requirements

- Outbound access to AWS APIs
- Port 8080 available (UI)
- Port 8089 available (API)
- Port 7700 available (Meilisearch)

## Installation Methods

### Method 1: Docker Compose (Recommended)

#### Development Setup

```bash
# Clone repository
git clone https://github.com/qburst/finala.git
cd finala

# Set AWS credentials
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key

# Start Finala
docker-compose -f docker-compose-hub.yaml up -d

# Or build from source
docker-compose up -d --build
```

#### Production Setup

```bash
# Clone repository
git clone https://github.com/qburst/finala.git
cd finala

# Create production configuration
cp configuration/api.yaml configuration/api.prod.yaml
cp configuration/collector.yaml configuration/collector.prod.yaml

# Edit configurations for production
# (see Configuration Guide for details)

# Start with production configs
docker-compose -f docker-compose-hub.yaml -f docker-compose.prod.yaml up -d
```

### Method 2: Custom Ports

If default ports are in use, create a custom compose file:

```yaml
# docker-compose.custom.yaml
version: '3.8'
services:
  ui:
    ports:
      - "9090:8080"  # Custom UI port
  
  api:
    ports:
      - "9099:8089"  # Custom API port
  
  meilisearch:
    ports:
      - "7701:7700"  # Custom Meilisearch port
```

```bash
docker-compose -f docker-compose-hub.yaml -f docker-compose.custom.yaml up -d
```

### Method 3: External Database

For production, you may want to use external Meilisearch:

```yaml
# docker-compose.external-db.yaml
version: '3.8'
services:
  api:
    environment:
      - OVERRIDE_STORAGE_ENDPOINT=http://your-meilisearch:7700
      - OVERRIDE_STORAGE_PASSWORD=your_master_key
  
  # Remove internal meilisearch service
  meilisearch:
    profiles:
      - skip
```

## Environment-Specific Configurations

### Development Environment

```yaml
# configuration/api.dev.yaml
---
log_level: debug
storage:
  meilisearch:
    password: "dev_master_key"
    endpoints: ["http://meilisearch:7700"]
auth:
  username: "dev"
  password: "dev_password"
```

```yaml
# configuration/collector.dev.yaml
---
name: development
log_level: debug
api_server:
  address: http://api:8081
  bulk_interval: 5s
providers:
  aws:
    accounts:
      - name: dev
        regions: [us-east-1]
        profile: dev-profile
```

### Staging Environment

```yaml
# configuration/api.staging.yaml
---
log_level: info
storage:
  meilisearch:
    password: "staging_master_key"
    endpoints: ["http://meilisearch:7700"]
auth:
  username: "staging"
  password: "secure_staging_password"
```

### Production Environment

```yaml
# configuration/api.prod.yaml
---
log_level: warn
storage:
  meilisearch:
    password: "${MEILI_MASTER_KEY}"
    endpoints: ["http://meilisearch:7700"]
auth:
  username: "${FINALA_USERNAME}"
  password: "${FINALA_PASSWORD}"
```

## Security Considerations

### Production Security

1. **Change Default Passwords**:
   ```bash
   # Generate secure passwords
   openssl rand -base64 32
   ```

2. **Use Environment Variables**:
   ```bash
   export MEILI_MASTER_KEY="your_secure_master_key"
   export FINALA_USERNAME="admin"
   export FINALA_PASSWORD="your_secure_password"
   ```

3. **Network Security**:
   ```yaml
   # docker-compose.secure.yaml
   version: '3.8'
   services:
     ui:
       networks:
         - finala_frontend
     
     api:
       networks:
         - finala_backend
         - finala_frontend
     
     collector:
       networks:
         - finala_backend
     
     meilisearch:
       networks:
         - finala_backend
   
   networks:
     finala_frontend:
       driver: bridge
     finala_backend:
       driver: bridge
       internal: true
   ```

4. **Resource Limits**:
   ```yaml
   # docker-compose.limits.yaml
   version: '3.8'
   services:
     api:
       deploy:
         resources:
           limits:
             memory: 1G
             cpus: '0.5'
     
     collector:
       deploy:
         resources:
           limits:
             memory: 2G
             cpus: '1.0'
   ```

## Data Management

### Backup and Restore

#### Backup Meilisearch Data

```bash
# Create backup
docker-compose exec meilisearch tar -czf /tmp/meilisearch-backup.tar.gz /data
docker cp finala_meilisearch_1:/tmp/meilisearch-backup.tar.gz ./backup/

# Restore backup
docker cp ./backup/meilisearch-backup.tar.gz finala_meilisearch_1:/tmp/
docker-compose exec meilisearch tar -xzf /tmp/meilisearch-backup.tar.gz -C /
```

#### Backup Configuration

```bash
# Backup configuration files
tar -czf finala-config-backup.tar.gz configuration/

# Restore configuration
tar -xzf finala-config-backup.tar.gz
```

### Data Persistence

```yaml
# docker-compose.persistent.yaml
version: '3.8'
services:
  meilisearch:
    volumes:
      - meilisearch_data:/data
      - ./backups:/backups
  
  api:
    volumes:
      - ./logs:/app/logs
  
  collector:
    volumes:
      - ./logs:/app/logs

volumes:
  meilisearch_data:
    driver: local
```

## Monitoring and Health Checks

### Health Check Script

```bash
#!/bin/bash
# health-check.sh

API_URL="http://localhost:8089/api/v1/health"
UI_URL="http://localhost:8080"
MEILI_URL="http://localhost:7700/health"

# Check API
if curl -f $API_URL > /dev/null 2>&1; then
    echo "✅ API: OK"
else
    echo "❌ API: FAILED"
    exit 1
fi

# Check UI
if curl -f $UI_URL > /dev/null 2>&1; then
    echo "✅ UI: OK"
else
    echo "❌ UI: FAILED"
    exit 1
fi

# Check Meilisearch
if curl -f $MEILI_URL > /dev/null 2>&1; then
    echo "✅ Meilisearch: OK"
else
    echo "❌ Meilisearch: FAILED"
    exit 1
fi

echo "🎉 All services healthy!"
```

### Docker Health Checks

```yaml
# docker-compose.health.yaml
version: '3.8'
services:
  api:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8089/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
  
  ui:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080"]
      interval: 30s
      timeout: 10s
      retries: 3
  
  meilisearch:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:7700/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Scaling Considerations

### Horizontal Scaling

```yaml
# docker-compose.scale.yaml
version: '3.8'
services:
  api:
    deploy:
      replicas: 3
      resources:
        limits:
          memory: 1G
          cpus: '0.5'
  
  collector:
    deploy:
      replicas: 2
      resources:
        limits:
          memory: 2G
          cpus: '1.0'
```

### Load Balancing

```yaml
# docker-compose.lb.yaml
version: '3.8'
services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - api
      - ui

  api:
    expose:
      - "8089"
  
  ui:
    expose:
      - "8080"
```

## Troubleshooting Installation

### Common Issues

**Port Conflicts**:
```bash
# Check what's using the ports
sudo lsof -i :8080
sudo lsof -i :8089
sudo lsof -i :7700

# Change ports in docker-compose.yaml
```

**Permission Issues**:
```bash
# Fix Docker permissions
sudo usermod -aG docker $USER
newgrp docker
```

**Resource Issues**:
```bash
# Check system resources
free -h
df -h
docker system df
```

### Installation Verification

```bash
# Check all services are running
docker-compose ps

# Check service logs
docker-compose logs api
docker-compose logs collector
docker-compose logs ui
docker-compose logs meilisearch

# Test API endpoint
curl http://localhost:8089/api/v1/health

# Test UI endpoint
curl http://localhost:8080
```

## Next Steps

After installation:

1. **[Quick Start Guide](quick-start.md)** - Get running quickly
2. **[Configuration Guide](configuration.md)** - Customize settings
3. **[AWS Setup Guide](aws-setup.md)** - Configure AWS access
4. **[Troubleshooting Guide](troubleshooting.md)** - Solve issues 

#### Alternative: Configure AWS Credentials in File

Instead of environment variables, you can configure AWS credentials directly in the configuration file:

```yaml
# configuration/collector.yaml
providers:
  aws:
    accounts: 
      - name: production
        regions:
          - us-east-1
          - us-west-2
        access_key: your_access_key_here
        secret_key: your_secret_key_here
```

**Note**: Using credentials in configuration files is less secure than environment variables. For production, prefer environment variables or IAM roles.

## AWS Authentication Methods

Finala supports multiple AWS authentication methods. Choose the one that best fits your security requirements:

### Method 1: Environment Variables (Recommended)

```bash
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_SESSION_TOKEN=your_session_token  # For temporary credentials
```

### Method 2: Configuration File

Edit `configuration/collector.yaml`:

```yaml
providers:
  aws:
    accounts: 
      - name: production
        regions: [us-east-1, us-west-2]
        access_key: your_access_key_here
        secret_key: your_secret_key_here
```

### Method 3: AWS Profile

```yaml
providers:
  aws:
    accounts: 
      - name: production
        regions: [us-east-1, us-west-2]
        profile: your_aws_profile_name
```

### Method 4: IAM Role (Production Recommended)

```yaml
providers:
  aws:
    accounts: 
      - name: production
        regions: [us-east-1, us-west-2]
        role: arn:aws:iam::123456789012:role/FinalaRole
```

### Method 5: Multiple Accounts

```yaml
providers:
  aws:
    accounts: 
      - name: production
        regions: [us-east-1, us-west-2]
        role: arn:aws:iam::123456789012:role/FinalaRole
      
      - name: development
        regions: [us-east-1]
        profile: dev-profile
      
      - name: staging
        regions: [eu-west-1]
        access_key: staging_access_key
        secret_key: staging_secret_key
```

**Security Best Practices**:
- Use IAM roles for production environments
- Use environment variables over configuration files
- Rotate credentials regularly
- Follow least privilege principle 