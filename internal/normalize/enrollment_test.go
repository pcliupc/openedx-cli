package normalize

import (
	"os"
	"testing"

	"github.com/openedx/cli/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnrollmentListFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/enrollment_list.json")
	require.NoError(t, err)

	enrollments, err := EnrollmentListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, enrollments, 2)
	expected := &model.Enrollment{
		Username: "alice",
		CourseID: "course-v1:demo+cs101+2026",
		Mode:     "audit",
		IsActive: true,
	}
	assert.Equal(t, expected, enrollments[0])
}

func TestEnrollmentFromJSONInvalidInput(t *testing.T) {
	_, err := EnrollmentListFromJSON([]byte("not json"))
	assert.Error(t, err)
}

func TestEnrollmentFromJSONEmptyObject(t *testing.T) {
	result, err := EnrollmentFromJSON([]byte("{}"))
	require.NoError(t, err)
	assert.Equal(t, &model.Enrollment{IsActive: true}, result)
}
