package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/openedx/cli/internal/model"
)

func TestRootCommandDefaultsToJSON(t *testing.T) {
	cmd := NewRootCmd()
	format, err := cmd.PersistentFlags().GetString("format")
	if err != nil {
		t.Fatalf("failed to get format flag: %v", err)
	}
	if format != "json" {
		t.Errorf("expected default format 'json', got '%s'", format)
	}
}

func TestRootCommandAcceptsProfileFlag(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--profile", "admin"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	profile, err := cmd.PersistentFlags().GetString("profile")
	if err != nil {
		t.Fatalf("failed to get profile flag: %v", err)
	}
	if profile != "admin" {
		t.Errorf("expected profile 'admin', got '%s'", profile)
	}
}

func TestRootCommandAcceptsAllFlags(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--profile", "ops", "--format", "table", "--config", "/tmp/cfg.yaml", "--verbose"})
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	profile, _ := cmd.PersistentFlags().GetString("profile")
	format, _ := cmd.PersistentFlags().GetString("format")
	config, _ := cmd.PersistentFlags().GetString("config")
	verbose, _ := cmd.PersistentFlags().GetBool("verbose")

	if profile != "ops" {
		t.Errorf("expected profile 'ops', got '%s'", profile)
	}
	if format != "table" {
		t.Errorf("expected format 'table', got '%s'", format)
	}
	if config != "/tmp/cfg.yaml" {
		t.Errorf("expected config '/tmp/cfg.yaml', got '%s'", config)
	}
	if !verbose {
		t.Error("expected verbose to be true")
	}
}

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"key": "value"}
	err := PrintJSON(&buf, data)
	if err != nil {
		t.Fatalf("PrintJSON failed: %v", err)
	}
	expected := `{
  "key": "value"
}`
	if strings.TrimSpace(buf.String()) != expected {
		t.Errorf("unexpected JSON output:\ngot:\n%s\nwant:\n%s", buf.String(), expected)
	}
}

func TestPrintTableWithSlice(t *testing.T) {
	users := []model.User{
		{Username: "alice", Email: "alice@example.com", Name: "Alice", IsActive: true},
		{Username: "bob", Email: "bob@example.com", Name: "Bob", IsActive: false},
	}
	var buf bytes.Buffer
	err := PrintTable(&buf, users)
	if err != nil {
		t.Fatalf("PrintTable failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "USERNAME") {
		t.Errorf("missing header: %s", output)
	}
	if !strings.Contains(output, "alice") {
		t.Errorf("missing data: %s", output)
	}
	if !strings.Contains(output, "bob") {
		t.Errorf("missing data: %s", output)
	}
}

func TestPrintTableWithNonSliceFallsBackToJSON(t *testing.T) {
	user := model.User{Username: "alice", Email: "alice@example.com"}
	var buf bytes.Buffer
	err := PrintTable(&buf, user)
	if err != nil {
		t.Fatalf("PrintTable failed: %v", err)
	}
	if !strings.Contains(buf.String(), `"username": "alice"`) {
		t.Errorf("expected JSON fallback, got: %s", buf.String())
	}
}

func TestPrintTableEmptySlice(t *testing.T) {
	var courses []model.Course
	var buf bytes.Buffer
	err := PrintTable(&buf, courses)
	if err != nil {
		t.Fatalf("PrintTable failed: %v", err)
	}
	if !strings.Contains(buf.String(), "(no results)") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestPrintOutputJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"hello": "world"}
	err := PrintOutput(&buf, "json", data)
	if err != nil {
		t.Fatalf("PrintOutput failed: %v", err)
	}
	if !strings.Contains(buf.String(), `"hello": "world"`) {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestPrintOutputUnknownFormatDefaultsToJSON(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"hello": "world"}
	err := PrintOutput(&buf, "unknown", data)
	if err != nil {
		t.Fatalf("PrintOutput failed: %v", err)
	}
	if !strings.Contains(buf.String(), `"hello": "world"`) {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestCLIError(t *testing.T) {
	err := NewCLIError("permission_denied", "profile 'ops' cannot create course")
	if err.Error() != "permission_denied: profile 'ops' cannot create course" {
		t.Errorf("unexpected Error() string: %s", err.Error())
	}
	if err.Code != "permission_denied" {
		t.Errorf("unexpected error code: %s", err.Code)
	}
	if err.Message != "profile 'ops' cannot create course" {
		t.Errorf("unexpected message: %s", err.Message)
	}
}

func TestCLIErrorPrint(t *testing.T) {
	var buf bytes.Buffer
	err := NewCLIError("not_found", "course not found")
	err.Resource = "course.get"
	err.Suggestion = "check the course ID"
	if printErr := err.Print(&buf); printErr != nil {
		t.Fatalf("Print failed: %v", printErr)
	}
	output := buf.String()
	if !strings.Contains(output, `"error": "not_found"`) {
		t.Errorf("missing error code in output: %s", output)
	}
	if !strings.Contains(output, `"resource": "course.get"`) {
		t.Errorf("missing resource in output: %s", output)
	}
	if !strings.Contains(output, `"suggestion": "check the course ID"`) {
		t.Errorf("missing suggestion in output: %s", output)
	}
}
