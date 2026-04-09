// Package registry defines the built-in command registry that maps CLI command
// keys to their backend Open edX API endpoints. The registry is versioned with
// the CLI code and is not user-configurable.
package registry

// CommandMeta describes a single CLI command's backend API mapping.
type CommandMeta struct {
	Key          string   // e.g. "course.list"
	Method       string   // HTTP method: GET, POST, etc.
	Path         string   // URL path template, e.g. "/api/courses/v1/courses"
	RequiredArgs []string // required argument names, e.g. ["course_id"]
	OutputModel  string   // name of the output resource model, e.g. "Course"
}
