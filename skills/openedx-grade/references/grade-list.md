# openedx grade list

List grades for a course, optionally filtered by student.

## Usage

```bash
openedx grade list --course-id <course_id> [--username <user>] [--page <n>] [--page-size <n>]
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --course-id | Yes | — | Course identifier |
| --username | No | — | Filter by student username |
| --page | No | 1 | Page number |
| --page-size | No | 50 | Results per page |

## Examples

### List all grades for a course
```bash
openedx grade list --course-id "course-v1:MyOrg+CS101+2026_Spring"
```

### Get grades for a specific student
```bash
openedx grade list --course-id "course-v1:MyOrg+CS101+2026_Spring" --username alice
```

## Output

JSON array of Grade objects:

```json
[
  {
    "username": "alice",
    "course_id": "course-v1:MyOrg+CS101+2026_Spring",
    "percent": 0.85,
    "letter_grade": "B+",
    "passed": true,
    "section": "Week 1: Introduction"
  }
]
```

## Notes

- `--course-id` is required to specify which course's grades to query
- Use `--username` to narrow results to a single student
- `percent` is a decimal (0.85 = 85%)
- `passed` indicates whether the student met the passing threshold
```