package normalize

import (
	"encoding/json"

	"github.com/openedx/cli/internal/model"
)

// rawUser is a lax intermediate struct that accepts field names from both
// public and extension API payloads.
type rawUser struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
	IsActive   *bool  `json:"is_active"`
	Active     *bool  `json:"active"`
	Created    string `json:"created"`
	DateJoined string `json:"date_joined"`
	CreatedAt  string `json:"created_at"`
}

// UserFromJSON parses a user from raw JSON bytes and returns the stable model.
// It handles both public API (fields: name, is_active, date_joined) and
// extension API (fields: full_name, active, created_at) naming conventions.
func UserFromJSON(data []byte) (*model.User, error) {
	var raw rawUser
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	name := firstNonEmpty(raw.Name, raw.FullName)
	active := resolveActive(raw.IsActive, raw.Active)
	createdAt := firstNonEmpty(raw.Created, raw.DateJoined, raw.CreatedAt)

	return &model.User{
		Username:  raw.Username,
		Email:     raw.Email,
		Name:      name,
		IsActive:  active,
		CreatedAt: createdAt,
	}, nil
}

// resolveActive returns the active status from whichever boolean pointer is
// non-nil, defaulting to true when neither is set.
func resolveActive(isActive, active *bool) bool {
	if isActive != nil {
		return *isActive
	}
	if active != nil {
		return *active
	}
	return true
}
