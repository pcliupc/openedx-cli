package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryContainsV1Commands(t *testing.T) {
	reg := LatestRegistry()

	expectedKeys := []string{
		"course.list",
		"course.get",
		"course.create",
		"course.import",
		"course.export",
		"course.outline.get",
		"user.create",
		"enrollment.add",
		"role.assign",
	}

	require.Len(t, reg, len(expectedKeys),
		"registry should contain exactly %d commands", len(expectedKeys))

	for _, key := range expectedKeys {
		assert.Contains(t, reg, key, "registry should contain command %q", key)
	}
}

func TestRegistryReturnsCourseListMapping(t *testing.T) {
	reg := LatestRegistry()

	cmd, ok := reg["course.list"]
	require.True(t, ok, "registry should contain course.list")

	assert.Equal(t, "course.list", cmd.Key)
	assert.Equal(t, "GET", cmd.Method)
	assert.Equal(t, "/api/courses/v1/courses", cmd.Path)
	assert.Empty(t, cmd.RequiredArgs, "course.list should have no required args")
	assert.Equal(t, "Course", cmd.OutputModel)
}

func TestRegistryReturnsCourseCreateMapping(t *testing.T) {
	reg := LatestRegistry()

	cmd, ok := reg["course.create"]
	require.True(t, ok, "registry should contain course.create")

	assert.Equal(t, "course.create", cmd.Key)
	assert.Equal(t, "POST", cmd.Method)
	assert.Equal(t, "/api/courses/v1/courses", cmd.Path)
	assert.Equal(t, []string{"org", "number", "run", "title"}, cmd.RequiredArgs)
	assert.Equal(t, "Course", cmd.OutputModel)
}

func TestRegistryCommandHasRequiredFields(t *testing.T) {
	reg := LatestRegistry()
	require.NotEmpty(t, reg, "registry should not be empty")

	for key, cmd := range reg {
		assert.NotEmpty(t, cmd.Method, "command %q should have a Method", key)
		assert.NotEmpty(t, cmd.Path, "command %q should have a Path", key)
		assert.NotEmpty(t, cmd.OutputModel, "command %q should have an OutputModel", key)
	}
}
