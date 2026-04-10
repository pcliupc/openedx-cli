# Command Expansion Design

Date: 2026-04-10
Status: Draft for review

## 1. Purpose

Two deliverables:

1. **Extension guide** — a document explaining how to add new commands (for contributors) and configure extension APIs (for admins), referenced from CLAUDE.md
2. **New commands** — additional CLI commands based on official Open edX public APIs

## 2. Extension Guide

Written to `docs/extending-the-cli.md`, covering:

- For contributors: 5-step process (registry → model → normalizer → cmd → root), naming conventions, path parameter handling, test patterns
- For admins: extension config syntax, fallback behavior, diagnostic commands, endpoint requirements

Referenced from CLAUDE.md under the Configuration section.

## 3. New Commands (Batch 1)

7 new commands covering high-frequency platform management operations:

| Command Key | Method | API Path | Required Args | Output Model |
|---|---|---|---|---|
| `enrollment.list` | GET | `/api/enrollment/v1/enrollments` | (none) | Enrollment |
| `enrollment.remove` | POST | `/api/enrollment/v1/enrollments` | `course_id`, `username` | Enrollment |
| `user.list` | GET | `/api/user/v1/accounts` | (none) | User |
| `user.get` | GET | `/api/user/v1/accounts/{username}` | `username` | User |
| `grade.list` | GET | `/api/grades/v1/courses/{course_id}/` | `course_id` | Grade |
| `gradebook.get` | GET | `/api/grades/v1/gradebook/{course_id}/` | `course_id` | Gradebook |
| `certificate.list` | GET | `/api/certificates/v0/certificates/{username}/` | `username` | Certificate |

### Changes per command

Each command requires:

1. Registry entry in `internal/registry/public_registry.go`
2. Model struct in `internal/model/` (new file for new domains: grade.go, certificate.go)
3. Normalizer in `internal/normalize/` (new file for new domains)
4. Cobra command in `internal/cli/cmd/` (extend enrollment.go, user.go; new grade.go, certificate.go)
5. Registration in `internal/cli/root.go`

### New models

**Grade** (`internal/model/grade.go`):

```go
type Grade struct {
    Username  string  `json:"username"`
    CourseID  string  `json:"course_id"`
    Percent   float64 `json:"percent"`
    LetterGrade string `json:"letter_grade,omitempty"`
    Passed    bool    `json:"passed"`
    Section   string  `json:"section,omitempty"`
}
```

**Gradebook** (`internal/model/grade.go`):

```go
type Gradebook struct {
    CourseID string  `json:"course_id"`
    Grades   []Grade `json:"grades"`
}
```

**Certificate** (`internal/model/certificate.go`):

```go
type Certificate struct {
    Username   string `json:"username"`
    CourseID   string `json:"course_id"`
    Type       string `json:"certificate_type"`
    Status     string `json:"status"`
    DownloadURL string `json:"download_url,omitempty"`
    Grade      string `json:"grade,omitempty"`
}
```

### Enrollment enhancements

`enrollment.list` — supports optional filters: `--course-id`, `--username`, `--page`, `--page-size`

`enrollment.remove` — sends POST with `is_active: false` to deactivate enrollment

### User enhancements

`user.list` — supports optional filters: `--page`, `--page-size`

`user.get` — path parameter `username`

## 4. New Commands (Batch 2 — Future)

Documented in the extension guide as examples. Not implemented now.

| Command Key | Method | API Path |
|---|---|---|
| `cohort.list` | GET | `/api/cohorts/v1/courses/{course_key}/cohorts/` |
| `cohort.add` | POST | `/api/cohorts/v1/courses/{course_key}/cohorts/` |
| `org.list` | GET | `/api/organizations/v0/organizations/` |
| `block.list` | GET | `/api/courses/v2/blocks/` |
| `discussion.list` | GET | `/api/discussion/v1/courses/{course_id}` |
| `team.list` | GET | `/api/team/v0/teams/` |

## 5. Implementation Order

1. Write `docs/extending-the-cli.md` (done)
2. Update CLAUDE.md (done)
3. Add new model structs (grade, certificate)
4. Add registry entries for all 7 commands
5. Add normalizers with test fixtures
6. Add/extend Cobra commands (enrollment.go, user.go, new grade.go, new certificate.go)
7. Register in root.go
8. Write tests
9. Update schema/doctor if needed
