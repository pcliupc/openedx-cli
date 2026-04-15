# openedx user create

Create a new user on the platform.

## Usage

```bash
openedx user create --username <username> --email <email> [--name <full_name>]
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --username | Yes | — | Unique username for the new user |
| --email | Yes | — | Email address for the new user |
| --name | No | — | Full display name |

## Examples

### Create a user with all fields
```bash
openedx user create --username alice --email alice@example.com --name "Alice Smith"
```

### Create with username and email only
```bash
openedx user create --username bob --email bob@example.com
```

## Output

User object:

```json
{
  "username": "alice",
  "email": "alice@example.com",
  "name": "Alice Smith",
  "is_active": true,
  "created_at": "2026-04-15T10:00:00Z"
}
```

## Notes

- Username must be unique across the platform
- The new user is active by default (`is_active: true`)
- Save the username for subsequent operations (enrollment, role assignment)
