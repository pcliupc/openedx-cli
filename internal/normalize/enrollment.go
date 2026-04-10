package normalize

import (
	"encoding/json"

	"github.com/openedx/cli/internal/model"
)

// rawEnrollment is a lax intermediate struct that accepts field names from both
// public and extension API payloads.
type rawEnrollment struct {
	Username string `json:"username"`
	CourseID string `json:"course_id"`
	User     string `json:"user"`
	Course   string `json:"course"`
	Mode     string `json:"mode"`
	IsActive *bool  `json:"is_active"`
	Active   *bool  `json:"active"`
}

// rawEnrollmentList wraps the public API paginated response shape.
type rawEnrollmentList struct {
	Results []rawEnrollment `json:"results"`
}

func (r rawEnrollment) toModel() *model.Enrollment {
	courseID := firstNonEmpty(r.CourseID, r.Course)
	username := firstNonEmpty(r.Username, r.User)
	active := resolveActive(r.IsActive, r.Active)
	return &model.Enrollment{
		Username: username,
		CourseID: courseID,
		Mode:     r.Mode,
		IsActive: active,
	}
}

// EnrollmentListFromJSON parses a list of enrollments from raw JSON bytes.
// It supports the public API paginated response {"results": [...]} and the
// extension API top-level array [...].
func EnrollmentListFromJSON(data []byte) ([]*model.Enrollment, error) {
	// Try paginated public API shape first.
	var paginated rawEnrollmentList
	if err := json.Unmarshal(data, &paginated); err == nil && len(paginated.Results) > 0 {
		out := make([]*model.Enrollment, len(paginated.Results))
		for i, r := range paginated.Results {
			out[i] = r.toModel()
		}
		return out, nil
	}

	// Try top-level array (extension API shape).
	var items []rawEnrollment
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	out := make([]*model.Enrollment, len(items))
	for i, r := range items {
		out[i] = r.toModel()
	}
	return out, nil
}

// EnrollmentFromJSON parses a single enrollment from raw JSON bytes.
func EnrollmentFromJSON(data []byte) (*model.Enrollment, error) {
	var raw rawEnrollment
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return raw.toModel(), nil
}
