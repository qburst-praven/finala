# Configuration Guide

Finala uses YAML configuration files to control its behavior. This guide covers all configuration options for each component.

## Configuration Files Overview

Finala uses four main configuration files:

- `configuration/api.yaml` - API server settings
- `configuration/collector.yaml` - Resource collection rules and AWS configuration
- `configuration/ui.yaml` - Web interface settings
- `configuration/notifier.yaml` - Notification settings

## API Configuration (`configuration/api.yaml`)

The API configuration controls the REST API server, authentication, and storage settings.

### Basic Configuration

```yaml
---
log_level: info  # debug, info, warn, error

storage:
  meilisearch:
    username: ""  # Usually empty for Meilisearch
    password: "your_meilisearch_master_key"
    endpoints: 
      - http://meilisearch:7700

smtp:
  username: "your_email@example.com"
  password: "your_email_password"
  smtpServer: "smtp.gmail.com"
  smtpPort: 587

auth:
  username: "admin"
  password: "your_secure_password"
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `log_level` | string | `info` | Logging level (debug, info, warn, error) |
| `storage.meilisearch.username` | string | `""` | Meilisearch username (usually empty) |
| `storage.meilisearch.password` | string | - | Meilisearch master key |
| `storage.meilisearch.endpoints` | array | - | List of Meilisearch endpoints |
| `smtp.username` | string | - | SMTP username for email notifications |
| `smtp.password` | string | - | SMTP password |
| `smtp.smtpServer` | string | - | SMTP server address |
| `smtp.smtpPort` | int | - | SMTP server port |
| `auth.username` | string | `admin` | Web interface username |
| `auth.password` | string | - | Web interface password |

## Collector Configuration (`configuration/collector.yaml`)

The collector configuration defines AWS accounts, regions, and resource detection rules.

### Basic Structure

```yaml
---
name: general
log_level: info
api_server: 
  address: http://127.0.0.1:8081
  bulk_interval: 5s

providers:
  aws:
    accounts: 
      - name: production
        regions:
          - us-east-1
          - us-west-2
        # Authentication methods (see AWS Setup Guide for details):
        # access_key: your_access_key
        # secret_key: your_secret_key
        # profile: your_aws_profile
        # role: arn:aws:iam::123456789012:role/FinalaRole
    metrics:
      # Resource-specific metrics configuration
```

**Note**: For detailed AWS authentication setup, see the [AWS Setup Guide](aws-setup.md).

### Resource Metrics Configuration

The metrics section defines detection rules for each AWS service. Here are examples for common services:

#### EC2 Instances

```yaml
ec2:
  - description: CPU utilization 
    enable: true
    metrics:
      - name: CPUUtilization
        statistic: Maximum
    period: 24h 
    start_time: 168h  # 7 days
    constraint:
      operator: "<"
      value: 40  # Less than 40% CPU utilization
```

**How EC2 CPU Check Works:**
- **Time Period**: Looks at the last 7 days (168 hours)
- **Aggregation**: Collects 24-hour periods within those 7 days
- **Statistic**: Uses `Maximum` to find the highest CPU usage in each 24-hour period
- **Detection**: If the maximum CPU utilization in any 24-hour period was **less than 40%**, the instance is considered underutilized
- **Result**: Instances that never exceeded 40% CPU in the last 7 days are flagged for potential downsizing or termination

**Example**: If an EC2 instance had 40% maximum CPU under over the last 7 days, it would be detected as underutilized.

### Metric Configuration Options

| Option | Type | Description |
|--------|------|-------------|
| `description` | string | Human-readable description of the rule |
| `enable` | boolean | Whether this rule is active |
| `metrics` | array | List of CloudWatch metrics to collect |
| `period` | string | Time period for metric aggregation (e.g., "24h") |
| `start_time` | string | How far back to look for metrics (e.g., "168h" = 7 days) |
| `constraint.operator` | string | Comparison operator ("<", ">", "==", ">=", etc.) |
| `constraint.value` | number | Threshold value for comparison |
| `constraint.formula` | string | Mathematical formula for complex calculations |

### Supported Operators

- `==` - Equal to
- `!=` - Not equal to
- `<` - Less than
- `<=` - Less than or equal to
- `>` - Greater than
- `>=` - Greater than or equal to

## UI Configuration (`configuration/ui.yaml`)

The UI configuration controls the web interface settings and API connection. When properly configured, you'll see the login screen and dashboard as shown below:

![Finala Login Screen](images/login-page.png)

After logging in, you'll see the main dashboard showing resource metrics and potential savings:

![Finala Dashboard](images/dashboard-metrics.png)

The UI configuration is minimal and primarily controls the connection to the API.

```yaml
---
log_level: info
api_server:
  address: http://127.0.0.1:8089  # API server address
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `log_level` | string | `info` | Logging level |
| `api_server.address` | string | - | API server URL |

## Notifier Configuration (`configuration/notifier.yaml`)

The notifier configuration controls automated notifications via Slack and email.

### Basic Configuration

```yaml
---
log_level: info
api_server_address: "http://127.0.0.1:8089"
ui_address: "http://127.0.0.1:8080"
notifiers:
  slack:
    token: your_slack_token
    default_channels:
      - '#general'
    notify_by_tags:
      production:
        minimum_cost_to_present: 100
        tags:
          - name: Environment
            value: production
        notify_to:
          - "#cost-alerts"
          - "@username"
```

### Slack Configuration

#### Basic Slack Setup

```yaml
notifiers:
  slack:
    token: xoxb-your-slack-bot-token
    default_channels:
      - '#general'
```

#### Tag-Based Notifications

```yaml
notify_by_tags:
  production:
    minimum_cost_to_present: 100  # Minimum cost to trigger notification
    tags:
      - name: Environment
        value: production
      - name: Team
        value: engineering
    notify_to:
      - "#cost-alerts"
      - "@username"
      - "#team-engineering"
  
  development:
    minimum_cost_to_present: 50
    tags:
      - name: Environment
        value: development
    notify_to:
      - "#dev-alerts"
```

### Email Configuration

Email notifications are configured in the API configuration file:

```yaml
smtp:
  username: "your_email@example.com"
  password: "your_app_password"
  smtpServer: "smtp.gmail.com"
  smtpPort: 587
```

## Environment Variables

You can override configuration values using environment variables:

### API Configuration

| Environment Variable | Configuration Path | Description |
|---------------------|-------------------|-------------|
| `OVERRIDE_STORAGE_ENDPOINT` | `storage.meilisearch.endpoints[0]` | Meilisearch endpoint |
| `OVERRIDE_STORAGE_PASSWORD` | `storage.meilisearch.password` | Meilisearch master key |

### Collector Configuration

| Environment Variable | Description |
|---------------------|-------------|
| `OVERRIDE_API_ENDPOINT` | API server endpoint |

### Meilisearch Configuration

| Environment Variable | Description |
|---------------------|-------------|
| `MEILI_MASTER_KEY` | Meilisearch master key |

## Configuration Best Practices

### Security

1. **Never commit credentials to version control**
2. **Use environment variables for sensitive data**
3. **Secure Meilisearch master key**
4. **Use strong authentication passwords**

### Performance

1. **Adjust `bulk_interval` based on your infrastructure size**
2. **Configure appropriate time periods for metrics**
3. **Use specific regions to reduce API calls**

### Monitoring

1. **Set up log aggregation for all services**
2. **Monitor API response times**
3. **Track resource collection success rates**

## Configuration Validation

Finala validates configuration files on startup. Common validation errors:

- **Invalid YAML syntax**
- **Missing required fields**
- **Invalid metric names**
- **Unsupported operators**

Check service logs for specific validation errors:

```bash
docker-compose logs api
docker-compose logs collector
```

## Example Configurations

### Minimal Production Setup

```yaml
# api.yaml
---
log_level: info
storage:
  meilisearch:
    password: "your_secure_master_key"
    endpoints: ["http://meilisearch:7700"]
auth:
  username: "admin"
  password: "your_secure_password"

# collector.yaml
---
name: production
log_level: info
api_server:
  address: http://api:8081
  bulk_interval: 10s
providers:
  aws:
    accounts:
      - name: production
        regions: [us-east-1, us-west-2]
        # Configure authentication (see AWS Setup Guide)
```

### Development Setup

```yaml
# api.yaml
---
log_level: debug
storage:
  meilisearch:
    password: "dev_master_key"
    endpoints: ["http://meilisearch:7700"]
auth:
  username: "dev"
  password: "dev_password"

# collector.yaml
---
name: development
log_level: debug
api_server:
  address: http://api:8081
  bulk_interval: 5s
providers:
  aws:
    accounts:
      - name: development
        regions: [us-east-1]
        # Configure authentication (see AWS Setup Guide)
``` 