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
