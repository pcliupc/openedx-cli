# openedx role assign

Assign a course-level role to a user.

## Usage

```bash
openedx role assign --course-id <course_id> --username <user> --role <role>
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --course-id | Yes | — | Course identifier |
| --username | Yes | — | Username to assign the role to |
| --role | Yes | — | Role to assign: instructor or staff |

## Examples

### Assign instructor role
```bash
openedx role assign --course-id "course-v1:MyOrg+CS101+2026_Spring" --username prof_smith --role instructor
```

### Assign staff role
```bash
openedx role assign --course-id "course-v1:MyOrg+CS101+2026_Spring" --username ta_jones --role staff
```

## Output

Raw JSON response from the course roles API.

## Notes

- Common roles: `instructor` (full course control) and `staff` (course management without some admin actions)
- The user must already exist on the platform (create with `openedx user create` if needed)
- Role changes take effect immediately