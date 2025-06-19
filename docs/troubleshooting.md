# Troubleshooting Guide

This guide helps you diagnose and resolve common issues with Finala. If you encounter a problem not covered here, please check the service logs and open an issue on GitHub.

## Quick Diagnostics

### Check Service Status

```bash
# Check if all services are running
docker-compose ps

# Check service health
docker-compose -f docker-compose-hub.yaml ps
```

### View Service Logs

```bash
# View logs for all services
docker-compose logs

# View logs for specific service
docker-compose logs api
docker-compose logs collector
docker-compose logs ui
docker-compose logs meilisearch

# Follow logs in real-time
docker-compose logs -f collector
```

### Check API Health

```bash
# Test API health endpoint
curl http://localhost:8089/api/v1/health

# Test with authentication
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8089/api/v1/health
```

## Common Issues

### 1. Services Not Starting

**Symptoms**: Services show as "Exit" or "Restarting" status

**Check Service Logs**:
```bash
docker-compose logs [service-name]
```

**Common Causes and Solutions**:

#### Port Conflicts
**Error**: `Bind for 0.0.0.0:8080 failed: port is already allocated`

**Solution**:
```bash
# Check what's using the port
sudo lsof -i :8080
sudo netstat -tulpn | grep :8080

# Stop conflicting service or change port in docker-compose.yaml
```

#### Insufficient Resources
**Error**: `Cannot start service: OOMKilled`

**Solution**:
- Increase Docker memory limit (minimum 4GB recommended)
- Reduce resource limits in docker-compose.yaml
- Close other applications using system resources

#### Configuration Errors
**Error**: `yaml: unmarshal errors`

**Solution**:
```bash
# Validate YAML syntax
yamllint configuration/*.yaml

# Check for indentation issues
cat -A configuration/collector.yaml
```

### 2. Authentication Issues

**Symptoms**: Cannot log in to web interface, 401 errors

#### Invalid Credentials
**Error**: `Invalid username or password`

**Solution**:
1. Check credentials in `configuration/api.yaml`:
```yaml
auth:
  username: admin
  password: your_password
```

2. Restart API service:
```bash
docker-compose restart api
```

3. Clear browser cache and cookies

#### Auto-Generated Password
**Check logs for auto-generated password**:
```bash
docker-compose logs api | grep "Generated random password"
```

**Expected output**:
```
INFO: Generated random password for admin user: XyzT7q2PwC8rLzV5
```

### 3. AWS Connection Issues

**Symptoms**: Collector fails to connect to AWS, no resources found

#### Invalid AWS Credentials
**Error**: `NoCredentialProviders: no valid providers in chain`

**Solution**:
1. Verify AWS credentials:
```bash
# Test AWS CLI
aws sts get-caller-identity

# Check environment variables
echo $AWS_ACCESS_KEY_ID
echo $AWS_SECRET_ACCESS_KEY
```

2. Update configuration:
```yaml
# configuration/collector.yaml
providers:
  aws:
    accounts:
      - name: production
        regions: [us-east-1]
        access_key: YOUR_ACCESS_KEY
        secret_key: YOUR_SECRET_KEY
```

#### Insufficient IAM Permissions
**Error**: `AccessDenied: User is not authorized to perform`

**Solution**:
1. Check IAM permissions (see [AWS Setup Guide](aws-setup.md))
2. Test specific permissions:
```bash
aws ec2 describe-instances --region us-east-1
aws cloudwatch get-metric-statistics --namespace AWS/EC2 --metric-name CPUUtilization
```

3. Verify required policies are attached

#### Region Access Issues
**Error**: `Could not connect to the endpoint URL`

**Solution**:
1. Check region availability:
```bash
aws ec2 describe-regions
```

2. Update configuration with accessible regions:
```yaml
regions:
  - us-east-1
  - us-west-2
```

### 4. CloudWatch Metrics Issues

**Symptoms**: Resources found but no utilization data

#### No Metrics Available
**Error**: `No metrics found for the specified time range`

**Solution**:
1. Enable detailed monitoring for EC2 instances
2. Check metric retention period (default: 15 months)
3. Verify time range in configuration:
```yaml
period: 24h
start_time: 168h  # 7 days
```

#### Metric Access Denied
**Error**: `AccessDenied: User is not authorized to perform: cloudwatch:GetMetricStatistics`

**Solution**:
Add CloudWatch permissions to IAM policy:
```json
{
  "Effect": "Allow",
  "Action": [
    "cloudwatch:GetMetricStatistics",
    "cloudwatch:ListMetrics"
  ],
  "Resource": "*"
}
```

### 5. Meilisearch Issues

**Symptoms**: Search not working, API errors

#### Connection Issues
**Error**: `Failed to connect to Meilisearch`

**Solution**:
1. Check Meilisearch service:
```bash
docker-compose logs meilisearch
```

2. Verify configuration:
```yaml
# configuration/api.yaml
storage:
  meilisearch:
    endpoints: 
      - http://meilisearch:7700
    password: "your_master_key"
```

3. Restart Meilisearch:
```bash
docker-compose restart meilisearch
```

#### Index Issues
**Error**: `Index not found`

**Solution**:
1. Recreate index:
```bash
# Access Meilisearch console
curl -X DELETE "http://localhost:7700/indexes/resources" \
  -H "Authorization: Bearer your_master_key"

# Restart collector to recreate index
docker-compose restart collector
```

### 6. Resource Collection Issues

**Symptoms**: No resources found, collection fails

#### No Resources Detected
**Check collector logs**:
```bash
docker-compose logs collector | grep "Found"
```

**Common causes**:
1. **No resources in specified regions**
2. **Insufficient permissions**
3. **Resources filtered out by detection rules**

**Debug steps**:
```bash
# Test AWS resource discovery
aws ec2 describe-instances --region us-east-1
aws rds describe-db-instances --region us-east-1

# Check collector configuration
cat configuration/collector.yaml
```

#### Collection Timeout
**Error**: `context deadline exceeded`

**Solution**:
1. Increase timeout in configuration:
```yaml
api_server:
  bulk_interval: 10s  # Increase from 5s
```

2. Reduce regions or accounts being scanned
3. Check network connectivity to AWS

### 7. Web Interface Issues

**Symptoms**: UI not loading, API errors

#### UI Not Accessible
**Error**: `Connection refused`

**Solution**:
1. Check UI service:
```bash
docker-compose logs ui
```

2. Verify port mapping:
```yaml
# docker-compose.yaml
ui:
  ports:
    - "8080:8080"
```

3. Check firewall settings:
```bash
sudo ufw status
```

#### API Connection Issues
**Error**: `Failed to fetch from API`

**Solution**:
1. Check API configuration:
```yaml
# configuration/ui.yaml
api_server:
  address: http://127.0.0.1:8089
```

2. Verify API is running:
```bash
curl http://localhost:8089/api/v1/health
```

### 8. Notification Issues

**Symptoms**: Slack/email notifications not working

#### Slack Integration
**Error**: `Invalid token`

**Solution**:
1. Verify Slack token:
```yaml
# configuration/notifier.yaml
notifiers:
  slack:
    token: xoxb-your-bot-token
```

2. Check bot permissions in Slack workspace
3. Verify channel names and user mentions

#### Email Notifications
**Error**: `SMTP authentication failed`

**Solution**:
1. Check SMTP configuration:
```yaml
# configuration/api.yaml
smtp:
  username: "your_email@example.com"
  password: "your_app_password"
  smtpServer: "smtp.gmail.com"
  smtpPort: 587
```

2. Use app passwords for Gmail
3. Check firewall blocking SMTP ports

## Performance Issues

### 1. Slow Resource Collection

**Symptoms**: Collection takes too long

**Optimization**:
1. **Reduce regions**: Only scan necessary regions
2. **Increase batch size**: Adjust `bulk_interval`
3. **Parallel processing**: Use multiple collector instances
4. **Caching**: Enable CloudWatch metric caching

### 2. High Memory Usage

**Symptoms**: Services using excessive memory

**Solution**:
1. **Monitor memory usage**:
```bash
docker stats
```

2. **Adjust memory limits**:
```yaml
# docker-compose.yaml
services:
  collector:
    deploy:
      resources:
        limits:
          memory: 2G
```

3. **Optimize configuration**: Reduce concurrent operations

### 3. API Response Slow

**Symptoms**: Web interface slow to load

**Solution**:
1. **Check Meilisearch performance**:
```bash
curl "http://localhost:7700/health"
```

2. **Optimize search queries**: Use specific filters
3. **Enable caching**: Configure response caching
4. **Scale API service**: Add more API instances

## Debugging Steps

### 1. Enable Debug Logging

**Update configuration**:
```yaml
# configuration/collector.yaml
log_level: debug

# configuration/api.yaml
log_level: debug
```

**Restart services**:
```bash
docker-compose restart collector api
```

### 2. Test AWS Connectivity

```bash
# Test basic connectivity
aws sts get-caller-identity

# Test specific services
aws ec2 describe-regions
aws cloudwatch list-metrics --namespace AWS/EC2

# Test pricing API
aws pricing get-products --service-code AmazonEC2
```

### 3. Validate Configuration

```bash
# Check YAML syntax
yamllint configuration/*.yaml

# Validate collector configuration
docker-compose exec collector finala collector --validate-config

# Test API configuration
curl -X POST http://localhost:8089/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "test"}'
```

### 4. Check Network Connectivity

```bash
# Test internal communication
docker-compose exec api ping meilisearch
docker-compose exec collector ping api

# Test external connectivity
docker-compose exec collector ping ec2.us-east-1.amazonaws.com
```

## Recovery Procedures

### 1. Reset Meilisearch Data

**If search data is corrupted**:
```bash
# Stop services
docker-compose down

# Remove Meilisearch volume
docker volume rm finala_meilisearch_data

# Restart services
docker-compose up -d
```

### 2. Reset Configuration

**If configuration is corrupted**:
```bash
# Backup current configuration
cp -r configuration configuration.backup

# Restore from backup or recreate
# Edit configuration files as needed

# Restart services
docker-compose restart
```

### 3. Complete Reset

**For complete system reset**:
```bash
# Stop and remove everything
docker-compose down -v
docker system prune -f

# Reclone repository
cd ..
rm -rf finala
git clone https://github.com/similarweb/finala.git
cd finala

# Start fresh
docker-compose up -d
```

## Monitoring and Alerts

### 1. Service Health Monitoring

**Create health check script**:
```bash
#!/bin/bash
# health-check.sh

API_URL="http://localhost:8089/api/v1/health"
UI_URL="http://localhost:8080"

# Check API
if curl -f $API_URL > /dev/null 2>&1; then
    echo "API: OK"
else
    echo "API: FAILED"
    exit 1
fi

# Check UI
if curl -f $UI_URL > /dev/null 2>&1; then
    echo "UI: OK"
else
    echo "UI: FAILED"
    exit 1
fi
```

### 2. Log Monitoring

**Set up log aggregation**:
```bash
# Forward logs to external system
docker-compose logs -f | tee /var/log/finala/combined.log

# Monitor for errors
docker-compose logs -f | grep -i error
```

### 3. Resource Monitoring

**Monitor system resources**:
```bash
# Check Docker resource usage
docker stats --no-stream

# Monitor disk space
df -h

# Check memory usage
free -h
```

## Getting Help

### 1. Collect Debug Information

**Before opening an issue, collect**:
```bash
# System information
docker version
docker-compose version
uname -a

# Service status
docker-compose ps
docker-compose logs > finala-logs.txt

# Configuration (remove sensitive data)
tar -czf finala-config.tar.gz configuration/

# AWS connectivity test
aws sts get-caller-identity > aws-test.txt
```

### 2. Open GitHub Issue

**Include in your issue**:
- **Description**: What you're trying to do
- **Expected behavior**: What should happen
- **Actual behavior**: What's happening
- **Steps to reproduce**: How to trigger the issue
- **Environment**: OS, Docker version, AWS setup
- **Logs**: Relevant error messages
- **Configuration**: Sanitized config files

### 3. Community Support

- **GitHub Discussions**: Ask questions and share solutions
- **GitHub Issues**: Report bugs and request features
- **Code Review**: Contribute fixes and improvements

## Prevention

### 1. Regular Maintenance

- **Update dependencies**: Keep Docker images current
- **Monitor logs**: Check for warnings and errors
- **Backup configuration**: Version control your configs
- **Test regularly**: Run health checks periodically

### 2. Best Practices

- **Use IAM roles**: Avoid hardcoded credentials
- **Limit permissions**: Follow least privilege principle
- **Monitor costs**: Track AWS API usage
- **Document changes**: Keep configuration changes documented

### 3. Security Considerations

- **Rotate credentials**: Change passwords and keys regularly
- **Network security**: Use VPC endpoints when possible
- **Access control**: Limit who can access Finala
- **Audit logs**: Monitor access and changes 