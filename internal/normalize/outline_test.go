package normalize

import (
	"os"
	"testing"

	"github.com/openedx/cli/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeOutlineFromPublicPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/outline.json")
	require.NoError(t, err)

	outline, err := OutlineFromJSON(data)
	require.NoError(t, err)

	expected := &model.CourseOutline{
		CourseID: "block-v1:demo+cs101+2026+type@course+block@course",
		Chapters: []model.OutlineBlock{
			{
				ID:    "block-v1:demo+cs101+2026+type@chapter+block@ch1",
				Title: "Chapter 1: Introduction",
				Type:  "chapter",
				Children: []model.OutlineBlock{
					{
						ID:    "block-v1:demo+cs101+2026+type@sequential+block@seq1",
						Title: "Lesson 1: Welcome",
						Type:  "sequential",
						Children: []model.OutlineBlock{
							{
								ID:    "block-v1:demo+cs101+2026+type@vertical+block@vert1",
								Title: "Welcome Video",
								Type:  "vertical",
							},
						},
					},
				},
			},
			{
				ID:       "block-v1:demo+cs101+2026+type@chapter+block@ch2",
				Title:    "Chapter 2: Basics",
				Type:     "chapter",
				Children: nil,
			},
		},
	}
	assert.Equal(t, expected, outline)
}

func TestNormalizeOutlineFromExtensionPayload(t *testing.T) {
	data, err := os.ReadFile("../../testdata/extension/outline.json")
	require.NoError(t, err)

	outline, err := OutlineFromJSON(data)
	require.NoError(t, err)

	expected := &model.CourseOutline{
		CourseID: "course-v1:demo+cs101+2026",
		Chapters: []model.OutlineBlock{
			{
				ID:    "block-v1:demo+cs101+2026+type@chapter+block@ch1",
				Title: "Chapter 1: Introduction",
				Type:  "chapter",
				Children: []model.OutlineBlock{
					{
						ID:    "block-v1:demo+cs101+2026+type@sequential+block@seq1",
						Title: "Lesson 1: Welcome",
						Type:  "sequential",
						Children: []model.OutlineBlock{
							{
								ID:    "block-v1:demo+cs101+2026+type@vertical+block@vert1",
								Title: "Welcome Video",
								Type:  "vertical",
							},
						},
					},
				},
			},
			{
				ID:       "block-v1:demo+cs101+2026+type@chapter+block@ch2",
				Title:    "Chapter 2: Basics",
				Type:     "chapter",
				Children: nil,
			},
		},
	}
	assert.Equal(t, expected, outline)
}

func TestNormalizeOutlineProducesStableTree(t *testing.T) {
	data, err := os.ReadFile("../../testdata/public/outline.json")
	require.NoError(t, err)

	outline, err := OutlineFromJSON(data)
	require.NoError(t, err)

	// Verify the tree structure is correctly reconstructed from flat blocks map.
	require.Len(t, outline.Chapters, 2)

	ch1 := outline.Chapters[0]
	assert.Equal(t, "chapter", ch1.Type)
	require.Len(t, ch1.Children, 1)

	seq1 := ch1.Children[0]
	assert.Equal(t, "sequential", seq1.Type)
	require.Len(t, seq1.Children, 1)

	vert1 := seq1.Children[0]
	assert.Equal(t, "vertical", vert1.Type)
	assert.Empty(t, vert1.Children)
}

func TestOutlineFromJSONInvalidInput(t *testing.T) {
	_, err := OutlineFromJSON([]byte("not json"))
	assert.Error(t, err)
}
