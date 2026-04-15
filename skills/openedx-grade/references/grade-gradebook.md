# openedx gradebook get

Get the complete gradebook for a course containing all student grades.

## Usage

```bash
openedx gradebook get --course-id <course_id>
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --course-id | Yes | — | Course identifier |

## Examples

### Get full gradebook
```bash
openedx gradebook get --course-id "course-v1:MyOrg+CS101+2026_Spring"
```

### Save gradebook to file
```bash
openedx gradebook get --course-id "course-v1:MyOrg+CS101+2026_Spring" > gradebook.json
```

## Output

Gradebook object:

```json
{
  "course_id": "course-v1:MyOrg+CS101+2026_Spring",
  "grades": [
    {
      "username": "alice",
      "course_id": "course-v1:MyOrg+CS101+2026_Spring",
      "percent": 0.85,
      "letter_grade": "B+",
      "passed": true,
      "section": ""
    }
  ]
}
```

## Notes

- Returns all students' grades in a single response
- Suitable for export and analysis — redirect to a file for storage
- Use `jq` for analysis: `openedx gradebook get --course-id <id> | jq '.grades | length'` to count graded students
```