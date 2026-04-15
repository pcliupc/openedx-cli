# openedx certificate list

List certificates earned by a user.

## Usage

```bash
openedx certificate list --username <user> [--page <n>] [--page-size <n>]
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --username | Yes | — | Username to look up certificates for |
| --page | No | 1 | Page number |
| --page-size | No | 50 | Results per page |

## Examples

### List all certificates for a user
```bash
openedx certificate list --username alice
```

### Paginate results
```bash
openedx certificate list --username alice --page-size 10
```

## Output

JSON array of Certificate objects:

```json
[
  {
    "username": "alice",
    "course_id": "course-v1:MyOrg+CS101+2026_Spring",
    "certificate_type": "verified",
    "status": "downloadable",
    "download_url": "https://openedx.example.com/certificates/abc123",
    "grade": "B+"
  }
]
```

## Notes

- `--username` is required
- `status` can be: `downloadable`, `generating`, `notpassing`, `unverified`
- `download_url` is only populated when status is `downloadable`
- Use this to verify a student has earned a certificate
```