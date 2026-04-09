package normalize

import (
	"encoding/json"

	"github.com/openedx/cli/internal/model"
)

// rawOutlineBlock is a lax intermediate struct that accepts field names from
// both public API (block_id, display_name) and extension API (id, title) shapes.
type rawOutlineBlock struct {
	BlockID     string            `json:"block_id"`
	ID          string            `json:"id"`
	DisplayName string            `json:"display_name"`
	Title       string            `json:"title"`
	Type        string            `json:"type"`
	Children    []rawOutlineBlock `json:"children"`
}

// rawPublicOutline represents the public API outline response shape.
type rawPublicOutline struct {
	CourseID string            `json:"course_id"`
	Blocks   []rawOutlineBlock `json:"blocks"`
}

// rawExtensionOutline represents the extension API outline response shape.
type rawExtensionOutline struct {
	CourseID string            `json:"course_id"`
	Chapters []rawOutlineBlock `json:"chapters"`
}

func (r rawOutlineBlock) toModel() model.OutlineBlock {
	id := firstNonEmpty(r.BlockID, r.ID)
	title := firstNonEmpty(r.DisplayName, r.Title)
	var children []model.OutlineBlock
	for _, c := range r.Children {
		children = append(children, c.toModel())
	}
	return model.OutlineBlock{
		ID:       id,
		Title:    title,
		Type:     r.Type,
		Children: children,
	}
}

// OutlineFromJSON parses a course outline from raw JSON bytes. It supports the
// public API shape {"course_id": "...", "blocks": [...]} where top-level blocks
// are treated as chapters, and the extension API shape {"course_id": "...",
// "chapters": [...]}.
func OutlineFromJSON(data []byte) (*model.CourseOutline, error) {
	// Try public API shape first (blocks field).
	var pub rawPublicOutline
	if err := json.Unmarshal(data, &pub); err == nil && len(pub.Blocks) > 0 {
		chapters := make([]model.OutlineBlock, 0, len(pub.Blocks))
		for _, b := range pub.Blocks {
			chapters = append(chapters, b.toModel())
		}
		return &model.CourseOutline{
			CourseID: pub.CourseID,
			Chapters: chapters,
		}, nil
	}

	// Try extension API shape (chapters field).
	var ext rawExtensionOutline
	if err := json.Unmarshal(data, &ext); err != nil {
		return nil, err
	}
	chapters := make([]model.OutlineBlock, 0, len(ext.Chapters))
	for _, c := range ext.Chapters {
		chapters = append(chapters, c.toModel())
	}
	return &model.CourseOutline{
		CourseID: ext.CourseID,
		Chapters: chapters,
	}, nil
}
