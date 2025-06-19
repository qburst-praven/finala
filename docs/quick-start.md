# Quick Start Guide

Get Finala up and running in minutes.

## Prerequisites

- Docker & Docker Compose
- AWS credentials with appropriate permissions
- 2GB RAM available (4GB recommended for large AWS environments)

## Quick Setup

### 1. Clone and Start

```bash
git clone https://github.com/qburst/finala.git
cd finala

# Set AWS credentials
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key

# Start Finala
docker-compose -f docker-compose-hub.yaml up -d
```

### 2. Access the Interface

- Open http://localhost:8080
- Login: `admin` / `test`

### 3. Configure AWS (Optional)

Edit `configuration/collector.yaml` to add your AWS accounts:

```yaml
providers:
  aws:
    accounts:
      - name: production
        regions: [us-east-1, us-west-2]
        # Uncomment and add your credentials:
        # access_key: your_access_key
        # secret_key: your_secret_key
```

## Alternative: Build from Source

```bash
git clone https://github.com/qburst/finala.git
cd finala
docker-compose up -d --build
```

## Troubleshooting

**Services not starting?**
```bash
docker-compose -f docker-compose-hub.yaml ps
docker-compose -f docker-compose-hub.yaml logs [service-name]
```

**Need help?** See the [Troubleshooting Guide](troubleshooting.md)

## Next Steps

- **[Configuration Guide](configuration.md)** - Customize detection rules
- **[AWS Setup Guide](aws-setup.md)** - Configure AWS permissions
- **[Architecture Overview](architecture.md)** - Understand the system 