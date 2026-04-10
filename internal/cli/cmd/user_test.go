package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCreateRequiresUsernameAndEmail(t *testing.T) {
	noop := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		return nil, nil
	}

	t.Run("missing both", func(t *testing.T) {
		cmd := NewUserCmd(noop)
		cmd.SetArgs([]string{"create"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("missing email", func(t *testing.T) {
		cmd := NewUserCmd(noop)
		cmd.SetArgs([]string{"create", "--username", "alice"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("missing username", func(t *testing.T) {
		cmd := NewUserCmd(noop)
		cmd.SetArgs([]string{"create", "--email", "alice@example.com"})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}

func TestUserCreateCallsWithCorrectArgs(t *testing.T) {
	fixture := loadFixture(t, "user.json")
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
	cmd.SetArgs([]string{"create", "--username", "alice", "--email", "alice@example.com"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "user.create", capturedKey)
	assert.Equal(t, "alice", capturedArgs["username"])
	assert.Equal(t, "alice@example.com", capturedArgs["email"])

	// Verify normalized user output.
	output := buf.String()
	assert.Contains(t, output, `"username": "alice"`)
	assert.Contains(t, output, `"email": "alice@example.com"`)
	assert.Contains(t, output, `"name": "Alice Smith"`)
}

func TestUserCreateWithOptionalName(t *testing.T) {
	fixture := loadFixture(t, "user.json")
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewUserCmd(execFn)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"create", "--username", "alice", "--email", "alice@example.com", "--name", "Alice Smith"})

	err := cmd.Execute()
	require.NoError(t, err)
	assert.Equal(t, "Alice Smith", capturedArgs["name"])

	// Verify it's valid JSON output.
	output := buf.String()
	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(output), &parsed))
}

func TestUserCreateWithoutOptionalName(t *testing.T) {
	fixture := loadFixture(t, "user.json")
	var capturedArgs map[string]string
	execFn := func(ctx context.Context, cmdKey string, args map[string]string) ([]byte, error) {
		capturedArgs = args
		return fixture, nil
	}

	cmd := NewUserCmd(execFn)
	cmd.SetArgs([]string{"create", "--username", "alice", "--email", "alice@example.com"})

	err := cmd.Execute()
	require.NoError(t, err)
	_, hasName := capturedArgs["name"]
	assert.False(t, hasName, "name should not be set when --name flag is omitted")
}

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
