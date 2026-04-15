# openedx user get

Get details for a specific user.

## Usage

```bash
openedx user get --username <username>
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --username | Yes | — | Username to look up |

## Examples

### Get user details
```bash
openedx user get --username alice
```

### Get user on a specific profile
```bash
openedx user get --username alice --profile staging
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

- Returns a 404 error if the username does not exist
- Use this to verify a user exists before enrolling or assigning roles
