package normalize

import (
	"encoding/json"

	"github.com/openedx/cli/internal/model"
)

// rawBlock represents a single block in the Open edX Blocks API flat-map response.
type rawBlock struct {
	ID          string   `json:"id"`
	BlockID     string   `json:"block_id"`
	DisplayName string   `json:"display_name"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Children    []string `json:"children"`
}

// rawBlocksAPIResponse represents the Open edX Blocks API response shape.
// The API returns a flat map of blocks with string IDs as keys and a "root"
// field pointing to the root block ID.
type rawBlocksAPIResponse struct {
	Root   string              `json:"root"`
	Blocks map[string]rawBlock `json:"blocks"`
}

// rawOutlineBlock is a lax intermediate struct for extension API payloads
// that already provide a nested tree structure.
type rawOutlineBlock struct {
	BlockID     string            `json:"block_id"`
	ID          string            `json:"id"`
	DisplayName string            `json:"display_name"`
	Title       string            `json:"title"`
	Type        string            `json:"type"`
	Children    []rawOutlineBlock `json:"children"`
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

// buildTree recursively reconstructs the block tree from the flat blocks map.
func buildTree(blockID string, blocks map[string]rawBlock) model.OutlineBlock {
	block, ok := blocks[blockID]
	if !ok {
		return model.OutlineBlock{ID: blockID}
	}

	id := firstNonEmpty(block.BlockID, block.ID)
	title := firstNonEmpty(block.DisplayName, block.Title)

	var children []model.OutlineBlock
	for _, childID := range block.Children {
		children = append(children, buildTree(childID, blocks))
	}

	return model.OutlineBlock{
		ID:       id,
		Title:    title,
		Type:     block.Type,
		Children: children,
	}
}

// OutlineFromJSON parses a course outline from raw JSON bytes. It supports
// three shapes:
//  1. The Open edX Blocks API response: {"root": "...", "blocks": {...}} — a flat
//     map of blocks that is reconstructed into a tree.
//  2. The extension API shape: {"course_id": "...", "chapters": [...]} — already
//     a nested tree.
func OutlineFromJSON(data []byte) (*model.CourseOutline, error) {
	// Try Blocks API shape first (flat map with root reference).
	var blocksResp rawBlocksAPIResponse
	if err := json.Unmarshal(data, &blocksResp); err == nil && len(blocksResp.Blocks) > 0 && blocksResp.Root != "" {
		root := blocksResp.Blocks[blocksResp.Root]
		courseID := firstNonEmpty(root.BlockID, root.ID)

		var chapters []model.OutlineBlock
		for _, childID := range root.Children {
			chapters = append(chapters, buildTree(childID, blocksResp.Blocks))
		}

		return &model.CourseOutline{
			CourseID: courseID,
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
