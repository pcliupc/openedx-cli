package model

// CourseOutline represents the hierarchical structure of a course.
type CourseOutline struct {
	CourseID string         `json:"course_id"`
	Chapters []OutlineBlock `json:"chapters"`
}

// OutlineBlock represents a single block (chapter, sequential, vertical, etc.)
// in a course outline tree.
type OutlineBlock struct {
	ID       string         `json:"id"`
	Title    string         `json:"title"`
	Type     string         `json:"type"`
	Children []OutlineBlock `json:"children,omitempty"`
}
