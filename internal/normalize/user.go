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

func (r rawUser) toModel() *model.User {
	name := firstNonEmpty(r.Name, r.FullName)
	active := resolveActive(r.IsActive, r.Active)
	createdAt := firstNonEmpty(r.Created, r.DateJoined, r.CreatedAt)

	return &model.User{
		Username:  r.Username,
		Email:     r.Email,
		Name:      name,
		IsActive:  active,
		CreatedAt: createdAt,
	}
}

// UserFromJSON parses a user from raw JSON bytes and returns the stable model.
// It handles both public API (fields: name, is_active, date_joined) and
// extension API (fields: full_name, active, created_at) naming conventions.
func UserFromJSON(data []byte) (*model.User, error) {
	var raw rawUser
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return raw.toModel(), nil
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

// rawUserList wraps the public API paginated response shape.
type rawUserList struct {
	Results []rawUser `json:"results"`
}

// UserListFromJSON parses a list of users from raw JSON bytes.
// It supports the public API paginated response {"results": [...]} and the
// extension API top-level array [...].
func UserListFromJSON(data []byte) ([]*model.User, error) {
	// Try paginated public API shape first.
	var paginated rawUserList
	if err := json.Unmarshal(data, &paginated); err == nil && len(paginated.Results) > 0 {
		out := make([]*model.User, len(paginated.Results))
		for i, r := range paginated.Results {
			out[i] = r.toModel()
		}
		return out, nil
	}

	// Try top-level array (extension API shape).
	var items []rawUser
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	out := make([]*model.User, len(items))
	for i, r := range items {
		out[i] = r.toModel()
	}
	return out, nil
}
