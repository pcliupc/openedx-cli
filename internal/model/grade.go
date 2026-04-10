package model

// Grade represents a student's grade in a course section.
type Grade struct {
	Username    string  `json:"username"`
	CourseID    string  `json:"course_id"`
	Percent     float64 `json:"percent"`
	LetterGrade string  `json:"letter_grade,omitempty"`
	Passed      bool    `json:"passed"`
	Section     string  `json:"section,omitempty"`
}

// Gradebook represents the full gradebook for a course.
type Gradebook struct {
	CourseID string  `json:"course_id"`
	Grades   []Grade `json:"grades"`
}
