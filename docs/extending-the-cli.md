# Extending the Open edX CLI

This guide explains how to add new commands to the Open edX CLI and how to configure extension APIs for custom endpoints.

## Audience

- **Contributors** — developers adding new commands to the CLI source code
- **Platform admins / DevOps** — operators configuring extension mappings in YAML to connect custom or backported APIs

---

## Part 1: For Contributors — Adding a New Command

Adding a command touches five areas of the codebase, always in this order:

1. **Registry** — define the API endpoint metadata
2. **Model** — define the output resource struct
3. **Normalizer** — convert raw API response to model
4. **Command** — build the Cobra command with flags
5. **Root** — register the command group

### Step-by-step Example: Adding `grade list`

#### 1. Add a registry entry

File: `internal/registry/public_registry.go`

```go
"grade.list": {
    Key:          "grade.list",
    Method:       "GET",
    Path:         "/api/grades/v1/courses/{course_id}/",
    RequiredArgs: []string{"course_id"},
    OutputModel:  "Grade",
},
```

Fields:

- `Key` — domain.verb identifier used everywhere (config, extensions, logging)
- `Method` — HTTP method
- `Path` — URL path template with `{param}` placeholders for path parameters
- `RequiredArgs` — argument names that must be provided
- `OutputModel` — name of the model struct for documentation

#### 2. Define the model

File: `internal/model/grade.go`

```go
package model

type Grade struct {
    Username string  `json:"username"`
    CourseID string  `json:"course_id"`
    Percent  float64 `json:"percent"`
    Letter   string  `json:"letter_grade,omitempty"`
    Passed   bool    `json:"passed"`
}
```

Keep models flat and stable. Add `omitempty` for optional fields. These structs define the JSON output contract.

#### 3. Write the normalizer

File: `internal/normalize/grade.go`

```go
package normalize

import "github.com/openedx/cli/internal/model"

func Grades(raw []byte) ([]model.Grade, error) {
    // Parse raw JSON from provider into []model.Grade.
    // Both public and extension payloads should produce the same output.
}
```

Write a test with fixture files under `testdata/public/` and `testdata/extension/`.

#### 4. Create the Cobra command

File: `internal/cli/cmd/grade.go`

Follow the pattern from existing command files (e.g. `course.go`):

```go
func NewGradeCmd(execFn ExecuteFunc) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "grade",
        Short: "Manage course grades",
    }

    cmd.AddCommand(newGradeListCmd(execFn))
    return cmd
}

func newGradeListCmd(execFn ExecuteFunc) *cobra.Command {
    var courseID string

    cmd := &cobra.Command{
        Use:   "list",
        Short: "List grades for a course",
        RunE: func(cmd *cobra.Command, args []string) error {
            result, err := execFn(cmd.Context(), "grade.list", map[string]string{
                "course_id": courseID,
            })
            // ... output handling
        },
    }

    cmd.Flags().StringVar(&courseID, "course-id", "", "Course ID (required)")
    cmd.MarkFlagRequired("course-id")

    return cmd
}
```

Key conventions:

- Noun-first: `grade list`, not `list grades`
- Long flags: `--course-id`, not `-c`
- Use `ExecuteFunc` for testability (injected in root.go)
- Required args get `MarkFlagRequired`
- Optional filter/pagination args: `--username`, `--page`, `--page-size`

#### 5. Register in root

File: `internal/cli/root.go`

Add the new command to the `rootCmd.AddCommand(...)` block:

```go
rootCmd.AddCommand(
    cmd.NewCourseCmd(execFn),
    cmd.NewUserCmd(execFn),
    cmd.NewEnrollmentCmd(execFn),
    cmd.NewRoleCmd(execFn),
    cmd.NewGradeCmd(execFn),    // add this
    // ...
)
```

#### 6. Write tests

- `internal/normalize/grade_test.go` — test with golden fixtures
- `internal/cli/cmd/grade_test.go` — test flag parsing and execFn invocation
- Provider-level tests are only needed if the command uses special HTTP behavior

### Naming Conventions

- Command keys use `domain.verb`: `grade.list`, `enrollment.remove`, `certificate.list`
- Domains align with Open edX API prefixes: grade → `/api/grades/`, certificate → `/api/certificates/`
- Verbs are consistent: `list`, `get`, `create`, `add`, `remove`, `assign`
- CLI surface: `openedx grade list`, `openedx enrollment remove`

### Path Parameter Handling

Path templates use `{param}` syntax. The `resolvePath` function replaces placeholders with values from the args map. Remaining args become query parameters for GET requests, or JSON body for POST requests.

Example for `GET /api/grades/v1/courses/{course_id}/?username=alice`:

```
Path: /api/grades/v1/courses/{course_id}/
Args: {course_id: "course-v1:Org+Num+Run", username: "alice"}
→ resolves to: /api/grades/v1/courses/course-v1:Org+Num+Run/?username=alice
```

### When to Add vs. Extend

- New domain (grade, certificate, org) → new cmd file, new model, new normalizer
- New verb on existing domain (enrollment.list alongside enrollment.add) → extend existing cmd file

---

## Part 2: For Platform Admins — Configuring Extension APIs

Extensions let you map CLI commands to custom API endpoints without modifying CLI source code. This is useful when:

- Your Open edX version lacks a public API that newer versions have
- You have custom internal APIs that serve the same purpose
- You backported an API to a different URL path

### Configuration Syntax

Add extension mappings under the `extensions` key in your config file (`openedx.yaml` or `~/.openedx/config.yaml`):

```yaml
version: 1

profiles:
  admin:
    base_url: https://openedx.example.com
    token_url: https://openedx.example.com/oauth2/access_token
    client_id_env: OPENEDX_ADMIN_CLIENT_ID
    client_secret_env: OPENEDX_ADMIN_CLIENT_SECRET

extensions:
  grade.list:
    method: GET
    url: https://openedx.example.com/api/custom/v1/grades
  certificate.list:
    method: GET
    url: https://openedx.example.com/api/custom/v1/certificates
  course.create:
    method: POST
    url: https://openedx.example.com/api/cli-ext/course/create
```

Each extension mapping has:

- `method` — HTTP method (GET, POST, etc.)
- `url` — full URL to the custom endpoint

### How Fallback Works

For every command:

1. CLI tries the built-in public API first
2. If the public API returns **404 Not Found**, **405 Method Not Allowed**, or **501 Not Implemented**, and an extension mapping exists for that command key, CLI retries through the extension endpoint
3. If no extension mapping exists, the original error is returned

**These errors NEVER trigger fallback:**

- `400 Bad Request` — your input is wrong
- `401 Unauthorized` — credentials are wrong
- `403 Forbidden` — your profile lacks permission

This ensures auth and validation issues are surfaced clearly, not hidden behind a provider switch.

### Verifying Your Configuration

Use built-in diagnostic commands:

```bash
# See all commands and whether extensions are configured
openedx schema all

# Check a specific command
openedx schema grade list

# Run full health check
openedx doctor

# Verify a specific command's endpoint
openedx doctor verify grade.list
```

### Extension Endpoint Requirements

Your extension endpoint must:

- Accept the same HTTP method as configured
- Accept a Bearer token in the `Authorization` header
- For GET requests: accept args as URL query parameters
- For POST requests: accept args as a JSON body (`Content-Type: application/json`)
- Return JSON responses with standard HTTP status codes
- Return 2xx on success, 4xx/5xx on failure

The CLI does not normalize extension responses differently — it passes them through the same normalizer pipeline. For best results, structure extension responses to match the public API response shape documented in the CLI models.

---

## Summary

| Task | Who | Where |
|------|-----|-------|
| Add a new command | Contributor | registry → model → normalizer → cmd → root |
| Add a verb to existing domain | Contributor | extend existing cmd file, update registry |
| Connect a custom API | Admin | `extensions:` block in config YAML |
| Check what's configured | Anyone | `openedx schema` |
| Verify endpoint health | Anyone | `openedx doctor` |
