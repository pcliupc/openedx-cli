---
name: openedx-shared
version: 1.0.0
description: "OpenEdX CLI 基础：配置初始化、认证、安全规则、诊断命令"
metadata:
  requires:
    skills: []
    bins: ["openedx"]
---

## CLI Installation

Install the OpenEdX CLI before using any skill:

```bash
# npm (recommended)
npm install -g @openedx/cli

# or build from source
git clone https://github.com/your-org/openedx-cli.git
cd openedx-cli && make build
```

Verify installation: `openedx --help`

## Configuration

The CLI uses a YAML config file (`openedx.yaml` or `~/.openedx/config.yaml`) with named profiles:

```yaml
version: 1
profiles:
  production:
    base_url: https://openedx.example.com
    token_url: https://openedx.example.com/oauth2/access_token
    client_id_env: OPENEDX_PROD_CLIENT_ID
    client_secret_env: OPENEDX_PROD_CLIENT_SECRET
    default_format: json
```

Key concepts:
- **Profile** — a named set of connection settings. Switch with `--profile <name>`.
- **Environment variables** — credential values are never stored in config. The config stores only the *env var names* that hold the actual secrets.
- **default_format** — output format for commands (json, table). Default is json.

Set environment variables before running commands:

```bash
export OPENEDX_PROD_CLIENT_ID="your-client-id"
export OPENEDX_PROD_CLIENT_SECRET="your-client-secret"
```

## Authentication

The CLI uses OAuth 2.0 client credentials flow:

1. Reads `client_id_env` and `client_secret_env` from the active profile
2. Resolves the actual credential values from environment variables
3. POSTs to `token_url` with `grant_type=client_credentials`
4. Receives and caches the access token in memory (not on disk)
5. Sends the token as `Authorization: Bearer <token>` on subsequent requests

Token caching is in-memory only. Tokens are never written to the filesystem.

## Security Rules

MUST follow these rules when using the CLI:

1. **Credentials never stored in config** — only env var names appear in YAML
2. **Non-interactive by default** — all commands run without prompting
3. **stdout = data, stderr = logs** — JSON output goes to stdout; diagnostic/log messages go to stderr
4. **No plaintext secrets in commands** — never pass `--client-id` or `--client-secret` as flags

## Diagnostics

Before using business commands, verify connectivity:

```bash
# Full health check: base URL → token acquisition → API availability
openedx doctor

# Check a specific command's endpoint
openedx doctor verify course.list

# View all command-to-endpoint mappings
openedx schema all

# View a specific command's mapping
openedx schema course.list
```

## Extension APIs

When the public API does not support a command, configure an extension endpoint:

```yaml
extensions:
  grade.list:
    method: GET
    url: https://openedx.example.com/api/custom/v1/grades
```

Fallback behavior:
- CLI tries the public API first
- If the public API returns **404**, **405**, or **501**, the CLI retries through the extension endpoint
- **400**, **401**, **403** errors are NEVER retried through extension — these indicate input or permission problems

Extension endpoint requirements:
- Accept `Authorization: Bearer <token>` header
- GET: args as query parameters
- POST: args as JSON body
- Return JSON with standard HTTP status codes
```