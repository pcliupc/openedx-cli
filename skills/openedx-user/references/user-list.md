# openedx user list

List users on the platform.

## Usage

```bash
openedx user list [--page <n>] [--page-size <n>]
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --page | No | 1 | Page number |
| --page-size | No | 50 | Results per page |

## Examples

### List users (first page)
```bash
openedx user list
```

### Paginate manually
```bash
openedx user list --page 2 --page-size 20
```

## Output

JSON array of User objects:

```json
[
  {
    "username": "alice",
    "email": "alice@example.com",
    "name": "Alice Smith",
    "is_active": true,
    "created_at": "2026-04-15T10:00:00Z"
  }
]
```

## Notes

- Use `--page-size` to control result count per page
- Combine with `jq` for filtering: `openedx user list | jq '.[] | select(.is_active == true)'`
