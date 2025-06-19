# Quick Start Guide

This guide will help you get Finala up and running quickly. Choose the deployment method that best fits your needs.

## Prerequisites

Before you begin, ensure you have:

- **Docker & Docker Compose**: Version 1.29 or higher
- **AWS Access**: Valid AWS credentials with appropriate permissions
- **System Resources**: At least 4GB RAM and 2GB free disk space
- **Network Access**: Ability to access AWS APIs and download Docker images

## Deployment Options

### Option 1: Pre-built Images (Recommended for Production)

This option uses pre-built Docker images from Docker Hub, making it faster to deploy and more suitable for production environments.

#### Step 1: Clone the Repository

```bash
git clone https://github.com/qburst/finala.git
cd finala
```

#### Step 2: Configure AWS Credentials

Set your AWS credentials as environment variables:

```bash
export AWS_ACCESS_KEY_ID=your_access_key_here
export AWS_SECRET_ACCESS_KEY=your_secret_key_here
# Optional: For temporary credentials
export AWS_SESSION_TOKEN=your_session_token_here
```

**Note**: Ensure your AWS credentials have the necessary permissions. See [AWS Setup Guide](aws-setup.md) for detailed permission requirements.

#### Step 3: Start Finala

```bash
docker-compose -f docker-compose-hub.yaml up -d
```

This command will:
- Pull pre-built images from Docker Hub
- Start all required services (Meilisearch, API, UI, Collector)
- Set up networking between services
- Mount configuration files

#### Step 4: Verify Deployment

Check that all services are running:

```bash
docker-compose -f docker-compose-hub.yaml ps
```

You should see all services in the "Up" state:
- `meilisearch`
- `api`
- `ui`
- `collector`

#### Step 5: Access the Web Interface

1. Open your browser and navigate to `http://localhost:8080`
2. Use the default credentials:
   - **Username**: `admin`
   - **Password**: `test`

**Security Note**: For production deployments, change the default password in the configuration file.

### Option 2: Build from Source (Development)

This option builds Finala from source code, which is useful for development or when you need to modify the codebase.

#### Step 1: Clone and Build

```bash
git clone https://github.com/qburst/finala.git
cd finala
docker-compose up -d --build
```

This will:
- Build all Docker images from source
- Install dependencies
- Compile the Go applications
- Start all services

#### Step 2: Configure and Access

Follow the same steps as Option 1 for AWS credentials and web interface access.

## First Run

### 1. Start Resource Collection

The collector will automatically start scanning your AWS resources. You can monitor the progress:

```bash
docker-compose logs -f collector
```

### 2. View Results

1. Access the web interface at `http://localhost:8080`
2. Navigate to the Dashboard to see discovered resources
3. Use filters to analyze specific resource types or regions

### 3. Enable Notifications (Optional)

After the first successful collection, you can enable notifications:

1. Edit `docker-compose-hub.yaml` and uncomment the notifier service
2. Configure notification settings in `configuration/notifier.yaml`
3. Restart the stack:
```bash
docker-compose -f docker-compose-hub.yaml down
docker-compose -f docker-compose-hub.yaml up -d
```

## Troubleshooting

If you encounter issues during setup or first run:

1. **Check service status**:
   ```bash
   docker-compose -f docker-compose-hub.yaml ps
   ```

2. **View service logs**:
   ```bash
   docker-compose -f docker-compose-hub.yaml logs [service-name]
   ```

3. **Common quick fixes**:
   - **Services not starting**: Check port conflicts and system resources
   - **Authentication issues**: Verify credentials in `configuration/api.yaml`
   - **AWS access problems**: Check credentials and permissions
   - **No resources found**: Verify AWS account configuration

For detailed troubleshooting and solutions, see the **[Troubleshooting Guide](troubleshooting.md)**.

## Next Steps

Now that Finala is running, explore these resources:

- **[Configuration Guide](configuration.md)** - Customize detection rules and settings
- **[Architecture Overview](architecture.md)** - Understand how Finala works
- **[API Reference](api-reference.md)** - Integrate with external systems
- **[AWS Setup Guide](aws-setup.md)** - Configure AWS permissions and accounts

## Production Considerations

For production deployments, consider:

- **Security**: Change default passwords and use secure authentication
- **Monitoring**: Set up monitoring for Finala services
- **Backup**: Configure backups for Meilisearch data
- **Scaling**: Adjust resource limits based on your infrastructure size
- **Networking**: Use proper network security groups and firewalls 