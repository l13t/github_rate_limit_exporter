# Helm Installation Guide

Complete guide for deploying GitHub Rate Limit Exporter to Kubernetes using Helm with Prometheus Operator integration.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator) installed (for ServiceMonitor and PrometheusRule)
- GitHub Personal Access Token(s)

## Quick Start

### 1. Create Configuration

Create a values file with your GitHub tokens:

```bash
cat > my-values.yaml <<EOF
githubTokens:
  - name: "ci-bot"
    token: "ghp_your_token_here"
  - name: "team-shared"
    token: "ghp_another_token"

serviceMonitor:
  enabled: true

prometheusRule:
  enabled: true
EOF
```

### 2. Install the Chart

```bash
# Install from local chart directory
helm install github-rate-limit-exporter ./helm/github-rate-limit-exporter \
  --namespace monitoring \
  --create-namespace \
  --values my-values.yaml
```

### 3. Verify Installation

```bash
# Check pods
kubectl get pods -n monitoring -l app.kubernetes.io/name=github-rate-limit-exporter

# Check metrics
kubectl port-forward -n monitoring svc/github-rate-limit-exporter 9101:9101
curl http://localhost:9101/metrics | grep github_rate_limit
```

## Installation Options

### Option 1: Simple Development Setup

Best for development and testing:

```bash
# Create simple values file
cat > dev-values.yaml <<EOF
replicaCount: 1

githubTokens:
  - name: "dev-user"
    token: "ghp_dev_token"

serviceMonitor:
  enabled: true
  
prometheusRule:
  enabled: true

resources:
  requests:
    cpu: 50m
    memory: 64Mi
  limits:
    cpu: 200m
    memory: 128Mi
EOF

# Install
helm install github-exporter ./helm/github-rate-limit-exporter \
  -f dev-values.yaml \
  -n monitoring \
  --create-namespace
```

### Option 2: Production with External Secrets

Best for production - keeps tokens out of Helm values:

**Step 1: Create Secret**

Using Sealed Secrets:
```bash
# Create config file
cat > config.yaml <<EOF
users:
  - name: "prod-bot"
    token: "ghp_production_token"
  - name: "ci-pipeline"
    token: "ghp_ci_token"
listen_addr: ":9101"
metrics_path: "/metrics"
poll_interval: 60
EOF

# Create sealed secret
kubectl create secret generic github-tokens \
  --from-file=config.yaml=config.yaml \
  --namespace monitoring \
  --dry-run=client -o yaml | \
  kubeseal -o yaml > github-tokens-sealed.yaml

# Apply sealed secret
kubectl apply -f github-tokens-sealed.yaml
```

Using External Secrets Operator:
```yaml
# external-secret.yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: github-tokens
  namespace: monitoring
spec:
  secretStoreRef:
    name: vault-backend
    kind: SecretStore
  target:
    name: github-tokens
  data:
  - secretKey: config.yaml
    remoteRef:
      key: github-exporter/config
```

**Step 2: Install Chart**

```bash
# Use production values
helm install github-exporter ./helm/github-rate-limit-exporter \
  -f ./helm/github-rate-limit-exporter/values-production.yaml \
  --set config.existingSecret=github-tokens \
  -n monitoring
```

### Option 3: Multi-Team Setup

For monitoring multiple teams:

```yaml
# multi-team-values.yaml
replicaCount: 2

githubTokens:
  - name: "team-platform"
    token: "ghp_platform_token"
  - name: "team-backend"
    token: "ghp_backend_token"
  - name: "team-frontend"
    token: "ghp_frontend_token"
  - name: "team-mobile"
    token: "ghp_mobile_token"
  - name: "ci-shared"
    token: "ghp_ci_token"

exporter:
  pollInterval: 30  # More frequent polling

resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 300m
    memory: 256Mi

podDisruptionBudget:
  enabled: true
  minAvailable: 1

serviceMonitor:
  enabled: true
  interval: 30s
  additionalLabels:
    prometheus: kube-prometheus
    release: prometheus-operator

prometheusRule:
  enabled: true
  additionalLabels:
    prometheus: kube-prometheus
```

```bash
helm install github-exporter ./helm/github-rate-limit-exporter \
  -f multi-team-values.yaml \
  -n monitoring
```

## Configuration

### GitHub Tokens

Three ways to configure tokens:

**1. Inline (simplest, development only):**
```yaml
githubTokens:
  - name: "user1"
    token: "ghp_token1"
```

**2. Full config inline:**
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

**3. Existing secret (production):**
```yaml
config:
  existingSecret: "github-tokens"
  existingSecretKey: "config.yaml"
```

### Prometheus Operator Integration

#### ServiceMonitor Configuration

```yaml
serviceMonitor:
  enabled: true
  # Labels must match your Prometheus selector
  additionalLabels:
    prometheus: kube-prometheus
    release: prometheus-operator  # Match your Prometheus Operator release name
  namespace: ""  # Empty = same as release namespace
  interval: 60s
  scrapeTimeout: 30s
  honorLabels: true
```

**Find your Prometheus selector:**
```bash
kubectl get prometheus -n monitoring -o yaml | grep -A 5 serviceMonitorSelector
```

#### PrometheusRule Configuration

```yaml
prometheusRule:
  enabled: true
  additionalLabels:
    prometheus: kube-prometheus
    release: prometheus-operator
  namespace: ""  # Empty = same as release namespace
  rules:
    # Default rules are pre-configured
    # Add custom rules here
    - alert: CustomAlert
      expr: custom_expression
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Custom alert"
```

**Default alerts included:**
- `GitHubRateLimitCritical` - <100 requests remaining
- `GitHubRateLimitLow` - <1000 requests remaining
- `GitHubRateLimitHighUsage` - >80% usage
- `GitHubSearchRateLimitLow` - Search API limits
- `GitHubRateLimitExporterDown` - Exporter unavailable

### Security Settings

#### Pod Security

```yaml
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000
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

#### Network Policies

```yaml
networkPolicy:
  enabled: true
  policyTypes:
    - Ingress
    - Egress
  ingress:
    # Allow Prometheus to scrape
    - from:
      - namespaceSelector:
          matchLabels:
            name: monitoring
        podSelector:
          matchLabels:
            app.kubernetes.io/name: prometheus
  egress:
    # Allow DNS
    - to:
      - namespaceSelector:
          matchLabels:
            name: kube-system
      ports:
      - protocol: UDP
        port: 53
    # Allow GitHub API
    - to:
      - podSelector: {}
      ports:
      - protocol: TCP
        port: 443
```

### High Availability

```yaml
# Multiple replicas
replicaCount: 2

# Pod Disruption Budget
podDisruptionBudget:
  enabled: true
  minAvailable: 1

# Anti-affinity for spreading across nodes
affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchLabels:
            app.kubernetes.io/name: github-rate-limit-exporter
        topologyKey: topology.kubernetes.io/zone

# Auto-scaling (optional)
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 5
  targetCPUUtilizationPercentage: 80
```

## Verification

### Check Deployment

```bash
# Pod status
kubectl get pods -n monitoring -l app.kubernetes.io/name=github-rate-limit-exporter

# View logs
kubectl logs -n monitoring -l app.kubernetes.io/name=github-rate-limit-exporter

# Describe deployment
kubectl describe deployment -n monitoring github-rate-limit-exporter
```

### Test Metrics Endpoint

```bash
# Port forward
kubectl port-forward -n monitoring svc/github-rate-limit-exporter 9101:9101

# In another terminal
curl http://localhost:9101/metrics | grep github_rate_limit
curl http://localhost:9101/health
```

### Verify ServiceMonitor

```bash
# Check if ServiceMonitor exists
kubectl get servicemonitor -n monitoring

# Check ServiceMonitor details
kubectl get servicemonitor -n monitoring github-rate-limit-exporter -o yaml

# Check if Prometheus is scraping (requires Prometheus port-forward)
kubectl port-forward -n monitoring svc/prometheus-operated 9090:9090
# Open http://localhost:9090/targets and look for github-rate-limit-exporter
```

### Verify PrometheusRule

```bash
# Check if PrometheusRule exists
kubectl get prometheusrule -n monitoring

# View configured alerts
kubectl get prometheusrule -n monitoring github-rate-limit-exporter -o yaml

# Check alerts in Prometheus
# Open http://localhost:9090/alerts
```

## Upgrading

### Standard Upgrade

```bash
# Update values if needed
vim my-values.yaml

# Upgrade release
helm upgrade github-rate-limit-exporter ./helm/github-rate-limit-exporter \
  --namespace monitoring \
  --values my-values.yaml
```

### Preview Changes

```bash
# See what will change
helm diff upgrade github-rate-limit-exporter ./helm/github-rate-limit-exporter \
  --namespace monitoring \
  --values my-values.yaml

# Or use Helm's built-in dry-run
helm upgrade github-rate-limit-exporter ./helm/github-rate-limit-exporter \
  --namespace monitoring \
  --values my-values.yaml \
  --dry-run --debug
```

### Upgrade to New Chart Version

```bash
# Pull new chart version
git pull

# Upgrade with new version
helm upgrade github-rate-limit-exporter ./helm/github-rate-limit-exporter \
  --namespace monitoring \
  --reuse-values
```

## Uninstalling

### Remove the Release

```bash
# Uninstall Helm release
helm uninstall github-rate-limit-exporter --namespace monitoring

# Optionally delete the namespace
kubectl delete namespace monitoring
```

### Clean Up Secrets

```bash
# If you created secrets manually
kubectl delete secret github-tokens -n monitoring
```

## Troubleshooting

### Pods Not Starting

**Check pod status:**
```bash
kubectl describe pod -n monitoring -l app.kubernetes.io/name=github-rate-limit-exporter
```

**Common issues:**
- Secret not found → Create secret first
- Image pull errors → Check image name and pull policy
- Resource constraints → Adjust resource limits

### Metrics Not in Prometheus

**1. Check ServiceMonitor labels:**
```bash
# Get Prometheus selector
kubectl get prometheus -n monitoring -o yaml | grep -A 5 serviceMonitorSelector

# Compare with ServiceMonitor labels
kubectl get servicemonitor -n monitoring github-rate-limit-exporter -o yaml | grep -A 5 labels
```

**2. Check Prometheus logs:**
```bash
kubectl logs -n monitoring prometheus-prometheus-0 | grep github-rate-limit
```

**3. Verify service selector:**
```bash
kubectl get svc -n monitoring github-rate-limit-exporter -o yaml
kubectl get pods -n monitoring -l app.kubernetes.io/name=github-rate-limit-exporter --show-labels
```

### Alerts Not Firing

**1. Check PrometheusRule:**
```bash
kubectl get prometheusrule -n monitoring github-rate-limit-exporter -o yaml
```

**2. Check Prometheus rules:**
```bash
# Port forward Prometheus
kubectl port-forward -n monitoring svc/prometheus-operated 9090:9090
# Visit http://localhost:9090/rules
```

**3. Test alert expression:**
```bash
# In Prometheus UI, test the expression
github_rate_limit_core_remaining < 100
```

### Configuration Errors

**Validate template:**
```bash
helm template ./helm/github-rate-limit-exporter \
  --values my-values.yaml \
  --debug
```

**Check secret content:**
```bash
kubectl get secret -n monitoring github-rate-limit-exporter -o yaml
# Decode config
kubectl get secret -n monitoring github-rate-limit-exporter -o jsonpath='{.data.config\.yaml}' | base64 -d
```

## Best Practices

### 1. Secret Management

 **Do:**
- Use external secret management (Sealed Secrets, External Secrets, Vault)
- Rotate tokens regularly
- Use least-privilege tokens (read-only)

 **Don't:**
- Commit tokens to git
- Use inline token configuration in production
- Share tokens across environments

### 2. Resource Management

```yaml
# Set appropriate resource limits
resources:
  requests:
    cpu: 50m      # Start small
    memory: 64Mi
  limits:
    cpu: 200m     # Prevent runaway
    memory: 128Mi # Prevent OOM
```

### 3. High Availability

```yaml
# For production
replicaCount: 2
podDisruptionBudget:
  enabled: true
  minAvailable: 1
```

### 4. Monitoring

```yaml
# Enable ServiceMonitor and alerts
serviceMonitor:
  enabled: true
prometheusRule:
  enabled: true
```

### 5. Security

```yaml
# Enable network policies
networkPolicy:
  enabled: true

# Run as non-root
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000
```

## Examples

### Complete Production Example

```yaml
# production.yaml
replicaCount: 2

image:
  tag: "1.0.0"  # Pin version

# Use external secret
config:
  existingSecret: "github-tokens-vault"
  existingSecretKey: "config.yaml"

exporter:
  pollInterval: 60
  logLevel: info

resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 200m
    memory: 256Mi

podDisruptionBudget:
  enabled: true
  minAvailable: 1

affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchLabels:
            app.kubernetes.io/name: github-rate-limit-exporter
        topologyKey: topology.kubernetes.io/zone

networkPolicy:
  enabled: true

serviceMonitor:
  enabled: true
  additionalLabels:
    prometheus: kube-prometheus
    release: prometheus-operator
  interval: 60s

prometheusRule:
  enabled: true
  additionalLabels:
    prometheus: kube-prometheus

priorityClassName: "system-cluster-critical"
```

Installation:
```bash
helm install github-exporter ./helm/github-rate-limit-exporter \
  -f production.yaml \
  -n monitoring \
  --create-namespace
```

## Additional Resources

- [Helm Chart README](helm/github-rate-limit-exporter/README.md)
- [Values Reference](helm/github-rate-limit-exporter/values.yaml)
- [Production Values Example](helm/github-rate-limit-exporter/values-production.yaml)
- [Prometheus Operator Docs](https://prometheus-operator.dev/)
- [External Secrets Operator](https://external-secrets.io/)

## Support

- **Issues**: https://github.com/l13t/github_rate_limit_exporter/issues
- **Documentation**: https://github.com/l13t/github_rate_limit_exporter
- **Discussions**: https://github.com/l13t/github_rate_limit_exporter/discussions