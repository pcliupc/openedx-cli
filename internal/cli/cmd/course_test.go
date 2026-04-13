package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadFixture is a test helper that reads a JSON fixture from testdata/public/.
func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(fmt.Sprintf("../../../testdata/public/%s", name))
	require.NoError(t, err, "failed to load fixture %s", name)
	return data
}

// mockExecFn returns an ExecuteFunc that returns the given fixture data.
func mockExecFn(fixtureData []byte) ExecuteFunc {
	return func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return fixtureData, nil
	}
}

// mockExecFnWithKey returns an ExecuteFunc that records the cmdKey it was
// called with and returns the given fixture data.
func mockExecFnWithKey(fixtureData []byte, captured *string) ExecuteFunc {
	return func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		*captured = cmdKey
		return fixtureData, nil
	}
}

func TestCourseListCallsRegistryCommand(t *testing.T) {
	fixture := loadFixture(t, "course_list.json")
	var captured string
	execFn := mockExecFnWithKey(fixture, &captured)

	cmd := NewCourseCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"list"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "course.list", captured)

	// Verify the output contains normalized course data.
	output := buf.String()
	assert.Contains(t, output, `"course_id": "course-v1:demo+cs101+2026"`)
	assert.Contains(t, output, `"org": "demo"`)
	assert.Contains(t, output, `"title": "Intro to AI"`)
}

func TestCourseGetRequiresCourseID(t *testing.T) {
	cmd := NewCourseCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})
	cmd.SetArgs([]string{"get"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestCourseGetCallsWithCorrectArgs(t *testing.T) {
	fixture := loadFixture(t, "course_get.json")
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewCourseCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"get", "--course-id", "course-v1:demo+cs101+2026"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "course-v1:demo+cs101+2026", capturedArgs["course_id"])

	output := buf.String()
	assert.Contains(t, output, `"course_id": "course-v1:demo+cs101+2026"`)
	assert.Contains(t, output, `"title": "Intro to AI"`)
}

func TestCourseCreateSupportsOrgNumberRunTitle(t *testing.T) {
	fixture := loadFixture(t, "course_get.json")
	var capturedKey string
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewCourseCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"create", "--org", "demo", "--number", "cs101", "--run", "2026", "--title", "Intro to AI"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "course.create", capturedKey)
	assert.Equal(t, "demo", capturedArgs["org"])
	assert.Equal(t, "cs101", capturedArgs["number"])
	assert.Equal(t, "2026", capturedArgs["run"])
	assert.Equal(t, "Intro to AI", capturedArgs["title"])
}

func TestCourseCreateRequiresAllFlags(t *testing.T) {
	cmd := NewCourseCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})
	cmd.SetArgs([]string{"create"})

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestCourseOutlineGetFormatsStableJSON(t *testing.T) {
	fixture := loadFixture(t, "outline.json")

	cmd := NewCourseCmd(mockExecFn(fixture))
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"outline", "get", "--course-id", "course-v1:demo+cs101+2026"})

	err := cmd.Execute()
	require.NoError(t, err)

	output := buf.String()

	// Verify it is valid JSON.
	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(output), &parsed))

	// Verify normalized fields present.
	assert.Contains(t, output, `"title": "Chapter 1: Introduction"`)
	assert.Contains(t, output, `"title": "Welcome Video"`)
}

func TestCourseExportRequiresCourseID(t *testing.T) {
	cmd := NewCourseCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})
	cmd.SetArgs([]string{"export"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestCourseExportCallsWithCorrectArgs(t *testing.T) {
	jobFixture := []byte(`{"job_id":"job-123","operation":"export","status":"submitted"}`)
	var capturedKey string
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		capturedArgs = args
		return jobFixture, nil
	}

	cmd := NewCourseCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"export", "--course-id", "course-v1:demo+cs101+2026"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "course.export", capturedKey)
	assert.Equal(t, "course-v1:demo+cs101+2026", capturedArgs["course_id"])

	output := buf.String()
	assert.Contains(t, output, `"job_id": "job-123"`)
	assert.Contains(t, output, `"operation": "export"`)
}

func TestCourseImportRequiresFile(t *testing.T) {
	cmd := NewCourseCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})
	cmd.SetArgs([]string{"import", "--course-id", "course-v1:demo+cs101+2026"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestCourseImportRequiresCourseID(t *testing.T) {
	cmd := NewCourseCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})
	cmd.SetArgs([]string{"import", "--file", "/tmp/course.tar.gz"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestCourseImportCallsWithCorrectArgs(t *testing.T) {
	jobFixture := []byte(`{"job_id":"job-456","operation":"import","status":"submitted"}`)
	var capturedKey string
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		capturedArgs = args
		return jobFixture, nil
	}

	cmd := NewCourseCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"import", "--course-id", "course-v1:demo+cs101+2026", "--file", "/tmp/course.tar.gz"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "course.import", capturedKey)
	assert.Equal(t, "course-v1:demo+cs101+2026", capturedArgs["course_id"])
	assert.Equal(t, "/tmp/course.tar.gz", capturedArgs["file"])

	output := buf.String()
	assert.Contains(t, output, `"job_id": "job-456"`)
}

func TestCourseListPageFlags(t *testing.T) {
	fixture := loadFixture(t, "course_list.json")
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewCourseCmd(execFn)
	cmd.SetArgs([]string{"list", "--page", "2", "--page-size", "10"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "2", capturedArgs["page"])
	assert.Equal(t, "10", capturedArgs["page_size"])
}

func TestCourseListAllFlag(t *testing.T) {
	fixture := loadFixture(t, "course_list.json")
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewCourseCmd(execFn)
	cmd.SetArgs([]string{"list", "--all"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "true", capturedArgs["all"])
	_, hasPage := capturedArgs["page"]
	assert.False(t, hasPage, "page should not be set when --all is used")
}

func TestCourseExportWithOutput(t *testing.T) {
	jobFixture := []byte(`{"job_id":"job-789","operation":"export","status":"submitted"}`)
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedArgs = args
		return jobFixture, nil
	}

	cmd := NewCourseCmd(execFn)
	cmd.SetArgs([]string{"export", "--course-id", "course-v1:demo+cs101+2026", "--output", "/tmp/exports"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "/tmp/exports", capturedArgs["output"])
}

func TestCourseCommandStructure(t *testing.T) {
	cmd := NewCourseCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})

	// Verify command structure.
	assert.Equal(t, "course", cmd.Use)
	subcommands := cmd.Commands()
	subNames := make([]string, len(subcommands))
	for i, sub := range subcommands {
		subNames[i] = sub.Use
	}
	assert.Contains(t, subNames, "list")
	assert.Contains(t, subNames, "get")
	assert.Contains(t, subNames, "create")
	assert.Contains(t, subNames, "import")
	assert.Contains(t, subNames, "export")
	assert.Contains(t, subNames, "outline")
}

func TestCourseOutlineGetRequiresCourseID(t *testing.T) {
	cmd := NewCourseCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})
	cmd.SetArgs([]string{"outline", "get"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "required"), "expected 'required' in error, got: %s", err.Error())
}

func TestCourseExecutorError(t *testing.T) {
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, fmt.Errorf("api server unreachable")
	}

	cmd := NewCourseCmd(execFn)
	cmd.SetArgs([]string{"list"})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "api server unreachable")
}
