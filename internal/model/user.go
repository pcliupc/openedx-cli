package model

// User represents a normalized OpenEdX user resource.
type User struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at,omitempty"`
}
