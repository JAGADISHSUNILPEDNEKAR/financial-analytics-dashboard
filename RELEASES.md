# Release and Distribution Guide

This repository uses GitHub Actions to automate the creation of releases and the distribution of packages (Docker images and binaries).

## Versioning Policy

We follow [Semantic Versioning (SemVer)](https://semver.org/).
- **Major**: Breaking changes.
- **Minor**: New features (backwards-compatible).
- **Patch**: Bug fixes (backwards-compatible).

Pre-releases use tags like `v1.0.0-alpha.1` or `v1.0.0-rc.1`.

## Release Process

1.  **Draft Release**: Create a new tag following the `v*` pattern (e.g., `git tag v1.0.0`).
2.  **Push Tag**: Push the tag to GitHub (`git push origin v1.0.0`).
3.  **Automated Workflow**: The `Release` workflow will trigger automatically.
    - Builds Go binaries for Linux, macOS, and Windows.
    - Builds the Rust Analytics Engine.
    - Packages the Flutter web application.
    - Creates a GitHub Release with these artifacts attached.

## Package Distribution

### Docker Images (GHCR)

Docker images are automatically built and pushed to the GitHub Container Registry (GHCR) on:
- Every push to the `main` branch (tagged with `latest` and `sha-*`).
- Tag creation (tagged with the version name).

To pull an image:
```bash
docker pull ghcr.io/jagadishsunilpednekar/financial-analytics-api-gateway:latest
```

Available images:
- `financial-analytics-api-gateway`
- `financial-analytics-analytics-engine`
- `financial-analytics-auth-service`
- `financial-analytics-user-service`
- `financial-analytics-dashboard-service`
- `financial-analytics-notification-service`
- `financial-analytics-ml-services`

### Binaries and Web Build

Static binaries and the Flutter web build are available in the [GitHub Releases](https://github.com/jagadishsunilpednekar/financial-analytics-dashboard/releases) section.

- **Go Services**: Binaries for `api-gateway`, `auth-service`, `user-service`, and `dashboard-service`.
- **Rust Service**: Binary for `analytics-engine` (Linux).
- **Frontend**: `flutter-web.zip` containing the production-ready web application.

---

## Technical Details

- **CI Workflow**: [.github/workflows/ci.yml](file:///Users/jagadishsunilpednekar/financial-analytics-dashboard/.github/workflows/ci.yml)
- **Release Workflow**: [.github/workflows/release.yml](file:///Users/jagadishsunilpednekar/financial-analytics-dashboard/.github/workflows/release.yml)
