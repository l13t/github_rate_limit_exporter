# Quick Reference

## 🚀 Quick Start

```bash
# 1. Get GitHub token from: https://github.com/settings/tokens
# 2. Download binary or use Docker
# 3. Create config.yaml:
cat > config.yaml <<EOF
users:
  - name: "my-user"
    token: "ghp_your_token_here"
listen_addr: ":9101"
EOF

# 4. Run
./github_rate_limit_exporter -config config.yaml

# 5. Test
curl http://localhost:9101/metrics | grep github_rate_limit
```

## 📋 Common Commands

### Using Task

```bash
task build          # Build binary
task test           # Run tests
task run            # Build and run
task docker:build   # Build Docker image
task --list         # Show all tasks
```

### Using Go Directly

```bash
go build -o github_rate_limit_exporter ./cmd/exporter
go test ./...
go run ./cmd/exporter -config config.yaml
```

## 🐳 Docker

```bash
# Run
docker run -d -p 9101:9101 \
  -v $(pwd)/config.yaml:/config.yaml:ro \
  ghcr.io/l13t/github_rate_limit_exporter:latest \
  -config /config.yaml

# With Docker Compose
docker-compose up -d
docker-compose logs -f
docker-compose down
```

## 🔧 Configuration

```yaml
users:
  - name: "user1"
    token: "ghp_token1"
  - name: "user2"
    token: "ghp_token2"

listen_addr: ":9101"      # Server address
metrics_path: "/metrics"  # Metrics endpoint
poll_interval: 60         # Seconds between polls
```

## 📊 Key Metrics

```promql
# Remaining requests
github_rate_limit_core_remaining{user="username"}

# Usage percentage
(github_rate_limit_core_used / github_rate_limit_core_limit) * 100

# Time until reset (minutes)
(github_rate_limit_core_reset_timestamp - time()) / 60
```

## 🔗 Default URLs

- **Metrics**: http://localhost:9101/metrics
- **Health**: http://localhost:9101/health
- **Prometheus**: http://localhost:9090 (if using docker-compose)
- **Grafana**: http://localhost:3000 (if using docker-compose)

## 🆘 Troubleshooting

```bash
# Test token
curl -H "Authorization: token ghp_YOUR_TOKEN" \
  https://api.github.com/rate_limit

# Check health
curl http://localhost:9101/health

# View logs (systemd)
journalctl -u github_rate_limit_exporter -f

# View logs (docker)
docker logs -f github_rate_limit_exporter
```

## 🔒 Security Checklist

- [ ] Never commit tokens to git
- [ ] Use `chmod 600 config.yaml`
- [ ] Run as non-root user
- [ ] Rotate tokens regularly
- [ ] Use secrets management in production

## 📚 Full Documentation

- [README.md](README.md) - Complete guide
- [CONTRIBUTING.md](CONTRIBUTING.md) - How to contribute
- [CHANGELOG.md](CHANGELOG.md) - Version history