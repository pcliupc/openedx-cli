# openedx enrollment remove

Remove a user's enrollment from a course.

## Usage

```bash
openedx enrollment remove --course-id <course_id> --username <user>
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --course-id | Yes | — | Course identifier |
| --username | Yes | — | Username to unenroll |

## Examples

### Remove a user's enrollment
```bash
openedx enrollment remove --course-id "course-v1:MyOrg+CS101+2026_Spring" --username alice
```

## Output

Raw JSON response from the enrollment API.

## Notes

- This action sets `is_active` to false (may not fully delete the enrollment record depending on platform configuration)
- Verify removal with `openedx enrollment list --course-id <id> --username <user>`