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
