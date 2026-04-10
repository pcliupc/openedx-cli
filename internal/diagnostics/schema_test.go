package diagnostics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/openedx/cli/internal/config"
)

func TestSchemaShowsPublicMapping(t *testing.T) {
	schema, err := GetCommandSchema("course.list", nil)
	require.NoError(t, err)

	assert.Equal(t, "course.list", schema.CommandKey)
	assert.Equal(t, "GET", schema.PublicMethod)
	assert.Equal(t, "/api/courses/v1/courses", schema.PublicPath)
	assert.False(t, schema.HasExtension)
	assert.Empty(t, schema.ExtensionURL)
	assert.Equal(t, []string{}, schema.RequiredArgs)
	assert.Equal(t, "Course", schema.OutputModel)
}

func TestSchemaShowsExtensionWhenConfigured(t *testing.T) {
	extensions := map[string]config.ExtensionMapping{
		"course.create": {
			Method: "POST",
			URL:    "https://custom.example.com/api/courses",
		},
	}

	schema, err := GetCommandSchema("course.create", extensions)
	require.NoError(t, err)

	assert.Equal(t, "course.create", schema.CommandKey)
	assert.Equal(t, "POST", schema.PublicMethod)
	assert.Equal(t, "/api/courses/v1/courses", schema.PublicPath)
	assert.True(t, schema.HasExtension)
	assert.Equal(t, "https://custom.example.com/api/courses", schema.ExtensionURL)
	assert.Equal(t, []string{"org", "number", "run", "title"}, schema.RequiredArgs)
	assert.Equal(t, "Course", schema.OutputModel)
}

func TestSchemaAllListsV1Commands(t *testing.T) {
	schemas, err := GetAllCommandSchemas(nil)
	require.NoError(t, err)

	// The registry contains 16 commands.
	assert.Len(t, schemas, 16)

	// Verify output is sorted by command key.
	for i := 1; i < len(schemas); i++ {
		assert.True(t, schemas[i-1].CommandKey < schemas[i].CommandKey,
			"schemas should be sorted alphabetically: %s < %s",
			schemas[i-1].CommandKey, schemas[i].CommandKey)
	}

	// Spot-check a known command.
	found := false
	for _, s := range schemas {
		if s.CommandKey == "course.list" {
			found = true
			assert.Equal(t, "GET", s.PublicMethod)
			assert.Equal(t, "/api/courses/v1/courses", s.PublicPath)
			assert.Equal(t, "Course", s.OutputModel)
			break
		}
	}
	assert.True(t, found, "course.list should be present in all schemas")
}

func TestSchemaReturnsErrorForUnknownCommand(t *testing.T) {
	_, err := GetCommandSchema("nonexistent.command", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown command")
	assert.Contains(t, err.Error(), "nonexistent.command")
}
