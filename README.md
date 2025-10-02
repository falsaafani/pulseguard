# PulseGuard - Uptime & Anomaly Monitoring Platform

A distributed monitoring platform built with Go, TypeScript, Python, Docker, and Kubernetes for learning and hands-on experience with modern microservices architecture.

## What it does

- **Checks** a list of URLs/services on a schedule
- **Stores** latency and status history
- **Detects anomalies** in latency and error rates using statistical methods
- **Visualizes** results on a web dashboard
- **Notifies** via Slack when incidents happen

## Architecture

### Services

| Service | Tech | Purpose |
|---------|------|---------|
| **probe-api** | Go | REST API, probes targets, publishes to NATS, stores in PostgreSQL |
| **anomaly-worker** | Python | Subscribes to probe results, detects anomalies using z-score |
| **webapp** | TypeScript/Next.js | Dashboard UI for monitoring and administration |
| **notifier** | Go | Listens for incidents, sends Slack notifications |

### Infrastructure

- **PostgreSQL** - Data storage
- **NATS** - Message bus for inter-service communication
- **Docker** - Containerization
- **Kubernetes** - Orchestration

## Project Structure

```
pulseguard/
├── probe-api/              # Go service for probing targets
│   ├── cmd/
│   ├── internal/
│   ├── Dockerfile
│   └── go.mod
├── anomaly-worker/         # Python service for anomaly detection
│   ├── app/
│   ├── Dockerfile
│   └── requirements.txt
├── webapp/                 # Next.js dashboard
│   ├── src/
│   ├── Dockerfile
│   └── package.json
├── notifier/               # Go service for notifications
│   ├── cmd/
│   ├── internal/
│   ├── Dockerfile
│   └── go.mod
├── deploy/
│   └── k8s/
│       └── base/          # Kubernetes manifests
├── scripts/
│   └── init.sql           # Database schema
├── docker-compose.yml     # Local development setup
└── Makefile              # Common commands
```

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.22+
- Node.js 20+
- Python 3.11+
- kubectl (for K8s deployment)
- kind (for local K8s cluster)

### Local Development with Docker Compose

1. **Start all services:**
   ```bash
   make up
   ```

2. **View logs:**
   ```bash
   make logs
   ```

3. **Access services:**
   - Webapp: http://localhost:3000
   - Probe API: http://localhost:8080
   - NATS Monitor: http://localhost:8222

4. **Stop services:**
   ```bash
   make down
   ```

### Kubernetes Deployment

1. **Create kind cluster:**
   ```bash
   make kind-up
   ```

2. **Build and load images:**
   ```bash
   make kind-load
   ```

3. **Deploy to cluster:**
   ```bash
   make k8s-deploy
   ```

4. **Check status:**
   ```bash
   make k8s-status
   ```

5. **Port forward to access services:**
   ```bash
   make port-forward-webapp  # Access at localhost:3000
   make port-forward-api     # Access at localhost:8080
   ```

## Database Schema

```sql
targets (
  id, name, url, type, enabled, created_at
)

probes (
  id, target_id, ts, latency_ms, status_code, ok
)

incidents (
  id, target_id, started_at, ended_at, kind, details
)
```

## API Endpoints

- `POST /targets` - Create new monitoring target
- `GET /targets` - List all targets
- `GET /status?target_id=N` - Get probe results for target
- `GET /metrics` - Prometheus metrics
- `GET /health` - Health check

## Development Roadmap

### Phase 1: Core Functionality ✅
- [x] Project scaffolding
- [ ] Implement probe-api endpoints
- [ ] HTTP/TCP probing logic
- [ ] Store probe results in PostgreSQL
- [ ] Publish results to NATS

### Phase 2: Anomaly Detection
- [ ] Z-score calculation in anomaly-worker
- [ ] Incident detection and storage
- [ ] Publish incident events to NATS

### Phase 3: UI & Notifications
- [ ] Webapp dashboard with real-time data
- [ ] Target management UI
- [ ] Incident timeline visualization
- [ ] Slack notification integration

### Phase 4: Kubernetes & Production
- [ ] Deploy to local K8s cluster
- [ ] Configure ingress and TLS
- [ ] Horizontal Pod Autoscaling
- [ ] Prometheus & Grafana observability

### Stretch Goals
- [ ] gRPC between services
- [ ] Multi-region probes
- [ ] Advanced anomaly detection (Prophet/ARIMA)
- [ ] SLO tracking and burn-rate alerts

## Configuration

Copy `.env.example` files in each service directory and configure:

```bash
# Database
DATABASE_URL=postgres://pulseguard:password@localhost:5432/pulseguard

# Message Bus
NATS_URL=nats://localhost:4222

# Slack (optional)
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

## Learning Goals

This project provides hands-on experience with:

- **Go**: REST APIs, goroutines, database interactions
- **TypeScript**: Next.js, React, API integration
- **Python**: Async programming, statistical analysis, message queues
- **Docker**: Multi-stage builds, containerization
- **Kubernetes**: Deployments, services, ConfigMaps, Secrets, Ingress
- **System Design**: Microservices, message buses, data storage

## Contributing

This is a learning project! Feel free to experiment, break things, and learn.

## License

MIT
