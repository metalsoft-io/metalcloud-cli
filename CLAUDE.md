# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

metalcloud-cli is a Go CLI tool for managing MetalCloud infrastructure. It wraps the `metalcloud-sdk-go` SDK and uses Cobra for command structure and Viper for configuration.

## Build & Test Commands

```bash
# Build
go build ./cmd/metalcloud-cli/

# Build for development (skips server version validation)
go build -ldflags="-X main.version=v7.0.0 -X main.allowDevelop=true" ./cmd/metalcloud-cli/

# Run all tests
go test ./...

# Run a single package's tests
go test ./pkg/formatter/
go test ./internal/firmware_catalog/

# Tidy dependencies
go mod tidy
```

## Architecture

### Entry Point & Command Flow

`cmd/metalcloud-cli/main.go` → `cmd/metalcloud-cli/cmd/root.go` (Execute) → Cobra command tree.

Before every command runs, `rootPersistentPreRun` initializes the logger, validates endpoint/API key, creates the SDK API client, validates CLI version compatibility (min 7.0, max 7.1), fetches user permissions, and hides commands the user lacks permissions for.

### Directory Layout

- `cmd/metalcloud-cli/cmd/` — Cobra command definitions (~45 files). Each file registers commands in `init()` and delegates to `internal/`.
- `cmd/metalcloud-cli/system/` — Constants (config keys, 100+ permission strings), version validation.
- `internal/<resource>/` — Command implementations. Each module calls the SDK, inspects responses, and formats output.
- `pkg/api/` — SDK client creation and context-based dependency injection (`api.GetApiClient(ctx)`, `api.GetUserId(ctx)`).
- `pkg/formatter/` — Multi-format output (text/csv/md/json/yaml) using `PrintConfig` structs and go-pretty tables.
- `pkg/response_inspector/` — Centralized HTTP response error handling.
- `pkg/logger/` — Zerolog-based logging.

### Adding a New Command

1. Create `cmd/metalcloud-cli/cmd/<resource>.go`:
   - Define a flags struct for command-specific options.
   - Define `cobra.Command` vars with `Use`, `Aliases`, `Short`, `Long`, `RunE`, and permission `Annotations`.
   - Register with `rootCmd.AddCommand()` in `init()`.

2. Create `internal/<resource>/<resource>.go`:
   - Define a `formatter.PrintConfig` with field titles, widths, order, and optional transformers.
   - Implement functions that get the client via `api.GetApiClient(ctx)`, call SDK APIs, check errors via `response_inspector.InspectResponse()`, and print via `formatter.PrintResult()`.

### Permission System

Commands use `Annotations: map[string]string{system.REQUIRED_PERMISSION: system.PERMISSION_X}` to declare required permissions. The root pre-run hook hides commands the user can't access.

### Configuration

Viper reads `metalcloud.yaml` from `.`, `$HOME/.metalcloud/`, or `/etc/metalcloud/`. Environment variables use prefix `METALCLOUD_` (e.g., `METALCLOUD_ENDPOINT`, `METALCLOUD_API_KEY`). Hyphens in flag names become underscores in env vars.

### Key Dependencies

- `github.com/metalsoft-io/metalcloud-sdk-go` — API SDK
- `github.com/spf13/cobra` + `github.com/spf13/viper` — CLI framework + config
- `github.com/jediv0t/go-pretty/v6` — Table rendering
- `github.com/rs/zerolog` — Structured logging

### Release Process

Tag-based releases via GoReleaser (`.goreleaser.yml`). Push a new tag (`git tag v7.0.x && git push --tags`) to trigger GitHub Actions that build cross-platform binaries and publish to Homebrew, deb, and rpm repos.
