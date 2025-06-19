# AWS Setup Guide

This guide covers everything you need to set up AWS integration with Finala, including IAM permissions, authentication methods, and best practices.

## Prerequisites

Before setting up AWS integration, ensure you have:

- An AWS account with appropriate permissions
- Access to AWS IAM console
- Understanding of AWS services and regions
- Finala deployed and running

## Authentication Methods

Finala supports multiple AWS authentication methods. Choose the one that best fits your security requirements.

### 1. AWS Access Keys (Quick Start)

**Use Case**: Development, testing, or single-account setups

**Setup Steps**:

1. **Create IAM User**:
   ```bash
   # Using AWS CLI
   aws iam create-user --user-name finala-collector
   ```

2. **Attach Required Policies**:
   ```bash
   # Attach the Finala policy (see below)
   aws iam attach-user-policy --user-name finala-collector --policy-arn arn:aws:iam::YOUR_ACCOUNT:policy/FinalaCollectorPolicy
   ```

3. **Create Access Keys**:
   ```bash
   aws iam create-access-key --user-name finala-collector
   ```

4. **Configure Finala**:
   ```yaml
   # configuration/collector.yaml
   providers:
     aws:
       accounts:
         - name: production
           regions: [us-east-1, us-west-2]
           access_key: YOUR_ACCESS_KEY
           secret_key: YOUR_SECRET_KEY
   ```

**Security Considerations**:
- Rotate keys regularly
- Use least privilege principle
- Monitor key usage
- Consider using temporary credentials

### 2. AWS IAM Role (Recommended for Production)

**Use Case**: Production environments, cross-account access, enhanced security

**Setup Steps**:

1. **Create IAM Role**:
   ```json
   {
     "Version": "2012-10-17",
     "Statement": [
       {
         "Effect": "Allow",
         "Principal": {
           "Service": "ec2.amazonaws.com"
         },
         "Action": "sts:AssumeRole"
       }
     ]
   }
   ```

2. **Attach Required Policies**:
   ```bash
   aws iam attach-role-policy --role-name FinalaCollectorRole --policy-arn arn:aws:iam::YOUR_ACCOUNT:policy/FinalaCollectorPolicy
   ```

3. **Configure Finala**:
   ```yaml
   # configuration/collector.yaml
   providers:
     aws:
       accounts:
         - name: production
           regions: [us-east-1, us-west-2]
           role: arn:aws:iam::123456789012:role/FinalaCollectorRole
   ```

### 3. AWS Profile (Local Development)

**Use Case**: Local development, multiple AWS profiles

**Setup Steps**:

1. **Configure AWS CLI**:
   ```bash
   aws configure --profile finala
   ```

2. **Configure Finala**:
   ```yaml
   # configuration/collector.yaml
   providers:
     aws:
       accounts:
         - name: production
           regions: [us-east-1]
           profile: finala
   ```

## Required IAM Permissions

### Minimum Required Permissions

Create an IAM policy with the following permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeInstances",
        "ec2:DescribeVolumes",
        "ec2:DescribeAddresses",
        "ec2:DescribeLoadBalancers",
        "ec2:DescribeLoadBalancerAttributes",
        "ec2:DescribeNatGateways",
        "rds:DescribeDBInstances",
        "rds:DescribeDBClusters",
        "dynamodb:ListTables",
        "dynamodb:DescribeTable",
        "elasticache:DescribeCacheClusters",
        "es:ListDomainNames",
        "es:DescribeElasticsearchDomain",
        "lambda:ListFunctions",
        "lambda:GetFunction",
        "kinesis:ListStreams",
        "kinesis:DescribeStream",
        "redshift:DescribeClusters",
        "neptune:DescribeDBClusters",
        "apigateway:GET",
        "iam:ListUsers",
        "iam:GetUser",
        "iam:GetAccessKeyLastUsed",
        "cloudwatch:GetMetricStatistics",
        "pricing:GetProducts",
        "sts:GetCallerIdentity"
      ],
      "Resource": "*"
    }
  ]
}
```

### Granular Permissions (Recommended)

For enhanced security, use more granular permissions:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeInstances",
        "ec2:DescribeVolumes",
        "ec2:DescribeAddresses",
        "ec2:DescribeLoadBalancers",
        "ec2:DescribeLoadBalancerAttributes",
        "ec2:DescribeNatGateways"
      ],
      "Resource": "*",
      "Condition": {
        "StringEquals": {
          "aws:RequestedRegion": ["us-east-1", "us-west-2"]
        }
      }
    },
    {
      "Effect": "Allow",
      "Action": [
        "rds:DescribeDBInstances",
        "rds:DescribeDBClusters"
      ],
      "Resource": "*",
      "Condition": {
        "StringEquals": {
          "aws:RequestedRegion": ["us-east-1", "us-west-2"]
        }
      }
    },
    {
      "Effect": "Allow",
      "Action": [
        "cloudwatch:GetMetricStatistics"
      ],
      "Resource": "*",
      "Condition": {
        "StringEquals": {
          "aws:RequestedRegion": ["us-east-1", "us-west-2"]
        }
      }
    },
    {
      "Effect": "Allow",
      "Action": [
        "pricing:GetProducts"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "sts:GetCallerIdentity"
      ],
      "Resource": "*"
    }
  ]
}
```

### Cross-Account Access

For scanning multiple AWS accounts, set up cross-account access:

1. **In the target account**, create a role that can be assumed:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::COLLECTOR_ACCOUNT:root"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "finala-collector"
        }
      }
    }
  ]
}
```

2. **In the collector account**, grant permission to assume the role:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Resource": "arn:aws:iam::TARGET_ACCOUNT:role/FinalaCollectorRole"
    }
  ]
}
```

3. **Configure Finala**:
```yaml
providers:
  aws:
    accounts:
      - name: production
        regions: [us-east-1, us-west-2]
        role: arn:aws:iam::TARGET_ACCOUNT:role/FinalaCollectorRole
        external_id: finala-collector
```

## CloudWatch Metrics Access

Finala requires CloudWatch metrics for resource analysis. Ensure the following:

### Required CloudWatch Permissions

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "cloudwatch:GetMetricStatistics",
        "cloudwatch:ListMetrics"
      ],
      "Resource": "*"
    }
  ]
}
```

### CloudWatch Metrics Configuration

Some metrics may require additional setup:

1. **EC2 Detailed Monitoring**: Enable for instances you want to monitor
2. **Custom Metrics**: Ensure custom metrics are accessible
3. **Metric Retention**: Verify metrics are available for the configured time period

## AWS Pricing API Access

For cost calculations, Finala uses the AWS Pricing API:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "pricing:GetProducts"
      ],
      "Resource": "*"
    }
  ]
}
```

## Region Configuration

### Supported Regions

Finala supports all AWS regions. Configure regions based on your infrastructure:

```yaml
providers:
  aws:
    accounts:
      - name: production
        regions:
          - us-east-1      # US East (N. Virginia)
          - us-west-2      # US West (Oregon)
          - eu-west-1      # Europe (Ireland)
          - ap-southeast-1 # Asia Pacific (Singapore)
```

### Region-Specific Considerations

1. **Service Availability**: Not all services are available in all regions
2. **Pricing**: Costs may vary by region
3. **Latency**: Consider network latency for API calls
4. **Compliance**: Ensure regions meet compliance requirements

## Security Best Practices

### 1. Least Privilege Principle

- Grant only necessary permissions
- Use specific resource ARNs when possible
- Regularly review and audit permissions

### 2. Credential Management

- Use IAM roles instead of access keys
- Rotate credentials regularly
- Use temporary credentials when possible
- Store credentials securely

### 3. Network Security

- Use VPC endpoints for AWS services
- Restrict access to specific IP ranges
- Use security groups and NACLs

### 4. Monitoring and Auditing

- Enable CloudTrail for API logging
- Monitor IAM user activity
- Set up alerts for unusual activity
- Regular security assessments

## Troubleshooting

### Common Issues

**Access Denied Errors**:
```bash
# Check current identity
aws sts get-caller-identity

# Test specific permissions
aws ec2 describe-instances --region us-east-1
```

**CloudWatch Metrics Not Available**:
- Verify instance has detailed monitoring enabled
- Check metric retention period
- Ensure metrics exist for the specified time range

**Cross-Account Access Issues**:
- Verify trust relationship is configured correctly
- Check external ID matches
- Ensure role ARN is correct

**Pricing API Errors**:
- Verify pricing API is available in your region
- Check API rate limits
- Ensure pricing data is accessible

### Debugging Steps

1. **Enable Debug Logging**:
   ```yaml
   # configuration/collector.yaml
   log_level: debug
   ```

2. **Test AWS Connectivity**:
   ```bash
   # Test basic connectivity
   aws sts get-caller-identity
   
   # Test specific service access
   aws ec2 describe-regions
   ```

3. **Check Service Logs**:
   ```bash
   docker-compose logs collector
   ```

4. **Verify Configuration**:
   - Check YAML syntax
   - Verify account configuration
   - Ensure regions are accessible

## Performance Optimization

### API Rate Limits

AWS has rate limits for API calls. Optimize by:

1. **Batch Operations**: Use bulk operations when possible
2. **Caching**: Cache frequently accessed data
3. **Parallel Processing**: Process multiple regions concurrently
4. **Retry Logic**: Implement exponential backoff

### Resource Limits

Monitor and adjust based on your infrastructure:

1. **API Throttling**: Implement rate limiting
2. **Memory Usage**: Monitor collector memory consumption
3. **Network Bandwidth**: Consider bandwidth limitations
4. **Storage**: Monitor Meilisearch storage usage

## Compliance and Governance

### Data Residency

- Ensure data processing meets compliance requirements
- Configure regions based on data residency needs
- Review data retention policies

### Audit Requirements

- Enable CloudTrail for API logging
- Implement log retention policies
- Regular security assessments
- Compliance reporting

### Cost Management

- Monitor AWS API costs
- Implement cost alerts
- Regular cost optimization reviews
- Budget tracking and reporting 