# openedx course list

List all courses on the platform.

## Usage

```bash
openedx course list [--page <n>] [--page-size <n>] [--all]
```

## Parameters

| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| --page | No | 1 | Page number for pagination |
| --page-size | No | 50 | Results per page |
| --all | No | false | Fetch all pages automatically |

## Examples

### List courses (first page, JSON)
```bash
openedx course list
```

### Fetch all courses across pages
```bash
openedx course list --all
```

### Paginate manually
```bash
openedx course list --page 2 --page-size 20
```

## Output

JSON array of Course objects:

```json
[
  {
    "course_id": "course-v1:MyOrg+CS101+2026_Spring",
    "org": "MyOrg",
    "number": "CS101",
    "run": "2026_Spring",
    "title": "Intro to CS",
    "pacing": "self",
    "start": "2026-02-01T00:00:00Z",
    "end": "2026-06-30T23:59:59Z",
    "status": "active"
  }
]
```

## Notes

- Default output is JSON to stdout for pipeline processing
- Use `--all` to avoid manual pagination when processing all courses
- Combine with `jq` for filtering: `openedx course list --all | jq '.[] | select(.org == "MyOrg")'`