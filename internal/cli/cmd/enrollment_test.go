package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnrollmentAddRequiresCourseAndUsername(t *testing.T) {
	noop := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	}

	t.Run("missing both", func(t *testing.T) {
		cmd := NewEnrollmentCmd(noop)
		cmd.SetArgs([]string{"add"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("missing username", func(t *testing.T) {
		cmd := NewEnrollmentCmd(noop)
		cmd.SetArgs([]string{"add", "--course-id", "course-v1:demo+cs101+2026"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("missing course-id", func(t *testing.T) {
		cmd := NewEnrollmentCmd(noop)
		cmd.SetArgs([]string{"add", "--username", "alice"})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}

func TestEnrollmentAddCallsWithCorrectArgs(t *testing.T) {
	enrollmentFixture := []byte(`{"username":"alice","course_id":"course-v1:demo+cs101+2026","mode":"audit","is_active":true}`)
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
	cmd.SetArgs([]string{"add", "--course-id", "course-v1:demo+cs101+2026", "--username", "alice", "--mode", "verified"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "enrollment.add", capturedKey)
	assert.Equal(t, "course-v1:demo+cs101+2026", capturedArgs["course_id"])
	assert.Equal(t, "alice", capturedArgs["username"])
	assert.Equal(t, "verified", capturedArgs["mode"])

	// Verify raw JSON output.
	output := buf.String()
	assert.Contains(t, output, `"username": "alice"`)
	assert.Contains(t, output, `"course_id": "course-v1:demo+cs101+2026"`)
}

func TestEnrollmentAddDefaultMode(t *testing.T) {
	enrollmentFixture := []byte(`{"username":"bob","course_id":"course-v1:demo+cs101+2026","mode":"audit","is_active":true}`)
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedArgs = args
		return enrollmentFixture, nil
	}

	cmd := NewEnrollmentCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"add", "--course-id", "course-v1:demo+cs101+2026", "--username", "bob"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "audit", capturedArgs["mode"], "mode should default to 'audit'")

	// Verify output is valid JSON.
	output := buf.String()
	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(output), &parsed))
}

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
