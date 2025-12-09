# GitHub Rate Limit Exporter Helm Chart

A Helm chart for deploying the GitHub Rate Limit Exporter to Kubernetes with Prometheus Operator integration.

## Features

-  Easy deployment to Kubernetes
-  Automatic ServiceMonitor creation for Prometheus Operator
-  Pre-configured PrometheusRule with alerts
-  Secure by default (non-root, read-only filesystem)
-  High availability support with PodDisruptionBudget
-  Network policies for enhanced security
-  Auto-scaling support with HPA

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator) (for ServiceMonitor and PrometheusRule)

## Installing the Chart

### Quick Start

```bash
# Add the repository (if published to a Helm repo)
helm repo add github-rate-limit-exporter https://l13t.github.io/github_rate_limit_exporter

# Or install from local chart
cd helm/github-rate-limit-exporter

# Create a values file with your GitHub tokens
cat > my-values.yaml <<EOF
githubTokens:
  - name: "ci-bot"
    token: "ghp_your_token_here"
  - name: "personal"
    token: "ghp_another_token"
EOF

# Install the chart
helm install github-rate-limit-exporter . \
  --namespace monitoring \
  --create-namespace \
  --values my-values.yaml
```

### Production Installation

For production, use external secret management:

```bash
# 1. Create secret separately (using sealed-secrets, external-secrets, or vault)
kubectl create secret generic github-tokens \
  --from-file=config.yaml=config.yaml \
  --namespace monitoring \
  --dry-run=client -o yaml | kubeseal -o yaml > sealed-secret.yaml

kubectl apply -f sealed-secret.yaml

# 2. Install chart referencing the secret
helm install github-rate-limit-exporter . \
  --namespace monitoring \
  --set config.existingSecret=github-tokens \
  --values values-production.yaml
```

## Configuration

### Basic Configuration

The following table lists the main configurable parameters:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `1` |
| `image.repository` | Image repository | `ghcr.io/l13t/github_rate_limit_exporter` |
| `image.tag` | Image tag | Chart appVersion |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `githubTokens` | List of GitHub tokens to monitor | `[]` |
| `exporter.port` | Exporter port | `9101` |
| `exporter.metricsPath` | Metrics endpoint path | `/metrics` |
| `exporter.pollInterval` | Polling interval in seconds | `60` |

### GitHub Tokens Configuration

#### Option 1: Inline Configuration (Development)

```yaml
githubTokens:
  - name: "user1"
    token: "ghp_token1"
  - name: "user2"
    token: "ghp_token2"
```

#### Option 2: Full Inline Config

```yaml
config:
  inline:
    users:
      - name: "user1"
        token: "ghp_token1"
    listen_addr: ":9101"
    metrics_path: "/metrics"
    poll_interval: 60
```

#### Option 3: Existing Secret (Production)

```yaml
config:
  existingSecret: "github-tokens"
  existingSecretKey: "config.yaml"
```

### Prometheus Operator Integration

#### ServiceMonitor

Automatically creates a ServiceMonitor for Prometheus Operator:

```yaml
serviceMonitor:
  enabled: true
  additionalLabels:
    prometheus: kube-prometheus
  interval: 60s
  scrapeTimeout: 30s
```

#### PrometheusRule (Alerts)

Pre-configured alerts for monitoring:

```yaml
prometheusRule:
  enabled: true
  additionalLabels:
    prometheus: kube-prometheus
  rules:
    # Default rules included, customize as needed
    - alert: GitHubRateLimitCritical
      expr: github_rate_limit_core_remaining < 100
      for: 1m
      labels:
        severity: critical
```

**Default Alerts:**
- `GitHubRateLimitLow` - Warning when <1000 requests remaining
- `GitHubRateLimitCritical` - Critical when <100 requests remaining
- `GitHubRateLimitHighUsage` - Warning when >80% usage
- `GitHubSearchRateLimitLow` - Warning for search API limits
- `GitHubRateLimitExporterDown` - Exporter availability

### Security Configuration

#### Pod Security Context

```yaml
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000
  runAsGroup: 1000
  fsGroup: 1000
  seccompProfile:
    type: RuntimeDefault

securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop:
    - ALL
  readOnlyRootFilesystem: true
```

#### Network Policy

```yaml
networkPolicy:
  enabled: true
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
      - namespaceSelector:
          matchLabels:
            name: monitoring
```

### High Availability

#### Pod Disruption Budget

```yaml
podDisruptionBudget:
  enabled: true
  minAvailable: 1
```

#### Horizontal Pod Autoscaler

```yaml
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 5
  targetCPUUtilizationPercentage: 80
```

#### Anti-Affinity

```yaml
affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchLabels:
            app.kubernetes.io/name: github-rate-limit-exporter
        topologyKey: topology.kubernetes.io/zone
```

## Examples

### Example 1: Simple Development Setup

```yaml
# dev-values.yaml
replicaCount: 1

githubTokens:
  - name: "dev-bot"
    token: "ghp_dev_token"

serviceMonitor:
  enabled: true
  
prometheusRule:
  enabled: true
```

```bash
helm install github-exporter . -f dev-values.yaml -n monitoring
```

### Example 2: Production with External Secrets

```yaml
# prod-values.yaml
replicaCount: 2

config:
  existingSecret: "github-tokens-sealed"

resources:
  limits:
    cpu: 200m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi

podDisruptionBudget:
  enabled: true
  minAvailable: 1

networkPolicy:
  enabled: true

serviceMonitor:
  enabled: true
  additionalLabels:
    prometheus: kube-prometheus
    release: prometheus-operator

prometheusRule:
  enabled: true
  additionalLabels:
    prometheus: kube-prometheus
```

```bash
helm install github-exporter . -f prod-values.yaml -n monitoring
```

### Example 3: Multi-Team Setup

```yaml
# multi-team-values.yaml
githubTokens:
  - name: "team-platform"
    token: "ghp_platform_token"
  - name: "team-backend"
    token: "ghp_backend_token"
  - name: "team-frontend"
    token: "ghp_frontend_token"
  - name: "ci-shared"
    token: "ghp_ci_token"

exporter:
  pollInterval: 30  # More frequent for busy teams

replicaCount: 2

resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 300m
    memory: 512Mi
```

## Upgrading

### Upgrade the Chart

```bash
# Update repository
helm repo update

# Upgrade release
helm upgrade github-rate-limit-exporter . \
  --namespace monitoring \
  --values my-values.yaml
```

### Check What Will Change

```bash
helm diff upgrade github-rate-limit-exporter . \
  --namespace monitoring \
  --values my-values.yaml
```

## Uninstalling

```bash
helm uninstall github-rate-limit-exporter --namespace monitoring
```

To delete the namespace as well:

```bash
kubectl delete namespace monitoring
```

## Verifying the Installation

### Check Deployment Status

```bash
# Check pods
kubectl get pods -n monitoring -l app.kubernetes.io/name=github-rate-limit-exporter

# Check logs
kubectl logs -n monitoring -l app.kubernetes.io/name=github-rate-limit-exporter

# Check metrics endpoint
kubectl port-forward -n monitoring svc/github-rate-limit-exporter 9101:9101
curl http://localhost:9101/metrics
```

### Verify ServiceMonitor

```bash
# Check if ServiceMonitor is created
kubectl get servicemonitor -n monitoring

# Verify Prometheus is scraping
kubectl port-forward -n monitoring svc/prometheus-operated 9090:9090
# Open http://localhost:9090/targets
```

### Verify Alerts

```bash
# Check PrometheusRule
kubectl get prometheusrule -n monitoring

# Check alerts in Prometheus
# Open http://localhost:9090/alerts
```

## Troubleshooting

### Pods Not Starting

```bash
# Check pod status
kubectl describe pod -n monitoring -l app.kubernetes.io/name=github-rate-limit-exporter

# Check logs
kubectl logs -n monitoring -l app.kubernetes.io/name=github-rate-limit-exporter
```

### Metrics Not Appearing in Prometheus

1. **Check ServiceMonitor:**
   ```bash
   kubectl get servicemonitor -n monitoring -o yaml
   ```

2. **Verify labels match Prometheus selector:**
   ```bash
   kubectl get prometheus -n monitoring -o yaml | grep serviceMonitorSelector
   ```

3. **Check Prometheus logs:**
   ```bash
   kubectl logs -n monitoring prometheus-prometheus-0
   ```

### Configuration Errors

```bash
# Validate configuration
helm template . --values my-values.yaml --debug

# Test installation
helm install --dry-run --debug github-exporter . -f my-values.yaml
```

### Common Issues

**Issue:** Secret not found
```bash
# Solution: Ensure secret exists before installation
kubectl get secret github-tokens -n monitoring
```

**Issue:** ServiceMonitor not being picked up
```bash
# Solution: Add correct labels
serviceMonitor:
  additionalLabels:
    release: prometheus-operator  # Match your Prometheus Operator release
```

**Issue:** Network policies blocking scraping
```bash
# Solution: Adjust network policy or disable temporarily
networkPolicy:
  enabled: false
```

## Values File Reference

See the following example values files:
- `values.yaml` - Default values with all options documented
- `values-production.yaml` - Production-ready configuration example

## Contributing

Contributions are welcome! Please see the main repository for contribution guidelines.

## License

MIT License - see LICENSE file in the main repository.

## Support

- **Documentation**: https://github.com/l13t/github_rate_limit_exporter
- **Issues**: https://github.com/l13t/github_rate_limit_exporter/issues
- **Discussions**: https://github.com/l13t/github_rate_limit_exporter/discussions

## Related Projects

- [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator)
- [kube-prometheus-stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)
- [External Secrets Operator](https://external-secrets.io/)