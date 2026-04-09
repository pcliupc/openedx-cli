// Package normalize converts raw JSON payloads from provider responses into
// stable model structs. It handles both public API and extension API field
// naming conventions, mapping them to a unified output shape.
package normalize

import (
	"encoding/json"

	"github.com/openedx/cli/internal/model"
)

// rawCourse is a lax intermediate struct that accepts field names from both
// public and extension API payloads.
type rawCourse struct {
	ID          string `json:"id"`
	CourseID    string `json:"course_id"`
	Org         string `json:"org"`
	Number      string `json:"number"`
	Run         string `json:"run"`
	DisplayName string `json:"display_name"`
	Title       string `json:"title"`
	Pacing      string `json:"pacing"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Status      string `json:"status"`
}

// rawCourseList wraps the public API paginated response shape.
type rawCourseList struct {
	Results []rawCourse `json:"results"`
}

func (r rawCourse) toModel() *model.Course {
	courseID := firstNonEmpty(r.ID, r.CourseID)
	title := firstNonEmpty(r.DisplayName, r.Title)
	return &model.Course{
		CourseID: courseID,
		Org:      r.Org,
		Number:   r.Number,
		Run:      r.Run,
		Title:    title,
		Pacing:   r.Pacing,
		Start:    r.Start,
		End:      r.End,
		Status:   r.Status,
	}
}

// CourseFromJSON parses a single course from raw JSON bytes, handling both
// public API (fields: id, display_name) and extension API (fields: course_id,
// title) naming conventions.
func CourseFromJSON(data []byte) (*model.Course, error) {
	var raw rawCourse
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return raw.toModel(), nil
}

// CourseListFromJSON parses a list of courses from raw JSON bytes. It supports
// two shapes: the public API paginated response {"results": [...]} and the
// extension API top-level array [...].
func CourseListFromJSON(data []byte) ([]*model.Course, error) {
	// Try paginated public API shape first.
	var paginated rawCourseList
	if err := json.Unmarshal(data, &paginated); err == nil && len(paginated.Results) > 0 {
		out := make([]*model.Course, len(paginated.Results))
		for i, r := range paginated.Results {
			out[i] = r.toModel()
		}
		return out, nil
	}

	// Try top-level array (extension API shape).
	var items []rawCourse
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	out := make([]*model.Course, len(items))
	for i, r := range items {
		out[i] = r.toModel()
	}
	return out, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
