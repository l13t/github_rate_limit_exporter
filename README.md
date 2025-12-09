# GitHub Rate Limit Exporter

[![CI](https://github.com/l13t/github_rate_limit_exporter/workflows/CI/badge.svg)](https://github.com/l13t/github_rate_limit_exporter/actions/workflows/ci.yml)
[![Docker](https://github.com/l13t/github_rate_limit_exporter/workflows/Docker/badge.svg)](https://github.com/l13t/github_rate_limit_exporter/actions/workflows/docker.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/l13t/github_rate_limit_exporter)](https://goreportcard.com/report/github.com/l13t/github_rate_limit_exporter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Prometheus exporter for monitoring GitHub API rate limits across multiple users and tokens.

## Features

- 📊 Export GitHub API rate limits as Prometheus metrics
- 👥 Monitor multiple users/tokens simultaneously
- 📝 Support for YAML, TOML, and HCL configuration
- 🎯 Track Core, Search, GraphQL, and Integration Manifest limits
- 🐳 Multi-arch Docker images (amd64, arm64, armv7)
- 🔒 Secure, non-root execution

## Quick Start

### 1. Get a GitHub Token

Create a [Personal Access Token](https://github.com/settings/tokens) - no scopes required.

### 2. Install

**Pre-built binary:**
```bash
wget https://github.com/l13t/github_rate_limit_exporter/releases/latest/download/github_rate_limit_exporter-linux-amd64.tar.gz
tar -xzf github_rate_limit_exporter-linux-amd64.tar.gz
sudo mv github_rate_limit_exporter /usr/local/bin/
```

**Docker:**
```bash
docker pull ghcr.io/l13t/github_rate_limit_exporter:latest
```

**Build from source:**
```bash
git clone https://github.com/l13t/github_rate_limit_exporter.git
cd github_rate_limit_exporter
go build -o github_rate_limit_exporter ./cmd/exporter
```

### 3. Configure

Create `config.yaml`:

```yaml
users:
  - name: "my-user"
    token: "ghp_your_token_here"

listen_addr: ":9101"
metrics_path: "/metrics"
poll_interval: 60
```

### 4. Run

```bash
./github_rate_limit_exporter -config config.yaml
```

### 5. Verify

```bash
curl http://localhost:9101/metrics | grep github_rate_limit
```

## Configuration

### Formats

Supports YAML, TOML, and HCL. See [examples](config.yaml.example).

### Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `users` | array | *required* | GitHub users to monitor |
| `users[].name` | string | *required* | User identifier (metric label) |
| `users[].token` | string | *required* | GitHub PAT |
| `listen_addr` | string | `:9101` | Server address |
| `metrics_path` | string | `/metrics` | Metrics endpoint |
| `poll_interval` | int | `60` | Poll interval (seconds) |

### Multiple Users

```yaml
users:
  - name: "personal"
    token: "ghp_personal_token"
  - name: "ci-bot"
    token: "ghp_bot_token"
  - name: "team-shared"
    token: "ghp_team_token"
```

## Metrics

For each user, the following metrics are exported:

### Core API
```
github_rate_limit_core_limit{user="username"}
github_rate_limit_core_remaining{user="username"}
github_rate_limit_core_used{user="username"}
github_rate_limit_core_reset_timestamp{user="username"}
```

### Search API
```
github_rate_limit_search_limit{user="username"}
github_rate_limit_search_remaining{user="username"}
github_rate_limit_search_used{user="username"}
github_rate_limit_search_reset_timestamp{user="username"}
```

### GraphQL API
```
github_rate_limit_graphql_limit{user="username"}
github_rate_limit_graphql_remaining{user="username"}
github_rate_limit_graphql_used{user="username"}
github_rate_limit_graphql_reset_timestamp{user="username"}
```

## Docker

### Run Container

```bash
docker run -d \
  -p 9101:9101 \
  -v $(pwd)/config.yaml:/config.yaml:ro \
  ghcr.io/l13t/github_rate_limit_exporter:latest \
  -config /config.yaml
```

### Docker Compose

```bash
docker-compose up -d
```

Includes Prometheus (port 9090) and Grafana (port 3000).

## Prometheus Integration

Add to `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'github_rate_limits'
    static_configs:
      - targets: ['localhost:9101']
    scrape_interval: 60s
```

### Example Queries

```promql
# Remaining requests
github_rate_limit_core_remaining

# Usage percentage
(github_rate_limit_core_used / github_rate_limit_core_limit) * 100

# Time until reset (hours)
(github_rate_limit_core_reset_timestamp - time()) / 3600
```

### Alerts

```yaml
- alert: GitHubRateLimitLow
  expr: github_rate_limit_core_remaining < 1000
  for: 5m
  annotations:
    summary: "Rate limit low for {{ $labels.user }}"

- alert: GitHubRateLimitCritical
  expr: github_rate_limit_core_remaining < 100
  for: 1m
  annotations:
    summary: "Rate limit critical for {{ $labels.user }}"
```

See [alerts.yml](alerts.yml) for complete examples.

## Deployment

### Systemd

```bash
# Install binary
sudo cp github_rate_limit_exporter /usr/local/bin/

# Install service
sudo cp systemd/github_rate_limit_exporter.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now github_rate_limit_exporter
```

### Kubernetes

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: github-tokens
stringData:
  config.yaml: |
    users:
      - name: "bot"
        token: "ghp_token"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: github-rate-limit-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: github-rate-limit-exporter
  template:
    metadata:
      labels:
        app: github-rate-limit-exporter
    spec:
      containers:
      - name: exporter
        image: ghcr.io/l13t/github_rate_limit_exporter:latest
        args: ["-config", "/config/config.yaml"]
        ports:
        - containerPort: 9101
        volumeMounts:
        - name: config
          mountPath: /config
      volumes:
      - name: config
        secret:
          secretName: github-tokens
```

## Development

### Build with Task

```bash
# Install Task
brew install go-task/tap/go-task

# Build
task build

# Run tests
task test

# Run locally
task run
```

### Without Task

```bash
# Build
go build -o build/github_rate_limit_exporter ./cmd/exporter

# Test
go test ./...

# Run
go run ./cmd/exporter -config config.yaml
```

## Security

- ⚠️ Never commit tokens to version control
- Use secrets management in production (Vault, AWS Secrets Manager)
- Set restrictive permissions: `chmod 600 config.yaml`
- Run as non-root user
- Rotate tokens regularly

## Troubleshooting

**No metrics appearing:**
```bash
# Test token validity
curl -H "Authorization: token ghp_YOUR_TOKEN" \
  https://api.github.com/rate_limit
```

**Config errors:**
```bash
# Validate config syntax
./github_rate_limit_exporter -config config.yaml
```

**Connection issues:**
```bash
# Check if running
curl http://localhost:9101/health

# Check logs
journalctl -u github_rate_limit_exporter -f  # systemd
docker logs github_rate_limit_exporter  # docker
```

## Contributing

Contributions welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT License - see [LICENSE](LICENSE) file.

## Resources

- [Changelog](CHANGELOG.md) - Version history
- [Quick Reference](QUICK_REFERENCE.md) - One-page command reference
- [Alert Examples](alerts.yml) - Prometheus alerting rules
- [GitHub Releases](https://github.com/l13t/github_rate_limit_exporter/releases) - Download binaries