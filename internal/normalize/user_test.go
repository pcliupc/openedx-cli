package normalize

import (
	"os"
	"testing"

	"github.com/openedx/cli/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeUserFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/user.json")
	require.NoError(t, err)

	user, err := UserFromJSON(data)
	require.NoError(t, err)

	expected := &model.User{
		Username:  "alice",
		Email:     "alice@example.com",
		Name:      "Alice Smith",
		IsActive:  true,
		CreatedAt: "2026-01-15T10:30:00Z",
	}
	assert.Equal(t, expected, user)
}

func TestNormalizeUserFromExtensionPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/extension/user.json")
	require.NoError(t, err)

	user, err := UserFromJSON(data)
	require.NoError(t, err)

	expected := &model.User{
		Username:  "alice",
		Email:     "alice@example.com",
		Name:      "Alice Smith",
		IsActive:  true,
		CreatedAt: "2026-01-15T10:30:00Z",
	}
	assert.Equal(t, expected, user)
}

func TestNormalizeUserInactive(t *testing.T) {
	data := []byte(`{
		"username": "bob",
		"email": "bob@example.com",
		"name": "Bob Jones",
		"is_active": false,
		"date_joined": "2025-12-01T00:00:00Z"
	}`)

	user, err := UserFromJSON(data)
	require.NoError(t, err)

	assert.Equal(t, "bob", user.Username)
	assert.False(t, user.IsActive)
}

func TestNormalizeUserDefaultsActiveToTrue(t *testing.T) {
	data := []byte(`{
		"username": "charlie",
		"email": "charlie@example.com",
		"name": "Charlie"
	}`)

	user, err := UserFromJSON(data)
	require.NoError(t, err)

	assert.True(t, user.IsActive)
	assert.Empty(t, user.CreatedAt)
}

func TestUserFromJSONInvalidInput(t *testing.T) {
	_, err := UserFromJSON([]byte("not json"))
	assert.Error(t, err)
}
