# openedx enrollment list

List enrollments, optionally filtered by course or user.

## Usage

```bash
openedx enrollment list [--course-id <id>] [--username <user>] [--page <n>] [--page-size <n>]
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --course-id | No | — | Filter by course identifier |
| --username | No | — | Filter by username |
| --page | No | 1 | Page number |
| --page-size | No | 50 | Results per page |

## Examples

### List all enrollments
```bash
openedx enrollment list
```

### List enrollments for a course
```bash
openedx enrollment list --course-id "course-v1:MyOrg+CS101+2026_Spring"
```

### Check a specific user's enrollments
```bash
openedx enrollment list --username alice
```

### Check if a user is enrolled in a course
```bash
openedx enrollment list --course-id "course-v1:MyOrg+CS101+2026_Spring" --username alice
```

## Output

JSON array of Enrollment objects:

```json
[
  {
    "username": "alice",
    "course_id": "course-v1:MyOrg+CS101+2026_Spring",
    "mode": "audit",
    "is_active": true
  }
]
```

## Notes

- At least one filter (--course-id or --username) is recommended to avoid large result sets
- Combine with `jq` to check specific enrollment status: `openedx enrollment list --username alice | jq '.[] | select(.is_active == true)'`