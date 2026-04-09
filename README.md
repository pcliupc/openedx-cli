# Open edX CLI

A CLI for Open edX designed for CI pipelines and coding agents. Uses official public APIs by default, with automatic fallback to configured extension APIs when official endpoints are unavailable.

## Build

```bash
make build        # builds bin/openedx
make test         # runs unit tests
make clean        # removes build artifacts
```

## Commands (v1)

```bash
openedx --profile admin course list
openedx --profile admin course get --course-id course-v1:demo+cs101+2026
openedx --profile admin course create --org demo --number cs101 --run 2026 --title "Intro to AI"
openedx --profile admin course import --course-id course-v1:demo+cs101+2026 --file ./course.tar.gz
openedx --profile admin course export --course-id course-v1:demo+cs101+2026
openedx --profile admin course outline get --course-id course-v1:demo+cs101+2026

openedx --profile ops user create --username alice --email alice@example.com --name "Alice"
openedx --profile ops enrollment add --course-id course-v1:demo+cs101+2026 --username alice --mode audit
openedx --profile admin role assign --course-id course-v1:demo+cs101+2026 --username alice --role staff

openedx schema all
openedx schema course.create
openedx doctor
openedx doctor verify course.list
```

## Configuration

Create `openedx.yaml` in the current directory or `~/.openedx/config.yaml`:

```yaml
version: 1

profiles:
  admin:
    base_url: https://openedx.example.com
    token_url: https://openedx.example.com/oauth2/access_token
    client_id_env: OPENEDX_ADMIN_CLIENT_ID
    client_secret_env: OPENEDX_ADMIN_CLIENT_SECRET
    default_format: json

extensions:
  course.create:
    method: POST
    url: https://openedx.example.com/api/cli-ext/course/create
```

Secrets are referenced by environment variable name, not stored directly. Set the env vars before running commands:

```bash
export OPENEDX_ADMIN_CLIENT_ID="your-client-id"
export OPENEDX_ADMIN_CLIENT_SECRET="your-client-secret"
```

## Extension Fallback

The CLI tries official public APIs first. If an endpoint returns 404, 405, or 501, and an extension mapping exists, the CLI retries through the extension provider. Auth and validation errors (400, 401, 403) never trigger fallback.

## Integration Tests

```bash
make test-integration   # requires OPENEDX_INTEGRATION=1 and a live Open edX instance
```
