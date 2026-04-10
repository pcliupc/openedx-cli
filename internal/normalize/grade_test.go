package normalize

import (
	"os"
	"testing"

	"github.com/openedx/cli/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGradeListFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/grade_list.json")
	require.NoError(t, err)

	grades, err := GradeListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, grades, 2)
	expected := &model.Grade{
		Username:    "alice",
		CourseID:    "course-v1:demo+cs101+2026",
		Percent:     0.85,
		LetterGrade: "B+",
		Passed:      true,
		Section:     "Week 1",
	}
	assert.Equal(t, expected, grades[0])
}

func TestGradeListFromExtensionPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/extension/grade_list.json")
	require.NoError(t, err)

	grades, err := GradeListFromJSON(data)
	require.NoError(t, err)

	require.Len(t, grades, 1)
	expected := &model.Grade{
		Username:    "alice",
		CourseID:    "course-v1:demo+cs101+2026",
		Percent:     0.85,
		LetterGrade: "B+",
		Passed:      true,
	}
	assert.Equal(t, expected, grades[0])
}

func TestGradebookFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/gradebook_get.json")
	require.NoError(t, err)

	gradebook, err := GradebookFromJSON(data)
	require.NoError(t, err)

	assert.Equal(t, "course-v1:demo+cs101+2026", gradebook.CourseID)
	require.Len(t, gradebook.Grades, 2)
	assert.Equal(t, "alice", gradebook.Grades[0].Username)
	assert.Equal(t, 0.85, gradebook.Grades[0].Percent)
}

func TestGradeListFromJSONInvalidInput(t *testing.T) {
	_, err := GradeListFromJSON([]byte("not json"))
	assert.Error(t, err)
}

func TestGradebookFromJSONInvalidInput(t *testing.T) {
	_, err := GradebookFromJSON([]byte("not json"))
	assert.Error(t, err)
}
