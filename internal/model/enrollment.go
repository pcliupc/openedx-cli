package model

// Enrollment represents a user's enrollment in a course.
type Enrollment struct {
	Username string `json:"username"`
	CourseID string `json:"course_id"`
	Mode     string `json:"mode"`
	IsActive bool   `json:"is_active"`
}
