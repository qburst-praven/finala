# API Reference

Finala provides a RESTful API for accessing resource data, managing configurations, and integrating with external systems.

## Base URL

The API is available at:
- **Local Development**: `http://localhost:8089`
- **Docker Compose**: `http://api:8081` (internal)
- **Production**: Configure based on your deployment

## Authentication

All API endpoints require authentication using JWT tokens.

### Login

**Endpoint**: `POST /api/v1/auth/login`

**Request Body**:
```json
{
  "username": "admin",
  "password": "your_password"
}
```

**Response**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "user": {
    "username": "admin"
  }
}
```

**Usage**:
```bash
curl -X POST http://localhost:8089/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "your_password"}'
```

### Using the Token

Include the JWT token in the Authorization header:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8089/api/v1/resources
```

## Resources Endpoints

### List Resources

**Endpoint**: `GET /api/v1/resources`

**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50, max: 100)
- `search` (optional): Search query
- `service` (optional): Filter by service (e.g., "ec2", "rds")
- `region` (optional): Filter by region
- `account` (optional): Filter by account name
- `tags` (optional): Filter by tags (comma-separated)
- `sort_by` (optional): Sort field (default: "created_at")
- `sort_order` (optional): Sort order ("asc" or "desc", default: "desc")

**Response**:
```json
{
  "data": [
    {
      "id": "i-1234567890abcdef0",
      "name": "web-server-01",
      "service": "ec2",
      "region": "us-east-1",
      "account": "production",
      "cost": 45.67,
      "tags": {
        "Environment": "production",
        "Team": "engineering"
      },
      "metrics": {
        "cpu_utilization": 15.5,
        "connection_count": 0
      },
      "detection_rules": [
        {
          "description": "Low CPU utilization",
          "status": "detected"
        }
      ],
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1250,
    "pages": 25
  }
}
```

**Usage Examples**:

```bash
# Get all resources
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8089/api/v1/resources

# Search for specific resources
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8089/api/v1/resources?search=web-server&service=ec2"

# Filter by region and tags
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8089/api/v1/resources?region=us-east-1&tags=production,engineering"

# Pagination
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8089/api/v1/resources?page=2&limit=25"
```

### Get Resource Details

**Endpoint**: `GET /api/v1/resources/{id}`

**Response**:
```json
{
  "id": "i-1234567890abcdef0",
  "name": "web-server-01",
  "service": "ec2",
  "region": "us-east-1",
  "account": "production",
  "instance_type": "t3.medium",
  "state": "running",
  "cost": 45.67,
  "monthly_cost": 1368.90,
  "tags": {
    "Environment": "production",
    "Team": "engineering",
    "Project": "web-app"
  },
  "metrics": {
    "cpu_utilization": {
      "current": 15.5,
      "average": 12.3,
      "maximum": 45.2
    },
    "network_in": {
      "current": 1024,
      "average": 2048
    }
  },
  "detection_rules": [
    {
      "description": "Low CPU utilization",
      "status": "detected",
      "threshold": 40,
      "current_value": 15.5
    }
  ],
  "recommendations": [
    {
      "type": "downsize",
      "description": "Consider downsizing to t3.small",
      "potential_savings": 22.84
    }
  ],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Usage**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8089/api/v1/resources/i-1234567890abcdef0
```

### Update Resource Tags

**Endpoint**: `PUT /api/v1/resources/{id}/tags`

**Request Body**:
```json
{
  "tags": {
    "Environment": "staging",
    "Team": "devops",
    "CostCenter": "IT-001"
  }
}
```

**Response**:
```json
{
  "message": "Tags updated successfully",
  "resource_id": "i-1234567890abcdef0"
}
```

**Usage**:
```bash
curl -X PUT \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"tags": {"Environment": "staging", "Team": "devops"}}' \
  http://localhost:8089/api/v1/resources/i-1234567890abcdef0/tags
```

## Statistics Endpoints

### Dashboard Statistics

**Endpoint**: `GET /api/v1/statistics`

**Response**:
```json
{
  "total_resources": 1250,
  "total_cost": 45678.90,
  "potential_savings": 12345.67,
  "resources_by_service": {
    "ec2": 450,
    "rds": 120,
    "dynamodb": 80,
    "lambda": 200,
    "elasticache": 50
  },
  "resources_by_region": {
    "us-east-1": 800,
    "us-west-2": 300,
    "eu-west-1": 150
  },
  "detection_summary": {
    "underutilized": 234,
    "unused": 56,
    "over_provisioned": 89
  },
  "cost_trends": {
    "current_month": 45678.90,
    "previous_month": 43210.50,
    "change_percentage": 5.7
  }
}
```

**Usage**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8089/api/v1/statistics
```

### Service Statistics

**Endpoint**: `GET /api/v1/statistics/services`

**Query Parameters**:
- `service` (optional): Specific service name

**Response**:
```json
{
  "ec2": {
    "total_resources": 450,
    "total_cost": 23456.78,
    "potential_savings": 5678.90,
    "detection_rules": {
      "low_cpu": 123,
      "unused_volumes": 45,
      "unattached_ips": 12
    }
  },
  "rds": {
    "total_resources": 120,
    "total_cost": 12345.67,
    "potential_savings": 2345.67,
    "detection_rules": {
      "low_connections": 34,
      "over_provisioned": 23
    }
  }
}
```

**Usage**:
```bash
# All services
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8089/api/v1/statistics/services

# Specific service
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8089/api/v1/statistics/services?service=ec2"
```

## Tags Endpoints

### List All Tags

**Endpoint**: `GET /api/v1/tags`

**Response**:
```json
{
  "tags": {
    "Environment": ["production", "staging", "development"],
    "Team": ["engineering", "devops", "data"],
    "Project": ["web-app", "api-service", "data-pipeline"]
  }
}
```

**Usage**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8089/api/v1/tags
```

### Get Resources by Tag

**Endpoint**: `GET /api/v1/tags/{tag_name}/{tag_value}/resources`

**Query Parameters**:
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response**:
```json
{
  "data": [
    {
      "id": "i-1234567890abcdef0",
      "name": "web-server-01",
      "service": "ec2",
      "cost": 45.67
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 150,
    "pages": 3
  }
}
```

**Usage**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8089/api/v1/tags/Environment/production/resources
```

## Executions Endpoints

### List Collection Executions

**Endpoint**: `GET /api/v1/executions`

**Query Parameters**:
- `page` (optional): Page number
- `limit` (optional): Items per page
- `status` (optional): Filter by status ("running", "completed", "failed")

**Response**:
```json
{
  "data": [
    {
      "id": "exec_1234567890",
      "started_at": "2024-01-15T10:00:00Z",
      "completed_at": "2024-01-15T10:15:00Z",
      "status": "completed",
      "resources_found": 1250,
      "resources_analyzed": 1250,
      "errors": [],
      "duration_seconds": 900
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 25,
    "pages": 1
  }
}
```

**Usage**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8089/api/v1/executions
```

### Get Execution Details

**Endpoint**: `GET /api/v1/executions/{id}`

**Response**:
```json
{
  "id": "exec_1234567890",
  "started_at": "2024-01-15T10:00:00Z",
  "completed_at": "2024-01-15T10:15:00Z",
  "status": "completed",
  "resources_found": 1250,
  "resources_analyzed": 1250,
  "errors": [],
  "duration_seconds": 900,
  "accounts_scanned": ["production", "development"],
  "regions_scanned": ["us-east-1", "us-west-2"],
  "services_scanned": ["ec2", "rds", "dynamodb"],
  "summary": {
    "underutilized_resources": 234,
    "unused_resources": 56,
    "potential_savings": 12345.67
  }
}
```

**Usage**:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8089/api/v1/executions/exec_1234567890
```

## Search Endpoints

### Advanced Search

**Endpoint**: `POST /api/v1/search`

**Request Body**:
```json
{
  "query": "web-server",
  "filters": {
    "service": ["ec2", "rds"],
    "region": ["us-east-1"],
    "tags": {
      "Environment": "production"
    },
    "cost_range": {
      "min": 10,
      "max": 100
    }
  },
  "sort": {
    "field": "cost",
    "order": "desc"
  },
  "pagination": {
    "page": 1,
    "limit": 50
  }
}
```

**Response**:
```json
{
  "data": [
    {
      "id": "i-1234567890abcdef0",
      "name": "web-server-01",
      "service": "ec2",
      "cost": 45.67,
      "score": 0.95
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 25,
    "pages": 1
  },
  "facets": {
    "service": {
      "ec2": 20,
      "rds": 5
    },
    "region": {
      "us-east-1": 25
    }
  }
}
```

**Usage**:
```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"query": "web-server", "filters": {"service": ["ec2"]}}' \
  http://localhost:8089/api/v1/search
```

## Health Check

### API Health

**Endpoint**: `GET /api/v1/health`

**Response**:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "2.0.0",
  "services": {
    "meilisearch": "healthy",
    "collector": "healthy"
  }
}
```

**Usage**:
```bash
curl http://localhost:8089/api/v1/health
```

## Error Handling

### Error Response Format

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "Resource with ID 'i-1234567890abcdef0' not found",
    "details": {
      "resource_id": "i-1234567890abcdef0"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Invalid or missing authentication |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `RESOURCE_NOT_FOUND` | 404 | Resource not found |
| `VALIDATION_ERROR` | 400 | Invalid request parameters |
| `INTERNAL_ERROR` | 500 | Internal server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

### Rate Limiting

The API implements rate limiting to prevent abuse:

- **Default Limit**: 100 requests per minute per IP
- **Response Headers**:
  - `X-RateLimit-Limit`: Request limit
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Reset time

**Rate Limit Exceeded Response**:
```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Try again in 60 seconds."
  }
}
```

## SDK Examples

### Python

```python
import requests
import json

class FinalaAPI:
    def __init__(self, base_url, username, password):
        self.base_url = base_url
        self.token = self._authenticate(username, password)
    
    def _authenticate(self, username, password):
        response = requests.post(
            f"{self.base_url}/api/v1/auth/login",
            json={"username": username, "password": password}
        )
        return response.json()["token"]
    
    def get_resources(self, **params):
        headers = {"Authorization": f"Bearer {self.token}"}
        response = requests.get(
            f"{self.base_url}/api/v1/resources",
            headers=headers,
            params=params
        )
        return response.json()
    
    def get_statistics(self):
        headers = {"Authorization": f"Bearer {self.token}"}
        response = requests.get(
            f"{self.base_url}/api/v1/statistics",
            headers=headers
        )
        return response.json()

# Usage
api = FinalaAPI("http://localhost:8089", "admin", "password")
resources = api.get_resources(service="ec2", region="us-east-1")
stats = api.get_statistics()
```

### JavaScript/Node.js

```javascript
class FinalaAPI {
    constructor(baseUrl, username, password) {
        this.baseUrl = baseUrl;
        this.token = null;
        this.authenticate(username, password);
    }
    
    async authenticate(username, password) {
        const response = await fetch(`${this.baseUrl}/api/v1/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });
        const data = await response.json();
        this.token = data.token;
    }
    
    async getResources(params = {}) {
        const queryString = new URLSearchParams(params).toString();
        const response = await fetch(
            `${this.baseUrl}/api/v1/resources?${queryString}`,
            {
                headers: {
                    'Authorization': `Bearer ${this.token}`
                }
            }
        );
        return response.json();
    }
    
    async getStatistics() {
        const response = await fetch(
            `${this.baseUrl}/api/v1/statistics`,
            {
                headers: {
                    'Authorization': `Bearer ${this.token}`
                }
            }
        );
        return response.json();
    }
}

// Usage
const api = new FinalaAPI('http://localhost:8089', 'admin', 'password');
api.getResources({ service: 'ec2', region: 'us-east-1' })
    .then(resources => console.log(resources));
```

## Webhook Integration

### Configure Webhooks

**Endpoint**: `POST /api/v1/webhooks`

**Request Body**:
```json
{
  "url": "https://your-app.com/webhook",
  "events": ["resource.detected", "execution.completed"],
  "secret": "your-webhook-secret"
}
```

**Response**:
```json
{
  "id": "webhook_1234567890",
  "url": "https://your-app.com/webhook",
  "events": ["resource.detected", "execution.completed"],
  "status": "active",
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Webhook Events

**Resource Detected**:
```json
{
  "event": "resource.detected",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "resource_id": "i-1234567890abcdef0",
    "service": "ec2",
    "cost": 45.67,
    "detection_rules": ["low_cpu_utilization"]
  }
}
```

**Execution Completed**:
```json
{
  "event": "execution.completed",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "execution_id": "exec_1234567890",
    "resources_found": 1250,
    "potential_savings": 12345.67
  }
}
```

## API Versioning

The API uses URL versioning (`/api/v1/`). Future versions will be available at `/api/v2/`, etc.

**Version Deprecation Policy**:
- Major versions are supported for at least 12 months
- Deprecation notices are sent 6 months in advance
- Breaking changes only occur in major versions

## Support

For API support:
- Check the [Troubleshooting Guide](troubleshooting.md)
- Review service logs for detailed error information
- Open an issue on GitHub with API-specific problems 