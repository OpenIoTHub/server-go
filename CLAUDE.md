# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the OpenIoTHub server written in Go. It acts as a central server for IoT gateway connections, providing session management, multiple protocol support (TCP, KCP/UDP, TLS, gRPC), and HTTP/HTTPS proxy services. The server manages connections from gateways and handles proxying requests to IoT devices.

## Common Development Tasks

### Building
- `go build` - Build the binary locally
- `go test ./...` - Run all tests
- GoReleaser is used for multi-platform releases (see `.goreleaser.yml`)

### Testing
- Tests are in `*_test.go` files. The main test is in `main_test.go`.
- CI runs `go test ./...` on pushes and PRs via GitHub Actions (`.github/workflows/test.yml`).

### Running
- Default command: `./server-go` (requires config file `server-go.yaml`)
- Generate config: `./server-go init`
- Generate tokens: `./server-go generate`
- Test command: `./server-go test`
- Configuration path can be set via `--config` flag or `ServerConfigFilePath` env var.

### Development Workflow
1. Make changes to Go files
2. Run `go test ./...` to ensure tests pass
3. Run `go build` to verify compilation
4. Use `go mod tidy` to clean up dependencies (if needed)

## Architecture Overview

### Key Components
- **Main Entry Point** (`main.go`): CLI setup using `urfave/cli`. Commands: `generate`, `init`, `test`, and default server start.
- **Configuration** (`config/`): YAML-based config (`server-go.yaml`). Supports Redis for persistent storage. Config constants in `const.go`.
- **Session Management** (`manager/`, `session/`):
  - `SessionsManager` (`manager/SessionsManager.go`): Central manager for all client sessions.
  - `Session` (`session/Session.go`): Represents individual gateway connections.
  - Supports multiplexed (yamux) and non-multiplexed connections.
- **Network Services** (launched in `main.go:run()`):
  - TCP server on port 34320
  - KCP/UDP server on port 34320 (UDP)
  - TLS server on port 34321
  - gRPC server on port 34322
  - HTTP proxy on port 80
  - HTTPS proxy on port 443
  - UDP API server on port 34321
  - KCP API server on port 34322
- **Storage Abstraction** (`iface/`, `imp/`):
  - Interface `RuntimeStorageIfce` (`iface/runtimeStorage/runtimeStorage.go`) for runtime storage.
  - Implementations: in-memory (`memImp/`) and Redis (`redisImp/`).
  - Used for HTTP proxy configurations and session data.
- **Protocol Support**: Custom messaging via `github.com/OpenIoTHub/utils/msg`, gRPC, KCP, HTTP/HTTPS.

### Configuration Details
Default config file `server-go.yaml` includes:
- `my_public_ip_or_domain`: Public IP/domain for token generation.
- `common`: Bind address and port settings.
- `security`: Login key, TLS certificates.
- `redisconfig`: Redis persistence settings (enable for HTTP proxy config persistence).
- `logconfig`: Log output settings.

Important ports (defaults):
- 34320: TCP and KCP/UDP
- 34321: TLS and UDP API
- 34322: gRPC and KCP API
- 80: HTTP proxy
- 443: HTTPS proxy

### Storage
- By default, HTTP proxy configs are stored in memory.
- Enable Redis in config (`redisconfig.enabled: true`) for persistence across restarts.
- Storage interface allows swapping implementations.

### Session Management
- Each gateway connection is a `Session` with unique RunId.
- Sessions map in `SessionsManager` keyed by RunId.
- Supports multiplexing via yamux for multiple streams over a single connection.

## Notes for Development

- The server uses vendored dependencies (`vendor/`).
- The project uses Go modules (`go.mod` with Go 1.21).
- Dockerfile builds a minimal image with entrypoint script.
- Releases are automated via GoReleaser and GitHub Actions, producing binaries for multiple platforms, Docker images, and package manager releases (Homebrew, Snapcraft, Scoop, DEB/RPM).
- The codebase is primarily in Chinese (comments, logs, README), but the code and architecture follow standard Go practices.