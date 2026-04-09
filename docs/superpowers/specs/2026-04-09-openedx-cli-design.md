# Open edX CLI Design

Date: 2026-04-09
Status: Draft for review

## 1. Purpose

Design a CLI for Open edX that is easy to use from CI and by coding agents.

The CLI should:

- prefer official Open edX public APIs
- support the latest official Open edX API shape as the built-in standard
- allow custom Open edX APIs through extension configuration when official APIs are missing or unavailable
- keep one stable command surface for users, CI pipelines, and agents
- support Tutor deployments on recent Open edX releases

This CLI is not intended to be a complete multi-version compatibility layer for every historical Open edX release. It is a product interface aligned to the latest supported official API surface, with an extension fallback mechanism for older or customized deployments.

## 2. Goals

### Primary goals

- Make common Open edX platform operations easy to run in CI.
- Give agents a stable, structured, non-interactive command interface.
- Keep the command surface compatible with future custom APIs without changing command names.
- Preserve compatibility with official Open edX where possible.

### Non-goals

- Supporting every Open edX release with separate built-in API maps
- Hiding authorization failures behind CLI-side role simulation
- Building a full Open edX replacement SDK before proving the core workflow

## 3. Initial Scope

The first supported capability set is:

- create course
- list courses
- get course details
- import course
- export course
- get course outline or chapter structure
- create user account
- enroll user into course
- assign course role to user

Course outline support should use official public APIs if available. If the latest public API shape is not available on a target deployment, the same command may be backed by an extension API.

## 4. Design Principles

The CLI should follow agent-friendly CLI design principles:

- noun-first command hierarchy
- long flags first
- JSON-first output for automation
- strict separation of stdout and stderr
- non-interactive by default in CI use
- stable command naming independent of backend provider
- machine-readable errors

The design takes inspiration from recent enterprise CLIs such as `lark-cli` and `dws`, especially:

- stable command tree
- automation-first behavior
- schema or doctor style inspection tools
- support for a simple command surface over a more complex backend

## 5. Architecture Summary

The CLI has one unified command surface and two backend provider sources:

- built-in public provider
- configured extension provider

The CLI always tries the built-in public provider first. If that call fails because the API is unavailable and the same command has an extension mapping configured, the CLI retries through the extension provider.

The user should not need to care whether a command is served by official Open edX APIs or by a custom extension API. The command surface stays stable.

## 6. Command Model

### Top-level command philosophy

Use a noun-verb structure:

```bash
openedx course create
openedx course list
openedx course get
openedx course import
openedx course export
openedx course outline get

openedx user create
openedx enrollment add
openedx role assign
```

### Recommended top-level domains

- `auth`
- `course`
- `user`
- `enrollment`
- `role`
- `api`
- `schema`
- `doctor`

### Example commands

```bash
openedx --profile admin course create \
  --org demo \
  --number cs101 \
  --run 2026 \
  --title "Intro to AI"

openedx --profile admin course list --format json

openedx --profile admin course get \
  --course-id course-v1:demo+cs101+2026

openedx --profile admin course import \
  --course-id course-v1:demo+cs101+2026 \
  --file ./course.tar.gz

openedx --profile admin course export \
  --course-id course-v1:demo+cs101+2026 \
  --output ./exports/

openedx --profile ops course outline get \
  --course-id course-v1:demo+cs101+2026

openedx --profile ops user create \
  --username alice \
  --email alice@example.com \
  --name "Alice"

openedx --profile ops enrollment add \
  --course-id course-v1:demo+cs101+2026 \
  --username alice \
  --mode audit

openedx --profile admin role assign \
  --course-id course-v1:demo+cs101+2026 \
  --username alice \
  --role staff
```

## 7. Provider Model

There are only two backend provider sources:

- `public`
- `extension`

### Public provider

The CLI ships with a built-in registry for the latest supported official Open edX API mappings.

This registry is maintained in code and versioned with the CLI. It is not user configuration.

Conceptually:

```yaml
commands:
  course.list:
    method: GET
    path: /api/...
  course.get:
    method: GET
    path: /api/...
  enrollment.add:
    method: POST
    path: /api/...
```

The CLI only maintains one built-in mapping set: the latest supported official API shape.

### Extension provider

Extensions are user-configured mappings for custom Open edX APIs.

These are meant for deployments that:

- add their own APIs
- backport latest API behavior into older versions
- expose custom authoring or operational endpoints not present in official Open edX

Extensions are not a different kind of product feature. They are just a different provider source.

## 8. Provider Resolution Rules

For any command:

1. The CLI resolves the command to its built-in latest public mapping.
2. It tries the public provider first.
3. If the public call succeeds, it returns the normalized result.
4. If the public call fails with an API-unavailable class of error and an extension mapping exists for the same command, the CLI retries using the extension provider.
5. If no extension mapping exists, the public error is returned.

### Errors that should trigger extension fallback

- `404 Not Found`
- `405 Method Not Allowed`
- `501 Not Implemented`
- explicit endpoint-not-available responses from the platform

### Errors that should not trigger extension fallback

- `400 Validation Error`
- `401 Unauthorized`
- `403 Permission Denied`
- semantic business errors where the command exists but the request is invalid

This avoids hiding auth or validation issues behind a provider switch.

## 9. Configuration Model

User configuration should contain:

- profiles
- base URL
- token endpoint
- environment variable names for credentials
- extension mappings
- default output format

Sensitive secrets should not be stored directly in the config file body.

Example:

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

## 10. Authentication and Authorization

The CLI should default to machine-to-machine authentication using service credentials.

Profiles represent operational identities such as:

- platform administrator
- course operations administrator

The CLI should not simulate RBAC locally. Authorization remains the responsibility of Open edX or the extension API.

If a profile lacks permission, the CLI should surface structured permission errors clearly rather than trying alternate auth modes.

Example structured error:

```json
{
  "error": "permission_denied",
  "message": "profile 'ops' cannot create course",
  "resource": "course.create",
  "suggestion": "retry with --profile admin or grant the required platform permission"
}
```

## 11. Output Contract

### Default output

Default output should be JSON for all automation-oriented commands.

Supported output formats:

- `json`
- `table`
- `yaml`

JSON is the only recommended format for CI.

### stdout and stderr

- `stdout`: machine-readable result
- `stderr`: logs, warnings, progress, diagnostics

### Error format

Errors should be structured and stable.

Suggested top-level error codes:

- `auth_error`
- `permission_denied`
- `not_found`
- `validation_error`
- `remote_error`
- `provider_unavailable`
- `extension_not_configured`

## 12. Resource Models

The CLI should normalize backend responses into stable internal resource shapes.

### Course

- `course_id`
- `org`
- `number`
- `run`
- `title`
- `pacing`
- `start`
- `end`
- `status`

### CourseOutline

- `course_id`
- `chapters[]`
- `sequentials[]`
- `verticals[]`
- `blocks[]`

The goal is not to expose every raw block detail. The goal is to return a stable, CI- and agent-friendly course structure tree.

### User

- `username`
- `email`
- `name`
- `is_active`
- `created_at`

### Enrollment

- `username`
- `course_id`
- `mode`
- `is_active`

### RoleAssignment

- `username`
- `course_id`
- `role`
- `assigned_by`
- `assigned_at`

### Job

Used for long-running operations such as import and export:

- `job_id`
- `operation`
- `status`
- `submitted_at`
- `finished_at`
- `result`
- `artifacts[]`

## 13. Introspection and Diagnostics

The CLI should expose how commands are resolved.

### Schema command

Examples:

```bash
openedx schema course create
openedx schema course list
openedx schema all
```

The schema output should show:

- command name
- public mapping
- whether an extension mapping exists
- required arguments
- output resource shape

### Doctor command

Examples:

```bash
openedx doctor
openedx doctor verify course.list
openedx doctor verify course.create
```

Doctor should verify:

- base URL reachability
- token acquisition
- baseline public API availability
- extension configuration health when present

It is a diagnostic tool, not a dynamic capability discovery system.

## 14. Internal Module Breakdown

The implementation should be split into small modules with clear boundaries.

### 1. Command Layer

Responsibilities:

- command tree
- argument parsing
- output selection
- help text

### 2. Capability Registry

Responsibilities:

- latest public API command map
- command metadata
- required parameter definitions

### 3. Provider Layer

Responsibilities:

- execute public API calls
- execute extension API calls
- provider fallback logic

### 4. Auth Layer

Responsibilities:

- profile loading
- token acquisition
- token caching
- credential lookup from environment

### 5. Normalizer Layer

Responsibilities:

- convert provider responses into stable CLI resource models

### 6. Diagnostics Layer

Responsibilities:

- doctor
- schema
- capability visibility

## 15. Command Behavior Conventions

The CLI should follow these conventions:

- non-interactive by default
- long flags preferred
- JSON-first output
- stable command names
- structured errors
- explicit `--dry-run` for side-effecting commands where supported
- support `--page`, `--page-size`, and `--all` for list commands

### Future-friendly flags

These should be considered early, even if not all are implemented in v1:

- `--dry-run`
- `--wait`
- `--idempotency-key`
- `--if-exists return-existing`

## 16. CI Usage Model

The CLI must be easy to use in pipelines.

Example CI flows:

### Create course

```bash
openedx --profile admin course create \
  --org demo \
  --number cs101 \
  --run 2026 \
  --title "Intro to AI" \
  --format json
```

### Import course

```bash
openedx --profile admin course import \
  --course-id course-v1:demo+cs101+2026 \
  --file ./artifacts/course.tar.gz \
  --format json
```

### Create user and enroll

```bash
openedx --profile ops user create \
  --username alice \
  --email alice@example.com \
  --name "Alice"

openedx --profile ops enrollment add \
  --course-id course-v1:demo+cs101+2026 \
  --username alice \
  --mode audit
```

### Fetch course outline

```bash
openedx --profile ops course outline get \
  --course-id course-v1:demo+cs101+2026 \
  --format json
```

## 17. Testing Strategy

Testing should be split into four layers.

### 1. Registry contract tests

Verify:

- command definitions exist
- required argument names are stable
- output contracts remain stable

### 2. Provider behavior tests

Simulate:

- public success
- public unavailable and extension success
- public auth failure
- public validation failure
- extension timeout
- extension malformed response

### 3. Tutor integration tests

Use Tutor with a recent Open edX deployment and verify:

- token acquisition
- course list or get
- user create
- enrollment add
- role assign
- outline get

### 4. Extension compatibility tests

Use a fake extension service to verify:

- fallback triggers only on API-unavailable errors
- extension responses normalize correctly
- command output shape remains stable across providers

## 18. v1 Roadmap

The first implementation should stay narrow and high-value.

### v1 command set

- `course list`
- `course get`
- `course create`
- `course import`
- `course export`
- `course outline get`
- `user create`
- `enrollment add`
- `role assign`
- `schema`
- `doctor`

### v1 implementation priorities

1. profiles and auth
2. built-in latest public registry
3. public provider
4. extension fallback
5. normalized JSON output
6. doctor and schema

## 19. v2 Roadmap

After the command and provider model is proven, v2 can add:

- asset upload
- richer course authoring operations
- publish operations
- block create or update operations
- bulk operations
- NDJSON output
- better dry-run support
- plan or apply style workflows

## 20. Key Decisions

This design makes the following explicit decisions:

- only one built-in public API mapping set is maintained
- that mapping set aligns to the latest supported official Open edX API shape
- extension APIs are first-class provider sources, not second-class hacks
- command names are stable regardless of provider
- fallback is runtime behavior, not a separate user-facing mode
- auth and permission failures do not trigger fallback
- CI and agents are first-class users of the CLI

## 21. Open Questions

These should be validated during implementation planning:

- which target commands are fully implementable with current latest official public APIs
- whether course create should ship in v1 as public, extension-backed, or provisional
- which course outline endpoint shape is the most stable for normalization
- whether import and export should be modeled as synchronous commands or job-based commands from day one

## 22. Recommendation

Proceed with a v1 implementation that:

- fixes the command surface now
- implements the latest public provider first
- supports extension fallback from the beginning
- proves the CI workflow on Tutor-backed recent Open edX deployments

This yields a stable CLI product surface without forcing full multi-version Open edX compatibility into the first release.
