# DocStore API Monitoring Setup

This document describes the monitoring stack for the DocStore API in production.

## Components

### Grafana (Port 3000)
- **URL**: http://localhost:3000
- **Default Login**: admin / admin123
- **Purpose**: Visualization and dashboards for metrics and logs

### Prometheus (Port 9090)
- **URL**: http://localhost:9090
- **Purpose**: Metrics collection and storage
- **Scrapes**:
  - API metrics from `/metrics` endpoint
  - API health from `/health` endpoint

### Loki (Port 3100)
- **URL**: http://localhost:3100
- **Purpose**: Log aggregation and storage

### Promtail
- **Purpose**: Log shipping from Docker containers to Loki
- **Collects**: All container logs with proper labeling

## Available Dashboards

### DocStore API Dashboard
- API health status monitoring
- Log volume analysis by service
- Error rate tracking
- Recent API logs table

## Metrics Endpoints

### Health Check: `/health`
Returns JSON health status:
```json
{
  "status": "ok",
  "timestamp": "2023-01-01T00:00:00Z",
  "service": "docstore-api",
  "version": "1.0.0",
  "environment": "production"
}
```

### Metrics: `/metrics`
Returns Prometheus-compatible metrics:
- `docstore_api_info` - Service information
- `docstore_api_uptime_seconds` - Service uptime
- `docstore_api_memory_usage_bytes` - Current memory usage
- `docstore_api_memory_allocated_bytes` - Total allocated memory
- `docstore_api_goroutines` - Current goroutine count
- `docstore_api_health_status` - Health status (1=healthy, 0=unhealthy)

## Log Labels

All logs are automatically labeled with:
- `service` - Service name (docstore-api, nginx)
- `container_name` - Docker container name
- `environment` - Environment (production)

## Starting the Monitoring Stack

```bash
# Start all services including monitoring
make prod-up

# Or start monitoring services only
docker-compose -f docker/docker-compose.prod.yml up grafana prometheus loki promtail -d
```

## Accessing Services

1. **Grafana**: http://localhost:3000 (admin/admin123)
2. **Prometheus**: http://localhost:9090
3. **API Health**: http://localhost:8080/health
4. **API Metrics**: http://localhost:8080/metrics

## Troubleshooting

### No logs in Grafana
1. Check Promtail is running: `docker logs docstore-promtail-prod`
2. Verify Loki is accessible: `curl http://localhost:3100/ready`
3. Check container labels are correct

### No metrics in Prometheus
1. Check Prometheus targets: http://localhost:9090/targets
2. Verify API is responding: `curl http://localhost:8080/metrics`
3. Check Prometheus logs: `docker logs docstore-prometheus-prod`

### Grafana dashboard not loading
1. Verify datasources are configured correctly
2. Check Grafana logs: `docker logs docstore-grafana-prod`
3. Ensure dashboard JSON is valid
