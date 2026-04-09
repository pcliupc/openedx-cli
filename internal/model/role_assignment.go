package model

// RoleAssignment represents a role assigned to a user on a course.
type RoleAssignment struct {
	Username   string `json:"username"`
	CourseID   string `json:"course_id"`
	Role       string `json:"role"`
	AssignedBy string `json:"assigned_by,omitempty"`
	AssignedAt string `json:"assigned_at,omitempty"`
}
