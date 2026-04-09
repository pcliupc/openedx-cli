# Open edX CLI Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go-based Open edX CLI for CI and agents that uses the latest official public API mappings by default and falls back to configured extension APIs when official endpoints are unavailable.

**Architecture:** Ship a single-binary CLI with a noun-first command tree, built-in public command registry, profile-based OAuth client credentials auth, public-first provider execution, extension fallback for API-unavailable errors, normalized JSON output, and `schema`/`doctor` diagnostics. Keep command names stable regardless of provider.

**Tech Stack:** Go, Cobra, Viper, net/http or resty, testify, golden JSON fixtures, Tutor-backed integration tests

---

## Planned File Structure

### Repo bootstrap

- Create: `go.mod`
- Create: `go.sum`
- Create: `README.md`
- Create: `.gitignore`
- Create: `Makefile`

### CLI entrypoints

- Create: `cmd/openedx/main.go`
- Create: `internal/cli/root.go`
- Create: `internal/cli/output.go`
- Create: `internal/cli/errors.go`

### Command packages

- Create: `internal/cli/cmd/auth.go`
- Create: `internal/cli/cmd/course.go`
- Create: `internal/cli/cmd/user.go`
- Create: `internal/cli/cmd/enrollment.go`
- Create: `internal/cli/cmd/role.go`
- Create: `internal/cli/cmd/schema.go`
- Create: `internal/cli/cmd/doctor.go`

### Config and auth

- Create: `internal/config/config.go`
- Create: `internal/config/types.go`
- Create: `internal/auth/token_client.go`
- Create: `internal/auth/cache.go`

### Registry and provider execution

- Create: `internal/registry/commands.go`
- Create: `internal/registry/public_registry.go`
- Create: `internal/provider/provider.go`
- Create: `internal/provider/public_provider.go`
- Create: `internal/provider/extension_provider.go`
- Create: `internal/provider/fallback.go`

### Normalization and models

- Create: `internal/model/course.go`
- Create: `internal/model/outline.go`
- Create: `internal/model/user.go`
- Create: `internal/model/enrollment.go`
- Create: `internal/model/role_assignment.go`
- Create: `internal/model/job.go`
- Create: `internal/normalize/course.go`
- Create: `internal/normalize/outline.go`
- Create: `internal/normalize/user.go`

### Diagnostics

- Create: `internal/diagnostics/schema.go`
- Create: `internal/diagnostics/doctor.go`

### Test fixtures and tests

- Create: `testdata/config/example.yaml`
- Create: `testdata/public/*.json`
- Create: `testdata/extension/*.json`
- Create: `internal/config/config_test.go`
- Create: `internal/auth/token_client_test.go`
- Create: `internal/registry/public_registry_test.go`
- Create: `internal/provider/fallback_test.go`
- Create: `internal/normalize/course_test.go`
- Create: `internal/diagnostics/schema_test.go`
- Create: `internal/diagnostics/doctor_test.go`
- Create: `integration/tutor_smoke_test.go`

### Docs

- Modify: `docs/superpowers/specs/2026-04-09-openedx-cli-design.md`
- Create: `docs/openedx-cli-config.md`
- Create: `docs/openedx-cli-commands.md`

## Implementation Assumptions

- Use Go to produce a single portable binary similar to the reference CLIs.
- Use Cobra for command structure and help generation.
- Use Viper only for config file loading and env interpolation, not as a global state container.
- Use OAuth client credentials for service profiles.
- Treat current workspace as a new project root.
- Because the current workspace is not a git repository, commit steps are specified but cannot be executed until the code is moved into a git repo or initialized as one.

## Task 1: Bootstrap the CLI repository skeleton

**Files:**
- Create: `go.mod`
- Create: `README.md`
- Create: `.gitignore`
- Create: `Makefile`
- Create: `cmd/openedx/main.go`

- [ ] **Step 1: Write the failing smoke test for CLI startup**

```go
package main

import "testing"

func TestBinaryBuilds(t *testing.T) {
	t.Fatal("main package not implemented")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./...`
Expected: FAIL with missing packages or failing placeholder test

- [ ] **Step 3: Initialize module and minimal entrypoint**

```go
package main

import "fmt"

func main() {
	fmt.Println("openedx")
}
```

- [ ] **Step 4: Add build helpers**

Add a `Makefile` with:

```makefile
test:
	go test ./...

build:
	go build ./cmd/openedx
```

- [ ] **Step 5: Run tests and build**

Run: `go test ./... && go build ./cmd/openedx`
Expected: PASS and binary builds

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum README.md .gitignore Makefile cmd/openedx/main.go
git commit -m "chore: bootstrap openedx cli skeleton"
```

## Task 2: Build the root command and JSON-first output shell

**Files:**
- Create: `internal/cli/root.go`
- Create: `internal/cli/output.go`
- Create: `internal/cli/errors.go`
- Modify: `cmd/openedx/main.go`
- Test: `internal/cli/root_test.go`

- [ ] **Step 1: Write the failing root command tests**

```go
func TestRootCommandDefaultsToJSON(t *testing.T) {}
func TestRootCommandAcceptsProfileFlag(t *testing.T) {}
```

- [ ] **Step 2: Run the root command tests**

Run: `go test ./internal/cli -v`
Expected: FAIL with missing root command package

- [ ] **Step 3: Implement the root command with persistent flags**

Required persistent flags:
- `--profile`
- `--format`
- `--config`
- `--verbose`

- [ ] **Step 4: Implement output writer helpers**

Add helpers for:
- JSON output to stdout
- table output placeholder
- structured error output

- [ ] **Step 5: Wire `main.go` to execute the root command**

```go
if err := cli.NewRootCmd().Execute(); err != nil {
	os.Exit(1)
}
```

- [ ] **Step 6: Run tests**

Run: `go test ./internal/cli -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add cmd/openedx/main.go internal/cli/root.go internal/cli/output.go internal/cli/errors.go internal/cli/root_test.go
git commit -m "feat: add root command and output shell"
```

## Task 3: Implement configuration loading for profiles and extensions

**Files:**
- Create: `internal/config/types.go`
- Create: `internal/config/config.go`
- Create: `testdata/config/example.yaml`
- Test: `internal/config/config_test.go`

- [ ] **Step 1: Write failing config parser tests**

```go
func TestLoadConfigProfiles(t *testing.T) {}
func TestLoadConfigExtensions(t *testing.T) {}
func TestSecretsComeFromEnvNames(t *testing.T) {}
```

- [ ] **Step 2: Run config tests to verify failure**

Run: `go test ./internal/config -v`
Expected: FAIL

- [ ] **Step 3: Define config structs**

Required structs:
- `Config`
- `Profile`
- `ExtensionMapping`

- [ ] **Step 4: Implement config loading from YAML**

Support:
- file path override with `--config`
- env var field names for client credentials
- basic validation

- [ ] **Step 5: Add example config fixture**

Include:
- `admin` and `ops` profiles
- `course.create`, `course.import`, `course.export` extension examples

- [ ] **Step 6: Run config tests**

Run: `go test ./internal/config -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/config/types.go internal/config/config.go testdata/config/example.yaml internal/config/config_test.go
git commit -m "feat: add profile and extension config loading"
```

## Task 4: Implement OAuth token acquisition and caching

**Files:**
- Create: `internal/auth/token_client.go`
- Create: `internal/auth/cache.go`
- Test: `internal/auth/token_client_test.go`

- [ ] **Step 1: Write failing token client tests**

```go
func TestClientCredentialsTokenRequest(t *testing.T) {}
func TestTokenCacheReturnsUnexpiredToken(t *testing.T) {}
func TestTokenCacheRefreshesExpiredToken(t *testing.T) {}
```

- [ ] **Step 2: Run auth tests**

Run: `go test ./internal/auth -v`
Expected: FAIL

- [ ] **Step 3: Implement client credentials token request**

Include:
- profile lookup
- token endpoint call
- JSON token response parsing

- [ ] **Step 4: Implement in-memory token cache**

The first version only needs process-local caching.

- [ ] **Step 5: Normalize auth errors**

Return machine-readable auth failures for:
- missing env var
- token endpoint unavailable
- invalid token response

- [ ] **Step 6: Run auth tests**

Run: `go test ./internal/auth -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/auth/token_client.go internal/auth/cache.go internal/auth/token_client_test.go
git commit -m "feat: add oauth token client and cache"
```

## Task 5: Implement the built-in latest public command registry

**Files:**
- Create: `internal/registry/commands.go`
- Create: `internal/registry/public_registry.go`
- Test: `internal/registry/public_registry_test.go`

- [ ] **Step 1: Write failing registry tests**

```go
func TestRegistryContainsV1Commands(t *testing.T) {}
func TestRegistryReturnsCourseListMapping(t *testing.T) {}
func TestRegistryReturnsCourseCreateMapping(t *testing.T) {}
```

- [ ] **Step 2: Run registry tests**

Run: `go test ./internal/registry -v`
Expected: FAIL

- [ ] **Step 3: Define command metadata types**

Include:
- command key
- HTTP method
- path template
- required args
- output model

- [ ] **Step 4: Populate the initial built-in registry**

Include entries for:
- `course.list`
- `course.get`
- `course.create`
- `course.import`
- `course.export`
- `course.outline.get`
- `user.create`
- `enrollment.add`
- `role.assign`

- [ ] **Step 5: Run registry tests**

Run: `go test ./internal/registry -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/registry/commands.go internal/registry/public_registry.go internal/registry/public_registry_test.go
git commit -m "feat: add latest public api registry"
```

## Task 6: Implement provider execution and fallback policy

**Files:**
- Create: `internal/provider/provider.go`
- Create: `internal/provider/public_provider.go`
- Create: `internal/provider/extension_provider.go`
- Create: `internal/provider/fallback.go`
- Test: `internal/provider/fallback_test.go`

- [ ] **Step 1: Write failing fallback tests**

```go
func TestFallbackUsesExtensionOnNotFound(t *testing.T) {}
func TestFallbackDoesNotUseExtensionOnForbidden(t *testing.T) {}
func TestFallbackDoesNotUseExtensionOnValidationError(t *testing.T) {}
```

- [ ] **Step 2: Run provider tests**

Run: `go test ./internal/provider -v`
Expected: FAIL

- [ ] **Step 3: Implement public provider**

Responsibilities:
- build request from registry entry
- attach bearer token
- return raw provider response and typed error

- [ ] **Step 4: Implement extension provider**

Responsibilities:
- resolve extension mapping by command key
- send request to configured URL
- return raw provider response and typed error

- [ ] **Step 5: Implement fallback policy**

Fallback should trigger only on:
- 404
- 405
- 501
- explicit endpoint-not-available errors

- [ ] **Step 6: Run provider tests**

Run: `go test ./internal/provider -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/provider/provider.go internal/provider/public_provider.go internal/provider/extension_provider.go internal/provider/fallback.go internal/provider/fallback_test.go
git commit -m "feat: add public provider and extension fallback"
```

## Task 7: Implement normalized resource models

**Files:**
- Create: `internal/model/course.go`
- Create: `internal/model/outline.go`
- Create: `internal/model/user.go`
- Create: `internal/model/enrollment.go`
- Create: `internal/model/role_assignment.go`
- Create: `internal/model/job.go`
- Create: `internal/normalize/course.go`
- Create: `internal/normalize/outline.go`
- Create: `internal/normalize/user.go`
- Test: `internal/normalize/course_test.go`
- Test: `internal/normalize/outline_test.go`
- Test: `internal/normalize/user_test.go`

- [ ] **Step 1: Write failing normalization tests**

```go
func TestNormalizeCourseFromPublicPayload(t *testing.T) {}
func TestNormalizeCourseFromExtensionPayload(t *testing.T) {}
func TestNormalizeOutlineProducesStableTree(t *testing.T) {}
```

- [ ] **Step 2: Run normalization tests**

Run: `go test ./internal/normalize -v`
Expected: FAIL

- [ ] **Step 3: Define stable model structs**

Required:
- `Course`
- `CourseOutline`
- `User`
- `Enrollment`
- `RoleAssignment`
- `Job`

- [ ] **Step 4: Implement normalizers**

Support both:
- latest public API payload shape
- extension payload shape

- [ ] **Step 5: Add fixture payloads**

Create JSON fixtures under:
- `testdata/public/`
- `testdata/extension/`

- [ ] **Step 6: Run normalization tests**

Run: `go test ./internal/normalize -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/model internal/normalize testdata/public testdata/extension
git commit -m "feat: add normalized resource models"
```

## Task 8: Implement the course command group

**Files:**
- Create: `internal/cli/cmd/course.go`
- Modify: `internal/cli/root.go`
- Test: `internal/cli/cmd/course_test.go`

- [ ] **Step 1: Write failing course command tests**

```go
func TestCourseListCallsRegistryCommand(t *testing.T) {}
func TestCourseGetRequiresCourseID(t *testing.T) {}
func TestCourseCreateSupportsOrgNumberRunTitle(t *testing.T) {}
func TestCourseOutlineGetFormatsStableJSON(t *testing.T) {}
```

- [ ] **Step 2: Run course command tests**

Run: `go test ./internal/cli/cmd -run TestCourse -v`
Expected: FAIL

- [ ] **Step 3: Implement course subcommands**

Include:
- `course list`
- `course get`
- `course create`
- `course import`
- `course export`
- `course outline get`

- [ ] **Step 4: Wire command handlers to provider + normalizer**

Each command should:
- parse flags
- resolve command key
- call execution layer
- print normalized output

- [ ] **Step 5: Run course command tests**

Run: `go test ./internal/cli/cmd -run TestCourse -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/cli/cmd/course.go internal/cli/root.go internal/cli/cmd/course_test.go
git commit -m "feat: add course commands"
```

## Task 9: Implement user, enrollment, and role commands

**Files:**
- Create: `internal/cli/cmd/user.go`
- Create: `internal/cli/cmd/enrollment.go`
- Create: `internal/cli/cmd/role.go`
- Modify: `internal/cli/root.go`
- Test: `internal/cli/cmd/user_test.go`
- Test: `internal/cli/cmd/enrollment_test.go`
- Test: `internal/cli/cmd/role_test.go`

- [ ] **Step 1: Write failing tests for user and access commands**

```go
func TestUserCreateRequiresUsernameAndEmail(t *testing.T) {}
func TestEnrollmentAddRequiresCourseAndUsername(t *testing.T) {}
func TestRoleAssignRequiresCourseUsernameRole(t *testing.T) {}
```

- [ ] **Step 2: Run command tests**

Run: `go test ./internal/cli/cmd -run 'TestUser|TestEnrollment|TestRole' -v`
Expected: FAIL

- [ ] **Step 3: Implement subcommands**

Include:
- `user create`
- `enrollment add`
- `role assign`

- [ ] **Step 4: Reuse output and error helpers**

All command errors should return structured CLI errors.

- [ ] **Step 5: Run tests**

Run: `go test ./internal/cli/cmd -run 'TestUser|TestEnrollment|TestRole' -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/cli/cmd/user.go internal/cli/cmd/enrollment.go internal/cli/cmd/role.go internal/cli/root.go internal/cli/cmd/user_test.go internal/cli/cmd/enrollment_test.go internal/cli/cmd/role_test.go
git commit -m "feat: add user enrollment and role commands"
```

## Task 10: Implement `schema` command for command visibility

**Files:**
- Create: `internal/cli/cmd/schema.go`
- Create: `internal/diagnostics/schema.go`
- Test: `internal/diagnostics/schema_test.go`

- [ ] **Step 1: Write failing schema tests**

```go
func TestSchemaShowsPublicMapping(t *testing.T) {}
func TestSchemaShowsExtensionWhenConfigured(t *testing.T) {}
func TestSchemaAllListsV1Commands(t *testing.T) {}
```

- [ ] **Step 2: Run schema tests**

Run: `go test ./internal/diagnostics -run TestSchema -v`
Expected: FAIL

- [ ] **Step 3: Implement schema model and command**

Output should include:
- command key
- public method/path
- extension presence
- required args
- output model

- [ ] **Step 4: Run schema tests**

Run: `go test ./internal/diagnostics -run TestSchema -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/cli/cmd/schema.go internal/diagnostics/schema.go internal/diagnostics/schema_test.go
git commit -m "feat: add schema inspection command"
```

## Task 11: Implement `doctor` command for health checks

**Files:**
- Create: `internal/cli/cmd/doctor.go`
- Create: `internal/diagnostics/doctor.go`
- Test: `internal/diagnostics/doctor_test.go`

- [ ] **Step 1: Write failing doctor tests**

```go
func TestDoctorChecksBaseURL(t *testing.T) {}
func TestDoctorChecksTokenAcquisition(t *testing.T) {}
func TestDoctorVerifyCommandChecksExtensionAvailability(t *testing.T) {}
```

- [ ] **Step 2: Run doctor tests**

Run: `go test ./internal/diagnostics -run TestDoctor -v`
Expected: FAIL

- [ ] **Step 3: Implement baseline doctor checks**

Include:
- base URL reachable
- token request succeeds
- verify command health

- [ ] **Step 4: Implement `doctor verify <command>`**

The first version only needs to:
- validate public mapping exists
- validate extension mapping exists if configured
- optionally hit a lightweight endpoint check

- [ ] **Step 5: Run doctor tests**

Run: `go test ./internal/diagnostics -run TestDoctor -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/cli/cmd/doctor.go internal/diagnostics/doctor.go internal/diagnostics/doctor_test.go
git commit -m "feat: add doctor diagnostics"
```

## Task 12: Add Tutor-backed integration coverage

**Files:**
- Create: `integration/tutor_smoke_test.go`
- Modify: `Makefile`
- Create: `docs/openedx-cli-config.md`
- Create: `docs/openedx-cli-commands.md`

- [ ] **Step 1: Write the failing integration test shell**

```go
func TestTutorSmokeCourseList(t *testing.T) {}
func TestTutorSmokeUserCreate(t *testing.T) {}
```

- [ ] **Step 2: Mark integration tests with an explicit gate**

Use an env var such as:
- `OPENEDX_INTEGRATION=1`

- [ ] **Step 3: Implement Tutor smoke scenarios**

Cover:
- auth
- course list or get
- user create
- enrollment add

- [ ] **Step 4: Add integration test target**

Extend `Makefile`:

```makefile
test-integration:
	OPENEDX_INTEGRATION=1 go test ./integration -v
```

- [ ] **Step 5: Document local and CI setup**

Document:
- config file format
- required env vars
- sample commands
- integration test prerequisites

- [ ] **Step 6: Run unit tests**

Run: `go test ./...`
Expected: PASS for unit tests

- [ ] **Step 7: Run integration tests in a Tutor-enabled environment**

Run: `OPENEDX_INTEGRATION=1 go test ./integration -v`
Expected: PASS against a prepared Tutor instance

- [ ] **Step 8: Commit**

```bash
git add integration/tutor_smoke_test.go Makefile docs/openedx-cli-config.md docs/openedx-cli-commands.md
git commit -m "test: add tutor integration coverage and usage docs"
```

## Task 13: Final polish before execution handoff

**Files:**
- Modify: `README.md`
- Modify: `docs/superpowers/specs/2026-04-09-openedx-cli-design.md`
- Modify: `docs/superpowers/plans/2026-04-09-openedx-cli-implementation-plan.md`

- [ ] **Step 1: Update README with supported commands**

Include:
- v1 scope
- config example
- profile usage
- extension fallback summary

- [ ] **Step 2: Reconcile implementation notes with the spec**

Add any confirmed deviations discovered during implementation.

- [ ] **Step 3: Run full unit test suite**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add README.md docs/superpowers/specs/2026-04-09-openedx-cli-design.md docs/superpowers/plans/2026-04-09-openedx-cli-implementation-plan.md
git commit -m "docs: finalize openedx cli v1 implementation plan and docs"
```

## Execution Notes

- Keep tasks small and land them in order.
- Do not implement extension-only commands until the public provider and fallback policy are working.
- Prefer exact fixture-driven normalization tests over broad integration-first development.
- Do not hide 401, 403, or 400 errors behind extension fallback.
- Keep command names stable even if backend provider changes.

## Review Constraints

- The current workspace is not a git repository, so commit steps are part of the plan but cannot be executed here until a repo exists.
- A formal plan-review subagent loop was not run in this session because subagent delegation was not explicitly requested.

## Recommended Execution Order

1. Bootstrap repo and root CLI
2. Config and auth
3. Public registry
4. Provider fallback
5. Normalization
6. Course commands
7. User, enrollment, role commands
8. Schema and doctor
9. Integration tests and docs
