package registry

// LatestRegistry returns the built-in command map for the latest official
// Open edX public APIs. This map is versioned with the CLI code and is not
// user-configurable; extension mappings are handled separately by the
// provider layer.
func LatestRegistry() map[string]CommandMeta {
	return map[string]CommandMeta{
		"course.list": {
			Key:          "course.list",
			Method:       "GET",
			Path:         "/api/courses/v1/courses",
			RequiredArgs: []string{},
			OutputModel:  "Course",
		},
		"course.get": {
			Key:          "course.get",
			Method:       "GET",
			Path:         "/api/courses/v1/courses/{course_id}",
			RequiredArgs: []string{"course_id"},
			OutputModel:  "Course",
		},
		"course.create": {
			Key:          "course.create",
			Method:       "POST",
			Path:         "/api/courses/v1/courses",
			RequiredArgs: []string{"org", "number", "run", "title"},
			OutputModel:  "Course",
		},
		"course.import": {
			Key:          "course.import",
			Method:       "POST",
			Path:         "/api/courses/v1/courses/{course_id}/import",
			RequiredArgs: []string{"course_id", "file"},
			OutputModel:  "Job",
		},
		"course.export": {
			Key:          "course.export",
			Method:       "POST",
			Path:         "/api/courses/v1/courses/{course_id}/export",
			RequiredArgs: []string{"course_id"},
			OutputModel:  "Job",
		},
		"course.outline.get": {
			Key:          "course.outline.get",
			Method:       "GET",
			Path:         "/api/courses/v2/blocks/",
			RequiredArgs: []string{"course_id"},
			OutputModel:  "CourseOutline",
		},
		"user.create": {
			Key:          "user.create",
			Method:       "POST",
			Path:         "/api/user/v1/accounts",
			RequiredArgs: []string{"username", "email"},
			OutputModel:  "User",
		},
		"enrollment.add": {
			Key:          "enrollment.add",
			Method:       "POST",
			Path:         "/api/enrollment/v1/enrollments",
			RequiredArgs: []string{"course_id", "username"},
			OutputModel:  "Enrollment",
		},
		"role.assign": {
			Key:          "role.assign",
			Method:       "POST",
			Path:         "/api/courses/v1/courses/{course_id}/roles",
			RequiredArgs: []string{"course_id", "username", "role"},
			OutputModel:  "RoleAssignment",
		},
		"enrollment.list": {
			Key:          "enrollment.list",
			Method:       "GET",
			Path:         "/api/enrollment/v1/enrollments",
			RequiredArgs: []string{},
			OutputModel:  "Enrollment",
		},
		"enrollment.remove": {
			Key:          "enrollment.remove",
			Method:       "POST",
			Path:         "/api/enrollment/v1/enrollments",
			RequiredArgs: []string{"course_id", "username"},
			OutputModel:  "Enrollment",
		},
		"user.list": {
			Key:          "user.list",
			Method:       "GET",
			Path:         "/api/user/v1/accounts",
			RequiredArgs: []string{},
			OutputModel:  "User",
		},
		"user.get": {
			Key:          "user.get",
			Method:       "GET",
			Path:         "/api/user/v1/accounts/{username}",
			RequiredArgs: []string{"username"},
			OutputModel:  "User",
		},
		"grade.list": {
			Key:          "grade.list",
			Method:       "GET",
			Path:         "/api/grades/v1/courses/{course_id}/",
			RequiredArgs: []string{"course_id"},
			OutputModel:  "Grade",
		},
		"gradebook.get": {
			Key:          "gradebook.get",
			Method:       "GET",
			Path:         "/api/grades/v1/gradebook/{course_id}/",
			RequiredArgs: []string{"course_id"},
			OutputModel:  "Gradebook",
		},
		"certificate.list": {
			Key:          "certificate.list",
			Method:       "GET",
			Path:         "/api/certificates/v0/certificates/{username}/",
			RequiredArgs: []string{"username"},
			OutputModel:  "Certificate",
		},
	}
}
