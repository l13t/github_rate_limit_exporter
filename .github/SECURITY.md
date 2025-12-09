# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of GitHub Rate Limit Exporter seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Please Do NOT

- Open a public GitHub issue
- Discuss the vulnerability in public forums, chat rooms, or social media
- Exploit the vulnerability beyond what is necessary to demonstrate it

### Please DO

**Report security vulnerabilities via GitHub Security Advisories:**

1. Go to https://github.com/l13t/github_rate_limit_exporter/security/advisories
2. Click "New draft security advisory"
3. Fill in the details of the vulnerability
4. Submit the advisory

**Or report via email:**

Send an email to: security@example.com (replace with your actual security contact)

### What to Include

Please include the following information in your report:

- Type of vulnerability (e.g., authentication bypass, code injection, etc.)
- Full paths of source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the vulnerability
- Proof-of-concept or exploit code (if possible)
- Impact of the vulnerability, including how an attacker might exploit it

### Response Timeline

- **Initial Response**: Within 48 hours
- **Severity Assessment**: Within 5 business days
- **Fix Development**: Depends on complexity and severity
- **Public Disclosure**: After fix is released and users have had time to update

## Security Update Process

1. **Assessment**: We'll acknowledge receipt and assess the severity
2. **Fix Development**: We'll develop a fix in a private repository
3. **Testing**: The fix will be thoroughly tested
4. **Release**: We'll release a new version with the fix
5. **Advisory**: We'll publish a security advisory
6. **Credit**: We'll credit you in the advisory (if desired)

## Security Best Practices for Users

### Configuration Security

- **Never commit tokens to version control**
  ```bash
  # Add to .gitignore
  config.yaml
  config.toml
  config.hcl
  ```

- **Use restrictive file permissions**
  ```bash
  chmod 600 config.yaml
  chown exporter:exporter config.yaml
  ```

- **Use secrets management in production**
  - Kubernetes Secrets
  - HashiCorp Vault
  - AWS Secrets Manager
  - Azure Key Vault
  - Google Secret Manager

### Deployment Security

#### Docker

- **Run as non-root** (default in our images)
  ```dockerfile
  USER exporter
  ```

- **Use read-only file systems**
  ```bash
  docker run --read-only -v $(pwd)/config.yaml:/config.yaml:ro ...
  ```

- **Enable security options**
  ```bash
  docker run --security-opt=no-new-privileges:true \
    --cap-drop=ALL ...
  ```

#### Kubernetes

- **Use SecurityContext**
  ```yaml
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    readOnlyRootFilesystem: true
    allowPrivilegeEscalation: false
    capabilities:
      drop:
        - ALL
  ```

- **Use Secrets for tokens**
  ```yaml
  env:
    - name: GITHUB_TOKEN
      valueFrom:
        secretKeyRef:
          name: github-tokens
          key: token
  ```

- **Enable Network Policies**
  ```yaml
  apiVersion: networking.k8s.io/v1
  kind: NetworkPolicy
  metadata:
    name: github-rate-limit-exporter
  spec:
    podSelector:
      matchLabels:
        app: github-rate-limit-exporter
    policyTypes:
      - Ingress
      - Egress
    ingress:
      - from:
        - podSelector:
            matchLabels:
              app: prometheus
        ports:
          - protocol: TCP
            port: 9101
    egress:
      - to:
        - namespaceSelector: {}
        ports:
          - protocol: TCP
            port: 443  # GitHub API
  ```

#### Systemd

- **Use systemd hardening**
  ```ini
  [Service]
  NoNewPrivileges=true
  PrivateTmp=true
  ProtectSystem=strict
  ProtectHome=true
  ProtectKernelTunables=true
  ProtectKernelModules=true
  ProtectControlGroups=true
  RestrictRealtime=true
  RestrictNamespaces=true
  RestrictAddressFamilies=AF_INET AF_INET6
  LockPersonality=true
  MemoryDenyWriteExecute=true
  SystemCallFilter=@system-service
  ```

### GitHub Token Security

#### Minimal Permissions

- **Classic tokens**: No scopes needed for rate limit reading
- **Fine-grained tokens**: Read-only access is sufficient

#### Token Management

- **Rotate regularly** (every 90 days minimum)
- **Use unique tokens** for different environments
- **Revoke immediately** if compromised
- **Monitor token usage** via GitHub audit log

### Network Security

- **Use HTTPS** for external access
  ```nginx
  server {
    listen 443 ssl;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location /metrics {
      proxy_pass http://localhost:9101;
    }
  }
  ```

- **Restrict access** with firewall rules
  ```bash
  # UFW example
  sudo ufw allow from 10.0.0.0/8 to any port 9101
  ```

- **Use mutual TLS** for Prometheus scraping
  ```yaml
  # prometheus.yml
  scrape_configs:
    - job_name: 'github_rate_limits'
      scheme: https
      tls_config:
        ca_file: /path/to/ca.pem
        cert_file: /path/to/client-cert.pem
        key_file: /path/to/client-key.pem
  ```

## Known Security Considerations

### Rate Limit Information

The exporter exposes rate limit information which could be considered sensitive in some environments:
- Number of API calls made
- Remaining API quota
- Reset times

**Mitigation**: Restrict access to metrics endpoint to authorized systems only.

### Configuration Files

Configuration files contain GitHub tokens which provide API access:
- Could be used to make API calls as the token owner
- Tokens are limited by their configured permissions

**Mitigation**: 
- Use read-only or minimal permission tokens
- Protect configuration files with appropriate permissions
- Use secrets management systems

## Security Scanning

### Automated Scanning

We use the following tools in our CI/CD:

- **Trivy**: Container and dependency vulnerability scanning
- **Gosec**: Go security checker for code issues
- **Grype**: Additional vulnerability scanning
- **Dependabot**: Automated dependency updates
- **CodeQL**: Code analysis for security issues

### Manual Audits

We perform manual security audits:
- Before major releases
- After significant changes
- When new attack vectors are discovered

## Compliance

### SBOM (Software Bill of Materials)

We provide SBOM for:
- Every release (attached to GitHub releases)
- Docker images (via attestation)
- Format: SPDX JSON

### Supply Chain Security

- **Signed commits**: All releases are from verified commits
- **Attestation**: Docker images include provenance attestation
- **Reproducible builds**: Builds are reproducible from source
- **Pinned dependencies**: All actions and dependencies are version-pinned

## Security Headers

When exposing the exporter via reverse proxy, add security headers:

```nginx
add_header X-Content-Type-Options "nosniff" always;
add_header X-Frame-Options "DENY" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "no-referrer" always;
add_header Content-Security-Policy "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self';" always;
```

## Incident Response

If a security incident occurs:

1. **Containment**: Immediately revoke affected tokens
2. **Investigation**: Determine scope and impact
3. **Remediation**: Apply fixes and patches
4. **Communication**: Notify affected users
5. **Documentation**: Document incident and lessons learned

## Security Contacts

- **Security Issues**: GitHub Security Advisories (preferred)
- **Email**: security@example.com
- **GPG Key**: Available on request

## Acknowledgments

We appreciate the security research community and thank those who responsibly disclose vulnerabilities to us.

### Hall of Fame

Contributors who have helped improve security:
- (Contributors will be listed here with their permission)

## Additional Resources

- [GitHub Security Best Practices](https://docs.github.com/en/code-security)
- [Docker Security](https://docs.docker.com/engine/security/)
- [Kubernetes Security](https://kubernetes.io/docs/concepts/security/)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)

---

**Last Updated**: 2024-12-09  
**Version**: 1.0  

This security policy is subject to change. Please check back regularly for updates.