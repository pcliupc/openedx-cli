package normalize

import (
	"os"
	"testing"

	"github.com/openedx/cli/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeCourseFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/course_get.json")
	require.NoError(t, err)

	course, err := CourseFromJSON(data)
	require.NoError(t, err)

	expected := &model.Course{
		CourseID: "course-v1:demo+cs101+2026",
		Org:      "demo",
		Number:   "cs101",
		Run:      "2026",
		Title:    "Intro to AI",
		Pacing:   "instructor",
		Start:    "2026-01-01T00:00:00Z",
		End:      "2026-06-30T23:59:59Z",
	}
	assert.Equal(t, expected, course)
}

func TestNormalizeCourseFromExtensionPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/extension/course_list.json")
	require.NoError(t, err)

	courses, err := CourseListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, courses, 1)
	expected := &model.Course{
		CourseID: "course-v1:demo+cs101+2026",
		Org:      "demo",
		Number:   "cs101",
		Run:      "2026",
		Title:    "Intro to AI",
		Status:   "active",
	}
	assert.Equal(t, expected, courses[0])
}

func TestNormalizeCourseListFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/course_list.json")
	require.NoError(t, err)

	courses, err := CourseListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, courses, 1)
	expected := &model.Course{
		CourseID: "course-v1:demo+cs101+2026",
		Org:      "demo",
		Number:   "cs101",
		Run:      "2026",
		Title:    "Intro to AI",
		Pacing:   "instructor",
		Start:    "2026-01-01T00:00:00Z",
		End:      "2026-06-30T23:59:59Z",
	}
	assert.Equal(t, expected, courses[0])
}

func TestCourseFromJSONInvalidInput(t *testing.T) {
	_, err := CourseFromJSON([]byte("not json"))
	assert.Error(t, err)
}

func TestCourseListFromJSONInvalidInput(t *testing.T) {
	_, err := CourseListFromJSON([]byte("not json"))
	assert.Error(t, err)
}

func TestCourseFromJSONEmptyObject(t *testing.T) {
	course, err := CourseFromJSON([]byte("{}"))
	require.NoError(t, err)
	assert.Equal(t, &model.Course{}, course)
}
