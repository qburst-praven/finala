# Finala

[![Lint](https://github.com/similarweb/finala/workflows/Lint/badge.svg)](https://github.com/similarweb/finala/actions)
[![Fmt](https://github.com/similarweb/finala/workflows/Fmt/badge.svg)](https://github.com/similarweb/finala/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/similarweb/finala)](https://goreportcard.com/report/github.com/similarweb/finala)

![Finala Logo](docs/images/main-logo.png)

## Overview

Finala is an open-source cloud resource scanner that analyzes, discloses, presents, and notifies about wasteful and unused resources across your cloud infrastructure. It helps organizations achieve two primary objectives: **cost optimization** and **unused resource detection**.

Finala provides comprehensive visibility into your cloud spending by identifying resources that are either underutilized or completely unused, enabling you to make informed decisions about resource optimization and cost reduction.

## Features

### 🔍 **Resource Discovery & Analysis**
- **Multi-Cloud Support**: Currently supports AWS with extensible architecture for other cloud providers
- **Comprehensive Resource Coverage**: Analyzes 18+ AWS services including EC2, RDS, Lambda, DynamoDB, and more
- **Intelligent Detection**: Uses CloudWatch metrics and custom rules to identify underutilized resources
- **Cost Impact Analysis**: Calculates potential cost savings for each identified resource

### 🎯 **Smart Detection Engine**
- **YAML-Based Configuration**: Easy-to-understand resource definitions using high-level YAML syntax
- **Customizable Rules**: Tailor detection criteria to match your infrastructure patterns and usage habits
- **Metric-Based Analysis**: Leverages CloudWatch metrics with configurable thresholds and time periods
- **Formula Support**: Advanced mathematical expressions for complex resource evaluation

### 🖥️ **Modern Web Interface**
- **React-Based UI**: Modern, responsive web interface built with React 18 and Material-UI v5
- **Real-Time Dashboard**: Interactive charts and visualizations of resource utilization
- **Advanced Filtering**: Filter resources by tags, regions, services, and cost thresholds
- **Search Capabilities**: Powered by Meilisearch for fast, relevant resource discovery

### 🔐 **Security & Authentication**
- **JWT-Based Authentication**: Secure login system with token-based authentication
- **Protected Routes**: Role-based access control for sensitive resource information
- **Auto-Generated Credentials**: Secure default setup with customizable authentication

### 📊 **Reporting & Notifications**
- **Scheduled Notifications**: Configure automated alerts via Slack or email
- **Tag-Based Filtering**: Group and notify based on resource tags and cost thresholds
- **Customizable Reports**: Generate reports based on specific criteria and time periods

### 🚀 **Easy Deployment**
- **Docker Compose**: One-command deployment with pre-configured services
- **Production Ready**: Optimized Docker images for both development and production
- **Scalable Architecture**: Microservices-based design for horizontal scaling

## Supported AWS Services

| Service | Cost Optimization | Unused Detection |
|---------|------------------|------------------|
| API Gateway | ❌ | ✅ |
| DocumentDB | ✅ | ❌ |
| DynamoDB | ✅ | ❌ |
| EC2 ALB/NLB | ✅ | ❌ |
| EC2 Elastic IPs | ✅ | ❌ |
| EC2 ELB | ✅ | ❌ |
| EC2 NAT Gateways | ✅ | ❌ |
| EC2 Instances | ✅ | ❌ |
| EC2 Volumes | ✅ | ❌ |
| ElastiCache | ✅ | ❌ |
| Elasticsearch | ✅ | ❌ |
| IAM Users | ❌ | ✅ |
| Kinesis | ✅ | ❌ |
| Lambda | ❌ | ✅ |
| Neptune | ✅ | ❌ |
| RDS | ✅ | ❌ |
| Redshift | ✅ | ❌ |

## Quick Start

### Prerequisites
- Docker and Docker Compose
- AWS credentials with appropriate permissions
- At least 4GB RAM available

### Option 1: Using Pre-built Images (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/similarweb/finala.git
   cd finala
   ```

2. **Configure AWS credentials**
   ```bash
   export AWS_ACCESS_KEY_ID=your_access_key
   export AWS_SECRET_ACCESS_KEY=your_secret_key
   # Optional: AWS_SESSION_TOKEN for temporary credentials
   ```

3. **Start Finala**
   ```bash
   docker-compose -f docker-compose-hub.yaml up -d
   ```

4. **Access the web interface**
   - Open http://localhost:8080 in your browser
   - Default credentials: `admin` / `test`
   - Check the logs for auto-generated password if using default config

### Option 2: Building from Source

1. **Clone and build**
   ```bash
   git clone https://github.com/similarweb/finala.git
   cd finala
   docker-compose up -d
   ```

2. **Configure and access as above**

## Documentation

- **[Quick Start Guide](docs/quick-start.md)** - Detailed setup instructions
- **[Configuration Guide](docs/configuration.md)** - Complete configuration reference
- **[Architecture Overview](docs/architecture.md)** - System design and components
- **[AWS Setup](docs/aws-setup.md)** - AWS permissions and configuration
- **[API Reference](docs/api-reference.md)** - REST API documentation
- **[Troubleshooting](docs/troubleshooting.md)** - Common issues and solutions

## Architecture

Finala follows a microservices architecture with four main components:

- **Collector**: Scans AWS resources and analyzes utilization
- **API**: RESTful API for data access and management
- **UI**: React-based web interface
- **Notifier**: Handles scheduled notifications and alerts

For detailed architecture information, see the [Architecture Documentation](docs/architecture.md).

## Configuration

Finala uses YAML configuration files for each component:

- `configuration/api.yaml` - API server configuration
- `configuration/collector.yaml` - Resource collection rules
- `configuration/ui.yaml` - Web interface settings
- `configuration/notifier.yaml` - Notification settings

See the [Configuration Guide](docs/configuration.md) for detailed configuration options.

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Code style and standards
- Testing requirements
- Pull request process
- Development setup

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Recent Updates

### v2.0+ Major Improvements
- **Search Backend**: Migrated from Elasticsearch to Meilisearch for improved performance
- **Go Version**: Upgraded to latest Go version with enhanced security
- **Frontend Modernization**: React 18, Material-UI v5, React Router v6
- **Authentication**: Added JWT-based authentication system
- **Containerization**: Optimized Docker images and build processes
- **Dependencies**: Updated all packages to latest stable versions

## Support

- **Issues**: Report bugs and request features via [GitHub Issues](https://github.com/similarweb/finala/issues)
- **Discussions**: Join community discussions in [GitHub Discussions](https://github.com/similarweb/finala/discussions)
- **Security**: Report security vulnerabilities via [GitHub Security](https://github.com/similarweb/finala/security/advisories) 