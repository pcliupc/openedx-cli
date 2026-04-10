# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go-based CLI for Open edX designed for CI pipelines and coding agents. It uses official public APIs by default and falls back to user-configured extension APIs when official endpoints are unavailable. The command surface stays stable regardless of which backend provider serves a request.

## Build and Test Commands

```bash
make build          # go build ./cmd/openedx
make test           # go test ./...
make test-integration  # OPENEDX_INTEGRATION=1 go test ./integration -v
```

Run a single test:

```bash
go test ./internal/cli -run TestRootCommandDefaultsToJSON -v
go test ./internal/cli/cmd -run TestCourse -v
```

## Tech Stack

- Go (single portable binary)
- Cobra for CLI command structure and help
- Viper for config file loading only (not global state)
- `net/http` or resty for HTTP
- testify for assertions
- Golden JSON fixtures under `testdata/`

## Architecture

The CLI has six internal layers with clear boundaries:

1. **Command Layer** (`internal/cli/`, `internal/cli/cmd/`) — Cobra command tree, argument parsing, output formatting. Noun-verb structure: `openedx course list`, `openedx user create`.

2. **Capability Registry** (`internal/registry/`) — Built-in map of latest official Open edX API endpoints. One mapping set, versioned with the CLI, not user-configurable.

3. **Provider Layer** (`internal/provider/`) — Executes API calls. Two providers: `public` (official APIs via registry) and `extension` (user-configured custom APIs). Fallback logic only retries on 404/405/501; auth/validation errors are never retried through extension.

4. **Auth Layer** (`internal/auth/`) — Profile-based OAuth client credentials. Secrets referenced by env var names, never stored in config. In-memory token caching.

5. **Normalizer Layer** (`internal/normalize/`) — Converts provider responses into stable resource models (`Course`, `User`, `Enrollment`, `RoleAssignment`, `CourseOutline`, `Job`). Both public and extension payloads normalize to the same output shape.

6. **Diagnostics Layer** (`internal/diagnostics/`) — `schema` command shows command-to-endpoint mapping and extension presence. `doctor` command checks base URL reachability, token acquisition, and API availability.

## Key Design Decisions

- Command names are stable and provider-independent
- Public provider is always tried first; extension is fallback only
- Only API-unavailable errors (404, 405, 501) trigger extension fallback
- 400/401/403 errors are never retried through extension
- JSON is the default output format (stdout = data, stderr = logs)
- Non-interactive by default
- Long flags preferred over short flags

## Configuration

YAML config with profiles and extension mappings. Profiles contain base URL, token URL, and env var names for credentials. Extension mappings define custom API endpoints per command key. Example at `testdata/config/example.yaml`.

For adding new commands or configuring extension APIs, see [docs/extending-the-cli.md](docs/extending-the-cli.md).

## Command Domains

Current built-in commands cover these domains:

- `course` — list, get, create, import, export, outline get
- `user` — create, list, get
- `enrollment` — add, list, remove
- `role` — assign
- `grade` — list, gradebook
- `certificate` — list
- `schema` — inspect command-to-endpoint mappings
- `doctor` — health checks

## Entrypoint

`cmd/openedx/main.go` — wires root command and calls `cli.NewRootCmd().Execute()`.

## License

MIT
