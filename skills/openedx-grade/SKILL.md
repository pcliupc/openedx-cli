---
name: openedx-grade
version: 1.0.0
description: "成绩与证书：成绩查询、成绩册导出、证书列表"
metadata:
  requires:
    skills: ["openedx-shared"]
    bins: ["openedx"]
  cliHelp: "openedx grade --help"
---

## Prerequisites

MUST read `../openedx-shared/SKILL.md` first to understand configuration, authentication, and security rules.

## Core Concepts

- **Grade** — a student's grade in a course section. Includes percent score, letter grade, pass/fail status, and section breakdown.
- **Gradebook** — a complete grade overview for a course. Contains an array of all student grades.
- **Certificate** — a completion certificate for a user. Has a type (verified, honor, etc.), status, download URL, and associated grade.

## Resource Relationships

```
Course (course_id)
├── grades[] → Grade
│   └── (username, course_id, percent, letter_grade, passed, section)
├── gradebook → Gradebook
│   └── grades[] → Grade
└── certificates[] → Certificate
    └── (username, course_id, certificate_type, status, download_url, grade)
```

## Commands

| Command | Description | Reference |
|---------|-------------|-----------|
| `openedx grade list` | List grades for a course | [grade-list.md](references/grade-list.md) |
| `openedx gradebook get` | Get full gradebook for a course | [grade-gradebook.md](references/grade-gradebook.md) |
| `openedx certificate list` | List certificates for a user | [certificate-list.md](references/certificate-list.md) |

## Common Workflows

### Query course grades for a specific student

```bash
openedx grade list --course-id "course-v1:MyOrg+CS101+2026_Spring" --username alice
```

### Export full gradebook

```bash
openedx gradebook get --course-id "course-v1:MyOrg+CS101+2026_Spring" > gradebook.json
```

### Check certificate status for a student

```bash
openedx certificate list --username alice
```

### Full audit: grades + certificates for a course

```bash
# Get all grades
openedx grade list --course-id "course-v1:MyOrg+CS101+2026_Spring"

# Check certificates for a specific student
openedx certificate list --username alice
```

## Error Handling

| Error | Cause | Resolution |
|-------|-------|------------|
| No grade data | Course has no graded content or no enrollments | Verify course has graded sections and enrolled students |
| Certificate not generated | Student hasn't met certificate requirements | Check grade status with `openedx grade list` first |
| Course not found (404) | Invalid course_id | Verify with `openedx course list` |
| Permission denied (403) | Insufficient role | Ensure staff or instructor role on the course |
```