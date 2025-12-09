# GitHub Actions Workflows

This document describes the GitHub Actions workflows used in this project for CI/CD, building, testing, and releasing.

## Overview

We use GitHub Actions for:
- 🧪 **Continuous Integration** - Testing and validation on every push/PR
- 🐳 **Docker Builds** - Multi-architecture container images
- 📦 **Releases** - Automated binary and Docker image releases
- 🔒 **Security** - Vulnerability scanning and security checks
- 📚 **Dependencies** - Automated dependency updates

## Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Trigger**: Push to master/develop, Pull Requests

**Purpose**: Comprehensive testing and validation

#### Jobs

##### Test
- Runs unit tests with race detection
- Generates code coverage reports
- Uploads coverage to Codecov
- **Runs on**: Ubuntu Latest
- **Go Version**: 1.21

##### Lint
- Runs golangci-lint for code quality
- Checks for common issues and style violations
- **Timeout**: 5 minutes

##### Security
- **Gosec**: Security scanner for Go code
- **Trivy**: Vulnerability scanner for dependencies
- Uploads results to GitHub Security tab
- **Format**: SARIF (for GitHub integration)

##### Build
- Builds binaries for all supported platforms
- **Platforms**:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
- Tests cross-compilation works
- Uploads artifacts (7-day retention)

##### Docker Build Test
- Tests Docker image builds successfully
- Uses BuildKit cache
- Does not push images
- **Purpose**: Ensure Dockerfile is valid

##### Validate
- YAML linting for workflow files
- Dockerfile linting with Hadolint
- Validates go.mod tidiness
- **Checks**:
  - YAML syntax and style
  - Dockerfile best practices
  - Dependency consistency

#### Artifacts

- **Build Artifacts**: Platform-specific binaries (7 days)
- **Coverage Reports**: Uploaded to Codecov
- **Security Scans**: SARIF files in GitHub Security

#### When Does It Run?

```yaml
on:
  push:
    branches: [ master, develop ]
  pull_request:
    branches: [ master, develop ]
  workflow_dispatch:  # Manual trigger
```

---

### 2. Docker Workflow (`.github/workflows/docker.yml`)

**Trigger**: Push to master/develop (specific paths), PRs, Weekly schedule, Manual

**Purpose**: Build, test, and publish Docker images

#### Jobs

##### Build and Test
- Builds images for multiple platforms
- **Platforms**:
  - linux/amd64
  - linux/arm64
  - linux/arm/v7
- Uses QEMU for cross-platform builds
- Uses layer caching for faster builds
- **Strategy**: Matrix build for each platform

##### Test Image
- Tests container functionality
- Verifies non-root user execution
- Tests health endpoint
- Tests metrics endpoint
- Validates container structure
- **Tests**:
  - Container starts successfully
  - Runs as non-root user
  - Health check responds
  - Metrics endpoint works

##### Security Scan
- **Trivy**: Scans for vulnerabilities
- **Grype**: Alternative vulnerability scanner
- Uploads findings to GitHub Security
- **Severity**: CRITICAL and HIGH
- **Format**: SARIF + Table

##### Push Images
- Pushes to GitHub Container Registry (ghcr.io)
- Optionally pushes to Docker Hub
- Only runs on master branch (not PRs)
- Generates SBOM (Software Bill of Materials)
- Creates attestations for supply chain security

##### Cleanup
- Deletes old untagged images
- Keeps last 10 versions
- Runs regardless of job status
- Prevents registry bloat

#### Image Tags

When pushing to master:
- `edge` - Latest master branch build
- `latest` - Alias for latest stable
- `master` - Branch-based tag
- `master-<sha>` - Commit-specific tag

#### Registries

1. **GitHub Container Registry** (always)
   - `ghcr.io/l13t/github_rate_limit_exporter`
   
2. **Docker Hub** (if configured)
   - Requires `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN` secrets
   - `docker.io/<username>/github_rate_limit_exporter`

#### When Does It Run?

```yaml
on:
  push:
    branches: [ master, develop ]
    paths: [ 'Dockerfile', '**.go', 'go.mod', 'go.sum' ]
  pull_request:
    branches: [ master, develop ]
  schedule:
    - cron: '0 0 * * 0'  # Weekly on Sunday
  workflow_dispatch:
```

---

### 3. Release Workflow (`.github/workflows/release.yml`)

**Trigger**: Git tags (v*.*.*), Manual workflow dispatch

**Purpose**: Create production releases with binaries and Docker images

#### Jobs

##### Build Binaries
- Builds for **9 platform combinations**:
  - Linux: amd64, arm64, armv7
  - macOS: amd64, arm64
  - Windows: amd64, arm64
  - FreeBSD: amd64, arm64
- Creates compressed archives:
  - `.tar.gz` for Unix-like systems
  - `.zip` for Windows
- Generates SHA256 checksums
- **Static binaries**: CGO_ENABLED=0
- **Optimized**: `-s -w` ldflags

##### Build Docker
- Builds multi-arch Docker images
- Pushes to registries with version tags
- **Platforms**: linux/amd64, linux/arm64, linux/arm/v7
- Generates SBOM for security
- **Tags**:
  - `v1.2.3` - Exact version
  - `v1.2` - Minor version
  - `v1` - Major version
  - `latest` - Latest stable

##### Create Release
- Downloads all build artifacts
- Generates release notes from git history
- Creates GitHub Release
- Uploads all binaries and checksums
- Marks pre-releases (tags with `-alpha`, `-beta`, etc.)
- **Auto-generated**: Includes commit list since last tag

##### Update Homebrew
- Updates Homebrew tap (if configured)
- Only for stable releases (no pre-releases)
- Requires `HOMEBREW_TAP_TOKEN` secret
- **Automatic**: Formula version and checksum

##### Notify
- Sends Slack notification (if configured)
- Requires `SLACK_WEBHOOK_URL` secret
- Runs regardless of success/failure
- **Info**: Release version and Docker image

#### Release Versioning

```bash
# Create a release
git tag v1.0.0
git push origin v1.0.0

# Create a pre-release
git tag v1.0.0-beta.1
git push origin v1.0.0-beta.1
```

#### Release Assets

Each release includes:
- **Binaries**: 9 platform-specific archives
- **Checksums**: SHA256 for verification
- **Release Notes**: Auto-generated changelog
- **Docker Images**: Multi-arch containers
- **SBOM**: Software Bill of Materials

#### Manual Release

Trigger manually from Actions tab:
```yaml
workflow_dispatch:
  inputs:
    tag:
      description: 'Tag to release'
      required: true
      type: string
```

---

## Configuration

### Required Secrets

None! The workflows use GitHub's built-in `GITHUB_TOKEN`.

### Optional Secrets

Configure these in repository Settings → Secrets and variables → Actions:

#### Docker Hub (Optional)
```
DOCKERHUB_USERNAME=your-username
DOCKERHUB_TOKEN=your-access-token
```

#### Homebrew (Optional)
```
HOMEBREW_TAP_TOKEN=github-personal-access-token
```

#### Slack Notifications (Optional)
```
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...
```

#### Codecov (Optional)
```
CODECOV_TOKEN=your-codecov-token
```

---

## Dependabot

Automated dependency updates via `.github/dependabot.yml`:

### Go Modules
- **Schedule**: Weekly (Monday 09:00 UTC)
- **Groups**:
  - Prometheus dependencies
  - GitHub API dependencies
  - Config parser dependencies

### GitHub Actions
- **Schedule**: Weekly (Monday 09:00 UTC)
- **Groups**:
  - Core actions
  - Docker actions
  - Security actions

### Docker Base Images
- **Schedule**: Weekly (Monday 09:00 UTC)
- **Updates**: Go and Alpine base images

---

## Permissions

### CI Workflow
```yaml
permissions:
  contents: read
  pull-requests: read
```

### Release Workflow
```yaml
permissions:
  contents: write       # Create releases
  packages: write       # Push to GHCR
  id-token: write       # SBOM attestation
```

### Docker Workflow
```yaml
permissions:
  contents: read
  packages: write       # Push to GHCR
  security-events: write # Upload security scans
```

---

## Workflow Optimization

### Caching

**Go Modules**:
```yaml
- uses: actions/setup-go@v5
  with:
    cache: true  # Caches go mod and build cache
```

**Docker Layers**:
```yaml
cache-from: type=gha
cache-to: type=gha,mode=max
```

### Matrix Builds

Parallel builds for multiple platforms:
```yaml
strategy:
  matrix:
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
```

### Conditional Execution

Skip unnecessary jobs:
```yaml
if: github.event_name != 'pull_request'
if: secrets.DOCKERHUB_USERNAME != ''
if: always()  # Run even if previous jobs fail
```

---

## Security Features

### Supply Chain Security

1. **SBOM Generation**: Software Bill of Materials for Docker images
2. **Attestation**: Cryptographic proof of image provenance
3. **Vulnerability Scanning**: Trivy and Grype scanners
4. **Security Tab Integration**: SARIF uploads for centralized view

### Static Analysis

- **Gosec**: Go security checker
- **golangci-lint**: Multiple linters including security checks
- **Hadolint**: Dockerfile best practices

### Dependency Security

- **Dependabot**: Automated security updates
- **govulncheck**: Go vulnerability database checks
- **Trivy**: CVE scanning for dependencies

---

## Common Tasks

### Running Workflows Locally

Use [act](https://github.com/nektos/act) to test workflows locally:

```bash
# Install act
brew install act

# Run CI workflow
act push

# Run specific job
act -j test

# List available workflows
act -l
```

### Debugging Workflows

Enable debug logging:

1. Go to repository Settings → Secrets
2. Add secret: `ACTIONS_STEP_DEBUG = true`
3. Re-run workflow

### Canceling Redundant Runs

Workflows automatically cancel:
- Previous runs on the same branch (for pushes)
- Concurrent runs (configurable)

### Manual Workflow Triggers

All workflows support `workflow_dispatch` for manual triggering:

1. Go to Actions tab
2. Select workflow
3. Click "Run workflow"
4. Select branch and input parameters

---

## Monitoring

### Workflow Status

Check workflow status:
- **Badge**: Add to README
- **GitHub Status**: Repository Actions tab
- **Notifications**: Configure in personal settings

### Success/Failure Rates

View metrics:
1. Repository → Insights
2. Community → Workflow runs

### Artifact Storage

- **Build artifacts**: 7 days retention
- **SBOM**: 30-90 days retention
- **Docker images**: Unlimited (with cleanup)

---

## Best Practices

### ✅ DO

- Test workflows on feature branches before merging
- Use specific action versions (not `@main`)
- Cache dependencies for faster builds
- Use matrix builds for multiple platforms
- Upload artifacts for debugging
- Add meaningful workflow and job names
- Use secrets for sensitive data
- Enable security scanning

### ❌ DON'T

- Commit secrets or tokens
- Use `latest` tags for actions
- Skip security scans
- Ignore workflow failures
- Hard-code credentials
- Run untrusted code
- Skip version pinning

---

## Troubleshooting

### Workflow Fails on Go Build

**Issue**: Build fails with "module not found"

**Solution**:
```bash
task tidy
git commit -m "chore: tidy go.mod"
```

### Docker Build Timeout

**Issue**: Multi-arch build takes too long

**Solution**: Builds are parallelized by platform. Check:
- Network connectivity
- Docker Hub rate limits
- Cache availability

### Release Artifacts Missing

**Issue**: GitHub Release doesn't have all binaries

**Solution**: 
- Check individual job logs
- Verify matrix includes all platforms
- Check artifact upload step

### Security Scan False Positives

**Issue**: Trivy reports known false positives

**Solution**: Add to `.trivyignore`:
```
# False positive in dependency X
CVE-2023-12345
```

---

## Metrics

### Current Stats

- **Total Workflows**: 3
- **Jobs per CI Run**: 6
- **Build Platforms**: 9 (binaries) + 3 (Docker)
- **Average CI Time**: ~10 minutes
- **Average Release Time**: ~30 minutes
- **Cache Hit Rate**: ~80%

---

## Future Enhancements

Planned improvements:

- [ ] Add performance benchmarking workflow
- [ ] Implement canary releases
- [ ] Add smoke tests for releases
- [ ] Create nightly builds
- [ ] Add integration tests with real GitHub API
- [ ] Implement blue-green deployment
- [ ] Add release drafter for changelog
- [ ] Create separate security workflow

---

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Build Push Action](https://github.com/docker/build-push-action)
- [Go Setup Action](https://github.com/actions/setup-go)
- [Trivy Action](https://github.com/aquasecurity/trivy-action)
- [golangci-lint Action](https://github.com/golangci/golangci-lint-action)

---

## Support

For issues with workflows:
1. Check workflow logs in Actions tab
2. Review this documentation
3. Open an issue with workflow run link
4. Tag with `ci/cd` label

---

**Last Updated**: 2024-12-09
**Maintained By**: [@l13t](https://github.com/l13t)