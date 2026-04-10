package normalize

import (
	"encoding/json"

	"github.com/openedx/cli/internal/model"
)

// rawGrade is a lax intermediate struct that accepts field names from both
// public and extension API payloads.
type rawGrade struct {
	Username    string  `json:"username"`
	CourseID    string  `json:"course_id"`
	Percent     float64 `json:"percent"`
	LetterGrade string  `json:"letter_grade"`
	Grade       string  `json:"grade"`
	Passed      *bool   `json:"passed"`
	Section     string  `json:"section"`
	Subsection  string  `json:"subsection"`
}

// rawGradeList wraps the public API paginated response shape.
type rawGradeList struct {
	Results []rawGrade `json:"results"`
}

func (r rawGrade) toModel() *model.Grade {
	letter := firstNonEmpty(r.LetterGrade, r.Grade)
	passed := true
	if r.Passed != nil {
		passed = *r.Passed
	}
	section := firstNonEmpty(r.Section, r.Subsection)
	return &model.Grade{
		Username:    r.Username,
		CourseID:    r.CourseID,
		Percent:     r.Percent,
		LetterGrade: letter,
		Passed:      passed,
		Section:     section,
	}
}

// GradeListFromJSON parses a list of grades from raw JSON bytes.
// It supports the public API paginated response {"results": [...]} and the
// extension API top-level array [...].
func GradeListFromJSON(data []byte) ([]*model.Grade, error) {
	// Try paginated public API shape first.
	var paginated rawGradeList
	if err := json.Unmarshal(data, &paginated); err == nil && len(paginated.Results) > 0 {
		out := make([]*model.Grade, len(paginated.Results))
		for i, r := range paginated.Results {
			out[i] = r.toModel()
		}
		return out, nil
	}

	// Try top-level array (extension API shape).
	var items []rawGrade
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	out := make([]*model.Grade, len(items))
	for i, r := range items {
		out[i] = r.toModel()
	}
	return out, nil
}

// rawGradebook is a lax intermediate struct for the gradebook response.
type rawGradebook struct {
	CourseID  string     `json:"course_id"`
	CourseKey string     `json:"course_key"`
	Grades    []rawGrade `json:"grades"`
}

// GradebookFromJSON parses a gradebook from raw JSON bytes.
func GradebookFromJSON(data []byte) (*model.Gradebook, error) {
	var raw rawGradebook
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	courseID := firstNonEmpty(raw.CourseID, raw.CourseKey)
	grades := make([]model.Grade, len(raw.Grades))
	for i, r := range raw.Grades {
		grades[i] = *r.toModel()
	}

	return &model.Gradebook{
		CourseID: courseID,
		Grades:   grades,
	}, nil
}
