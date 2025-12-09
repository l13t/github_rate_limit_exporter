# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of GitHub Rate Limit Exporter
- Support for monitoring multiple GitHub users/tokens
- Prometheus metrics export for:
  - Core API rate limits
  - Search API rate limits
  - GraphQL API rate limits
  - Integration Manifest rate limits
- Configuration file support for YAML, TOML, and HCL formats
- Docker support with multi-stage builds
- Docker Compose setup with Prometheus and Grafana
- Systemd service file for Linux deployments
- Kubernetes deployment manifests
- Comprehensive documentation (README, INSTALL, CONTRIBUTING, GETTING_STARTED)
- Example configuration files for all supported formats
- Taskfile.yml for modern task running (replaces Makefile)
- Task usage guide (TASK_USAGE.md)
- Quick reference card (QUICK_REFERENCE.md)
- Quick start script for easy setup (supports both Task and direct Go builds)
- Health check endpoint
- Graceful shutdown handling
- Unit tests for configuration loading (100% coverage)
- Example Prometheus alerting rules
- Security hardening in Docker and systemd
- **GitHub Actions CI/CD**:
  - CI workflow for testing, linting, and security scanning
  - Docker workflow for multi-arch container builds (amd64, arm64, armv7)
  - Release workflow for automated binary and Docker image publishing
  - Support for 9 platform combinations (Linux, macOS, Windows, FreeBSD)
  - Automated security scanning with Trivy and Gosec
  - SBOM generation and attestation
  - Dependabot for automated dependency updates
  - Issue templates for bugs and feature requests
  - Pull request template for consistency
  - Comprehensive workflow documentation

### Changed
- Migrated from Makefile to Taskfile.yml for better cross-platform support
- Updated all documentation to reference Task instead of Make
- Enhanced quickstart script to detect and use Task when available
- Updated README with CI/CD badges and multi-arch Docker image info

### Features
- Automatic polling with configurable intervals
- Concurrent metric collection for multiple users
- Secure token handling
- Non-root Docker container execution
- Comprehensive error handling and logging
- Prometheus best practices compliance

## [1.0.0] - YYYY-MM-DD

### Added
- First stable release

---

## Release Notes Template

When creating a new release, use the following template:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security improvements
```

---

[Unreleased]: https://github.com/l13t/github_rate_limit_exporter/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/l13t/github_rate_limit_exporter/releases/tag/v1.0.0