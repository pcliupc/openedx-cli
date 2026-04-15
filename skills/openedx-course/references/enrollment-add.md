# openedx enrollment add

Enroll a user in a course.

## Usage

```bash
openedx enrollment add --course-id <course_id> --username <user> [--mode <mode>]
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --course-id | Yes | — | Course identifier |
| --username | Yes | — | Username to enroll |
| --mode | No | audit | Enrollment mode: audit, verified, professional, no-id-professional, honor, credit |

## Examples

### Enroll a student (audit mode)
```bash
openedx enrollment add --course-id "course-v1:MyOrg+CS101+2026_Spring" --username alice
```

### Enroll with verified mode
```bash
openedx enrollment add --course-id "course-v1:MyOrg+CS101+2026_Spring" --username bob --mode verified
```

### Enroll in a professional course
```bash
openedx enrollment add --course-id "course-v1:MyOrg+CS101+2026_Spring" --username charlie --mode professional
```

## Output

Raw JSON response from the enrollment API.

## Notes

- Default mode is `audit` if `--mode` is not specified
- If the user is already enrolled, this may return an error or update the mode depending on the platform configuration
- Verify enrollment with `openedx enrollment list --course-id <id> --username <user>`