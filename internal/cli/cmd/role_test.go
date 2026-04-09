package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleAssignRequiresCourseUsernameRole(t *testing.T) {
	noop := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	}

	t.Run("missing all", func(t *testing.T) {
		cmd := NewRoleCmd(noop)
		cmd.SetArgs([]string{"assign"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("missing role", func(t *testing.T) {
		cmd := NewRoleCmd(noop)
		cmd.SetArgs([]string{"assign", "--course-id", "course-v1:demo+cs101+2026", "--username", "alice"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("missing username", func(t *testing.T) {
		cmd := NewRoleCmd(noop)
		cmd.SetArgs([]string{"assign", "--course-id", "course-v1:demo+cs101+2026", "--role", "instructor"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("missing course-id", func(t *testing.T) {
		cmd := NewRoleCmd(noop)
		cmd.SetArgs([]string{"assign", "--username", "alice", "--role", "instructor"})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}

func TestRoleAssignCallsWithCorrectArgs(t *testing.T) {
	roleFixture := []byte(`{"username":"alice","course_id":"course-v1:demo+cs101+2026","role":"instructor","assigned_by":"admin"}`)
	var capturedKey string
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedKey = cmdKey
		capturedArgs = args
		return roleFixture, nil
	}

	cmd := NewRoleCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"assign", "--course-id", "course-v1:demo+cs101+2026", "--username", "alice", "--role", "instructor"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "role.assign", capturedKey)
	assert.Equal(t, "course-v1:demo+cs101+2026", capturedArgs["course_id"])
	assert.Equal(t, "alice", capturedArgs["username"])
	assert.Equal(t, "instructor", capturedArgs["role"])

	// Verify raw JSON output.
	output := buf.String()
	assert.Contains(t, output, `"username": "alice"`)
	assert.Contains(t, output, `"role": "instructor"`)

	// Verify output is valid JSON.
	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(output), &parsed))
}

func TestRoleCommandStructure(t *testing.T) {
	cmd := NewRoleCmd(func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	})

	assert.Equal(t, "role", cmd.Use)
	subcommands := cmd.Commands()
	subNames := make([]string, len(subcommands))
	for i, sub := range subcommands {
		subNames[i] = sub.Use
	}
	assert.Contains(t, subNames, "assign")
}
