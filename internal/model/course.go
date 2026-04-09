// Package model defines stable output structures for the CLI.
// These types represent the canonical shape that all provider responses
// are normalized into, insulating the command surface from API differences.
package model

// Course represents a normalized OpenEdX course resource.
type Course struct {
	CourseID string `json:"course_id"`
	Org      string `json:"org"`
	Number   string `json:"number"`
	Run      string `json:"run"`
	Title    string `json:"title"`
	Pacing   string `json:"pacing,omitempty"`
	Start    string `json:"start,omitempty"`
	End      string `json:"end,omitempty"`
	Status   string `json:"status,omitempty"`
}
