# Command Expansion Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add 7 new commands (enrollment.list, enrollment.remove, user.list, user.get, grade.list, gradebook.get, certificate.list) to the Open edX CLI with full model, normalizer, Cobra command, test, and fixture support.

**Architecture:** Follow the existing 5-layer pattern: registry entry → model struct → normalizer with dual-source JSON parsing → Cobra command with flags → root.go registration. Extend existing files for enrollment and user domains; create new files for grade and certificate domains.

**Tech Stack:** Go, Cobra, testify, JSON fixtures under testdata/

---

## Planned File Structure

### New files
- `internal/model/grade.go` — Grade and Gradebook structs
- `internal/model/certificate.go` — Certificate struct
- `internal/normalize/grade.go` — grade/gradebook normalizers
- `internal/normalize/grade_test.go` — grade normalizer tests
- `internal/normalize/certificate.go` — certificate normalizer
- `internal/normalize/certificate_test.go` — certificate normalizer tests
- `internal/normalize/enrollment.go` — enrollment list normalizer
- `internal/normalize/enrollment_test.go` — enrollment normalizer tests
- `internal/cli/cmd/grade.go` — grade and gradebook commands
- `internal/cli/cmd/grade_test.go` — grade command tests
- `internal/cli/cmd/certificate.go` — certificate command
- `internal/cli/cmd/certificate_test.go` — certificate command tests
- `testdata/public/grade_list.json` — public grades fixture
- `testdata/public/gradebook_get.json` — public gradebook fixture
- `testdata/public/certificate_list.json` — public certificates fixture
- `testdata/public/enrollment_list.json` — public enrollment list fixture
- `testdata/public/user_list.json` — public user list fixture
- `testdata/public/user_get.json` — public user get fixture
- `testdata/extension/grade_list.json` — extension grades fixture
- `testdata/extension/certificate_list.json` — extension certificates fixture

### Modified files
- `internal/registry/public_registry.go` — add 7 new registry entries
- `internal/cli/cmd/enrollment.go` — add list and remove subcommands
- `internal/cli/cmd/user.go` — add list and get subcommands
- `internal/cli/root.go` — register grade and certificate command groups
- `internal/cli/cmd/enrollment_test.go` — add tests for list and remove
- `internal/cli/cmd/user_test.go` — add tests for list and get

---

## Task 1: Add Grade and Certificate model structs

**Files:**
- Create: `internal/model/grade.go`
- Create: `internal/model/certificate.go`

- [ ] **Step 1: Create `internal/model/grade.go`**

```go
package model

// Grade represents a student's grade in a course section.
type Grade struct {
	Username    string  `json:"username"`
	CourseID    string  `json:"course_id"`
	Percent     float64 `json:"percent"`
	LetterGrade string  `json:"letter_grade,omitempty"`
	Passed      bool    `json:"passed"`
	Section     string  `json:"section,omitempty"`
}

// Gradebook represents the full gradebook for a course.
type Gradebook struct {
	CourseID string  `json:"course_id"`
	Grades   []Grade `json:"grades"`
}
```

- [ ] **Step 2: Create `internal/model/certificate.go`**

```go
package model

// Certificate represents a user's course completion certificate.
type Certificate struct {
	Username    string `json:"username"`
	CourseID    string `json:"course_id"`
	CertificateType string `json:"certificate_type"`
	Status      string `json:"status"`
	DownloadURL string `json:"download_url,omitempty"`
	Grade       string `json:"grade,omitempty"`
}
```

- [ ] **Step 3: Run tests to verify no breakage**

Run: `go build ./...`
Expected: PASS (models have no dependencies)

- [ ] **Step 4: Commit**

```bash
git add internal/model/grade.go internal/model/certificate.go
git commit -m "feat: add Grade, Gradebook, and Certificate model structs"
```

---

## Task 2: Add registry entries for all 7 new commands

**Files:**
- Modify: `internal/registry/public_registry.go`

- [ ] **Step 1: Add 7 new entries to the LatestRegistry map**

Append these entries inside the `map[string]CommandMeta{` literal in `internal/registry/public_registry.go`, after the existing `"role.assign"` entry:

```go
		"enrollment.list": {
			Key:          "enrollment.list",
			Method:       "GET",
			Path:         "/api/enrollment/v1/enrollments",
			RequiredArgs: []string{},
			OutputModel:  "Enrollment",
		},
		"enrollment.remove": {
			Key:          "enrollment.remove",
			Method:       "POST",
			Path:         "/api/enrollment/v1/enrollments",
			RequiredArgs: []string{"course_id", "username"},
			OutputModel:  "Enrollment",
		},
		"user.list": {
			Key:          "user.list",
			Method:       "GET",
			Path:         "/api/user/v1/accounts",
			RequiredArgs: []string{},
			OutputModel:  "User",
		},
		"user.get": {
			Key:          "user.get",
			Method:       "GET",
			Path:         "/api/user/v1/accounts/{username}",
			RequiredArgs: []string{"username"},
			OutputModel:  "User",
		},
		"grade.list": {
			Key:          "grade.list",
			Method:       "GET",
			Path:         "/api/grades/v1/courses/{course_id}/",
			RequiredArgs: []string{"course_id"},
			OutputModel:  "Grade",
		},
		"gradebook.get": {
			Key:          "gradebook.get",
			Method:       "GET",
			Path:         "/api/grades/v1/gradebook/{course_id}/",
			RequiredArgs: []string{"course_id"},
			OutputModel:  "Gradebook",
		},
		"certificate.list": {
			Key:          "certificate.list",
			Method:       "GET",
			Path:         "/api/certificates/v0/certificates/{username}/",
			RequiredArgs: []string{"username"},
			OutputModel:  "Certificate",
		},
```

- [ ] **Step 2: Run registry tests**

Run: `go test ./internal/registry -v`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add internal/registry/public_registry.go
git commit -m "feat: add registry entries for enrollment, user, grade, gradebook, certificate commands"
```

---

## Task 3: Create test fixtures for new commands

**Files:**
- Create: `testdata/public/grade_list.json`
- Create: `testdata/public/gradebook_get.json`
- Create: `testdata/public/certificate_list.json`
- Create: `testdata/public/enrollment_list.json`
- Create: `testdata/public/user_list.json`
- Create: `testdata/public/user_get.json`
- Create: `testdata/extension/grade_list.json`
- Create: `testdata/extension/certificate_list.json`

- [ ] **Step 1: Create `testdata/public/grade_list.json`**

Public API returns paginated results:

```json
{
  "results": [
    {
      "username": "alice",
      "course_id": "course-v1:demo+cs101+2026",
      "percent": 0.85,
      "letter_grade": "B+",
      "passed": true,
      "section": "Week 1"
    },
    {
      "username": "bob",
      "course_id": "course-v1:demo+cs101+2026",
      "percent": 0.92,
      "letter_grade": "A-",
      "passed": true,
      "section": "Week 1"
    }
  ]
}
```

- [ ] **Step 2: Create `testdata/public/gradebook_get.json`**

```json
{
  "course_id": "course-v1:demo+cs101+2026",
  "grades": [
    {
      "username": "alice",
      "course_id": "course-v1:demo+cs101+2026",
      "percent": 0.85,
      "letter_grade": "B+",
      "passed": true,
      "section": "Week 1"
    },
    {
      "username": "bob",
      "course_id": "course-v1:demo+cs101+2026",
      "percent": 0.92,
      "letter_grade": "A-",
      "passed": true,
      "section": "Week 1"
    }
  ]
}
```

- [ ] **Step 3: Create `testdata/public/certificate_list.json`**

Public API returns paginated results:

```json
{
  "results": [
    {
      "username": "alice",
      "course_id": "course-v1:demo+cs101+2026",
      "certificate_type": "verified",
      "status": "downloadable",
      "download_url": "https://openedx.example.com/certificates/abc123",
      "grade": "85%"
    }
  ]
}
```

- [ ] **Step 4: Create `testdata/public/enrollment_list.json`**

```json
{
  "results": [
    {
      "username": "alice",
      "course_id": "course-v1:demo+cs101+2026",
      "mode": "audit",
      "is_active": true
    },
    {
      "username": "bob",
      "course_id": "course-v1:demo+cs101+2026",
      "mode": "verified",
      "is_active": true
    }
  ]
}
```

- [ ] **Step 5: Create `testdata/public/user_list.json`**

```json
{
  "results": [
    {
      "username": "alice",
      "email": "alice@example.com",
      "name": "Alice Smith",
      "is_active": true,
      "date_joined": "2026-01-15T10:30:00Z"
    },
    {
      "username": "bob",
      "email": "bob@example.com",
      "name": "Bob Jones",
      "is_active": true,
      "date_joined": "2026-02-01T08:00:00Z"
    }
  ]
}
```

- [ ] **Step 6: Create `testdata/public/user_get.json`**

```json
{
  "username": "alice",
  "email": "alice@example.com",
  "name": "Alice Smith",
  "is_active": true,
  "date_joined": "2026-01-15T10:30:00Z"
}
```

- [ ] **Step 7: Create `testdata/extension/grade_list.json`**

Extension API returns a top-level array:

```json
[
  {
    "username": "alice",
    "course_id": "course-v1:demo+cs101+2026",
    "percent": 0.85,
    "letter_grade": "B+",
    "passed": true
  }
]
```

- [ ] **Step 8: Create `testdata/extension/certificate_list.json`**

Extension API returns a top-level array:

```json
[
  {
    "username": "alice",
    "course_id": "course-v1:demo+cs101+2026",
    "certificate_type": "verified",
    "status": "downloadable"
  }
]
```

- [ ] **Step 9: Commit**

```bash
git add testdata/public/grade_list.json testdata/public/gradebook_get.json testdata/public/certificate_list.json testdata/public/enrollment_list.json testdata/public/user_list.json testdata/public/user_get.json testdata/extension/grade_list.json testdata/extension/certificate_list.json
git commit -m "test: add fixtures for grade, gradebook, certificate, enrollment, user commands"
```

---

## Task 4: Add normalizers for enrollment, grade, and certificate

**Files:**
- Create: `internal/normalize/enrollment.go`
- Create: `internal/normalize/enrollment_test.go`
- Create: `internal/normalize/grade.go`
- Create: `internal/normalize/grade_test.go`
- Create: `internal/normalize/certificate.go`
- Create: `internal/normalize/certificate_test.go`

- [ ] **Step 1: Create `internal/normalize/enrollment.go`**

```go
package normalize

import (
	"encoding/json"

	"github.com/openedx/cli/internal/model"
)

// rawEnrollment is a lax intermediate struct that accepts field names from both
// public and extension API payloads.
type rawEnrollment struct {
	Username string `json:"username"`
	CourseID string `json:"course_id"`
	User     string `json:"user"`
	Course   string `json:"course"`
	Mode     string `json:"mode"`
	IsActive *bool  `json:"is_active"`
	Active   *bool  `json:"active"`
}

// rawEnrollmentList wraps the public API paginated response shape.
type rawEnrollmentList struct {
	Results []rawEnrollment `json:"results"`
}

func (r rawEnrollment) toModel() *model.Enrollment {
	courseID := firstNonEmpty(r.CourseID, r.Course)
	username := firstNonEmpty(r.Username, r.User)
	active := resolveActive(r.IsActive, r.Active)
	return &model.Enrollment{
		Username: username,
		CourseID: courseID,
		Mode:     r.Mode,
		IsActive: active,
	}
}

// EnrollmentListFromJSON parses a list of enrollments from raw JSON bytes.
// It supports the public API paginated response {"results": [...]} and the
// extension API top-level array [...].
func EnrollmentListFromJSON(data []byte) ([]*model.Enrollment, error) {
	// Try paginated public API shape first.
	var paginated rawEnrollmentList
	if err := json.Unmarshal(data, &paginated); err == nil && len(paginated.Results) > 0 {
		out := make([]*model.Enrollment, len(paginated.Results))
		for i, r := range paginated.Results {
			out[i] = r.toModel()
		}
		return out, nil
	}

	// Try top-level array (extension API shape).
	var items []rawEnrollment
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	out := make([]*model.Enrollment, len(items))
	for i, r := range items {
		out[i] = r.toModel()
	}
	return out, nil
}

// EnrollmentFromJSON parses a single enrollment from raw JSON bytes.
func EnrollmentFromJSON(data []byte) (*model.Enrollment, error) {
	var raw rawEnrollment
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return raw.toModel(), nil
}
```

- [ ] **Step 2: Create `internal/normalize/enrollment_test.go`**

```go
package normalize

import (
	"os"
	"testing"

	"github.com/openedx/cli/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnrollmentListFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/enrollment_list.json")
	require.NoError(t, err)

	enrollments, err := EnrollmentListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, enrollments, 2)
	expected := &model.Enrollment{
		Username: "alice",
		CourseID: "course-v1:demo+cs101+2026",
		Mode:     "audit",
		IsActive: true,
	}
	assert.Equal(t, expected, enrollments[0])
}

func TestEnrollmentFromJSONInvalidInput(t *testing.T) {
	_, err := EnrollmentListFromJSON([]byte("not json"))
	assert.Error(t, err)
}

func TestEnrollmentFromJSONEmptyObject(t *testing.T) {
	result, err := EnrollmentFromJSON([]byte("{}"))
	require.NoError(t, err)
	assert.Equal(t, &model.Enrollment{IsActive: true}, result)
}
```

- [ ] **Step 3: Create `internal/normalize/grade.go`**

```go
package normalize

import (
	"encoding/json"

	"github.com/openedx/cli/internal/model"
)

// rawGrade is a lax intermediate struct that accepts field names from both
// public and extension API payloads.
type rawGrade struct {
	Username    string  `json:"username"`
	CourseID    string  `json:"course_id"`
	Percent     float64 `json:"percent"`
	LetterGrade string  `json:"letter_grade"`
	Grade       string  `json:"grade"`
	Passed      *bool   `json:"passed"`
	Section     string  `json:"section"`
	Subsection  string  `json:"subsection"`
}

// rawGradeList wraps the public API paginated response shape.
type rawGradeList struct {
	Results []rawGrade `json:"results"`
}

func (r rawGrade) toModel() *model.Grade {
	letter := firstNonEmpty(r.LetterGrade, r.Grade)
	passed := true
	if r.Passed != nil {
		passed = *r.Passed
	}
	section := firstNonEmpty(r.Section, r.Subsection)
	return &model.Grade{
		Username:    r.Username,
		CourseID:    r.CourseID,
		Percent:     r.Percent,
		LetterGrade: letter,
		Passed:      passed,
		Section:     section,
	}
}

// GradeListFromJSON parses a list of grades from raw JSON bytes.
// It supports the public API paginated response {"results": [...]} and the
// extension API top-level array [...].
func GradeListFromJSON(data []byte) ([]*model.Grade, error) {
	// Try paginated public API shape first.
	var paginated rawGradeList
	if err := json.Unmarshal(data, &paginated); err == nil && len(paginated.Results) > 0 {
		out := make([]*model.Grade, len(paginated.Results))
		for i, r := range paginated.Results {
			out[i] = r.toModel()
		}
		return out, nil
	}

	// Try top-level array (extension API shape).
	var items []rawGrade
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	out := make([]*model.Grade, len(items))
	for i, r := range items {
		out[i] = r.toModel()
	}
	return out, nil
}

// rawGradebook is a lax intermediate struct for the gradebook response.
type rawGradebook struct {
	CourseID    string     `json:"course_id"`
	CourseKey   string     `json:"course_key"`
	Grades      []rawGrade `json:"grades"`
}

// GradebookFromJSON parses a gradebook from raw JSON bytes.
func GradebookFromJSON(data []byte) (*model.Gradebook, error) {
	var raw rawGradebook
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	courseID := firstNonEmpty(raw.CourseID, raw.CourseKey)
	grades := make([]model.Grade, len(raw.Grades))
	for i, r := range raw.Grades {
		grades[i] = *r.toModel()
	}

	return &model.Gradebook{
		CourseID: courseID,
		Grades:   grades,
	}, nil
}
```

- [ ] **Step 4: Create `internal/normalize/grade_test.go`**

```go
package normalize

import (
	"os"
	"testing"

	"github.com/openedx/cli/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGradeListFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/grade_list.json")
	require.NoError(t, err)

	grades, err := GradeListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, grades, 2)
	expected := &model.Grade{
		Username:    "alice",
		CourseID:    "course-v1:demo+cs101+2026",
		Percent:     0.85,
		LetterGrade: "B+",
		Passed:      true,
		Section:     "Week 1",
	}
	assert.Equal(t, expected, grades[0])
}

func TestGradeListFromExtensionPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/extension/grade_list.json")
	require.NoError(t, err)

	grades, err := GradeListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, grades, 1)
	expected := &model.Grade{
		Username:    "alice",
		CourseID:    "course-v1:demo+cs101+2026",
		Percent:     0.85,
		LetterGrade: "B+",
		Passed:      true,
	}
	assert.Equal(t, expected, grades[0])
}

func TestGradebookFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/gradebook_get.json")
	require.NoError(t, err)

	gradebook, err := GradebookFromJSON(data)
	require.NoError(t, err)

	assert.Equal(t, "course-v1:demo+cs101+2026", gradebook.CourseID)
	require.Len(t, gradebook.Grades, 2)
	assert.Equal(t, "alice", gradebook.Grades[0].Username)
	assert.Equal(t, 0.85, gradebook.Grades[0].Percent)
}

func TestGradeListFromJSONInvalidInput(t *testing.T) {
	_, err := GradeListFromJSON([]byte("not json"))
	assert.Error(t, err)
}

func TestGradebookFromJSONInvalidInput(t *testing.T) {
	_, err := GradebookFromJSON([]byte("not json"))
	assert.Error(t, err)
}
```

- [ ] **Step 5: Create `internal/normalize/certificate.go`**

```go
package normalize

import (
	"encoding/json"

	"github.com/openedx/cli/internal/model"
)

// rawCertificate is a lax intermediate struct that accepts field names from both
// public and extension API payloads.
type rawCertificate struct {
	Username        string `json:"username"`
	CourseID        string `json:"course_id"`
	CourseKey       string `json:"course_key"`
	CertificateType string `json:"certificate_type"`
	CertType        string `json:"cert_type"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	DownloadURL     string `json:"download_url"`
	Grade           string `json:"grade"`
}

// rawCertificateList wraps the public API paginated response shape.
type rawCertificateList struct {
	Results []rawCertificate `json:"results"`
}

func (r rawCertificate) toModel() *model.Certificate {
	courseID := firstNonEmpty(r.CourseID, r.CourseKey)
	certType := firstNonEmpty(r.CertificateType, r.CertType, r.Type)
	return &model.Certificate{
		Username:        r.Username,
		CourseID:        courseID,
		CertificateType: certType,
		Status:          r.Status,
		DownloadURL:     r.DownloadURL,
		Grade:           r.Grade,
	}
}

// CertificateListFromJSON parses a list of certificates from raw JSON bytes.
// It supports the public API paginated response {"results": [...]} and the
// extension API top-level array [...].
func CertificateListFromJSON(data []byte) ([]*model.Certificate, error) {
	// Try paginated public API shape first.
	var paginated rawCertificateList
	if err := json.Unmarshal(data, &paginated); err == nil && len(paginated.Results) > 0 {
		out := make([]*model.Certificate, len(paginated.Results))
		for i, r := range paginated.Results {
			out[i] = r.toModel()
		}
		return out, nil
	}

	// Try top-level array (extension API shape).
	var items []rawCertificate
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	out := make([]*model.Certificate, len(items))
	for i, r := range items {
		out[i] = r.toModel()
	}
	return out, nil
}
```

- [ ] **Step 6: Create `internal/normalize/certificate_test.go`**

```go
package normalize

import (
	"os"
	"testing"

	"github.com/openedx/cli/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCertificateListFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/certificate_list.json")
	require.NoError(t, err)

	certs, err := CertificateListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, certs, 1)
	expected := &model.Certificate{
		Username:        "alice",
		CourseID:        "course-v1:demo+cs101+2026",
		CertificateType: "verified",
		Status:          "downloadable",
		DownloadURL:     "https://openedx.example.com/certificates/abc123",
		Grade:           "85%",
	}
	assert.Equal(t, expected, certs[0])
}

func TestCertificateListFromExtensionPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/extension/certificate_list.json")
	require.NoError(t, err)

	certs, err := CertificateListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, certs, 1)
	expected := &model.Certificate{
		Username:        "alice",
		CourseID:        "course-v1:demo+cs101+2026",
		CertificateType: "verified",
		Status:          "downloadable",
	}
	assert.Equal(t, expected, certs[0])
}

func TestCertificateListFromJSONInvalidInput(t *testing.T) {
	_, err := CertificateListFromJSON([]byte("not json"))
	assert.Error(t, err)
}
```

- [ ] **Step 7: Run all normalizer tests**

Run: `go test ./internal/normalize -v`
Expected: PASS (all existing + new tests)

- [ ] **Step 8: Commit**

```bash
git add internal/normalize/enrollment.go internal/normalize/enrollment_test.go internal/normalize/grade.go internal/normalize/grade_test.go internal/normalize/certificate.go internal/normalize/certificate_test.go
git commit -m "feat: add normalizers for enrollment list, grades, gradebook, certificates"
```

---

## Task 5: Extend enrollment command with list and remove

**Files:**
- Modify: `internal/cli/cmd/enrollment.go`
- Modify: `internal/cli/cmd/enrollment_test.go`

- [ ] **Step 1: Add `list` and `remove` subcommands to `internal/cli/cmd/enrollment.go`**

Replace the entire file content with:

```go
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/normalize"
)

// NewEnrollmentCmd creates the "enrollment" command group with all its subcommands.
func NewEnrollmentCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enrollment",
		Short: "Manage Open edX enrollments",
		Long:  "Add, remove, list, and inspect course enrollments in an Open edX deployment.",
	}

	cmd.AddCommand(
		newEnrollmentAddCmd(execFn),
		newEnrollmentListCmd(execFn),
		newEnrollmentRemoveCmd(execFn),
	)

	return cmd
}

// --- enrollment add ---

func newEnrollmentAddCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, username, mode string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Enroll a user in a course",
		Long:  "Enroll a user in a course with the specified enrollment mode.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"course_id": courseID,
				"username":  username,
				"mode":      mode,
			}

			data, err := execFn(cmd.Context(), "enrollment.add", cmdArgs)
			if err != nil {
				return err
			}

			// No dedicated normalizer for enrollment yet; output raw JSON.
			var raw json.RawMessage = data
			return printOutput(cmd, &raw)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	cmd.Flags().StringVar(&username, "username", "", "username to enroll (required)")
	cmd.Flags().StringVar(&mode, "mode", "audit", "enrollment mode (e.g. audit, verified)")
	_ = cmd.MarkFlagRequired("course-id")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}

// --- enrollment list ---

func newEnrollmentListCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, username string
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List enrollments",
		Long:  "List course enrollments with optional filters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{}
			if courseID != "" {
				cmdArgs["course_id"] = courseID
			}
			if username != "" {
				cmdArgs["username"] = username
			}
			if page > 0 {
				cmdArgs["page"] = fmt.Sprintf("%d", page)
			}
			if pageSize > 0 {
				cmdArgs["page_size"] = fmt.Sprintf("%d", pageSize)
			}

			data, err := execFn(cmd.Context(), "enrollment.list", cmdArgs)
			if err != nil {
				return err
			}

			enrollments, err := normalize.EnrollmentListFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize enrollment list: %w", err)
			}

			return printOutput(cmd, enrollments)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "filter by course ID")
	cmd.Flags().StringVar(&username, "username", "", "filter by username")
	cmd.Flags().IntVar(&page, "page", 0, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "number of results per page")

	return cmd
}

// --- enrollment remove ---

func newEnrollmentRemoveCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, username string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove a user from a course",
		Long:  "Deactivate a user's enrollment in a course.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"course_id": courseID,
				"username":  username,
			}

			data, err := execFn(cmd.Context(), "enrollment.remove", cmdArgs)
			if err != nil {
				return err
			}

			var raw json.RawMessage = data
			return printOutput(cmd, &raw)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	cmd.Flags().StringVar(&username, "username", "", "username to remove (required)")
	_ = cmd.MarkFlagRequired("course-id")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}
```

- [ ] **Step 2: Add tests for list and remove to `internal/cli/cmd/enrollment_test.go`**

Append these tests to the existing file (do NOT remove existing tests):

```go
func TestEnrollmentListCallsWithCorrectArgs(t *testing.T) {
	fixture := loadFixture(t, "enrollment_list.json")
	var capturedKey string
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewEnrollmentCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"list", "--course-id", "course-v1:demo+cs101+2026", "--username", "alice"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "enrollment.list", capturedKey)
	assert.Equal(t, "course-v1:demo+cs101+2026", capturedArgs["course_id"])
	assert.Equal(t, "alice", capturedArgs["username"])

	output := buf.String()
	assert.Contains(t, output, `"username": "alice"`)
	assert.Contains(t, output, `"course_id": "course-v1:demo+cs101+2026"`)
}

func TestEnrollmentListPageFlags(t *testing.T) {
	fixture := loadFixture(t, "enrollment_list.json")
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewEnrollmentCmd(execFn)
	cmd.SetArgs([]string{"list", "--page", "2", "--page-size", "10"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "2", capturedArgs["page"])
	assert.Equal(t, "10", capturedArgs["page_size"])
}

func TestEnrollmentRemoveRequiresCourseAndUsername(t *testing.T) {
	noop := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	}

	t.Run("missing both", func(t *testing.T) {
		cmd := NewEnrollmentCmd(noop)
		cmd.SetArgs([]string{"remove"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("missing username", func(t *testing.T) {
		cmd := NewEnrollmentCmd(noop)
		cmd.SetArgs([]string{"remove", "--course-id", "course-v1:demo+cs101+2026"})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}

func TestEnrollmentRemoveCallsWithCorrectArgs(t *testing.T) {
	enrollmentFixture := []byte(`{"username":"alice","course_id":"course-v1:demo+cs101+2026","mode":"audit","is_active":false}`)
	var capturedKey string
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		capturedArgs = args
		return enrollmentFixture, nil
	}

	cmd := NewEnrollmentCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"remove", "--course-id", "course-v1:demo+cs101+2026", "--username", "alice"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "enrollment.remove", capturedKey)
	assert.Equal(t, "course-v1:demo+cs101+2026", capturedArgs["course_id"])
	assert.Equal(t, "alice", capturedArgs["username"])
}

func TestEnrollmentCommandStructure(t *testing.T) {
	cmd := NewEnrollmentCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})

	assert.Equal(t, "enrollment", cmd.Use)
	subcommands := cmd.Commands()
	subNames := make([]string, len(subcommands))
	for i, sub := range subcommands {
		subNames[i] = sub.Use
	}
	assert.Contains(t, subNames, "add")
	assert.Contains(t, subNames, "list")
	assert.Contains(t, subNames, "remove")
}
```

Note: The existing `TestEnrollmentCommandStructure` in the file must be removed and replaced by the one above, which now checks for "list" and "remove" in addition to "add".

- [ ] **Step 3: Run enrollment tests**

Run: `go test ./internal/cli/cmd -run TestEnrollment -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/cli/cmd/enrollment.go internal/cli/cmd/enrollment_test.go
git commit -m "feat: add enrollment list and remove subcommands"
```

---

## Task 6: Extend user command with list and get

**Files:**
- Modify: `internal/cli/cmd/user.go`
- Modify: `internal/cli/cmd/user_test.go`

- [ ] **Step 1: Add `list` and `get` subcommands to `internal/cli/cmd/user.go`**

Replace the entire file content with:

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/normalize"
)

// NewUserCmd creates the "user" command group with all its subcommands.
func NewUserCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage Open edX users",
		Long:  "Create, list, and inspect users in an Open edX deployment.",
	}

	cmd.AddCommand(
		newUserCreateCmd(execFn),
		newUserListCmd(execFn),
		newUserGetCmd(execFn),
	)

	return cmd
}

// --- user create ---

func newUserCreateCmd(execFn ExecuteFunc) *cobra.Command {
	var username, email, name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Long:  "Create a new user in the configured Open edX deployment.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"username": username,
				"email":    email,
			}
			if name != "" {
				cmdArgs["name"] = name
			}

			data, err := execFn(cmd.Context(), "user.create", cmdArgs)
			if err != nil {
				return err
			}

			user, err := normalize.UserFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize user: %w", err)
			}

			return printOutput(cmd, user)
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "username for the new user (required)")
	cmd.Flags().StringVar(&email, "email", "", "email address for the new user (required)")
	cmd.Flags().StringVar(&name, "name", "", "full name for the new user")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("email")

	return cmd
}

// --- user list ---

func newUserListCmd(execFn ExecuteFunc) *cobra.Command {
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		Long:  "List user accounts in the configured Open edX deployment.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{}
			if page > 0 {
				cmdArgs["page"] = fmt.Sprintf("%d", page)
			}
			if pageSize > 0 {
				cmdArgs["page_size"] = fmt.Sprintf("%d", pageSize)
			}

			data, err := execFn(cmd.Context(), "user.list", cmdArgs)
			if err != nil {
				return err
			}

			users, err := normalize.UserListFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize user list: %w", err)
			}

			return printOutput(cmd, users)
		},
	}

	cmd.Flags().IntVar(&page, "page", 0, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "number of results per page")

	return cmd
}

// --- user get ---

func newUserGetCmd(execFn ExecuteFunc) *cobra.Command {
	var username string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get user details",
		Long:  "Retrieve details for a specific user by username.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"username": username,
			}

			data, err := execFn(cmd.Context(), "user.get", cmdArgs)
			if err != nil {
				return err
			}

			user, err := normalize.UserFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize user: %w", err)
			}

			return printOutput(cmd, user)
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "username to look up (required)")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}
```

- [ ] **Step 2: Add `UserListFromJSON` to `internal/normalize/user.go`**

Append this function to the existing file:

```go
// rawUserList wraps the public API paginated response shape.
type rawUserList struct {
	Results []rawUser `json:"results"`
}

// UserListFromJSON parses a list of users from raw JSON bytes.
// It supports the public API paginated response {"results": [...]} and the
// extension API top-level array [...].
func UserListFromJSON(data []byte) ([]*model.User, error) {
	// Try paginated public API shape first.
	var paginated rawUserList
	if err := json.Unmarshal(data, &paginated); err == nil && len(paginated.Results) > 0 {
		out := make([]*model.User, len(paginated.Results))
		for i, r := range paginated.Results {
			out[i] = r.toModel()
		}
		return out, nil
	}

	// Try top-level array (extension API shape).
	var items []rawUser
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	out := make([]*model.User, len(items))
	for i, r := range items {
		out[i] = r.toModel()
	}
	return out, nil
}
```

Also add a `toModel()` method on `rawUser`. Insert after the `rawUser` struct definition, before `UserFromJSON`:

```go
func (r rawUser) toModel() *model.User {
	name := firstNonEmpty(r.Name, r.FullName)
	active := resolveActive(r.IsActive, r.Active)
	createdAt := firstNonEmpty(r.Created, r.DateJoined, r.CreatedAt)

	return &model.User{
		Username:  r.Username,
		Email:     r.Email,
		Name:      name,
		IsActive:  active,
		CreatedAt: createdAt,
	}
}
```

Then refactor `UserFromJSON` to use `toModel()`:

```go
func UserFromJSON(data []byte) (*model.User, error) {
	var raw rawUser
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return raw.toModel(), nil
}
```

Note: the `firstNonEmpty` function is already defined in `normalize/course.go` and is accessible within the same package.

- [ ] **Step 3: Add tests to `internal/cli/cmd/user_test.go`**

Append these tests to the existing file (do NOT remove existing tests):

```go
func TestUserListCallsWithCorrectArgs(t *testing.T) {
	fixture := loadFixture(t, "user_list.json")
	var capturedKey string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		return fixture, nil
	}

	cmd := NewUserCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"list"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "user.list", capturedKey)

	output := buf.String()
	assert.Contains(t, output, `"username": "alice"`)
	assert.Contains(t, output, `"username": "bob"`)
}

func TestUserListPageFlags(t *testing.T) {
	fixture := loadFixture(t, "user_list.json")
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewUserCmd(execFn)
	cmd.SetArgs([]string{"list", "--page", "2", "--page-size", "5"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "2", capturedArgs["page"])
	assert.Equal(t, "5", capturedArgs["page_size"])
}

func TestUserGetRequiresUsername(t *testing.T) {
	cmd := NewUserCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})
	cmd.SetArgs([]string{"get"})

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestUserGetCallsWithCorrectArgs(t *testing.T) {
	fixture := loadFixture(t, "user_get.json")
	var capturedKey string
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewUserCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"get", "--username", "alice"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "user.get", capturedKey)
	assert.Equal(t, "alice", capturedArgs["username"])

	output := buf.String()
	assert.Contains(t, output, `"username": "alice"`)
	assert.Contains(t, output, `"email": "alice@example.com"`)
}
```

Update the existing `TestUserCommandStructure` to check for "list" and "get":

```go
func TestUserCommandStructure(t *testing.T) {
	cmd := NewUserCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})

	assert.Equal(t, "user", cmd.Use)
	subcommands := cmd.Commands()
	subNames := make([]string, len(subcommands))
	for i, sub := range subcommands {
		subNames[i] = sub.Use
	}
	assert.Contains(t, subNames, "create")
	assert.Contains(t, subNames, "list")
	assert.Contains(t, subNames, "get")
}
```

- [ ] **Step 4: Run user tests**

Run: `go test ./internal/cli/cmd -run TestUser -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/cli/cmd/user.go internal/cli/cmd/user_test.go internal/normalize/user.go
git commit -m "feat: add user list and get subcommands"
```

---

## Task 7: Create grade and gradebook commands

**Files:**
- Create: `internal/cli/cmd/grade.go`
- Create: `internal/cli/cmd/grade_test.go`

- [ ] **Step 1: Create `internal/cli/cmd/grade.go`**

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/normalize"
)

// NewGradeCmd creates the "grade" command group with all its subcommands.
func NewGradeCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grade",
		Short: "Manage course grades",
		Long:  "List grades and view gradebooks for Open edX courses.",
	}

	cmd.AddCommand(
		newGradeListCmd(execFn),
		newGradebookCmd(execFn),
	)

	return cmd
}

// --- grade list ---

func newGradeListCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID, username string
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List grades for a course",
		Long:  "List student grades for a specific course with optional filters.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"course_id": courseID,
			}
			if username != "" {
				cmdArgs["username"] = username
			}
			if page > 0 {
				cmdArgs["page"] = fmt.Sprintf("%d", page)
			}
			if pageSize > 0 {
				cmdArgs["page_size"] = fmt.Sprintf("%d", pageSize)
			}

			data, err := execFn(cmd.Context(), "grade.list", cmdArgs)
			if err != nil {
				return err
			}

			grades, err := normalize.GradeListFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize grade list: %w", err)
			}

			return printOutput(cmd, grades)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	cmd.Flags().StringVar(&username, "username", "", "filter by username")
	cmd.Flags().IntVar(&page, "page", 0, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "number of results per page")
	_ = cmd.MarkFlagRequired("course-id")

	return cmd
}

// --- grade gradebook ---

func newGradebookCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gradebook",
		Short: "Gradebook operations",
		Long:  "View course gradebooks.",
	}

	cmd.AddCommand(newGradebookGetCmd(execFn))
	return cmd
}

func newGradebookGetCmd(execFn ExecuteFunc) *cobra.Command {
	var courseID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get course gradebook",
		Long:  "Retrieve the full gradebook for a specific course.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"course_id": courseID,
			}

			data, err := execFn(cmd.Context(), "gradebook.get", cmdArgs)
			if err != nil {
				return err
			}

			gradebook, err := normalize.GradebookFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize gradebook: %w", err)
			}

			return printOutput(cmd, gradebook)
		},
	}

	cmd.Flags().StringVar(&courseID, "course-id", "", "course identifier (required)")
	_ = cmd.MarkFlagRequired("course-id")

	return cmd
}
```

- [ ] **Step 2: Create `internal/cli/cmd/grade_test.go`**

```go
package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGradeListRequiresCourseID(t *testing.T) {
	cmd := NewGradeCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})
	cmd.SetArgs([]string{"list"})

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestGradeListCallsWithCorrectArgs(t *testing.T) {
	fixture := loadFixture(t, "grade_list.json")
	var capturedKey string
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewGradeCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"list", "--course-id", "course-v1:demo+cs101+2026", "--username", "alice"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "grade.list", capturedKey)
	assert.Equal(t, "course-v1:demo+cs101+2026", capturedArgs["course_id"])
	assert.Equal(t, "alice", capturedArgs["username"])

	output := buf.String()
	assert.Contains(t, output, `"username": "alice"`)
	assert.Contains(t, output, `"percent": 0.85`)
}

func TestGradebookGetRequiresCourseID(t *testing.T) {
	cmd := NewGradeCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})
	cmd.SetArgs([]string{"gradebook", "get"})

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestGradebookGetCallsWithCorrectArgs(t *testing.T) {
	fixture := loadFixture(t, "gradebook_get.json")
	var capturedKey string
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewGradeCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"gradebook", "get", "--course-id", "course-v1:demo+cs101+2026"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "gradebook.get", capturedKey)
	assert.Equal(t, "course-v1:demo+cs101+2026", capturedArgs["course_id"])

	output := buf.String()
	assert.Contains(t, output, `"course_id": "course-v1:demo+cs101+2026"`)
	assert.Contains(t, output, `"grades"`)
}

func TestGradeCommandStructure(t *testing.T) {
	cmd := NewGradeCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})

	assert.Equal(t, "grade", cmd.Use)
	subcommands := cmd.Commands()
	subNames := make([]string, len(subcommands))
	for i, sub := range subcommands {
		subNames[i] = sub.Use
	}
	assert.Contains(t, subNames, "list")
	assert.Contains(t, subNames, "gradebook")
}
```

- [ ] **Step 3: Run grade tests**

Run: `go test ./internal/cli/cmd -run TestGrade -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/cli/cmd/grade.go internal/cli/cmd/grade_test.go
git commit -m "feat: add grade list and gradebook get commands"
```

---

## Task 8: Create certificate command

**Files:**
- Create: `internal/cli/cmd/certificate.go`
- Create: `internal/cli/cmd/certificate_test.go`

- [ ] **Step 1: Create `internal/cli/cmd/certificate.go`**

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openedx/cli/internal/normalize"
)

// NewCertificateCmd creates the "certificate" command group with all its subcommands.
func NewCertificateCmd(execFn ExecuteFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "Manage course certificates",
		Long:  "List and inspect course completion certificates.",
	}

	cmd.AddCommand(
		newCertificateListCmd(execFn),
	)

	return cmd
}

// --- certificate list ---

func newCertificateListCmd(execFn ExecuteFunc) *cobra.Command {
	var username string
	var page, pageSize int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List certificates",
		Long:  "List certificates for a specific user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdArgs := map[string]string{
				"username": username,
			}
			if page > 0 {
				cmdArgs["page"] = fmt.Sprintf("%d", page)
			}
			if pageSize > 0 {
				cmdArgs["page_size"] = fmt.Sprintf("%d", pageSize)
			}

			data, err := execFn(cmd.Context(), "certificate.list", cmdArgs)
			if err != nil {
				return err
			}

			certs, err := normalize.CertificateListFromJSON(data)
			if err != nil {
				return fmt.Errorf("failed to normalize certificate list: %w", err)
			}

			return printOutput(cmd, certs)
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "username to look up certificates for (required)")
	cmd.Flags().IntVar(&page, "page", 0, "page number")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "number of results per page")
	_ = cmd.MarkFlagRequired("username")

	return cmd
}
```

- [ ] **Step 2: Create `internal/cli/cmd/certificate_test.go`**

```go
package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCertificateListRequiresUsername(t *testing.T) {
	cmd := NewCertificateCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})
	cmd.SetArgs([]string{"list"})

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestCertificateListCallsWithCorrectArgs(t *testing.T) {
	fixture := loadFixture(t, "certificate_list.json")
	var capturedKey string
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewCertificateCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"list", "--username", "alice"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "certificate.list", capturedKey)
	assert.Equal(t, "alice", capturedArgs["username"])

	output := buf.String()
	assert.Contains(t, output, `"username": "alice"`)
	assert.Contains(t, output, `"certificate_type": "verified"`)
}

func TestCertificateCommandStructure(t *testing.T) {
	cmd := NewCertificateCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})

	assert.Equal(t, "certificate", cmd.Use)
	subcommands := cmd.Commands()
	subNames := make([]string, len(subcommands))
	for i, sub := range subcommands {
		subNames[i] = sub.Use
	}
	assert.Contains(t, subNames, "list")
}
```

- [ ] **Step 3: Run certificate tests**

Run: `go test ./internal/cli/cmd -run TestCertificate -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/cli/cmd/certificate.go internal/cli/cmd/certificate_test.go
git commit -m "feat: add certificate list command"
```

---

## Task 9: Register new command groups in root and run full test suite

**Files:**
- Modify: `internal/cli/root.go`

- [ ] **Step 1: Add grade and certificate commands to root.go**

In `internal/cli/root.go`, add imports for the grade and certificate packages (they are in the same `cmd` package, so no import needed — just call the constructors).

Update the `rootCmd.AddCommand(...)` block to include:

```go
	rootCmd.AddCommand(
		cmd.NewCourseCmd(execFn),
		cmd.NewUserCmd(execFn),
		cmd.NewEnrollmentCmd(execFn),
		cmd.NewRoleCmd(execFn),
		cmd.NewGradeCmd(execFn),
		cmd.NewCertificateCmd(execFn),
		cmd.NewSchemaCmd(extProvider),
		cmd.NewDoctorCmd(doctorFn),
	)
```

- [ ] **Step 2: Run the full test suite**

Run: `go test ./... -v`
Expected: PASS (all existing + new tests)

- [ ] **Step 3: Commit**

```bash
git add internal/cli/root.go
git commit -m "feat: register grade and certificate command groups in root"
```

---

## Self-Review Checklist

- [x] **Spec coverage:** All 7 commands from batch 1 have registry entries, models, normalizers, Cobra commands, tests, and fixtures
- [x] **Placeholder scan:** No TBD, TODO, or placeholder code — every step has complete code
- [x] **Type consistency:** `rawUser.toModel()` uses the same `firstNonEmpty`, `resolveActive` helpers as the existing `UserFromJSON`; `rawGrade.toModel()` uses `*bool` for `Passed` matching the `resolveActive` pattern; `rawCertificate.toModel()` uses `firstNonEmpty` for course_id resolution
