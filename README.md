# Open edX CLI

A command-line tool for [Open edX](https://open.edx.org/) designed for CI pipelines and coding agents.

The CLI uses official Open edX public APIs by default. When an endpoint is not available on a deployment, it can automatically fall back to user-configured extension APIs. The command surface stays stable regardless of which backend serves the request.

## Why This Exists

Interacting with Open edX from automation typically means writing ad-hoc `curl` calls against undocumented or deployment-specific endpoints. This CLI provides:

- **A stable command surface** — command names do not change when backends differ.
- **JSON-first output** — structured, machine-readable, suitable for piping into `jq` or other tools.
- **Non-interactive by default** — no prompts, no TUI, works in headless CI.
- **Extension fallback** — when official APIs are missing, configured extension endpoints fill the gap transparently.

## Install

### npm (recommended)

```bash
npm install -g @pcliupc/openedx-cli
```

The postinstall script automatically downloads the correct pre-built binary for your platform from [GitHub Releases](https://github.com/pcliupc/openedx-cli/releases).

Requires Node.js 14+.

### Build from source

Requires Go 1.23+:

```bash
git clone https://github.com/pcliupc/openedx-cli.git
cd openedx-cli
make build
```

This produces a single binary at `bin/openedx`.

## Quick Start

Create a config file at `./openedx.yaml` or `~/.openedx/config.yaml`:

```yaml
version: 1

profiles:
  admin:
    base_url: https://openedx.example.com
    token_url: https://openedx.example.com/oauth2/access_token
    client_id_env: OPENEDX_ADMIN_CLIENT_ID
    client_secret_env: OPENEDX_ADMIN_CLIENT_SECRET
    default_format: json
```

Set the credential environment variables:

```bash
export OPENEDX_ADMIN_CLIENT_ID="your-client-id"
export OPENEDX_ADMIN_CLIENT_SECRET="your-client-secret"
```

Run a command:

```bash
openedx --profile admin course list
```

## Commands

All commands use a noun-verb structure: `openedx <noun> <verb>`.

### Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--profile` | `-p` | | Config profile to use |
| `--format` | `-f` | `json` | Output format (`json`, `table`) |
| `--config` | `-c` | | Path to config file |
| `--verbose` | `-v` | `false` | Enable verbose output |

### Course

```bash
# List all courses
openedx --profile admin course list
openedx --profile admin course list --all
openedx --profile admin course list --page 1 --page-size 25

# Get a single course
openedx --profile admin course get --course-id course-v1:DemoX+DemoCourse+2026

# Create a new course
openedx --profile admin course create \
  --org DemoX \
  --number CS101 \
  --run 2026 \
  --title "Introduction to Computer Science"

# Import a course from a tar.gz archive
openedx --profile admin course import \
  --course-id course-v1:DemoX+CS101+2026 \
  --file ./course.tar.gz

# Export a course
openedx --profile admin course export \
  --course-id course-v1:DemoX+CS101+2026

# Get course outline structure
openedx --profile admin course outline get \
  --course-id course-v1:DemoX+CS101+2026
```

### User

```bash
# Create a user
openedx --profile admin user create \
  --username alice \
  --email alice@example.com \
  --name "Alice Smith"
```

### Enrollment

```bash
# Enroll a user in a course
openedx --profile admin enrollment add \
  --course-id course-v1:DemoX+CS101+2026 \
  --username alice \
  --mode audit
```

The `--mode` flag defaults to `audit`. Other common values are `verified`, `honor`, and `professional`.

### Role

```bash
# Assign a role to a user on a course
openedx --profile admin role assign \
  --course-id course-v1:DemoX+CS101+2026 \
  --username alice \
  --role staff
```

### Schema

Inspect how commands map to backend API endpoints:

```bash
# Show schema for a single command
openedx schema course.create

# Show all command schemas
openedx schema all
```

Output includes the command key, public HTTP method and path, whether an extension override exists, required arguments, and the output model name.

### Doctor

Run health diagnostics against the configured Open edX deployment:

```bash
# Run all checks (base URL reachability, token acquisition)
openedx doctor

# Verify a specific command's mapping
openedx doctor verify course.list
```

## Configuration

### File Locations

The CLI looks for configuration in this order:

1. Path specified by `--config` flag
2. `./openedx.yaml` in the current directory
3. `~/.openedx/config.yaml` in the user home directory

### Profile

A profile defines how to connect to an Open edX deployment:

```yaml
profiles:
  admin:
    base_url: https://openedx.example.com
    token_url: https://openedx.example.com/oauth2/access_token
    client_id_env: OPENEDX_ADMIN_CLIENT_ID
    client_secret_env: OPENEDX_ADMIN_CLIENT_SECRET
    default_format: json
```

`client_id_env` and `client_secret_env` are **environment variable names**, not the actual credentials. The CLI reads the credential values from those env vars at runtime. This avoids storing secrets in config files.

### Extensions

Extension mappings define custom API endpoints that the CLI can fall back to when official endpoints are unavailable:

```yaml
extensions:
  course.create:
    method: POST
    url: https://openedx.example.com/api/cli-ext/course/create
  course.import:
    method: POST
    url: https://openedx.example.com/api/cli-ext/course/import
  course.export:
    method: POST
    url: https://openedx.example.com/api/cli-ext/course/export
```

See the [Extension Fallback](#extension-fallback) section for details on when extensions are used.

### Full Example

```yaml
version: 1

profiles:
  admin:
    base_url: https://openedx.example.com
    token_url: https://openedx.example.com/oauth2/access_token
    client_id_env: OPENEDX_ADMIN_CLIENT_ID
    client_secret_env: OPENEDX_ADMIN_CLIENT_SECRET
    default_format: json

  ops:
    base_url: https://openedx.example.com
    token_url: https://openedx.example.com/oauth2/access_token
    client_id_env: OPENEDX_OPS_CLIENT_ID
    client_secret_env: OPENEDX_OPS_CLIENT_SECRET
    default_format: json

extensions:
  course.create:
    method: POST
    url: https://openedx.example.com/api/cli-ext/course/create
  course.import:
    method: POST
    url: https://openedx.example.com/api/cli-ext/course/import
  course.export:
    method: POST
    url: https://openedx.example.com/api/cli-ext/course/export
```

## Extension Fallback

When the CLI executes a command, it always tries the official public API first. If that call fails with an "endpoint unavailable" error and an extension mapping exists for the same command, the CLI retries through the extension provider.

**Errors that trigger fallback:**

| HTTP Status | Meaning |
|-------------|---------|
| `404 Not Found` | Endpoint does not exist on this deployment |
| `405 Method Not Allowed` | Endpoint exists but does not support this HTTP method |
| `501 Not Implemented` | Server explicitly signals the endpoint is not available |

**Errors that do NOT trigger fallback:**

| HTTP Status | Meaning |
|-------------|---------|
| `400 Bad Request` | The request is invalid — fix the input |
| `401 Unauthorized` | Credentials are wrong — fix the auth config |
| `403 Forbidden` | The profile lacks permission — use a different profile |
| `5xx` (except 501) | Server error — retrying through extension would mask the problem |

This design ensures that auth and permission issues are surfaced clearly rather than hidden behind a provider switch.

## Output Format

### JSON (default)

All commands output JSON to stdout. Errors go to stderr as JSON.

Example `course list` output:

```json
[
  {
    "course_id": "course-v1:DemoX+CS101+2026",
    "org": "DemoX",
    "number": "CS101",
    "run": "2026",
    "title": "Introduction to Computer Science",
    "pacing": "instructor",
    "start": "2026-01-01T00:00:00Z",
    "end": "2026-06-30T23:59:59Z"
  }
]
```

Example error output:

```json
{
  "error": "permission_denied",
  "message": "profile 'ops' cannot create course",
  "resource": "course.create",
  "suggestion": "retry with --profile admin or grant the required platform permission"
}
```

### Resource Models

All commands normalize backend responses into these stable shapes:

| Model | Fields |
|-------|--------|
| **Course** | `course_id`, `org`, `number`, `run`, `title`, `pacing`, `start`, `end`, `status` |
| **CourseOutline** | `course_id`, `chapters[]` (recursive: `id`, `title`, `type`, `children[]`) |
| **User** | `username`, `email`, `name`, `is_active`, `created_at` |
| **Enrollment** | `username`, `course_id`, `mode`, `is_active` |
| **RoleAssignment** | `username`, `course_id`, `role`, `assigned_by`, `assigned_at` |
| **Job** | `job_id`, `operation`, `status`, `submitted_at`, `finished_at`, `result`, `artifacts[]` |

Both public API responses and extension API responses normalize to the same output shape, so the command output is always consistent.

## CI Usage

The CLI is designed for non-interactive use in CI pipelines:

```bash
# Create a course
openedx --profile admin course create \
  --org DemoX --number CS101 --run 2026 --title "New Course" \
  --format json

# Create a user and enroll them
openedx --profile admin user create --username bob --email bob@example.com
openedx --profile admin enrollment add \
  --course-id course-v1:DemoX+CS101+2026 --username bob --mode audit

# Export course content
openedx --profile admin course export \
  --course-id course-v1:DemoX+CS101+2026
```

Tips for CI:
- Use the `--config` flag to point to a config file in a known location
- Store credentials in CI secret variables, then export them as env vars before running commands
- Pipe JSON output through `jq` for specific field extraction: `openedx course list | jq '.[0].course_id'`

## Development

### Prerequisites

- Go 1.23 or later

### Build and Test

```bash
make build              # Build the binary
make test               # Run all unit tests
make test-integration   # Run integration tests (requires OPENEDX_INTEGRATION=1)
make clean              # Remove build artifacts
```

Run a specific test:

```bash
go test ./internal/cli/cmd -run TestCourseList -v
go test ./internal/provider -run TestFallback -v
```

### Adding a New Command

Adding a command touches five files, always in this order:

1. **Registry** (`internal/registry/public_registry.go`) — add the API endpoint metadata (key, method, path, required args, output model)
2. **Model** (`internal/model/`) — define the output struct (flat fields, stable JSON contract)
3. **Normalizer** (`internal/normalize/`) — convert raw API response into the model struct, with test fixtures under `testdata/`
4. **Command** (`internal/cli/cmd/`) — build the Cobra command with flags, call `execFn`, normalize, print output
5. **Root** (`internal/cli/root.go`) — register the new command group in `rootCmd.AddCommand(...)`

Example — adding `grade list`:

```
Registry:  "grade.list" → GET /api/grades/v1/courses/{course_id}/
Model:     model.Grade{Username, CourseID, Percent, Letter, Passed}
Normalize: GradesFromJSON(raw) → []model.Grade
Command:   NewGradeCmd(execFn) → newGradeListCmd → --course-id flag
Root:      rootCmd.AddCommand(cmd.NewGradeCmd(execFn))
```

For the full step-by-step guide with code examples, see [docs/extending-the-cli.md](docs/extending-the-cli.md).

### Architecture

```
cmd/openedx/main.go              Entry point
internal/
  cli/
    root.go                      Root command, global flags
    output.go                    JSON output helpers
    errors.go                    Structured error type
    cmd/
      course.go                  course list/get/create/import/export/outline
      user.go                    user create
      enrollment.go              enrollment add
      role.go                    role assign
      schema.go                  schema inspection
      doctor.go                  health diagnostics
  config/
    types.go                     Config, Profile, ExtensionMapping structs
    config.go                    YAML config loading via Viper
  auth/
    token_client.go              OAuth client credentials flow
    cache.go                     In-memory token cache with expiry
  registry/
    commands.go                  CommandMeta type
    public_registry.go           Built-in map of official API endpoints
  provider/
    provider.go                  Provider interface, ProviderError
    public_provider.go           Official API execution
    extension_provider.go        Extension API execution
    fallback.go                  Public-first with extension fallback
  normalize/
    course.go                    Course + CourseList normalizers
    outline.go                   CourseOutline normalizer
    user.go                      User normalizer
  model/
    course.go, outline.go, ...   Stable output structs
  diagnostics/
    schema.go                    Command schema introspection
    doctor.go                    Health check functions
testdata/
  config/example.yaml            Example configuration
  public/*.json                  Public API response fixtures
  extension/*.json               Extension API response fixtures
integration/
  tutor_smoke_test.go            Integration tests for Tutor deployments
```

### Data Flow

```
User runs: openedx --profile admin course list
       │
       ▼
  Command Layer (Cobra)
  Parses flags, builds args map
       │
       ▼
  Registry Layer
  Looks up "course.list" → GET /api/courses/v1/courses
       │
       ▼
  Auth Layer
  Loads profile, acquires OAuth token from env vars
       │
       ▼
  Provider Layer
  Tries public API first
  On 404/405/501: retries through extension (if configured)
       │
       ▼
  Normalizer Layer
  Converts raw JSON response into stable Course model
       │
       ▼
  Output Layer
  Prints JSON to stdout
```

## License

[MIT](LICENSE)
